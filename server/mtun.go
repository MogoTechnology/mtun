package server

import (
	"errors"
	"github.com/icechen128/mtun/common"
	"github.com/icechen128/mtun/netutil"
	"github.com/icechen128/mtun/tun"
	"github.com/icechen128/mtun/util"
	"github.com/net-byte/water"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"net"
	"time"
)

const DefaultPort = 20231
const DefaultBufferSize = 65507

type Server struct {
	conn      *net.UDPConn
	iface     *water.Interface
	connCache *cache.Cache
}

func New() *Server {
	return &Server{
		connCache: cache.New(30*time.Minute, 10*time.Minute),
	}
}

func (s *Server) Start() error {
	iface, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		return err
	}
	tun.SetRoute("10.10.0.1/24", iface, "", "", true)
	s.iface = iface

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: DefaultPort,
	})
	if err != nil {
		return err
	}
	s.conn = conn
	zap.L().Info("Server started", zap.String("address", conn.LocalAddr().String()))

	go s.handleUDPRequest()
	s.handleTUNRequest()
	return errors.New("server stopped")
}

func (s *Server) handleUDPRequest() {
	zap.L().Info("Handling UDP request")

	var buf [DefaultBufferSize]byte
	for {
		n, remoteAddr, err := s.conn.ReadFromUDP(buf[0:])
		if err != nil {
			zap.L().Error("Read UDP data failed", zap.Error(err))
			break
		}

		zap.L().Info("Read UDP data", zap.Int("data_size", n))
		zap.L().Info("Data", zap.String("data", netutil.BytesToHexStr(buf[:n])))

		if common.IsControlReqIPHeader(buf[:n]) {
			zap.L().Info("Control request received")
			_, err := s.conn.WriteToUDP(append(common.ControlHeaderResIP, []byte(NextIP())...), remoteAddr)
			if err != nil {
				zap.L().Error("Write UDP data failed", zap.Error(err))
				continue
			}
			zap.L().Info("Control request sent")
			continue
		} else if common.IsControlDataHeader(buf[:n]) {
			copy(buf[0:], buf[common.ControlHeaderSize:])
		} else {
			zap.L().Error("Invalid data received")
			continue
		}

		decrypt, err := util.Decrypt(buf[:n-common.ControlHeaderSize], util.DefaultKey)
		if err != nil {
			zap.L().Error("Decrypt data failed", zap.Error(err))
			continue
		}
		zap.L().Info("Read UDP data", zap.Int("data_size", n))
		zap.L().Info("Data", zap.String("data", netutil.BytesToHexStr(decrypt)))

		if key := netutil.GetSrcKey(decrypt); key != "" {
			zap.L().Info("Write to TUN", zap.String("key", key))
			_, err := s.iface.Write(decrypt)
			if err != nil {
				zap.L().Error("Write to TUN failed", zap.Error(err))
				continue
			}
			s.connCache.Set(key, remoteAddr, cache.DefaultExpiration)
			zap.L().Info("Write to TUN success")
		}
	}
}

func (s *Server) handleTUNRequest() {
	zap.L().Info("Handling TUN request")

	var buf [DefaultBufferSize]byte
	for {
		n, err := s.iface.Read(buf[0:])
		if err != nil {
			zap.L().Error("Read TUN interface failed", zap.Error(err))
			break
		}
		zap.L().Info("Read TUN data", zap.Int("data_size", n))
		zap.L().Info("Data", zap.String("data", netutil.BytesToHexStr(buf[:n])))

		cryptBuf, err := util.Encrypt(buf[:n], util.DefaultKey)
		if err != nil {
			zap.L().Error("Encrypt data failed", zap.Error(err))
			continue
		}
		cryptBuf = append(common.ControlHeaderData, cryptBuf...)

		if key := netutil.GetDstKey(buf[:n]); key != "" {
			zap.L().Info("Write to UDP", zap.String("key", key))
			if v, ok := s.connCache.Get(key); ok {
				writeSize, err := s.conn.WriteToUDP(cryptBuf, v.(*net.UDPAddr))
				if err != nil {
					s.connCache.Delete(key)
					zap.L().Error("Write to UDP failed", zap.Error(err))
					continue
				}
				zap.L().Info("Write to UDP", zap.Int("data_size", writeSize))
			}
		}
	}
}

func (s *Server) Stop() error {
	return nil
}
