package data

// Data structs for the API / UI

// Structs for API Requests

type ApiSimpleAction struct {
	Action     string `json:"action"`
	Identifier string `json:"identifier"`
}

type ApiInterfaceAction struct {
	Action string `json:"action"`
	Time   int    `json:"time"`
	Target Target `json:"target"`
	Deauth bool   `json:"deauth"`
}

// Structs for API Responses

type ApiInterface struct {
	Name   string `json:"name"`
	State  string `json:"mode"`
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
