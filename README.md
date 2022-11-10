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
2. 将`peer.NewGenericPeer("tcp.Connector", "name", addr, queue)`改为`peer.NewGenericPeer("kcp.Connector", "name", addr, queue)`，
3. 但`import _ "github.com/davyxu/cellnet/proc/tcp"`和下面的`"tcp.ltv"`无需改动
