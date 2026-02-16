package hy

import (
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/apernet/hysteria/core/v2/client"
	"github.com/xjasonlyu/tun2socks/v2/core/adapter"
	"gvisor.dev/gvisor/pkg/tcpip/stack"
)

type MogoHysteria struct {
	client client.Client
	device stack.LinkEndpoint
	stack  *stack.Stack
	flow   PacketFlow
	IP     string
}

var defaultMogoHysteria *MogoHysteria
var _ adapter.TransportHandler = (*MogoHysteria)(nil)

type HyConfig struct {
	Server      string
	Port        int
	Uuid        string
	Obfs        string
	IsDebug     bool
	LimitMemory int
	Bandwidth   string
	Token       string
}

type PacketFlow interface {
	// WritePacket 向 tun 设备写入 IP 包。
	WritePacket(packet []byte)
	// ReadPacket 从 tun 设备读取 IP 包。
	// ios 应该已废弃使用。ios 须主动调用 Send() 将 IP 包发送到 hy 服务器。 
	// TOOD：Android 须额外实现从 tun 读取 IP 包并调用 Send()
	ReadPacket() []byte
	Log(msg string)
}

// StartTunnel 启动hysteria隧道。
//
// Android 系统可使用 StartTunnelWithAndroidTunFd(), 更简单。
func StartTunnel(flow PacketFlow, cfg *HyConfig) (*MogoHysteria, error) {
	//cfg = &HyConfig{
	//	//Server: "127.0.0.1",
	//	Server: "47.101.36.120",
	//	Port:   9996,
	//	Uuid:   cfg.Uuid,
	//	Token:  cfg.Token,
	//	//Obfs:               "mogo2022",
	//	IsDebug:            false,
	//	LimitMemory:        1000,
	//	Bandwidth:          "",
	//}
	//if cfg.LimitMemory != 0 {
	//	debug.SetMemoryLimit(int64(cfg.LimitMemory * 1 << 20))
	//} else {
	//	debug.SetMemoryLimit(20 * 1 << 20)
	//}
	flow.Log("start tunnel...")
	if len(cfg.Server) == 0 {
		return defaultMogoHysteria, errors.New("configured server is empty")
	}
	if cfg.Port == 0 {
		return defaultMogoHysteria, errors.New("configured port is 0")
	}
	if len(cfg.Uuid) == 0 {
		return defaultMogoHysteria, errors.New("configured uuid is empty")
	}
	if len(cfg.Obfs) != 0 && len(cfg.Obfs) < 4 {
		return defaultMogoHysteria, errors.New("configured obfs is too short")
	}

	//if cfg.Bandwidth == "" {
	//	cfg.Bandwidth = "80mbps"
	//}

	defaultMogoHysteria = &MogoHysteria{
		flow: flow,
	}

	config := &clientConfig{
		//Server: "47.95.31.127:7865",
		//Server: "192.144.225.219:4433",
		Server: cfg.Server + ":" + strconv.Itoa(cfg.Port),
		Auth:   cfg.Uuid + "|" + cfg.Token,
		TLS: clientConfigTLS{
			Insecure: true,
			//SNI:      "n1234.platovpn.com",
		},
		QUIC: clientConfigQUIC{
			InitStreamReceiveWindow:     2097152,
			MaxStreamReceiveWindow:      2097152,
			InitConnectionReceiveWindow: 5242880,
			MaxConnectionReceiveWindow:  5242880,
			MaxIdleTimeout:              time.Second * 30,
			KeepAlivePeriod:             time.Second * 5,
			DisablePathMTUDiscovery:     false,
		},
		Bandwidth: clientConfigBandwidth{
			Up:   "30mbps",
			Down: "60mbps",
		},
	}
	if cfg.Obfs != "" {
		config.Obfs.Type = "salamander"
		config.Obfs.Salamander.Password = cfg.Obfs
	}

	flow.Log("before create client")
	hyClient, err := client.NewReconnectableClient(config.Config, func(c client.Client, info *client.HandshakeInfo, i int) {
		flow.Log(fmt.Sprintf("connected, count: %d", i))
	}, false)

	if err != nil {
		err = fmt.Errorf("create client error: %w", err)
		flow.Log(err.Error())
		return defaultMogoHysteria, err
	}
	flow.Log("after create client")

	defaultMogoHysteria.client = hyClient

	//defaultMogoHysteria.IP = hyClient.ClientIP()

	flow.Log("before create stack")
	err = defaultMogoHysteria.serve()
	flow.Log("after create stack")

	//err = defaultMogoHysteria.serverTun()

	//go Free()

	//go logLoop(flow)

	if defaultMogoHysteria.IP == "" {
		defaultMogoHysteria.IP = "10.20.0.1"
	}
	return defaultMogoHysteria, err
}

// Send 是 tun 设备向 Hysteria 服务器发送 IP 包数据。
// 仅用于 ios 平台，Android 平台使用 StartTunnelWithAndroidTunFd(), 直接从 tun fd 读取 IP 包数据。
func Send(data []byte) error {
	// TODO(jinq): check closed
	// if defaultMogoHysteria.client.IsClose() {
	// 	return errors.New("closed")
	// }
	buf := make([]byte, len(data))
	copy(buf, data)
	//atomic.AddInt64(&waitSendCount, 1)
	waitSend <- buf  // tunnel.Read() 将从 waitSend 读取数据
	return nil
}

//func BatchSend(data [][]byte) error {
//	var err error
//	for _, d := range data {
//		e := Send(d)
//		if err != nil {
//			err = e
//			return err
//		}
//	}
//	return err
//}

// 应该没用。waitReceive 没人写入。
func Receive() ([]byte, error) {
	//timeoutTicker := time.NewTicker(time.Second * 10)
	//defer timeoutTicker.Stop()
	//select {
	//case data := <-waitReceive:
	//	return data, nil
	//case <-timeoutTicker.C:
	//	return nil, errors.New("timeout")
	//}

	data := <-waitReceive
	//atomic.AddInt64(&waitReceiveCount, -1)
	return data, nil
}

func (mhy *MogoHysteria) StopTunnel() error {
	//go defaultMogoHysteria.stack.Close()
	if mhy == nil {
		return errors.New("mogo hysteria nil")
	}
	if mhy.flow == nil {
		return errors.New("package flow nil")
	}
	if mhy.client == nil {
		return errors.New("mogo hysteria client nil")
	}
	if mhy.stack == nil {
		return errors.New("mogo hysteria stack nil")
	}

	mhy.flow.Log("start stop")
	mhy.client.Close()

	if androidFlow, ok := mhy.flow.(*androidPacketFlow); ok {
		androidFlow.close()
	}

	return nil
}

//var log = make(chan string)
//
//func Log() string {
//	return <-log
//}

func logLoop(flow PacketFlow) {
	for {
		time.Sleep(time.Second)
		if defaultMogoHysteria == nil {
			flow.Log("mogo hysteria nil")
		}
		if defaultMogoHysteria.flow == nil {
			flow.Log("package flow nil")
		}
		if defaultMogoHysteria.client == nil {
			flow.Log("mogo hysteria client nil")
		}
		if defaultMogoHysteria.stack == nil {
			flow.Log("mogo hysteria stack nil")
		}
		flow.Log("mogo hysteria")
	}
}

func Free() {
	for {
		time.Sleep(time.Second * 5)
		runtime.GC()
		debug.FreeOSMemory()
	}
}
