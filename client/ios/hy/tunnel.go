package hy

import (
	"github.com/xjasonlyu/tun2socks/v2/core/device/iobased"
)

type tunnel struct{}

var waitSend = make(chan []byte, 1024)
var waitReceive = make(chan []byte, 1024)

var DefaultTunnel = tunnel{}

func (t tunnel) Read(p []byte) (n int, err error) {
	b := <-waitSend
	//atomic.AddInt64(&waitSendCount, -1)
	n = copy(p, b)
	//defaultPool.Put(b)
	return n, nil
}

func (t tunnel) Write(p []byte) (n int, err error) {
	if defaultMogoHysteria.flow != nil {
		defaultMogoHysteria.flow.WritePacket(p)
	}
	return len(p), nil
}

const (
	offset     = 0
	defaultMTU = 1500
)

type Device struct {
	*iobased.Endpoint
}

func WarpTun() (*Device, error) {
	d := &Device{}
	ep, err := iobased.New(d, defaultMTU, offset)
	if err != nil {
		return nil, err
	}
	d.Endpoint = ep
	return d, nil
}

func (d *Device) Write(b []byte) (int, error) {
	return DefaultTunnel.Write(b)
}

func (d *Device) Read(b []byte) (int, error) {
	n, err := DefaultTunnel.Read(b)
	return n, err
}
