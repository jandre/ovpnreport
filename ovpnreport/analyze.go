package ovpnreport

import (
	"strings"
	"time"
)

type counts struct {
	UniqueIps       map[string]bool
	UniqueLocations map[string]bool
	Count           int
	LastSeenAt      time.Time
}

//
// LoginsByHostname() takes a list of OpenVpnLogin objects,
// and arranges it into a map of hostname -> []OpenVpnLogin
//
func LoginsByHostname(input []*OpenVpnLogin) OpenVpnLogins {
	var logins OpenVpnLogins = make(OpenVpnLogins)

	for _, login := range input {
		if logins[login.Hostname] == nil {
			logins[login.Hostname] = make([]*OpenVpnLogin, 0, 5)
		}
		logins[login.Hostname] = append(logins[login.Hostname], login)
	}

	return logins
}

func (c *counts) UniqueIpsString() string {
	keys := make([]string, 0, len(c.UniqueIps))
	for k := range c.UniqueIps {
		keys = append(keys, k)
	}

	return strings.Join(keys, ",")
}

func (c *counts) UniqueLocationsString() string {

	if len(c.UniqueLocations) == 0 {
		return ""
	}

	keys := make([]string, 0, len(c.UniqueLocations))
	for k := range c.UniqueLocations {
		keys = append(keys, k)
	}

	return strings.Join(keys, "\n\n")
}

func newCounts(logs []*OpenVpnLogin) map[string]*counts {
	var loginsByUser map[string]*counts = make(map[string]*counts)

	for i := range logs {
		log := logs[i]
		if loginsByUser[log.User] != nil {
			loginsByUser[log.User].Count++
			loginsByUser[log.User].UniqueIps[log.IpAddress.String()] = true
			if log.City != "" {
				loginsByUser[log.User].UniqueLocations[log.City+", "+log.Country] = true
			}
			if loginsByUser[log.User].LastSeenAt.Before(log.Timestamp) {
				loginsByUser[log.User].LastSeenAt = log.Timestamp
			}
		} else {
			ips := make(map[string]bool)
			ips[log.IpAddress.String()] = true
			locs := make(map[string]bool)
			if log.City != "" {
				locs[log.City+", "+log.Country] = true
			}

			loginsByUser[log.User] = &counts{
				Count:           1,
				LastSeenAt:      log.Timestamp,
				UniqueIps:       ips,
				UniqueLocations: locs,
			}

		}
	}
	return loginsByUser
}
