package main

import (
	"log"
	"os"

	"./ovpnreport"
	"github.com/visionmedia/go-flags"
)

var opts struct {
	Config string `short:"c" long:"config" description:"Path to configuration file." value-type:"config.json"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
		os.Exit(1)
	}

	log.Printf("Using config file: %s", opts.Config)

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
			log.Printf("Got papertrail")
			papertrail := ovpnreport.NewPapertrail(val["token"])
			papertrail.Fetch()
		default:
			log.Fatalf("unrecognized config %s", val)
		}
	}

}
