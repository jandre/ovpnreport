package ovpnreport

import (
	"time"

	"github.com/jandre/go-papertrail/papertrail"
)

type Papertrail struct {
	Client  *papertrail.Client
	Token   string
	MinID   string
	MaxID   string
	MinTime time.Time
	MaxTime time.Time
	Query   string
}

const VPN_QUERY = "openvpn.log Peer Connection Initiated with "

//
// NewPapertrail(string) will create a new Papertrail
// connection object with the given `token` as
// auth credentials.
//
func NewPapertrail(token string) *Papertrail {
	t := &papertrail.TokenTransport{Token: token}
	client := papertrail.NewClient(t.Client())
	return &Papertrail{Client: client}
}

//
// fetchWithMax(string) fetches a batch of logs with the maxId set
//
func (p *Papertrail) fetchWithMax(id string) ([]*OpenVpnLogin, *papertrail.SearchResponse, error) {
	var logins []*OpenVpnLogin

	if p.Query == "" {
		p.Query = VPN_QUERY
	}

	options := papertrail.SearchOptions{
		Query:   p.Query,
		MinTime: p.MinTime,
		MaxTime: p.MaxTime,
	}

	if id != "" {
		options.MaxID = id
	}

	debug("querying papertrail with options: %q", options)

	response, _, err := p.Client.Search(options)

	if response != nil {
		debug("got response: %q", response)
	}

	logins = make([]*OpenVpnLogin, 0, len(response.Events))

	if err != nil {
		return nil, nil, err
	}

	if response != nil {
		for i := range response.Events {
			ovpn := parseLog(response.Events[i].Message)
			ovpn.Hostname = response.Events[i].SourceName
			if ovpn != nil {
				logins = append(logins, ovpn)
			}
		}
	}

	return logins, response, nil
}

//
// Fetch() will fetch logs from PaperTrail, and return
// an array of OpenVpnLogin logs returned from parseLog()
//
// If it cannot parse a log, it will simply skip it.
//
func (p *Papertrail) Fetch() ([]*OpenVpnLogin, error) {
	var logins []*OpenVpnLogin
	var continueSearch bool = true
	var maxId string = p.MaxID

	for continueSearch {

		logs, response, err := p.fetchWithMax(maxId)

		if err != nil {
			return nil, err
		}

		if response.ReachedTimeLimit && !response.ReachedBeginning {
			debug("continuing with search: %s, %s", response.MaxID, response)
			continueSearch = true
			maxId = response.MaxID
		} else {
			continueSearch = false
		}
		logins = append(logins, logs...)

	}

	return logins, nil
}
