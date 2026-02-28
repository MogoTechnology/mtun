package hy

import (
	"fmt"
	"io"

	"github.com/xjasonlyu/tun2socks/v2/core/device/iobased"
	"gvisor.dev/gvisor/pkg/tcpip/stack"
)

// tunReadWriter 将 tun 设备封装成 io.ReadWriter, 从 waitSend 读取数据，写入到 defaultMogoHysteria.flow
// 其数据是IP包。
type tunReadWriter struct{}

var _ io.ReadWriter = (*tunReadWriter)(nil)

// waitSend 是从 tun 到 server 发送 IP 包数据的通道。
// waitSend 由 Send() 写入，(tunReadWriter).Read() 读取，也即 (*device).Read()
// 其数据是IP包。
var waitSend = make(chan []byte, 1024)

// Read implements io.ReadWriter.Read.
func (t tunReadWriter) Read(p []byte) (n int, err error) {
	b := <-waitSend
	//atomic.AddInt64(&waitSendCount, -1)
	n = copy(p, b)
	//defaultPool.Put(b)
	return n, nil
}

// Write implements io.ReadWriter.Write.
func (t tunReadWriter) Write(p []byte) (n int, err error) {
	// TODO: add flow WritePacket() into tunReadWriter
	if defaultMogoHysteria.flow != nil {
		defaultMogoHysteria.flow.WritePacket(p)
	}
	return len(p), nil
}

const (
	offset     = 0
	defaultMTU = 1500
)

// device 是一个 stack.LinkEndpoint, 实际读写在 tunnel。
//
// 在 tun2socks 中， LinkEndpoint 的实现通常是一个 TUN 设备包装器，它：
// - 从 TUN 设备读取 IP 数据包（来自操作系统的网络流量）
// - 将数据包传递给 gVisor 网络栈进行处理
// - 将处理后的数据包写回 TUN 设备
type device struct {
	// iobased.Endpoint 实现了 stack.LinkEndpoint 接口
	*iobased.Endpoint
}

var _ stack.LinkEndpoint = (*iobased.Endpoint)(nil)
var _ stack.LinkEndpoint = (*device)(nil)

// warpTun 创建一个 device, 其中内嵌 stack.LinkEndpoint(.*Endpoint)
// device 是空的，实际读写在 DefaultTunnel。
func warpTun() (*device, error) {
	// TODO: use tunReadWriter instead of device d
	ep, err := iobased.New(&tunReadWriter{}, defaultMTU, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to new EndPoint: %w", err)
	}
	return &device{
		Endpoint: ep,
	}, nil
}
