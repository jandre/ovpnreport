package ovpnreport

import (
	"fmt"
	"os"
	"strconv"

	"github.com/karlbunch/tablewriter"
	"github.com/mgutz/ansi"
)

var (
	lime  = ansi.ColorCode("green+h:black")
	red   = ansi.ColorCode("red")
	green = ansi.ColorCode("green")
	reset = ansi.ColorCode("reset")
)

func NewLogins(db *Db, hostname string, logins []*OpenVpnLogin) map[string]*OpenVpnLogin {
	var newLogins map[string]*OpenVpnLogin

	for _, login := range logins {
		if newLogins[login.User] == nil {
			if db.IsNewLoginForUser(login.User, login.Timestamp, hostname) {
				newLogins[login.User] = login
			}
		}
	}

	return newLogins
}

func NewLoginsReport(db *Db, loginsByHostname OpenVpnLogins) {

	for hostname, hostLogins := range loginsByHostname {

		newLogins := NewLogins(db, hostname, hostLogins)

		if len(newLogins) > 0 {
			fmt.Printf("--- New Logins (Never Seen Before) on %s%s%s ---\n", green, hostname, reset)

			table := tablewriter.NewWriter(os.Stdout)

			table.SetHeader([]string{"User", "IP", "Port", "Location"})
			// table.SetBorder(false)

			for user, record := range newLogins {
				var loc string

				if record.City != "" {
					loc = record.City + ", " + record.Country
				}

				fmt.Printf("")
				table.Append([]string{
					lime + user + reset,
					record.Timestamp.String(),
					record.IpAddress.String(),
					loc,
				})
			}

			table.Render()
			fmt.Printf("\n")

		}
	}
}

func LoginsReportByHost(loginsByHostname OpenVpnLogins) {

	for hostname, hostLogins := range loginsByHostname {
		fmt.Printf("--- Logins report for %s%s%s ---\n", green, hostname, reset)
		LoginsReport(hostLogins)
	}
}

func LoginsReport(logs []*OpenVpnLogin) {

	logins := newCounts(logs)

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"User", "Count", "Last Seen At", "# Unique IPs", "Unique IPs", "Locations"})
	// table.SetBorder(false)

	for user, record := range logins {

		fmt.Printf("")
		table.Append([]string{
			lime + user + reset,
			strconv.Itoa(record.Count),
			record.LastSeenAt.String(),
			strconv.Itoa(len(record.UniqueIps)),
			record.UniqueIpsString(),
			record.UniqueLocationsString(),
		})
	}

	table.Render()
	fmt.Printf("\n")
}
