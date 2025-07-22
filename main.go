package main

import (
	"fmt"
	"github.com/icechen128/mtun/client/ios/hy"
	"github.com/xjasonlyu/tun2socks/v2/buffer"
	"golang.zx2c4.com/wireguard/tun"
	"log"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
)

//type RWDevice struct {
//	tunDevice tun.Device
//	rMutex    sync.Mutex
//	wMutex    sync.Mutex
//}
//
//func (device RWDevice) Read(p []byte) (n int, err error) {
//	device.rMutex.Lock()
//	defer device.rMutex.Unlock()
//	data := [][]byte{p}
//	dataLen := []int{0}
//	_, err = device.tunDevice.Read(data, dataLen, 4)
//	return dataLen[0], err
//}
//
//func (device RWDevice) Write(p []byte) (n int, err error) {
//	device.wMutex.Lock()
//	defer device.wMutex.Unlock()
//	np := gopacket.NewPacket(p[4:], layers.LayerTypeIPv4, gopacket.DecodeOptions{})
//	_ = np
//	return device.tunDevice.Write([][]byte{p}, 4)
//}

type PacketFlow struct {
	device tun.Device
}

func (p PacketFlow) ReadPacket() []byte {
	buf := buffer.Get(buffer.RelayBufferSize)
	defer buffer.Put(buf)
	var buff = [][]byte{buf}
	size := make([]int, 1)
	_, _ = p.device.Read(buff, size, 4)
	return buf[4 : size[0]+4]
}

func (p PacketFlow) ReconnectCallback(err error) {
	fmt.Printf("reconnect call back: %+v", err)
}

func (p PacketFlow) WritePacket(packet []byte) {
	packet = append([]byte{0, 0, 0, 0}, packet...)

	_, err := p.device.Write([][]byte{packet}, 4)
	if err != nil {
		fmt.Println("tun write error: ", err)
		return
	}
}

func (p PacketFlow) Log(msg string) {
	fmt.Println(msg)
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	d, err := tun.CreateTUN("utun99", 1500)
	if err != nil {
		panic(err)
	}

	//device := RWDevice{tunDevice: d}
	//hy.DefaultTunnel = device

	ifName, _ := d.Name()

	runtime.GOMAXPROCS(runtime.NumCPU())

	tunnel, err := hy.StartTunnel(PacketFlow{d}, &hy.HyConfig{
		Server: "45.76.158.147",
		Port:   2021,
		//Server: "127.0.0.1",
		//Port:   443,
		Uuid: "ice",
		//Obfs:               "mogo2022",
		IsDebug:            false,
		LimitMemory:        1000,
		MaxReconnectSecond: 15,
		MaxReconnectCount:  1,
		Bandwidth:          "80mbps",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tunnel)

	SetRoute(tunnel.IP+"/16", ifName, "45.76.158.147:2121", "10.20.0.1", false)

	go func() {
		buf := buffer.Get(buffer.RelayBufferSize)
		defer buffer.Put(buf)
		var buff = [][]byte{buf}
		size := make([]int, 1)
		for {
			_, err = d.Read(buff, size, 4)

			if err != nil {
				fmt.Println("read tun error: " + err.Error())
				return
			}

			err = hy.Send(buf[4 : size[0]+4])
			if err != nil {
				slog.Error("hy closed")
				_ = tunnel.StopTunnel()
				ResetRoute()
				os.Exit(0)
			}
		}
	}()

	select {}

	//go func() {
	//	for {
	//		bytes, err := hy.Receive()
	//		if err != nil {
	//			continue
	//		}
	//
	//		bytes = append([]byte{0, 0, 0, 0}, bytes...)
	//
	//		_, err = d.Write([][]byte{bytes}, 4)
	//		if err != nil {
	//			fmt.Println(err)
	//			wg.Done()
	//			return
	//		}
	//	}
	//}()
}
