package kcp

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/util"
)

// Socket会话
type kcpSession struct {
	peer.CoreContextSet
	peer.CoreSessionIdentify
	*peer.CoreProcBundle

	pInterface cellnet.Peer

	// Socket原始连接
	conn      net.Conn
	connGuard sync.RWMutex

	// 退出同步器
	exitSync sync.WaitGroup

	// 发送队列
	sendQueue *cellnet.Pipe

	//cleanupGuard sync.Mutex

	endNotify func()

	closing int64
}

func (self *kcpSession) setConn(conn net.Conn) {
	self.connGuard.Lock()
	self.conn = conn
	self.connGuard.Unlock()
}

func (self *kcpSession) Conn() net.Conn {
	self.connGuard.RLock()
	defer self.connGuard.RUnlock()
	return self.conn
}

func (self *kcpSession) Peer() cellnet.Peer {
	return self.pInterface
}

// 取原始连接
func (self *kcpSession) Raw() interface{} {
	return self.Conn()
}

func (self *kcpSession) Close() {

	closing := atomic.SwapInt64(&self.closing, 1)
	if closing != 0 {
		return
	}

	conn := self.Conn()

	if conn != nil {
		// 关闭读
		tcpConn := conn.(*net.TCPConn)
		// 关闭读
		_ = tcpConn.CloseRead()
		// 手动读超时
		_ = tcpConn.SetReadDeadline(time.Now())
	}
}

// 发送封包
func (self *kcpSession) Send(msg interface{}) {

	// 只能通过Close关闭连接
	if msg == nil {
		return
	}

	// 已经关闭，不再发送
	if self.IsManualClosed() {
		return
	}

	self.sendQueue.Add(msg)
}

func (self *kcpSession) IsManualClosed() bool {
	return atomic.LoadInt64(&self.closing) != 0
}

func (self *kcpSession) protectedReadMessage() (msg interface{}, err error) {

	defer func() {

		if err := recover(); err != nil {
			log.Errorf("IO panic: %s", err)
			_ = self.Conn().Close()
		}

	}()

	msg, err = self.ReadMessage(self)

	return
}

// 接收循环
func (self *kcpSession) recvLoop() {

	var capturePanic bool

	if i, ok := self.Peer().(cellnet.PeerCaptureIOPanic); ok {
		capturePanic = i.CaptureIOPanic()
	}

	for self.Conn() != nil {

		var msg interface{}
		var err error

		if capturePanic {
			msg, err = self.protectedReadMessage()
		} else {
			msg, err = self.ReadMessage(self)
		}

		if err != nil {
			if !util.IsEOFOrNetReadError(err) {
				log.Errorf("session closed, sesid: %d, err: %s", self.ID(), err)
			}

			self.sendQueue.Add(nil)

			// 标记为手动关闭原因
			closedMsg := &cellnet.SessionClosed{}
			if self.IsManualClosed() {
				closedMsg.Reason = cellnet.CloseReason_Manual
			}

			self.ProcEvent(&cellnet.RecvMsgEvent{Ses: self, Msg: closedMsg})
			break
		}

		self.ProcEvent(&cellnet.RecvMsgEvent{Ses: self, Msg: msg})
	}

	// 通知完成
	self.exitSync.Done()
}

// 发送循环
func (self *kcpSession) sendLoop() {

	var writeList []interface{}

	for {
		writeList = writeList[0:0]
		exit := self.sendQueue.Pick(&writeList)

		// 遍历要发送的数据
		for _, msg := range writeList {

			self.SendMessage(&cellnet.SendMsgEvent{Ses: self, Msg: msg})
		}

		if exit {
			break
		}
	}

	// 完整关闭
	conn := self.Conn()
	if conn != nil {
		_ = conn.Close()
	}

	// 通知完成
	self.exitSync.Done()
}

// 启动会话的各种资源
func (self *kcpSession) Start() {

	atomic.StoreInt64(&self.closing, 0)

	// connector复用session时，上一次发送队列未释放可能造成问题
	self.sendQueue.Reset()

	// 需要接收和发送线程同时完成时才算真正的完成
	self.exitSync.Add(2)

	// 将会话添加到管理器, 在线程处理前添加到管理器(分配id), 避免ID还未分配,就开始使用id的竞态问题
	self.Peer().(peer.SessionManager).Add(self)

	go func() {

		// 等待2个任务结束
		self.exitSync.Wait()

		// 将会话从管理器移除
		self.Peer().(peer.SessionManager).Remove(self)

		if self.endNotify != nil {
			self.endNotify()
		}

	}()

	// 启动并发接收goroutine
	go self.recvLoop()

	// 启动并发发送goroutine
	go self.sendLoop()
}

func newSession(conn net.Conn, p cellnet.Peer, endNotify func()) *kcpSession {
	self := &kcpSession{
		conn:       conn,
		endNotify:  endNotify,
		sendQueue:  cellnet.NewPipe(),
		pInterface: p,
		CoreProcBundle: p.(interface {
			GetBundle() *peer.CoreProcBundle
		}).GetBundle(),
	}

	return self
}
