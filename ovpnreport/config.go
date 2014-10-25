package ovpnreport

import (
	"encoding/json"
	"io/ioutil"
	"time"

	. "github.com/visionmedia/go-debug"
)

var debug = Debug("ovpnreport")

type Config struct {
	// Start     time.Time
	// End       time.Time
	Start   time.Time           `json:"start_time"`
	End     time.Time           `json:"end_time"`
	Inputs  []map[string]string `json:"inputs"`
	Db      string              `json:"db_path"`
	GeoIPDb string              `json:"geoip_db"`
	Save    bool
}

var (
	ONE_DAY_BACK time.Duration = time.Hour * time.Duration(-24)
)

func NewConfig(file string) (*Config, error) {
	var config Config

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &config)

	debug("Using configuration: %q", config)

	if err != nil {
		return nil, err
	}
	// look 24 hours back
	if config.End.IsZero() {
		config.End = time.Now()
	}

	if config.Start.IsZero() {
		config.Start = config.End.Add(ONE_DAY_BACK)
	}

	if config.Db == "" {
		config.Db = "./db.sqlite3"
	}

	return &config, nil
}
