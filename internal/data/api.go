package data

// Data structs for the API / UI

// Access Point list for the UI
type ApiAP struct {
	Essid string `json:"essid"`
	Bssid string `json:"bssid"`
	Priv  string `json:"priv"`
}

// Client list for the UI
type ApiClient struct {
	Bssid   string `json:"bssid"`
	Station string `json:"station"`
}

type ApiInterface struct {
	Name   string `json:"name"`
	State  string `json:"mode"`
	Deauth bool   `json:"deauth"`
}

type ApiAction struct {
	Identifier string `json:"identifier"`
	Action     string `json:"action"`
	Deauth     bool   `json:"deauth"`
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
