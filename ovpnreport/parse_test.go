package ovpnreport

import "testing"

const line = "Oct 15 06:14:10 testhost 41175573 1.2.4.5 User Notice 0xc2081d8090 Wed Oct 15 10:14:10 2014 98.216.104.173:64291 [byoung] Peer Connection Initiated with [AF_INET]98.216.104.173:64291"

func TestParseSyslogAfnet(t *testing.T) {

	result := parseLog(line)

	if result == nil {
		t.Fatal("no result returned")
	}
	t.Log("Result", result)

	if result.Port != 64291 {
		t.Fatal("expected port to be 64291")
	}

	if result.IpAddress.String() != "98.216.104.173" {
		t.Fatal("expected ip to be 98.216.104.173")
	}

	if result.User != "byoung" {
		t.Fatal("expected user to be byoung")
	}
}
