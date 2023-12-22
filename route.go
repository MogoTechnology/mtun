package main

import (
	"github.com/icechen128/mtun/netutil"
	"net"
	"runtime"
	"strconv"
)

const DefaultMTU = 1496

var localGateway = ""

// SetRoute sets the system routes
func SetRoute(cidr string, ifaceName, serverAddr string, serverIP string, isServer bool) {
	ip, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return
	}
	localGateway = netutil.DiscoverGateway(true)

	os := runtime.GOOS
	if os == "linux" {
		netutil.ExecCmd("/sbin/ip", "link", "set", "dev", ifaceName, "mtu", strconv.Itoa(DefaultMTU))
		netutil.ExecCmd("/sbin/ip", "addr", "add", cidr, "dev", ifaceName)
		netutil.ExecCmd("/sbin/ip", "link", "set", "dev", ifaceName, "up")
		if !isServer {
			physicalIface := netutil.GetInterface()
			serverAddrIP := netutil.LookupServerAddrIP(serverAddr)
			if physicalIface != "" && serverAddrIP != nil {
				netutil.ExecCmd("/sbin/ip", "route", "add", "0.0.0.0/1", "dev", ifaceName)
				netutil.ExecCmd("/sbin/ip", "route", "add", "128.0.0.0/1", "dev", ifaceName)
				if serverAddrIP.To4() != nil {
					netutil.ExecCmd("/sbin/ip", "route", "add", serverAddrIP.To4().String()+"/32", "via", localGateway, "dev", physicalIface)
				}
			}
		}

	} else if os == "darwin" {
		netutil.ExecCmd("ifconfig", ifaceName, "inet", ip.String(), serverIP, "up")
		if !isServer {
			physicalIface := netutil.GetInterface()
			serverAddrIP := netutil.LookupServerAddrIP(serverAddr)
			if physicalIface != "" && serverAddrIP != nil {
				netutil.ExecCmd("route", "add", "default", serverIP)
				netutil.ExecCmd("route", "change", "default", serverIP)
				netutil.ExecCmd("route", "add", "0.0.0.0/1", "-interface", ifaceName)
				netutil.ExecCmd("route", "add", "128.0.0.0/1", "-interface", ifaceName)
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
	}
}

// ResetRoute resets the system routes
func ResetRoute() {
	os := runtime.GOOS
	if os == "darwin" {
		netutil.ExecCmd("route", "add", "default", localGateway)
		netutil.ExecCmd("route", "change", "default", localGateway)
	} else if os == "windows" {
		netutil.ExecCmd("cmd", "/C", "route", "delete", "0.0.0.0", "mask", "0.0.0.0")
		netutil.ExecCmd("cmd", "/C", "route", "add", "0.0.0.0", "mask", "0.0.0.0", localGateway, "metric", "6")
	}
}
