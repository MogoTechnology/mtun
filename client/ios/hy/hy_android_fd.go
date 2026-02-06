package hy

import (
	"errors"
	"fmt"
	"os"

	"github.com/xjasonlyu/tun2socks/v2/buffer"
)

// StartTunnelWithAndroidTunFd 启动hysteria隧道，使用Android Tun FD。
//
// 仅用于 Android 系统，其他系统请使用 StartTunnel()。
// 传入的 fd 会在 StopTunnel() 时关闭，读写出错时也会关闭。
func StartTunnelWithAndroidTunFd(fd int, cfg *HyConfig) (*MogoHysteria, error) {
	tunFile, err := makeTunFile(fd)
	if err != nil {
		return defaultMogoHysteria, errors.New("failed to create the TUN device")
	}

	androidPacketFlow := &androidPacketFlow{
		tunFile: tunFile,
	}
	return StartTunnel(androidPacketFlow, cfg)
}

type androidPacketFlow struct {
	tunFile *os.File
}

// 确保 addroidPacketFlow 实现了 PacketFlow 接口（编译期检查）
var _ PacketFlow = (*androidPacketFlow)(nil)

func (a *androidPacketFlow) WritePacket(packet []byte) {
	_, err := a.tunFile.Write(packet)
	if err != nil {
		a.Log(fmt.Sprintf("tun write error: %v", err))
		a.close()
	}
}

func (a *androidPacketFlow) ReadPacket() []byte {
	buf := make([]byte, buffer.RelayBufferSize)
	n, err := a.tunFile.Read(buf)
	if err != nil {
		a.Log(fmt.Sprintf("tun read error: %v", err))
		a.close()
		return []byte{}
	}
	return buf[:n]
}

func (a *androidPacketFlow) Log(msg string) {
	fmt.Println(msg)
}

func (a *androidPacketFlow) close() {
	a.tunFile.Close()
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
