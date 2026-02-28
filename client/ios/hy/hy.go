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
	stack  *stack.Stack
	flow   PacketFlow
	IP     string

	// waitSend 是从 tun 到 server 发送 IP 包数据的通道。
	// waitSend 由 Send() 写入，(tunReadWriter).Read() 读取
	// 其数据是IP包。
	waitSend chan<- []byte
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
	// Deprecated: ios 须主动调用 Send() 将 IP 包发送到 hy 服务器。 Android 须额外实现从 tun 读取 IP 包并调用 Send()
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
		return nil, errors.New("configured server is empty")
	}
	if cfg.Port == 0 {
		return nil, errors.New("configured port is 0")
	}
	if len(cfg.Uuid) == 0 {
		return nil, errors.New("configured uuid is empty")
	}
	if len(cfg.Obfs) != 0 && len(cfg.Obfs) < 4 {
		return nil, errors.New("configured obfs is too short")
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
		return nil, err
	}
	flow.Log("after create client")

	defaultMogoHysteria.client = hyClient

	//defaultMogoHysteria.IP = hyClient.ClientIP()

	waitSend := make(chan []byte, 1024)
	defaultMogoHysteria.waitSend = waitSend

	flow.Log("before create stack")
	err = defaultMogoHysteria.createStack(waitSend)
	if err != nil {
		err = fmt.Errorf("create stack error: %w", err)
		flow.Log(err.Error())
		return nil, err
	}
	flow.Log("after create stack")

	//err = defaultMogoHysteria.serverTun()

	//go Free()

	//go logLoop(flow)

	if defaultMogoHysteria.IP == "" {
		defaultMogoHysteria.IP = "10.20.0.1"
	}
	return defaultMogoHysteria, nil
}

// Send 是 tun 设备向 Hysteria 服务器发送 IP 包数据。
// 仅用于 ios 平台调用。
func Send(data []byte) error {
	// TODO: rename defaultMogoHysteria to iosDefaultMogoHysteria
	defaultMogoHysteria.send(data)
	return nil
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

func Free() {
	for {
		time.Sleep(time.Second * 5)
		runtime.GC()
		debug.FreeOSMemory()
	}
}
