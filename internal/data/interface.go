package data

type Target struct {
	Bssid   string `json:"bssid"`
	Station string `json:"station"`
	Channel string `json:"channell"`
}

type Deauth struct {
	Running   bool `json:"deauth"`
	DeauthCan chan struct{}
}