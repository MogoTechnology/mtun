package hy

import (
	"fmt"

	"github.com/xjasonlyu/tun2socks/v2/core/device/iobased"
	"gvisor.dev/gvisor/pkg/tcpip/stack"
)

type tunnel struct{}

// waitSend 由 Send() 写入，(tunnel).Read() 读取，也即 (*device).Read()
// 其数据是IP包。
var waitSend = make(chan []byte, 1024)

// DefaultTunnel 从 waitSend 读取数据，写入到 defaultMogoHysteria.flow
var DefaultTunnel = tunnel{}

func (t tunnel) Read(p []byte) (n int, err error) {
	b := <-waitSend
	//atomic.AddInt64(&waitSendCount, -1)
	n = copy(p, b)
	//defaultPool.Put(b)
	return n, nil
}

func (t tunnel) Write(p []byte) (n int, err error) {
	// TODO: add flow WritePacket() into tunnel
	if defaultMogoHysteria.flow != nil {
		defaultMogoHysteria.flow.WritePacket(p)
	}
	return len(p), nil
}

const (
	offset     = 0
	defaultMTU = 1500
)

type device struct {
	*iobased.Endpoint
}

var _ stack.LinkEndpoint = (*device)(nil)

// warpTun 创建一个 device, 其中内嵌 stack.LinkEndpoint(.*Endpoint)
// device 是空的，实际读写在 DefaultTunnel。
func warpTun() (*device, error) {
	d := &device{}
	ep, err := iobased.New(d, defaultMTU, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to new EndPoint: %w", err)
	}
	d.Endpoint = ep
	return d, nil
}

func (d *device) Write(b []byte) (int, error) {
	return DefaultTunnel.Write(b)
}

func (d *device) Read(b []byte) (int, error) {
	n, err := DefaultTunnel.Read(b)
	return n, err
}
