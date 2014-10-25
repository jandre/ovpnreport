package ovpnreport

//
// "Input" defines an interface for an input mechanism.
//
type Input interface {
	//
	// Fetch() fetches logs and returns an array of OpenVpnLogins
	//
	Fetch() ([]*OpenVpnLogin, error)
}
