package hy

import (
	"io"
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
