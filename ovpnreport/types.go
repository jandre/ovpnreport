package ovpnreport

import (
	"net"
	"time"
)

type OpenVpnLogin struct {
	Timestamp time.Time
	User      string
	IpAddress *net.IP
	Port      int
	Hostname  string
}

type OpenVpnLogins map[string]([]*OpenVpnLogin)
