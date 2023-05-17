package tun

import (
	"github.com/icechen128/mtun/netutil"
	"github.com/net-byte/water"
	"go.uber.org/zap"
	"net"
	"runtime"
	"strconv"
)

const DefaultMTU = 1500

var localGateway = ""

// SetRoute sets the system routes
func SetRoute(cidr string, iface *water.Interface, serverAddr string, serverIP string, isServer bool) {
	ip, _, err := net.ParseCIDR(cidr)
	if err != nil {
		zap.L().Error("error cidr", zap.String("cidr", cidr))
		return
	}
	localGateway = netutil.DiscoverGateway(true)
	zap.L().With(zap.String("gateway", localGateway)).Info("reset route")

	os := runtime.GOOS
	if os == "linux" {
		netutil.ExecCmd("/sbin/ip", "link", "set", "dev", iface.Name(), "mtu", strconv.Itoa(DefaultMTU))
		netutil.ExecCmd("/sbin/ip", "addr", "add", cidr, "dev", iface.Name())
		netutil.ExecCmd("/sbin/ip", "link", "set", "dev", iface.Name(), "up")
		if !isServer {
			physicalIface := netutil.GetInterface()
			serverAddrIP := netutil.LookupServerAddrIP(serverAddr)
			if physicalIface != "" && serverAddrIP != nil {
				netutil.ExecCmd("/sbin/ip", "route", "add", "0.0.0.0/1", "dev", iface.Name())
				netutil.ExecCmd("/sbin/ip", "route", "add", "128.0.0.0/1", "dev", iface.Name())
				if serverAddrIP.To4() != nil {
					netutil.ExecCmd("/sbin/ip", "route", "add", serverAddrIP.To4().String()+"/32", "via", localGateway, "dev", physicalIface)
				}

			}
		}

	} else if os == "darwin" {
		netutil.ExecCmd("ifconfig", iface.Name(), "inet", ip.String(), serverIP, "up")
		if !isServer {
			physicalIface := netutil.GetInterface()
			serverAddrIP := netutil.LookupServerAddrIP(serverAddr)
			if physicalIface != "" && serverAddrIP != nil {
				netutil.ExecCmd("route", "add", "default", serverIP)
				netutil.ExecCmd("route", "change", "default", serverIP)
				netutil.ExecCmd("route", "add", "0.0.0.0/1", "-interface", iface.Name())
				netutil.ExecCmd("route", "add", "128.0.0.0/1", "-interface", iface.Name())
				if serverAddrIP.To4() != nil {
					netutil.ExecCmd("route", "add", serverAddrIP.To4().String(), localGateway)
				}
			}
		}
	} else if os == "windows" {
		if !isServer {
			serverAddrIP := netutil.LookupServerAddrIP(serverAddr)
			if serverAddrIP != nil {
				netutil.ExecCmd("cmd", "/C", "route", "delete", "0.0.0.0", "mask", "0.0.0.0")
				netutil.ExecCmd("cmd", "/C", "route", "add", "0.0.0.0", "mask", "0.0.0.0", serverIP, "metric", "6")
				if serverAddrIP.To4() != nil {
					netutil.ExecCmd("cmd", "/C", "route", "add", serverAddrIP.To4().String()+"/32", localGateway, "metric", "5")
				}
			}
		}
	} else {
		zap.L().Error("unsupported os", zap.String("os", os))
	}
	zap.L().Info("set route", zap.String("cidr", cidr), zap.String("iface", iface.Name()), zap.String("serverAddr", serverAddr), zap.String("serverIP", serverIP), zap.Bool("isServer", isServer))
}

// ResetRoute resets the system routes
func ResetRoute() {
	zap.L().With(zap.String("gateway", localGateway)).Info("reset route")
	os := runtime.GOOS

	if os == "darwin" {
		netutil.ExecCmd("route", "add", "default", localGateway)
		netutil.ExecCmd("route", "change", "default", localGateway)
	} else if os == "windows" {
		netutil.ExecCmd("cmd", "/C", "route", "delete", "0.0.0.0", "mask", "0.0.0.0")
		netutil.ExecCmd("cmd", "/C", "route", "add", "0.0.0.0", "mask", "0.0.0.0", localGateway, "metric", "6")

	}
}
