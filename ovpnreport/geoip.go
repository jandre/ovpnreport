package ovpnreport

import (
	"log"

	"github.com/oschwald/geoip2-golang"
)

func openGeoIPDb(dbPath string) *geoip2.Reader {

	db, err := geoip2.Open(dbPath)

	if err != nil {
		log.Printf("Unable to open GeoIP database: %s", dbPath)
		return nil
	}
	return db
}

//
// ApplyIPLocations() tries to do an IP location lookup
//
func ApplyIPLocations(dbPath string, loginsByHost OpenVpnLogins) {

	db := openGeoIPDb(dbPath)

	if db == nil {
		// no database, nothing to apply
		return
	}

	//
	// For each login entry, try to add the City/Country/Lat/Long data.
	//
	for _, logins := range loginsByHost {
		for _, login := range logins {
			if login.IpAddress != nil {
				record, err := db.City(*login.IpAddress)
				if err != nil {
					debug("Error loading city: %s, %s", err, login.IpAddress)
				} else if record != nil {
					login.City = record.City.Names["en"]
					login.Country = record.Country.Names["en"]
					login.Latitude = record.Location.Latitude
					login.Longitude = record.Location.Longitude
				}
			}
		}
	}

}
