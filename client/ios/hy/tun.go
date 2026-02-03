package hy

import (
	"github.com/xjasonlyu/tun2socks/v2/core"
	"github.com/xjasonlyu/tun2socks/v2/core/adapter"
	"github.com/xjasonlyu/tun2socks/v2/core/option"
	"gvisor.dev/gvisor/pkg/tcpip/stack"
)

var DefaultDevice stack.LinkEndpoint

func (mhy *MogoHysteria) serve() error {
	if DefaultDevice == nil {
		device, err := warpTun()
		if err != nil {
			return err
		}

		mhy.device = device
	} else {
		mhy.device = DefaultDevice
	}
	//
	var err error
	//mhy.device, err = tun.Open("utun123", 1500)
	//if err != nil {
	//	return err
	//}

	//dialer, err := proxy.NewSocks5("127.0.0.1:8123", "", "")
	//if err != nil {
	//	return err
	//}
	//proxy.SetDialer(dialer)

	var opts []option.Option
	opts = append(opts, option.WithTCPSendBufferSize(65536))
	opts = append(opts, option.WithTCPReceiveBufferSize(65536))
	mhy.stack, err = core.CreateStack(&core.Config{
		LinkEndpoint:     mhy.device,
		TransportHandler: mhy,
		Options:          opts,
	})

	return err
}

// HandleTCP 是 TransportHandler 要求的方法
func (mhy *MogoHysteria) HandleTCP(conn adapter.TCPConn) {
	go mhy.handleTCP(conn)
}

// HandleUDP 是 TransportHandler 要求的方法
func (mhy *MogoHysteria) HandleUDP(conn adapter.UDPConn) {
	go mhy.handleUDP(conn)
}

// 没用到，可删
func (mhy *MogoHysteria) serverTun() error {
	//go func() {
	//	for {
	//		if mhy.client.IsClose() {
	//			return
	//		}
	//
	//		hyTun, _, err := mhy.client.Tun()
	//		if err != nil {
	//			continue
	//		}
	//
	//		errCh := make(chan error, 2)
	//
	//		inTun := make(chan []byte, 3000)
	//
	//		go func() {
	//			buf := make([]byte, 64*1024)
	//			length := make([]byte, 2)
	//			var err error
	//			var size, count int
	//			for {
	//				_, err = hyTun.Read(length)
	//				if err != nil {
	//					fmt.Println(err)
	//					errCh <- err
	//					return
	//				}
	//
	//				size = utils.ReadLength(length)
	//
	//				count, err = utils.SplitRead(hyTun, size, buf)
	//				if err != nil {
	//					fmt.Println(err)
	//					errCh <- err
	//					return
	//				}
	//
	//				b := buf[:count]
	//				c := make([]byte, len(b))
	//				copy(c, b)
	//				inTun <- c
	//			}
	//		}()
	//
	//		go func() {
	//			for data := range inTun {
	//				fmt.Println("recv data")
	//				mhy.flow.WritePacket(data)
	//				fmt.Println("recv data success")
	//			}
	//		}()
	//
	//		count := atomic.Int64{}
	//		go func() {
	//			//var buf []byte
	//			head := make([]byte, 2)
	//			for {
	//				buf := <-waitSend
	//				//buf := mhy.flow.ReadPacket()
	//				utils.WriteLength(head, len(buf))
	//
	//				fmt.Println("send to data: ", count.Add(1))
	//				_, err = hyTun.Write(utils.Merge(head, buf))
	//				if err != nil {
	//					errCh <- err
	//					return
	//				}
	//				fmt.Println("send to data success")
	//			}
	//		}()
	//
	//		err = <-errCh
	//		slog.Error("serverTun error:", "error", err)
	//	}
	//}()

	return nil
}
