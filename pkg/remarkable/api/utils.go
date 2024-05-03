package api

func (a *authToken) Missing() bool {
	return a.deviceToken == "" || a.userToken == ""
}
