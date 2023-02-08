package tcp

import (
	"encoding/binary"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proc"
)

// InitProc 使用"tcp.lv"前需自行在init()中调用这个函数
func InitProc(order binary.ByteOrder, bodySize int) {
	proc.RegisterProcessor("tcp.lv", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {
		bundle.SetTransmitter(&TCPMessageTransmitter{order: order, bodySize: bodySize})
		bundle.SetHooker(new(MsgHooker))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
