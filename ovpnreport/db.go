package ovpnreport

import (
	"log"
	"net"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type DbOpenVpnLogin struct {
	Id        int64
	Timestamp time.Time
	User      string
	IpAddress string
	Port      int
	Hostname  string
	City      string
	Country   string
	Latitude  float64
	Longitude float64
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

	if !db.connection.HasTable(&DbOpenVpnLogin{}) {
		db.connection.CreateTable(&DbOpenVpnLogin{})
		db.connection.Model(&DbOpenVpnLogin{}).AddIndex("user")
	}
}

//
// IsNewLoginForUser(string,ip,time) returns true
// if we haven't seen this ip for this user before the max time
//
func (db *Db) IsNewLoginForUser(user string, before time.Time, hostname string) bool {
	var login DbOpenVpnLogin
	db.connection.
		Model(&DbOpenVpnLogin{}).
		Where("user = ? and timestamp < ? and hostname = ?", user, before, hostname).
		First(&login)

	if login.Id != 0 {
		return false
	}
	return true
}

//
// IsNewIpForUser(string,ip,time) returns true
// if we haven't seen this ip for this user before the max time
//
func (db *Db) IsNewIpForUser(user string, ip *net.IP, since time.Time, hostname string) bool {
	var login DbOpenVpnLogin
	db.
		connection.Model(&DbOpenVpnLogin{}).
		Where("user = ? and ip_address = ? and timestamp < ? and hostname = ?", user, ip.String(), since, hostname).First(&login)
	log.Printf("XXX", login)
	if login.Id != 0 {
		return false
	}
	return true
}

//
// Save() saves a collection of OpenVpnLogins to the database.
//
func (db *Db) Save(loginsByHostname OpenVpnLogins) {

	for _, logins := range loginsByHostname {
		db.ImportLogs(logins)
	}
}

//
// ImportLogs(logs[]) imports an array of vpn logs into the database.
//
func (db *Db) ImportLogs(logs []*OpenVpnLogin) {
	debug("Populating Sqlite... (%s entries)", len(logs))
	for i := range logs {
		var login DbOpenVpnLogin
		db.connection.Model(&DbOpenVpnLogin{}).
			Where(&DbOpenVpnLogin{
			User:      logs[i].User,
			IpAddress: logs[i].IpAddress.String(),
			Hostname:  logs[i].Hostname,
			Timestamp: logs[i].Timestamp,
			Port:      logs[i].Port,
		}).First(&login)
		if login.Id == 0 {

			db.connection.Model(&DbOpenVpnLogin{}).Create(&DbOpenVpnLogin{
				User:      logs[i].User,
				IpAddress: logs[i].IpAddress.String(),
				Hostname:  logs[i].Hostname,
				Timestamp: logs[i].Timestamp,
				Port:      logs[i].Port,
				Latitude:  logs[i].Latitude,
				Longitude: logs[i].Longitude,
				City:      logs[i].City,
				Country:   logs[i].Country,
			})
		}
	}
	debug("Finished Populating Sqlite...")
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
