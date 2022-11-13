package raw

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/util"
	"reflect"
	"strings"
)

type Packet struct {
	Msg []byte
}

func (m *Packet) String() string {
	ret := make([]string, len(m.Msg))
	for i, b := range m.Msg {
		ret[i] = fmt.Sprintf("%02x", b)
	}
	return strings.Join(ret, " ")
}

func init() {
	cellnet.RegisterMessageMeta(&cellnet.MessageMeta{
		Codec: new(packetCodec),
		Type:  reflect.TypeOf((*Packet)(nil)).Elem(),
		ID:    int(util.StringHash("raw.Packet")),
	})
}
