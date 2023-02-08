package tcp

import (
	"encoding/binary"
	"errors"
	"github.com/CuteReimu/cellnet-plus/codec/raw"
	"io"

	"github.com/davyxu/cellnet"
)

var (
	ErrMaxPacket = errors.New("packet over size")
	ErrMinPacket = errors.New("packet short size")
)

// RecvLVPacket 接收Length-Value格式的封包流程
func RecvLVPacket(reader io.Reader, maxPacketSize int, order binary.ByteOrder, bodySize int) (msg interface{}, err error) {

	// Size为uint16，占2字节
	var sizeBuffer = make([]byte, bodySize)

	// 持续读取Size直到读到为止
	_, err = io.ReadFull(reader, sizeBuffer)

	// 发生错误时返回
	if err != nil {
		return
	}

	if len(sizeBuffer) < bodySize {
		return nil, ErrMinPacket
	}

	// 读取Size
	size := order.Uint32(sizeBuffer)

	if maxPacketSize > 0 && size >= uint32(maxPacketSize) {
		return nil, ErrMaxPacket
	}

	// 分配包体大小
	body := make([]byte, size)

	// 读取包体数据
	_, err = io.ReadFull(reader, body)

	// 发生错误时返回
	if err != nil {
		return
	}

	msg = &raw.Packet{Msg: body}

	return
}

// SendLVPacket 发送Length-Value格式的封包流程
func SendLVPacket(writer io.Writer, _ cellnet.ContextSet, data interface{}, order binary.ByteOrder, bodySize int) error {

	msgData := data.(*raw.Packet).Msg

	pkt := make([]byte, bodySize+len(msgData))

	// Length
	order.PutUint32(pkt, uint32(len(msgData)))

	// Value
	copy(pkt[bodySize:], msgData)

	// 将数据写入Socket
	err := WriteFull(writer, pkt)

	return err
}

// WriteFull 完整发送所有封包
func WriteFull(writer io.Writer, buf []byte) error {

	total := len(buf)

	for pos := 0; pos < total; {

		n, err := writer.Write(buf[pos:])

		if err != nil {
			return err
		}

		pos += n
	}

	return nil

}
