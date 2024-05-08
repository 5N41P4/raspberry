package modules

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"

	"github.com/5N41P4/raspberry/internal/data"
)

// Interface is a struct that contains the network interface name and the state of the interface

type Interface struct {
	Name          string `json:"name"`
	State         string `json:"state"`
	Deauth        bool   `json:"deauth"`
	process       *exec.Cmd
	deauth        chan struct{}
	TargetBssid   string
	TargetStation string
	TargetChannel string
}

// GetInterfaces is an initialization function to gat all the availabe network interfaces
func GetInterfaces(ifnames []string) map[string]*Interface {
	interfaces := make(map[string]*Interface)
	for _, ifname := range ifnames {
		inf, err := GetInterface(ifname)
		if err != nil {
			continue
		}
		interfaces[inf.Name] = &inf
	}

	if len(interfaces) == 0 {
		return nil
	}

	return interfaces
}

// GetInterface is a function that returns an Interface struct with the network interface name and the possible modules
func GetInterface(name string) (Interface, error) {
	// Check if the interface exists
	cinf := exec.Command("ip", "address", "show", "dev", name)
	out, err := cinf.Output()

	// If the interface does not exist return an error
	if err != nil {
		return Interface{}, err
	}

	inf := Interface{
		Name:  name,
		State: "up",
	}

	if strings.Contains(string(out), "inet") {
		inf.State = "inet"
	}

	return inf, nil

}

// StartCapture is a function that starts a capture on the network interface
func (inf *Interface) TryAction(action data.ApiAction, access *map[string]*data.AppAP, clients *map[string]*data.AppClient) (string, error) {
	if inf.State == "inet" {
		return "", errors.New("Interface is the internet access")
	}

	var err error

	switch action.Action {
	case "capture":
		inf.TargetBssid, inf.TargetStation, inf.TargetChannel = getTarget(action.Target, access, clients)
		go inf.captureStart()

	case "recon":
		go inf.reconStart()

	case "stop":
		inf.stop()
		return action.Action, nil

	default:
		err = errors.New("invalid action")
		return "", err
	}

	if action.Deauth {
		inf.Deauth = true
		inf.deauth = make(chan struct{})
		go inf.RunDeauth(access, clients)
	}

	return action.Action, err
}

func (i *Interface) stop() {
	// Go to the fitting cleanup function
	i.process.Process.Kill()

	close(i.deauth)

	// Stop the monitor mode
	mon := exec.Command("sudo", "airmon-ng", "stop", i.Name)
	_ = mon.Run()
}

func getTarget(target string, access *map[string]*data.AppAP, clients *map[string]*data.AppClient) (bssid string, station string, channel string) {
	// If target is a client station then fill in the target information from the client
	cl, ok := (*clients)[target]
	if ok {
		return cl.Bssid, cl.Station, strconv.Itoa((*access)[cl.Bssid].Channel)
	}

	// If the target is a BSSID then fill in the target with the information from the accesspoint
	ap, ok := (*access)[target]
	if ok {
		return ap.Bssid, "", strconv.Itoa(ap.Channel)
	}

	// If the target could not be found then we fill in empty strings as a default
	return "", "", ""
}
