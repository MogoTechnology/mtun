package hy

import (
	"io"
)

// tunReadWriter 将 tun 设备封装成 io.ReadWriter, 从 waitSend 读取数据，写入到 defaultMogoHysteria.flow
// 其数据是IP包。
type tunReadWriter struct{
	waitSend <-chan []byte;
	packetWriter packetWriter;
}

type packetWriter interface {
	WritePacket(packet []byte)
}

var _ io.ReadWriter = (*tunReadWriter)(nil)

// newTunReadWriter 创建一个 tunReadWriter, 它从 waitSend 读取数据，写入到 packetWriter。
func newTunReadWriter(waitSend <-chan []byte, packetWriter packetWriter) *tunReadWriter {
	return &tunReadWriter{
		waitSend:     waitSend,
		packetWriter: packetWriter,
	}
}

// Read implements io.ReadWriter.Read.
func (t *tunReadWriter) Read(p []byte) (n int, err error) {
	b := <-t.waitSend
	//atomic.AddInt64(&waitSendCount, -1)
	n = copy(p, b)
	//defaultPool.Put(b)
	return n, nil
}

// Write implements io.ReadWriter.Write.
func (t *tunReadWriter) Write(p []byte) (n int, err error) {
	// TODO: add flow WritePacket() into tunReadWriter
	if defaultMogoHysteria.flow != nil {
		defaultMogoHysteria.flow.WritePacket(p)
	}
	return len(p), nil
}
