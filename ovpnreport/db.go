package ovpnreport

import (
	"log"
	"net"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type LoginHistory struct {
	Id          int64
	User        string
	IpAddress   string
	Hostname    string
	Count       int64
	FirstSeenAt time.Time
	LastSeenAt  time.Time
}

//
// Stores all history OpenVPN logins here for analysis
//
type Db struct {
	connection *gorm.DB
}

//
// setup() will setup that sqlite database with an openvpnlogins
// table
//
func (db *Db) setup() {

	if !db.connection.HasTable(&LoginHistory{}) {
		db.connection.CreateTable(&LoginHistory{})
		db.connection.Model(&LoginHistory{}).AddIndex("user")
		db.connection.Model(&LoginHistory{}).AddIndex("ip_address")
	}
}

//
// IsNewLoginForUser(string,ip,time) returns true
// if we haven't seen this ip for this user before the max time
//
func (db *Db) IsNewLoginForUser(user string, before time.Time) bool {
	var login LoginHistory
	err := db.connection.
		Model(&LoginHistory{}).
		Where("user = ? and first_seen_at < ?", user, before).
		First(&login)
	if err != nil {
		log.Panicf("Unable to query IsNewLoginForUser", err)
	}
	if login.Id != 0 {
		return false
	}
	return true
}

//
// IsNewIpForUser(string,ip,time) returns true
// if we haven't seen this ip for this user before the max time
//
func (db *Db) IsNewIpForUser(user string, ip *net.IP, since time.Time) bool {
	var login LoginHistory
	err := db.connection.Model(&LoginHistory{}).Where("user = ? and ip_address = ? and last_seen_at < ?", user, ip.String(), since).First(&login)
	if err != nil {
		log.Panicf("Unable to query IsNewIpForUser", err)
	}
	log.Printf("XXX", login)
	if login.Id != 0 {
		return false
	}
	return true
}

//
// ImportLogs(logs[]) imports an array of vpn logs into the database.
//
func (db *Db) ImportLogs(logs []*OpenVpnLogin) {
	log.Printf("Populating Sqlite... (%s entries)", len(logs))
	for i := range logs {
		var login LoginHistory
		err := db.connection.Model(&LoginHistory{}).
			Where("user = ? and ip_address = ?", logs[i].User, logs[i].IpAddress.String()).
			First(&login)
		if err != nil {
			log.Panicf("Unable to query", err)
		}
		if login.Id != 0 {

			lastSeen := login.LastSeenAt
			firstSeen := login.FirstSeenAt

			if lastSeen.Before(logs[i].Timestamp) {
				lastSeen = logs[i].Timestamp
			}
			if firstSeen.After(logs[i].Timestamp) {
				firstSeen = logs[i].Timestamp
			}

			db.connection.Model(&login).Update(&LoginHistory{
				LastSeenAt:  lastSeen,
				FirstSeenAt: firstSeen,
				Count:       login.Count + 1,
			})

		} else {
			login.User = logs[i].User
			login.LastSeenAt = logs[i].Timestamp
			login.FirstSeenAt = logs[i].Timestamp
			login.IpAddress = logs[i].IpAddress.String()
			login.Count = 1
			db.connection.Table("login_histories").Create(login)
		}
	}
	log.Printf("Finished Populating Sqlite...")
}

//
// NewDb(config) creates a new database connection.
//
// It uses config.Db as the sqlite connection string.
//
// Panics if it cannot open the sqlite database
//
func NewDb(config *Config) *Db {
	var db Db

	connection, err := gorm.Open("sqlite3", config.Db)
	if err != nil {
		log.Panicf("Unable to open sqlitedb: %s", config.Db, err)
	}

	db = Db{connection: &connection}
	db.setup()

	return &db
}
