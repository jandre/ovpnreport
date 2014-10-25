package ovpnreport

import (
	"net"
	"os"
	"testing"
)

func TestGeoIp(t *testing.T) {

	filename := "../data/GeoIP2-City.mmdb"

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Log("-- Skipping TestGeoIp: no such file or directory: %s", filename)
		return
	}

	result := openGeoIPDb(filename)

	if result == nil {
		t.Fatal("no result returned")
	}

	loc, err := result.City(net.ParseIP("15.5.5.5"))
	t.Log("Location %s", loc.Subdivisions[0].Names["en"])
	t.Log("location %+v %s", loc, err)
}
