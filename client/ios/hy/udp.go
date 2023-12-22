package hy

import (
	"github.com/xjasonlyu/tun2socks/v2/core/adapter"
	"net"
	"strconv"
	"strings"
	"time"
)

var timeout = time.Second * 60

func (mhy *MogoHysteria) handleUDP(conn adapter.UDPConn) error {
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

	remoteConn, err := mhy.client.UDP()
	if err != nil {
		return err
	}
	defer remoteConn.Close()

	errChan := make(chan error, 2)
	go func() {
		buf := defaultPool.Get()
		defer defaultPool.Put(buf)
		for {
			_ = conn.SetDeadline(time.Now().Add(timeout))

			n, err := conn.Read(buf)
			if err != nil {
				errChan <- err
				return
			}

			err = remoteConn.Send(buf[:n], remoteAddr.String())
			if err != nil {
				errChan <- err
				return
			}
		}
	}()

	go func() {
		for {
			pkt, addr, err := remoteConn.Receive()
			if err != nil {
				errChan <- err
				return
			}

			if addrSlice := strings.Split(addr, ":"); len(addrSlice) == 2 {
				host := addrSlice[0]
				port, _ := strconv.Atoi(addrSlice[1])
				if host == remoteAddr.IP.String() && port == remoteAddr.Port {
					_, err = conn.Write(pkt)
					if err != nil {
						errChan <- err
						return
					}
				}
			}
		}
	}()

	err = <-errChan
	return err
}
