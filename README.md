# cellnet-plus

![](https://img.shields.io/github/languages/top/CuteReimu/cellnet-plus "语言")
[![](https://img.shields.io/github/actions/workflow/status/CuteReimu/cellnet-plus/golangci-lint.yml?branch=master)](https://github.com/CuteReimu/cellnet-plus/actions/workflows/golangci-lint.yml "代码分析")
[![](https://img.shields.io/github/contributors/CuteReimu/cellnet-plus)](https://github.com/CuteReimu/cellnet-plus/graphs/contributors "贡献者")
[![](https://img.shields.io/github/license/CuteReimu/cellnet-plus)](https://github.com/CuteReimu/cellnet-plus/blob/master/LICENSE "许可协议")

基于 [github.com/davyxu/cellnet](https://github.com/davyxu/cellnet) ，使其支持更多的协议。

## peer/kcp

见 [github.com/davyxu/cellnet/examples](https://github.com/davyxu/cellnet/tree/master/examples) ，但是需要一些改动：

1. 将`import _ "github.com/davyxu/cellnet/peer/tcp"`改为`import _ "github.com/CuteReimu/cellnet-plus/peer/kcp"`
2. 将`peer.NewGenericPeer("tcp.Acceptor", "name", addr, queue)`改为`peer.NewGenericPeer("kcp.Acceptor", "name", addr, queue)`，同理将`"tcp.Connector"`改为`"kcp.Connector"`
3. 但`import _ "github.com/davyxu/cellnet/proc/tcp"`和下面的`"tcp.ltv"`无需改动
4. **注意，在服务端使用`kcp.Acceptor`时，用户需要自行用心跳或者其它形式检测是否超时，超时后在服务端自行调用`Session.Close()`，以防内存泄漏**

## peer/tcp

length-value格式的tcp包

```go
package main

import (
	"encoding/binary"
	"github.com/CuteReimu/cellnet-plus/peer/tcp"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/tcp"
)

func init() {
	tcp.InitProc(binary.BigEndian, 4) // 下面用"tcp.lv"即可
}

func main() {
	const addr = "0.0.0.0:12345"
	queue := cellnet.NewEventQueue()
	p := peer.NewGenericPeer("tcp.Acceptor", "name", addr, queue)
	proc.BindProcessorHandler(p, "tcp.lv", func(ev cellnet.Event) {
		// ......
	})
}
```

## proc/udp

Encode和Decode纯udp包，不含length和type(id)数据

```go
package main

import (
	"fmt"
	"github.com/CuteReimu/cellnet-plus/codec/raw"
	_ "github.com/CuteReimu/cellnet-plus/proc/udp"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/udp"
)

func main() {
	const addr = "0.0.0.0:12345"
	queue := cellnet.NewEventQueue()
	p := peer.NewGenericPeer("udp.Acceptor", "name", addr, queue)
	proc.BindProcessorHandler(p, "udp.packet", func(ev cellnet.Event) {
		switch msg := ev.Message().(type) {
		case *raw.Packet:
			fmt.Println(msg.Msg)
		}
	})
}
```
