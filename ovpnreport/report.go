package ovpnreport

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"strings"

	"github.com/karlbunch/tablewriter"
	"github.com/mgutz/ansi"
)

var (
	lime  = ansi.ColorCode("green+h:black")
	red   = ansi.ColorCode("red")
	green = ansi.ColorCode("green")
	reset = ansi.ColorCode("reset")
)

type counts struct {
	UniqueIps  map[string]bool
	Count      int
	LastSeenAt time.Time
}

func (c *counts) UniqueIpsString() string {
	keys := make([]string, 0, len(c.UniqueIps))
	for k := range c.UniqueIps {
		keys = append(keys, k)
	}

	return strings.Join(keys, ",")
}

func newCounts(logs []*OpenVpnLogin) map[string]*counts {
	var loginsByUser map[string]*counts = make(map[string]*counts)

	for i := range logs {
		log := logs[i]
		if loginsByUser[log.User] != nil {
			loginsByUser[log.User].Count++
			loginsByUser[log.User].UniqueIps[log.IpAddress.String()] = true
			if loginsByUser[log.User].LastSeenAt.Before(log.Timestamp) {
				loginsByUser[log.User].LastSeenAt = log.Timestamp
			}
		} else {
			ips := make(map[string]bool)
			ips[log.IpAddress.String()] = true
			loginsByUser[log.User] = &counts{
				Count:      1,
				LastSeenAt: log.Timestamp,
				UniqueIps:  ips,
			}

		}
	}
	return loginsByUser
}

func LoginsReport(logs []*OpenVpnLogin) {

	logins := newCounts(logs)

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"User", "Count", "Last Seen At", "# Unique IPs", "Unique IPs"})
	// table.SetBorder(false)
	//

	for user, record := range logins {

		table.Append([]string{
			lime + user + reset,
			strconv.Itoa(record.Count),
			record.LastSeenAt.String(),
			strconv.Itoa(len(record.UniqueIps)),
			record.UniqueIpsString(),
		})
	}

	table.Render()
	fmt.Printf("\n")
}
