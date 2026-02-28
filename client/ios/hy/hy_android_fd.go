package hy

import (
	"errors"
	"fmt"
	"os"
	"sync/atomic"

	"github.com/xjasonlyu/tun2socks/v2/buffer"
)

// StartTunnelWithAndroidTunFd 启动hysteria隧道，使用Android Tun FD。
//
// 仅用于 Android 系统，其他系统请使用 StartTunnel()。
// 传入的 fd 由调用者管理，即由调用者负责关闭。
func StartTunnelWithAndroidTunFd(fd int, cfg *HyConfig) (*MogoHysteria, error) {
	tunFile, err := makeTunFile(fd)
	if err != nil {
		return nil, errors.New("failed to create the TUN device")
	}

	androidPacketFlow := &androidPacketFlow{
		tunFile: tunFile,
	}

	mogoHysteria, err := StartTunnel(androidPacketFlow, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to start tunnel: %w", err)
	}

	// ReadPacket() 已废弃，须开启协程读取 tun 并调用 Send()
	go mogoHysteria.readFromTunAndSend(tunFile, &androidPacketFlow.closed)

	return mogoHysteria, nil
}

type androidPacketFlow struct {
	tunFile *os.File
	closed  atomic.Bool
}

// 确保 addroidPacketFlow 实现了 PacketFlow 接口（编译期检查）
var _ PacketFlow = (*androidPacketFlow)(nil)

func (a *androidPacketFlow) WritePacket(packet []byte) {
	if a.isClosed() {
		return
	}

	_, err := a.tunFile.Write(packet)
	if err != nil {
		a.Log(fmt.Sprintf("tun write error: %v", err))
	}
}

// ReadPacket 已废弃，不应该被调用。
func (a *androidPacketFlow) ReadPacket() []byte {
	return nil
}

func (a *androidPacketFlow) Log(msg string) {
	fmt.Println(msg)
}

func (a *androidPacketFlow) close() {
	a.closed.Store(true)
}

func (a *androidPacketFlow) isClosed() bool {
	return a.closed.Load()
}

// makeTunFile returns an os.File object from a TUN file descriptor `fd`.
func makeTunFile(fd int) (*os.File, error) {
	if fd < 0 {
		return nil, errors.New("must provide a valid TUN file descriptor")
	}
	file := os.NewFile(uintptr(fd), "")
	if file == nil {
		return nil, errors.New("failed to open TUN file descriptor")
	}
	return file, nil
}

// readFromTunAndSend 从 tun 文件读取 IP 包数据，并调用 Send() 发送。
// 读取出错时会退出。
// 仅用于 Android 系统，其他系统请使用 StartTunnel()。
func (mhy *MogoHysteria) readFromTunAndSend(tunFile *os.File, closed *atomic.Bool) {
	buf := buffer.Get(buffer.RelayBufferSize)
	defer buffer.Put(buf)

	for {
		if closed.Load() {
			return
		}
		n, err := tunFile.Read(buf)
		if err != nil {
			fmt.Println("read tun error: " + err.Error())
			return
		}

		if closed.Load() {
			return
		}
		mhy.Send(buf[:n])
	}
}

// Send 是 tun 设备向 Hysteria 服务器发送 IP 包数据。
// ios 平台须主动调用 Send()。Android 平台使用 StartTunnelWithAndroidTunFd() 自动调用。
func (mhy *MogoHysteria) Send(data []byte) {
	// TODO(jinq): check closed
	// if mhy.client.IsClose() {
	// 	return errors.New("closed")
	// }
	buf := make([]byte, len(data))
	copy(buf, data)
	//atomic.AddInt64(&waitSendCount, 1)
	mhy.waitSend <- buf // tunReadWriter.Read() 将从 waitSend 读取数据
}
