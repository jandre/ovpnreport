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
	City      string
	Country   string
	Latitude  float64
	Longitude float64
}

type OpenVpnLogins map[string]([]*OpenVpnLogin)
