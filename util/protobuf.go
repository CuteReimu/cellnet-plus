package util

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/codec"
	"github.com/davyxu/cellnet/util"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"reflect"
)

// RegisterAllProtobuf 按照用hash值作ID注册所有的protobuf协议，调用前需要先 import _ "xxx/xxx" 确保协议的生成文件已经被导入
func RegisterAllProtobuf() {
	c := codec.MustGetCodec("protobuf")
	// 注册所有协议
	protoregistry.GlobalTypes.RangeMessages(func(messageType protoreflect.MessageType) bool {
		cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
			Codec: c,
			Type:  reflect.TypeOf(messageType.Zero().Interface()).Elem(),
			ID:    int(util.StringHash(string(messageType.Descriptor().Name()))),
		})
		return true
	})
}
