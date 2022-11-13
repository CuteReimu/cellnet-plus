package protobuf

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"google.golang.org/protobuf/proto"
)

type protoCodec struct {
}

// Name 编码器的名称
func (c *protoCodec) Name() string {
	return "protobuf"
}

func (c *protoCodec) MimeType() string {
	return "application/x-protobuf"
}

func (c *protoCodec) Encode(msgObj interface{}, _ cellnet.ContextSet) (data interface{}, err error) {
	return proto.Marshal(msgObj.(proto.Message))
}

func (c *protoCodec) Decode(data interface{}, msgObj interface{}) error {
	return proto.Unmarshal(data.([]byte), msgObj.(proto.Message))
}

// 将消息注册到系统
func init() {
	codec.RegisterCodec(new(protoCodec))
}
