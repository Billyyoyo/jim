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

var (
	socketServer *TcpServer
)

type event struct {
	action int
	data   interface{}
}

type ConnData struct {
	Uid      int64  // userId
	Did      int64  // deviceId
	Serial   string //device serial
	PongTime int64  //ping之后客户端发送pong消息记录收到时间
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
			cs.handlePong(c)
		} else if inPack.Type == rpc.PackType_PT_PING {
			cs.handlePing(c)
		} else if inPack.Type == rpc.PackType_PT_AUTH {
			regInfo := &rpc.RegInfo{}
			err = proto.Unmarshal(inPack.Data, regInfo)
			if err != nil {
				log.Error("register connection but parse data err:", err.Error())
				c.Close()
				return
			}
			err = cs.handleReg(c, regInfo)
			if err != nil {
				log.Error("register connection - reg failed:", err.Error())
				c.Close()
				return
			}
		} else if inPack.Type == rpc.PackType_PT_MSG {
			msg := &rpc.Message{}
			err = proto.Unmarshal(inPack.Data, msg)
			cs.handleMsg(&c, msg)
		} else if inPack.Type == rpc.PackType_PT_ACTION {
			act := &rpc.Action{}
			err = proto.Unmarshal(inPack.Data, act)
			cs.handleAct(&c, act)
		} else if inPack.Type == rpc.PackType_PT_ACK {
			ack := &rpc.Ack{}
			err = proto.Unmarshal(inPack.Data, ack)
			cs.handleAck(&c, ack)
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
	log.Info("someone connected:", c.RemoteAddr())
	// 作为无状态连接保存
	c.SetContext(ConnData{PongTime: time.Now().Unix(),})
	cs.uconns.Store(c.RemoteAddr().String(), c)
	return
}

func (cs *TcpServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	// todo 从Context()取出deviceId   再到uconns中取得用户信息  否则过调
	// todo 当连接关闭 调用logic的offline
	log.Info("someone connection closed ", c.RemoteAddr())
	cs.eventQueue <- event{action: EVENT_TYPE_KICKOFF, data: c.RemoteAddr().String()}
	if c.Context() != nil {
		data := c.Context().(ConnData)
		if data.Uid > 0 {
			tool.AsyncRun(func() {
				cs.handleOffline(&data)
			})
		}
	}
	return
}

func (cs *TcpServer) Tick() (delay time.Duration, action gnet.Action) {
	// 给所有连接发送心跳包
	cs.uconns.Range(func(k, v interface{}) bool {
		conn := v.(gnet.Conn)
		// 检查超时
		cData := conn.Context().(ConnData)
		if cData.Uid == 0 {
			if time.Now().Unix()-cData.PongTime > 5 {
				log.Info("this client register timeout, so kickoff", cData)
				conn.Close()
			}
			return true
		}
		if time.Now().Unix()-cData.PongTime > core.AppConfig.Socket.Timeout {
			log.Info("someone connect timeout, so kickoff")
			conn.Close()
			return true
		}
		pingPack := &rpc.Output{Type: rpc.PackType_PT_PING,}
		bs, err := proto.Marshal(pingPack)
		if err != nil {
			return true
		}
		conn.AsyncWrite(bs)
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
				remoteAddr := ev.data.(string)
				log.Info("delete ", remoteAddr, " from server")
				cs.uconns.Delete(remoteAddr)
			}
			break
		}
	}
}

func GetUserConn(remoteAddr string) *gnet.Conn {
	val, ok := socketServer.uconns.Load(remoteAddr)
	if ok {
		conn := val.(gnet.Conn)
		return &conn
	} else {
		return nil
	}
}

func StartUpSocketServer() {
	addr := fmt.Sprintf("tcp://%s:%d", core.AppConfig.Socket.Host, core.AppConfig.Socket.Port)
	socketServer = &TcpServer{
		addr:       addr,
		tick:       time.Second * time.Duration(core.AppConfig.Socket.Tick),
		eventQueue: make(chan event, 100),
	}
	go socketServer.watchEvent()
	// 使用自定义的编码解码
	err := gnet.Serve(socketServer, addr,
		gnet.WithMulticore(true),
		gnet.WithTCPKeepAlive(time.Minute*5),
		gnet.WithTicker(true),
		gnet.WithCodec(core.NewJimDataFrameCodec()))
	if err != nil {
		panic(err)
	}
}
