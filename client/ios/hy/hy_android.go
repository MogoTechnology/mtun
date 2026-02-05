package hy

import (
	"errors"
	"os"

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

	androidPacketFlow := &addroidPacketFlow{
		tunFile: tunFile,
	}
	return StartTunnel(androidPacketFlow, cfg)
}

type androidPacketFlow struct {
	tunFile *os.File
}
// 确保 addroidPacketFlow 实现了 PacketFlow 接口（编译期检查）
var _ PacketFlow = (*androidPacketFlow)(nil)

(a *androidPacketFlow) WritePacket(packet []byte) {
	// TODO
}

(a *androidPacketFlow) ReadPacket() []byte {
	// TODO
}

(a *androidPacketFlow) Log(msg string) {
	// TODO
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
