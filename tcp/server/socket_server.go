package server

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/panjf2000/gnet"
	log "github.com/sirupsen/logrus"
	"jim/common/rpc"
	"jim/common/tool"
	"jim/tcp/core"
	"sync"
	"time"
)

const (
	EVENT_TYPE_KICKOFF = iota + 1
)

type event struct {
	action int
	data   interface{}
}

type UConn struct {
	Uid      int64  // userId
	Did      int64  // deviceId
	Seq      int64  //Last receive message sequence
	Serial   string //device serial
	Conn     gnet.Conn
	PongTime int64 //ping之后客户端发送pong消息记录收到时间
}

type TcpServer struct {
	*gnet.EventServer
	addr       string
	uconns     sync.Map //key是deviceId， value是UConn
	tick       time.Duration
	eventQueue chan event
}

func (cs *TcpServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Info(fmt.Sprintf("Tcp server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop))
	// todo 将接入服务注册到zk
	return
}

func (cs *TcpServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	data := append([]byte{}, frame...)
	tool.AsyncRun(func() {
		inPack := &rpc.Input{}
		err := proto.Unmarshal(data, inPack)
		if err != nil {
			return
		}
		if inPack.Type == rpc.PackType_PT_PONG {
			did := c.Context().(int64)
			obj, ok := cs.uconns.Load(did)
			if !ok {
				fmt.Println("can not find device", did)
				return
			}
			uconn := obj.(UConn)
			uconn.PongTime = time.Now().Unix()
			cs.uconns.Store(did, uconn)
		} else if inPack.Type == rpc.PackType_PT_PING {
			pingPack := &rpc.Output{Type: rpc.PackType_PT_PONG,}
			bs, err := proto.Marshal(pingPack)
			if err != nil {
				return
			}
			c.AsyncWrite(bs)
		}

	})
	return
}

func (cs *TcpServer) OnShutdown(svr gnet.Server) {
	//	todo 怎么同步最后同步消息序列号 待思考
	if cs.eventQueue != nil {
		close(cs.eventQueue)
	}
	log.Info("server stopped")
}

func (cs *TcpServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	// todo 当连接建立 调用logic的register
	// todo 同步注册成功 将deviceId放到c.SetContext()中
	uconn := UConn{
		Uid:      0,
		Did:      1,
		Seq:      0,
		Serial:   "",
		Conn:     c,
		PongTime: time.Now().Unix(),
	}
	c.SetContext(uconn.Did)
	cs.uconns.Store(uconn.Did, uconn)
	return
}

func (cs *TcpServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	// todo 从Context()取出deviceId   再到uconns中取得用户信息  否则过调
	// todo 当连接关闭 调用logic的offline
	log.Info("someone connection closed")
	return
}

func (cs *TcpServer) Tick() (delay time.Duration, action gnet.Action) {
	// todo 给所有连接发送心跳包
	cs.uconns.Range(func(k, v interface{}) bool {
		uconn := v.(UConn)
		// 检查超时
		if time.Now().Unix()-uconn.PongTime > core.AppConfig.Socket.Timeout {
			log.Info("someone connect timeout, so kickoff")
			uconn.Conn.Close()
			cs.eventQueue <- event{action: EVENT_TYPE_KICKOFF, data: uconn.Did}
			return true
		}
		pingPack := &rpc.Output{Type: rpc.PackType_PT_PING,}
		bs, err := proto.Marshal(pingPack)
		if err != nil {
			return true
		}
		uconn.Conn.AsyncWrite(bs)
		return true
	})
	delay = cs.tick
	return
}

func (cs *TcpServer) watchEvent() {
	for {
		select {
		case ev := <-cs.eventQueue:
			if ev.action == EVENT_TYPE_KICKOFF {
				did := ev.data.(int64)
				cs.uconns.Delete(did)
				log.Info("delete uconn from server")
			}
			break
		}
	}
}

func StartUpSocketServer() {
	addr := fmt.Sprintf("tcp://%s:%d", core.AppConfig.Socket.Host, core.AppConfig.Socket.Port)
	cs := &TcpServer{
		addr:       addr,
		tick:       time.Second * time.Duration(core.AppConfig.Socket.Tick),
		eventQueue: make(chan event, 100),
	}
	go cs.watchEvent()
	// 使用自定义的编码解码
	err := gnet.Serve(cs, addr,
		gnet.WithMulticore(true),
		gnet.WithTCPKeepAlive(time.Minute*5),
		gnet.WithTicker(true),
		gnet.WithCodec(core.NewJimDataFrameCodec()))
	if err != nil {
		panic(err)
	}
}
