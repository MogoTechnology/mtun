module github.com/icechen128/mtun

go 1.23.2

require (
	github.com/apernet/hysteria/core/v2 v2.5.2
	github.com/apernet/hysteria/extras/v2 v2.5.2
	github.com/mdp/qrterminal/v3 v3.2.0
	github.com/net-byte/go-gateway v0.0.2
	github.com/oschwald/geoip2-golang v1.11.0
	github.com/stretchr/testify v1.9.0
	github.com/xjasonlyu/tun2socks/v2 v2.5.2
	golang.org/x/exp v0.0.0-20241009180824-f66d83c29e7c
	golang.org/x/sys v0.33.0
	golang.zx2c4.com/wireguard v0.0.0-20231211153847-12269c276173
	gvisor.dev/gvisor v0.0.0-20230927004350-cbd86285d259
)

require (
	github.com/apernet/quic-go v0.52.1-0.20250607183305-9320c9d14431 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/google/btree v1.1.3 // indirect
	github.com/google/pprof v0.0.0-20241029153458-d1b30febd7db // indirect
	github.com/onsi/ginkgo/v2 v2.21.0 // indirect
	github.com/oschwald/maxminddb-golang v1.13.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	go.uber.org/mock v0.5.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/term v0.32.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	golang.org/x/time v0.7.0 // indirect
	golang.org/x/tools v0.33.0 // indirect
	golang.zx2c4.com/wintun v0.0.0-20230126152724-0fa3db229ce2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	rsc.io/qr v0.2.0 // indirect
)

replace github.com/apernet/hysteria/core/v2 v2.5.2 => ../hysteria/core
