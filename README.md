# cellnet-plus

![](https://img.shields.io/github/languages/top/CuteReimu/cellnet-plus "语言")
[![](https://img.shields.io/github/workflow/status/CuteReimu/cellnet-plus/Go)](https://github.com/CuteReimu/cellnet-plus/actions/workflows/golangci-lint.yml "代码分析")
[![](https://img.shields.io/github/contributors/CuteReimu/cellnet-plus)](https://github.com/CuteReimu/cellnet-plus/graphs/contributors "贡献者")
[![](https://img.shields.io/github/license/CuteReimu/cellnet-plus)](https://github.com/CuteReimu/cellnet-plus/blob/master/LICENSE "许可协议")

基于 [github.com/davyxu/cellnet](https://github.com/davyxu/cellnet) ，使其支持更多的协议。

目前支持了kcp

## 使用方法

见 [github.com/davyxu/cellnet/examples](https://github.com/davyxu/cellnet/tree/master/examples) ，但是需要一些改动：

1. 将`import _ "github.com/davyxu/cellnet/peer/tcp"`改为`import _ "github.com/CuteReimu/cellnet-plus/kcp"`
2. 将`peer.NewGenericPeer("tcp.Acceptor", "name", addr, queue)`改为`peer.NewGenericPeer("kcp.Acceptor", "name", addr, queue)`，同理将`"tcp.Acceptor"`改为`"kcp.Acceptor"`
3. 但`import _ "github.com/davyxu/cellnet/proc/tcp"`和下面的`"tcp.ltv"`无需改动
4. 注意，在服务端使用`kcp.Acceptor`时，用户需要自行用心跳或者其它形式检测是否超时，超时后在服务端自行调用`Session.Close()`，以防内存泄漏
