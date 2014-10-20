# ovpnreport - A small OpenVPN reporting tool

This tool is designed to run periodically to send reports
on OpenVPN logins. 

It will track users logging in a configurable time period, and flag any:

 * New Users          (with --report-new-users)
 * New IP Addresses   (with --report-new-ips)
 * New Locations      (with --report-new-locations) 

The GeoIP MaxMind database is used to track historical login locations 
in a local SQLLite database.


# Install & Build 
```bash
export GOPATH=(desired gopath)
git clone https://github.com/jandre/ovpnreport.git
cd ovpnreport
make
cp config.json.template config.json
# edit config.json
./bin/ovpnreport --config=config.json
```

# Option Flags

TODO

# Inputs

## Papertrail

It uses the Papertrail API to fetch OpenVPN login logs. 

# FAQ

## Q. Why don't you track failed logins?

Because I don't care unless they have successful access to my system.
