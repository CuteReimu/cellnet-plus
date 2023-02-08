package udp

import (
	"github.com/CuteReimu/cellnet-plus/codec/raw"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/peer/udp"
	"github.com/davyxu/cellnet/proc"
)

type UDPMessageTransmitter struct {
}

func (UDPMessageTransmitter) OnRecvMessage(ses cellnet.Session) (msg interface{}, err error) {
	data := ses.Raw().(udp.DataReader).ReadData()
	m := &raw.Packet{}
	m.Msg = append(([]byte)(nil), data...)
	msg = m
	msglog.WriteRecvLogger(log, "udp", ses, msg)
	return
}

func (UDPMessageTransmitter) OnSendMessage(ses cellnet.Session, msg interface{}) error {
	writer := ses.(udp.DataWriter)
	message, ok := msg.(*raw.Packet)
	if !ok {
		log.Warnf("unsupported message type: %T", message)
		return nil
	}
	msglog.WriteSendLogger(log, "udp", ses, msg)
	writer.WriteData(message.Msg)
	return nil
}

func init() {
	proc.RegisterProcessor("udp.packet", func(bundle proc.ProcessorBundle, userCallback cellnet.EventCallback) {
		bundle.SetTransmitter(new(UDPMessageTransmitter))
		bundle.SetCallback(userCallback)
	})
}
