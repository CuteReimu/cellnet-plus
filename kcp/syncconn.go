package kcp

import (
	"net"
	"time"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/xtaci/kcp-go/v5"
)

type kcpSyncConnector struct {
	peer.SessionManager

	peer.CorePeerProperty
	peer.CoreContextSet
	peer.CoreProcBundle
	peer.CoreTCPSocketOption

	defaultSes *kcpSession
}

func (self *kcpSyncConnector) Port() int {
	conn := self.defaultSes.Conn()

	if conn == nil {
		return 0
	}

	return conn.LocalAddr().(*net.UDPAddr).Port
}

func (self *kcpSyncConnector) Start() cellnet.Peer {

	// 尝试用Socket连接地址
	conn, err := kcp.DialWithOptions(self.Address(), blockCrypto, 10, 3)

	// 发生错误时退出
	if err != nil {

		log.Debugf("#tcp.connect failed(%s)@%d address: %s", self.Name(), self.defaultSes.ID(), self.Address())

		self.ProcEvent(&cellnet.RecvMsgEvent{Ses: self.defaultSes, Msg: &cellnet.SessionConnectError{}})
		return self
	}

	self.defaultSes.setConn(conn)

	self.ApplySocketOption(conn)

	self.defaultSes.Start()

	self.ProcEvent(&cellnet.RecvMsgEvent{Ses: self.defaultSes, Msg: &cellnet.SessionConnected{}})

	return self
}

func (self *kcpSyncConnector) Session() cellnet.Session {
	return self.defaultSes
}

func (self *kcpSyncConnector) SetSessionManager(raw interface{}) {
	self.SessionManager = raw.(peer.SessionManager)
}

func (self *kcpSyncConnector) ReconnectDuration() time.Duration {
	return 0
}

func (self *kcpSyncConnector) SetReconnectDuration(_ time.Duration) {

}

func (self *kcpSyncConnector) Stop() {

	if self.defaultSes != nil {
		self.defaultSes.Close()
	}

}

func (self *kcpSyncConnector) IsReady() bool {

	return self.SessionCount() != 0
}

func (self *kcpSyncConnector) TypeName() string {
	return "kcp.SyncConnector"
}

func init() {

	peer.RegisterPeerCreator(func() cellnet.Peer {
		self := &kcpSyncConnector{
			SessionManager: new(peer.CoreSessionManager),
		}

		self.defaultSes = newSession(nil, self, nil)

		self.CoreTCPSocketOption.Init()

		return self
	})
}
