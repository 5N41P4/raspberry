package data

type Target struct {
	Bssid   string `json:"bssid"`
	Essid   string `json:"essid"`
	Station string `json:"station"`
	Channel string `json:"channel"`
	Cipher  string `json:"cipher"`
}

type Deauth struct {
	Running   bool `json:"deauth"`
	DeauthCan chan struct{}
}
