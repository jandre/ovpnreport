package main

import (
	"fmt"
	"log"
	"os"

	"./ovpnreport"
	. "github.com/visionmedia/go-debug"
	"github.com/visionmedia/go-flags"
)

var debug = Debug("ovpnreport")

var opts struct {
	Config string `short:"c" long:"config" description:"Path to configuration file." value-type:"config.json"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
		os.Exit(1)
	}

	debug("Using config file: %s", opts.Config)

	if _, err := os.Stat(opts.Config); os.IsNotExist(err) {
		log.Panicf("Config file not found: %s, did you specify the correct path?", opts.Config)
		os.Exit(1)
	}

	config, err := ovpnreport.NewConfig(opts.Config)
	if err != nil {
		panic(err)
		os.Exit(1)
	}

	for r := range config.Inputs {
		val := config.Inputs[r]

		switch val["type"] {
		case "papertrail":

			debug("using papertrail as an input: %q", val)

			papertrail := ovpnreport.NewPapertrail(val["token"])
			logins, _ := papertrail.Fetch()
			// log.Printf("XXX GOT %q", logins, len(logins))
			loginsByHostname := ovpnreport.LoginsByHostname(logins)

			for hostname, hostLogins := range loginsByHostname {
				fmt.Printf("--- Logins report for %s ---\n", hostname)
				ovpnreport.LoginsReport(hostLogins)
			}
		default:
			log.Fatalf("unrecognized config %s", val)
		}
	}

}
