package modules

import (
	"os/exec"
	"strings"

	"github.com/5N41P4/raspberry/internal/data"
)

// Interface is a struct that contains the network interface name and the state of the interface

type Interface struct {
	Name      string       `json:"name"`
	State     string       `json:"state"`
	Target    *data.Target `json:"target"`
	Deauth    *data.Deauth
	Scheduler *data.Scheduler
	process   *exec.Cmd
}

// GetInterfaces is an initialization function to gat all the availabe network interfaces
func GetInterfaces(ifnames []string) map[string]*Interface {
	interfaces := make(map[string]*Interface)
	for _, ifname := range ifnames {
		inf, err := getInterface(ifname)
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
func getInterface(name string) (Interface, error) {
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
		Deauth: &data.Deauth{
			Running: false,
		},
	}

	if strings.Contains(string(out), "inet") {
		inf.State = "inet"
	}

	return inf, err

}

func (i *Interface) Stop() {
	// Go to the fitting cleanup function
	if i.process != nil && i.process.Process != nil {
		i.process.Process.Kill()
	}

	if i.Deauth != nil && i.Deauth.DeauthCan != nil {
		close(i.Deauth.DeauthCan)
	}

	// Stop the monitor mode
	mon := exec.Command("sudo", "airmon-ng", "stop", i.Name)
	_ = mon.Run()
}
