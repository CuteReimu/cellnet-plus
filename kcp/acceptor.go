package kcp

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/util"
	"github.com/xtaci/kcp-go/v5"
	"net"
	"strings"
)

// 接受器
type kcpAcceptor struct {
	peer.SessionManager
	peer.CorePeerProperty
	peer.CoreContextSet
	peer.CoreRunningTag
	peer.CoreProcBundle
	peer.CoreTCPSocketOption
	peer.CoreCaptureIOPanic

	// 保存侦听器
	listener *kcp.Listener
}

func (self *kcpAcceptor) Port() int {
	if self.listener == nil {
		return 0
	}

	return self.listener.Addr().(*net.UDPAddr).Port
}

func (self *kcpAcceptor) IsReady() bool {

	return self.IsRunning()
}

// 异步开始侦听
func (self *kcpAcceptor) Start() cellnet.Peer {

	self.WaitStopFinished()

	if self.IsRunning() {
		return self
	}

	ln, err := util.DetectPort(self.Address(), func(a *util.Address, port int) (interface{}, error) {
		return kcp.ListenWithOptions(a.HostPortString(port), blockCrypto, 10, 3)
	})

	if err != nil {

		log.Errorf("#kcp.listen failed(%s) %v", self.Name(), err.Error())

		self.SetRunning(false)

		return self
	}

	self.listener = ln.(*kcp.Listener)

	log.Infof("#kcp.listen(%s) %s", self.Name(), self.ListenAddress())

	go self.accept()

	return self
}

func (self *kcpAcceptor) ListenAddress() string {

	pos := strings.Index(self.Address(), ":")
	if pos == -1 {
		return self.Address()
	}

	host := self.Address()[:pos]

	return util.JoinAddress(host, self.Port())
}

func (self *kcpAcceptor) accept() {
	self.SetRunning(true)

	for {
		conn, err := self.listener.AcceptKCP()

		if self.IsStopping() {
			break
		}

		if err != nil {

			// 调试状态时, 才打出accept的具体错误
			if log.IsDebugEnabled() {
				log.Errorf("#kcp.accept failed(%s) %v", self.Name(), err.Error())
			}

			continue
		}

		// 处理连接进入独立线程, 防止accept无法响应
		go self.onNewSession(conn)

	}

	self.SetRunning(false)

	self.EndStopping()

}

func (self *kcpAcceptor) onNewSession(conn net.Conn) {

	self.ApplySocketOption(conn)

	ses := newSession(conn, self, nil)

	ses.Start()

	self.ProcEvent(&cellnet.RecvMsgEvent{
		Ses: ses,
		Msg: &cellnet.SessionAccepted{},
	})
}

// 停止侦听器
func (self *kcpAcceptor) Stop() {
	if !self.IsRunning() {
		return
	}

	if self.IsStopping() {
		return
	}

	self.StartStopping()

	_ = self.listener.Close()

	// 断开所有连接
	self.CloseAllSession()

	// 等待线程结束
	self.WaitStopFinished()
}

func (self *kcpAcceptor) TypeName() string {
	return "kcp.Acceptor"
}

func init() {

	peer.RegisterPeerCreator(func() cellnet.Peer {
		p := &kcpAcceptor{
			SessionManager: new(peer.CoreSessionManager),
		}

		p.CoreTCPSocketOption.Init()

		return p
	})
}
