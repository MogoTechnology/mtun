package hy

import (
	"fmt"
	"github.com/xjasonlyu/tun2socks/v2/core/adapter"
)

func (mhy *MogoHysteria) handleTCP(conn adapter.TCPConn) error {
	defer conn.Close()

	id := conn.ID()

	remoteConn, err := mhy.client.TCP(fmt.Sprintf("%s:%d", id.LocalAddress.String(), id.LocalPort))
	if err != nil {
		return err
	}
	defer remoteConn.Close()

	errChan := make(chan error, 2)

	go func() {
		buf := defaultPool.Get()
		defer defaultPool.Put(buf)
		for {
			rn, err := conn.Read(buf)
			if rn > 0 {
				_, err := remoteConn.Write(buf[:rn])
				if err != nil {
					errChan <- err
					return
				}
			}
			if err != nil {
				errChan <- err
				return
			}
		}
	}()

	go func() {
		buf := defaultPool.Get()
		defer defaultPool.Put(buf)
		for {
			rn, err := remoteConn.Read(buf)
			if rn > 0 {
				_, err := conn.Write(buf[:rn])
				if err != nil {
					errChan <- err
					return
				}
			}
			if err != nil {
				errChan <- err
				return
			}
		}
	}()

	err = <-errChan
	return err
}
