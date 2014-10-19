package ovpnreport

import (
	"log"

	"github.com/sourcegraph/go-papertrail/papertrail"
)

type Papertrail struct {
	Client *papertrail.Client
}

const VPN_QUERY = "openvpn.log Peer Connection Initiated with "

func NewPapertrail(token string) *Papertrail {
	t := &papertrail.TokenTransport{Token: token}
	client := papertrail.NewClient(t.Client())
	return &Papertrail{Client: client}
}

func (p *Papertrail) Fetch() {
	options := papertrail.SearchOptions{
		Query: VPN_QUERY,
	}
	response, _, err := p.Client.Search(options)
	log.Printf("response: %s", response, err)
}
