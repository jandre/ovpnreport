package ovpnreport

//
// LoginsByHostname() takes a list of OpenVpnLogin objects,
// and arranges it into a map of hostname -> []OpenVpnLogin
//
func LoginsByHostname(input []*OpenVpnLogin) OpenVpnLogins {
	var logins OpenVpnLogins = make(OpenVpnLogins)

	for _, login := range input {
		if logins[login.Hostname] == nil {
			logins[login.Hostname] = make([]*OpenVpnLogin, 0, 5)
		}
		logins[login.Hostname] = append(logins[login.Hostname], login)
	}

	return logins
}
