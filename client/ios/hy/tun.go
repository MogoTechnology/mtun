package hy

import (
	"fmt"

	"github.com/xjasonlyu/tun2socks/v2/core"
	"github.com/xjasonlyu/tun2socks/v2/core/adapter"
	"github.com/xjasonlyu/tun2socks/v2/core/device/iobased"
	"github.com/xjasonlyu/tun2socks/v2/core/option"
)

const (
	offset     = 0
	defaultMTU = 1500
)

func (mhy *MogoHysteria) createStack() error {
	// 创建 Endpoint, 实际读写在 tunReadWriter. 
	// Endpoint 实现了 stack.LinkEndpoint 接口。
	// 在 tun2socks 中， LinkEndpoint 的实现通常是一个 TUN 设备包装器，它：
	// - 从 TUN 设备读取 IP 数据包（来自操作系统的网络流量）
	// - 将数据包传递给 gVisor 网络栈进行处理
	// - 将处理后的数据包写回 TUN 设备
	endpoint, err := iobased.New(&tunReadWriter{}, defaultMTU, offset)
	if err != nil {
		return fmt.Errorf("failed to new Endpoint: %w", err)
	}

	var opts []option.Option
	opts = append(opts, option.WithTCPSendBufferSize(65536))
	opts = append(opts, option.WithTCPReceiveBufferSize(65536))
	mhy.stack, err = core.CreateStack(&core.Config{
		LinkEndpoint:     endpoint,
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

//func (mhy *MogoHysteria) serverTun() error {
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

// 	return nil
// }
