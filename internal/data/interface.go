package data

type Target struct {
	Bssid   string `json:"bssid"`
	Essid   string `json:"essid"`
	Station string `json:"station"`
	Channel string `json:"channel"`
	Privacy string `json:"privacy"`
}

type Deauth struct {
	Running   bool `json:"deauth"`
	DeauthCan chan struct{}
}
