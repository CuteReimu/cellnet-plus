package protobuf

import (
	"bytes"
	"encoding/gob"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
)

type gobCodec struct {
}

// Name 编码器的名称
func (c *gobCodec) Name() string {
	return "gob"
}

func (c *gobCodec) MimeType() string {
	return "application/binary"
}

func (c *gobCodec) Encode(msgObj interface{}, _ cellnet.ContextSet) (data interface{}, err error) {
	b := &bytes.Buffer{}
	if err := gob.NewEncoder(b).Encode(msgObj); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (c *gobCodec) Decode(data interface{}, msgObj interface{}) error {
	return gob.NewDecoder(bytes.NewReader(data.([]byte))).Decode(msgObj)
}

// 将消息注册到系统
func init() {
	codec.RegisterCodec(new(gobCodec))
}
