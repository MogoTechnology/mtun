package hy

import (
	"github.com/xjasonlyu/tun2socks/v2/core/adapter"
	"net"
	"time"
)

func (mhy *MogoHysteria) handleDirectUDP(conn adapter.UDPConn) error {
	defer conn.Close()
	id := conn.ID()

	remoteAddr := net.UDPAddr{
		IP:   net.IP(id.LocalAddress.AsSlice()),
		Port: int(id.LocalPort),
	}
	//localAddr := net.UDPAddr{
	//	IP:   net.IP(id.RemoteAddress.AsSlice()),
	//	Port: int(id.RemotePort),
	//}

	remoteConn, err := net.DialUDP("udp", nil, &remoteAddr)
	if err != nil {
		return err
	}
	defer remoteConn.Close()

	errChan := make(chan error, 2)
	go func() {
		buf := make([]byte, 32*1024)

		for {
			_ = conn.SetDeadline(time.Now().Add(timeout))

			n, err := conn.Read(buf)
			if err != nil {
				errChan <- err
				return
			}

			_, err = remoteConn.Write(buf[:n])
			if err != nil {
				errChan <- err
				return
			}
		}
	}()

	go func() {
		buf := make([]byte, 32*1024)

		for {
			n, addr, err := remoteConn.ReadFromUDP(buf)
			if err != nil {
				errChan <- err
				return
			}

			if addr.IP.Equal(remoteAddr.IP) && addr.Port == remoteAddr.Port {
				_, err = conn.Write(buf[:n])
				if err != nil {
					errChan <- err
					return
				}
			}
		}
	}()

	err = <-errChan
	return err
}
