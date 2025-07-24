package hy

import (
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/apernet/hysteria/core/v2/client"
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
	WritePacket(packet []byte)
	ReadPacket() []byte
	Log(msg string)
}

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
		return defaultMogoHysteria, errors.New("server error")
	}
	if cfg.Port == 0 {
		return defaultMogoHysteria, errors.New("port error")
	}
	if len(cfg.Uuid) == 0 {
		return defaultMogoHysteria, errors.New("uuid error")
	}
	if len(cfg.Obfs) != 0 && len(cfg.Obfs) < 4 {
		return defaultMogoHysteria, errors.New("obfs error")
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
		flow.Log("create client error: " + err.Error())
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

func Send(data []byte) error {
	if defaultMogoHysteria.client.IsClose() {
		return errors.New("closed")
	}
	buf := make([]byte, len(data))
	copy(buf, data)
	//atomic.AddInt64(&waitSendCount, 1)
	waitSend <- buf
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
