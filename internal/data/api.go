package data

// Data structs for the API / UI

type ApiInterface struct {
	Name   string `json:"name"`
	State  string `json:"mode"`
	Deauth bool   `json:"deauth"`
}

type ApiAction struct {
	Identifier string `json:"identifier"`
	Action     string `json:"action"`
	Time       int    `json:"time"`
	Target     string `json:"target"`
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

type ApiCapture struct {
	Accesspoints []Accesspoint `json:"accesspoints"`
	Clients      []Client      `json:"clients"`
}
