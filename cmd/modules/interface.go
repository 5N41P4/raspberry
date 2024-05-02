package modules

import (
	"errors"
	"log"
	"os/exec"
	"strings"
)

// Interface is a struct that contains the network interface name and the state of the interface

type Interface struct {
	Name    string `json:"name"`
	State   string `json:"state"`
	process *exec.Cmd
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
func (inf *Interface) TryAction(action string) (string, error) {
	if inf.State == "inet" {
		return "", errors.New("Interface is the internet access")
	}

	var err error

	switch action {
	case "capture":
	  go inf.captureStart()
		break

	case "recon":
    go inf.reconStart()
		break

	case "stop":
		inf.stop()
		return action, nil

	default:
		err = errors.New("Invalid action")
	}

	if err != nil {
		log.Println(err)
		return "", err
	}

	return action, nil
}

func (i *Interface) stop() {
	// Go to the fitting cleanup function
	switch i.State {
	case "capture":
    i.captureStop()
    break

	case "recon":
		i.reconStop()
		break

	default:
		break
	}

	// Stop the monitor mode
	mon := exec.Command("sudo", "airmon-ng", "stop", i.Name)
	_ = mon.Run()
}
