module github.com/icechen128/mtun

go 1.24.0

require (
	github.com/apernet/hysteria/core/v2 v2.8.1
	github.com/apernet/hysteria/extras/v2 v2.8.1
	github.com/mdp/qrterminal/v3 v3.2.0
	github.com/net-byte/go-gateway v0.0.2
	github.com/oschwald/geoip2-golang v1.11.0
	github.com/stretchr/testify v1.11.1
	github.com/xjasonlyu/tun2socks/v2 v2.6.0
	golang.org/x/exp v0.0.0-20241009180824-f66d83c29e7c
	golang.org/x/sys v0.41.0
	golang.zx2c4.com/wireguard v0.0.0-20250521234502-f333402bd9cb
	gvisor.dev/gvisor v0.0.0-20250523182742-eede7a881b20
)

require (
	github.com/apernet/quic-go v0.59.1-0.20260330051153-c402ee641eb6 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/btree v1.1.3 // indirect
	github.com/oschwald/maxminddb-golang v1.13.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/quic-go/qpack v0.6.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	golang.org/x/crypto v0.47.0 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/term v0.39.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	golang.org/x/time v0.12.0 // indirect
	golang.zx2c4.com/wintun v0.0.0-20230126152724-0fa3db229ce2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	rsc.io/qr v0.2.0 // indirect
)

replace github.com/apernet/hysteria/core/v2 v2.8.1 => ../hysteria/core
