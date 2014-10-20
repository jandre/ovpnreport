package ovpnreport

import (
	"encoding/json"
	"io/ioutil"
	"time"

	. "github.com/visionmedia/go-debug"
)

var debug = Debug("ovpnreport")

type Config struct {
	Start  time.Time           `json: "start_time"` // time at which to start reporting on user logins
	End    time.Time           `json: "end_time"`   // time at which to end reporting on user logins
	Inputs []map[string]string `json: "inputs"`
	Db     string              `json: "db_path"`
}

func NewConfig(file string) (*Config, error) {
	var config Config
	var ONE_DAY_BACK time.Duration = time.Hour * time.Duration(-24)

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &config)
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

	cfg, _ := json.Marshal(config)
	debug("Using configuration: %q", cfg)

	return &config, nil
}
