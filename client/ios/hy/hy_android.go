package hy

import (
	"errors"
	"fmt"
	"os"

	"github.com/xjasonlyu/tun2socks/v2/buffer"
	"golang.org/x/sys/unix"
)

// StartTunnelWithAndroidTunFd 启动hysteria隧道，使用Android Tun FD。
//
// 仅用于 Android 系统，其他系统请使用 StartTunnel()。
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
		a.Close()
	}
}

func (a *androidPacketFlow) ReadPacket() []byte {
	buf := make([]byte, buffer.RelayBufferSize)
	n, err := a.tunFile.Read(buf)
	if err != nil {
		a.Log(fmt.Sprintf("tun read error: %v", err))
		a.Close()
		return []byte{}
	}
	return buf[:n]
}

func (a *androidPacketFlow) Log(msg string) {
	fmt.Println(msg)
}

func (a *androidPacketFlow) Close() {
	err := a.tunFile.Close()
	if err != nil {
		a.Log(fmt.Sprintf("tun close error: %v", err))
	}
}

// makeTunFile returns an os.File object from a TUN file descriptor `fd`.
// The returned os.File holds a separate reference to the underlying file,
// so the file will not be closed until both `fd` and the os.File are
// separately closed.  (UNIX only.)
func makeTunFile(fd int) (*os.File, error) {
	if fd < 0 {
		return nil, errors.New("must provide a valid TUN file descriptor")
	}
	// Make a copy of `fd` so that os.File's finalizer doesn't close `fd`.
	newfd, err := unix.Dup(fd)
	if err != nil {
		return nil, err
	}
	file := os.NewFile(uintptr(newfd), "")
	if file == nil {
		return nil, errors.New("failed to open TUN file descriptor")
	}
	return file, nil
}
