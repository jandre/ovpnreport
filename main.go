package main

import (
	"log"
	"os"
	"time"

	"./ovpnreport"
	. "github.com/visionmedia/go-debug"
	"github.com/visionmedia/go-flags"
)

var debug = Debug("ovpnreport")

var opts struct {
	Config             string `short:"c" long:"config" description:"Path to configuration file." value-type:"config.json"`
	DaysBack           int    `long:"days-back" description:"Look n days back for report instead of using start/end times in config file."`
	Save               bool   `long:"save" description:"Save output to database for future analysis" value-type:"false"`
	ReportNewUsers     bool   `long:"report-new-users"`
	ReportNewIps       bool   `long:"report-new-ips"`
	ReportNewLocations bool   `long:"report-new-locations"`
	ReportAll          bool   `long:"report-all"`
}

//
// applyOverrides() will apply any config overrides from command line
//
// currently supported:
//
// --days-back: use now() - `--days-back` as start/end time for report
// --save
//
func applyOverrides(config *ovpnreport.Config) {

	if opts.DaysBack > 0 {
		config.End = time.Now()
		config.Start = time.Now().Add(time.Hour * time.Duration(-24*opts.DaysBack))
	}
	config.Save = config.Save || opts.Save
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

	applyOverrides(config)

	if err != nil {
		panic(err)
		os.Exit(1)
	}

	for r := range config.Inputs {
		val := config.Inputs[r]

		switch val["type"] {
		case "papertrail":

			debug("using papertrail as an input: %s, %s, %s", val, config.Start, config.End)

			papertrail := ovpnreport.NewPapertrail(val["token"])
			papertrail.MinTime = config.Start
			papertrail.MaxTime = config.End
			logins, _ := papertrail.Fetch()
			loginsByHostname := ovpnreport.LoginsByHostname(logins)
			var db *ovpnreport.Db

			if config.GeoIPDb != "" {
				ovpnreport.ApplyIPLocations(config.GeoIPDb, loginsByHostname)
			}

			ovpnreport.LoginsReportByHost(loginsByHostname)

			if config.Db != "" {
				db = ovpnreport.NewDb(config)
			}

			if db != nil {

				// XXX: add opts to config
				if opts.ReportAll || opts.ReportNewUsers {
					ovpnreport.NewLoginsReport(db, loginsByHostname)
				}

				if opts.Save {
					debug("saving results to database")
					db.Save(loginsByHostname)
				}
			}

		default:
			log.Fatalf("unrecognized config %s", val)
		}
	}

}
