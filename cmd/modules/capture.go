package modules

import (
	"os"
	"os/exec"
	"time"
)

func (i *Interface) captureStart() error {
	// If the interface is not im monitor mode, try to set it.
	mon := exec.Command("sudo", "airmon-ng", "start", i.Name)
	_ = mon.Run()

	dir := exec.Command("sudo", "mkdir", "/usr/local/raspberry/captures/")
	_ = dir.Run()

	// Create a Capture and an ID
	id := i.Name + time.Now().Format("20060102150405")

	cmd := exec.Command("sudo", "airodump-ng", "-K", "1", "--write", "usr/local/raspberry/captures/"+id, "--output-format", "cap", "--wps", i.Name)
	err := cmd.Run()

	return err
}

func (i *Interface) captureStop() {
	// Dirty hack to stop the airmon-ng
	stop := exec.Command("sudo", "killall", "airodump-ng")
	_ = stop.Run()
}

func (i *Interface) captureDelete(capName string) {
	os.Remove("/usr/local/raspberry/captures/" + capName)
}
