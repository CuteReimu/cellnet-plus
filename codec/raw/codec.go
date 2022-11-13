package raw

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
)

type packetCodec struct {
}

func (c *packetCodec) Name() string {
	return "packet"
}

func (c *packetCodec) MimeType() string {
	return "application/binary"
}

func (c *packetCodec) Encode(msgObj interface{}, _ cellnet.ContextSet) (data interface{}, err error) {
	return msgObj.(*Packet).Msg, nil
}

func (c *packetCodec) Decode(data interface{}, msgObj interface{}) error {
	msgObj.(*Packet).Msg = data.([]byte)
	return nil
}

func init() {
	codec.RegisterCodec(new(packetCodec))
}
