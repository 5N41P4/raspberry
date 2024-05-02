package data

type AppAP struct {
	Bssid   string `json:"bssid"`
	First   string `json:"first seen at"`
	Last    string `json:"last seen at"`
	Channel int    `json:"channel"`
	Speed   int    `json:"speed"`
	Privacy string `json:"privacy"`
	Cipher  string `json:"cipher"`
	Auth    string `json:"auth"`
	Power   int    `json:"power"`
	Beacons int    `json:"beacons"`
	IVs     int    `json:"ivs"`
	Lan     string `json:"lan ip"`
	IdLen   int    `json:"id len"`
	Essid   string `json:"essid"`
	Key     string `json:"key"`
}

type AppClient struct {
	// MAC address
	Station string `json:"station"`
	First   string `json:"first seen at"`
	Last    string `json:"last seen at"`
	Power   int    `json:"power"`
	Packets int    `json:"packets"`
	Bssid   string `json:"bssid"`
	Probed  string `json:"probed essids"`
}
