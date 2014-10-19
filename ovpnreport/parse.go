package ovpnreport

import (
	"log"
	"net"
	"regexp"
	"strconv"
	"time"
)

type OpenVpnLogin struct {
	Timestamp time.Time
	User      string
	IpAddress *net.IP
	Port      int
}

var regexpMatchAfNet *regexp.Regexp = regexp.MustCompile(`(\w{3,4} \w{3}\s+\d+ \d+:\d+:\d+ \d+) [\d:\.]+ \[(\w+)\] Peer Connection Initiated with \[AF_INET\](\d+.\d+.\d+.\d+):(\d+)`)

func tryParseAfnet(input string) *OpenVpnLogin {
	var ovpn *OpenVpnLogin

	result := regexpMatchAfNet.FindStringSubmatch(input)
	//	log.Printf("RESULT %s \nmatch: %q", input, result)

	if result != nil && len(result) > 0 {
		user := result[2]
		ip := net.ParseIP(result[3])
		port, _ := strconv.Atoi(result[4])
		date, err := time.Parse(time.ANSIC, result[1])
		if err != nil {
			log.Printf("Unable to parse date %s: %q", result[1], err)
		}
		ovpn = &OpenVpnLogin{
			Timestamp: date,
			User:      user,
			IpAddress: &ip,
			Port:      port,
		}
	}

	return ovpn
}

func parseLog(input string) *OpenVpnLogin {
	return (tryParseAfnet(input))
}
