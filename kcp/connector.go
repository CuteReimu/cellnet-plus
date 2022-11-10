package kcp

import (
	"net"
	"sync"
	"time"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/xtaci/kcp-go/v5"
)

type kcpConnector struct {
	peer.SessionManager

	peer.CorePeerProperty
	peer.CoreContextSet
	peer.CoreRunningTag
	peer.CoreProcBundle
	peer.CoreTCPSocketOption

	defaultSes *kcpSession

	tryConnTimes int // 尝试连接次数

	sesEndSignal sync.WaitGroup

	reconDur time.Duration
}

func (self *kcpConnector) Start() cellnet.Peer {

	self.WaitStopFinished()

	if self.IsRunning() {
		return self
	}

	go self.connect(self.Address())

	return self
}

func (self *kcpConnector) Session() cellnet.Session {
	return self.defaultSes
}

func (self *kcpConnector) SetSessionManager(raw interface{}) {
	self.SessionManager = raw.(peer.SessionManager)
}

func (self *kcpConnector) Stop() {
	if !self.IsRunning() {
		return
	}

	if self.IsStopping() {
		return
	}

	self.StartStopping()

	// 通知发送关闭
	self.defaultSes.Close()

	// 等待线程结束
	self.WaitStopFinished()

}

func (self *kcpConnector) ReconnectDuration() time.Duration {

	return self.reconDur
}

func (self *kcpConnector) SetReconnectDuration(v time.Duration) {
	self.reconDur = v
}

func (self *kcpConnector) Port() int {

	conn := self.defaultSes.Conn()

	if conn == nil {
		return 0
	}

	return conn.LocalAddr().(*net.UDPAddr).Port
}

const reportConnectFailedLimitTimes = 3

// 连接器，传入连接地址和发送封包次数
func (self *kcpConnector) connect(address string) {

	self.SetRunning(true)

	for {
		self.tryConnTimes++

		// 尝试用Socket连接地址
		conn, err := kcp.DialWithOptions(address, blockCrypto, 10, 3)

		self.defaultSes.setConn(conn)

		// 发生错误时退出
		if err != nil {

			if self.tryConnTimes <= reportConnectFailedLimitTimes {
				log.Errorf("#kcp.connect failed(%s) %v", self.Name(), err.Error())

				if self.tryConnTimes == reportConnectFailedLimitTimes {
					log.Errorf("(%s) continue reconnecting, but mute log", self.Name())
				}
			}

			// 没重连就退出
			if self.ReconnectDuration() == 0 || self.IsStopping() {

				self.ProcEvent(&cellnet.RecvMsgEvent{
					Ses: self.defaultSes,
					Msg: &cellnet.SessionConnectError{},
				})
				break
			}

			// 有重连就等待
			time.Sleep(self.ReconnectDuration())

			// 继续连接
			continue
		}

		self.sesEndSignal.Add(1)

		self.ApplySocketOption(conn)

		self.defaultSes.Start()

		self.tryConnTimes = 0

		self.ProcEvent(&cellnet.RecvMsgEvent{Ses: self.defaultSes, Msg: &cellnet.SessionConnected{}})

		self.sesEndSignal.Wait()

		self.defaultSes.setConn(nil)

		// 没重连就退出/主动退出
		if self.IsStopping() || self.ReconnectDuration() == 0 {
			break
		}

		// 有重连就等待
		time.Sleep(self.ReconnectDuration())

		// 继续连接
		continue

	}

	self.SetRunning(false)

	self.EndStopping()
}

func (self *kcpConnector) IsReady() bool {

	return self.SessionCount() != 0
}

func (self *kcpConnector) TypeName() string {
	return "kcp.Connector"
}

func init() {

	peer.RegisterPeerCreator(func() cellnet.Peer {
		self := &kcpConnector{
			SessionManager: new(peer.CoreSessionManager),
		}

		self.defaultSes = newSession(nil, self, func() {
			self.sesEndSignal.Done()
		})

		self.CoreTCPSocketOption.Init()

		return self
	})
}
