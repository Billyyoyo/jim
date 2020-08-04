package tests

import (
	"fmt"
	"github.com/panjf2000/gnet"
	log "github.com/sirupsen/logrus"
	"jim/common/tool"
	"testing"
	"time"
)

var (
	socketServer *TestServer
)

type TestServer struct {
	*gnet.EventServer
}

func (cs *TestServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Info(fmt.Sprintf("Tcp server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop))
	return
}

func (cs *TestServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Info("someone connected:", c.RemoteAddr())
	return
}

func (cs *TestServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	data := append([]byte{}, frame...)
	tool.AsyncRun(func() {
		c.AsyncWrite(data)
	})
	return
}

func (cs *TestServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	log.Info("someone closed:", c.RemoteAddr())
	return
}

func TestSocketServer(t *testing.T){
	addr := fmt.Sprintf("tcp://%s:%d", "localhost", 4009)
	socketServer = &TestServer{}
	err := gnet.Serve(socketServer, addr,
		gnet.WithMulticore(true),
		gnet.WithTCPKeepAlive(time.Minute*5),
		//gnet.WithCodec(core.NewJimDataFrameCodec()),
		gnet.WithTicker(false))
	if err != nil {
		panic(err)
	}
}