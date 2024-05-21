package data

// Data structs for the API / UI

type ApiInterface struct {
	Name   string `json:"name"`
	State  string `json:"mode"`
	Deauth bool   `json:"deauth"`
}

type ApiAction struct {
	Action string `json:"action"`
	Time   int    `json:"time"`
	Target string `json:"target"`
	Deauth bool   `json:"deauth"`
}

type ApiSecurity struct {
	WEP  int `json:"wep"`
	WPA  int `json:"wpa"`
	WPA2 int `json:"wpa2"`
	WPA3 int `json:"wpa3"`
}

type ApiCaptures struct {
	Files []string `json:"files"`
}

type ApiCapture struct {
	Accesspoints []Accesspoint `json:"accesspoints"`
	Clients      []Client      `json:"clients"`
}

type ApiFakeAPStats struct {
	Name      string `json:"name"`
	Running   bool   `json:"running"`
	Handshake bool   `json:"handshake"`
	Key       bool   `json:"key"`
}

type ApiFakeAPAction struct {
	Bssid   string `json:"bssid"`
	Channel int    `json:"channel"`
	Essid   string `json:"essid"`
	Cipher  string `json:"cipher"`
}
