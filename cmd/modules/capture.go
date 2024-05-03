package modules

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func (i *Interface) captureStart() error {
	// If the interface is not im monitor mode, try to set it.
	mon := exec.Command("sudo", "airmon-ng", "start", i.Name)
	_ = mon.Run()

	fileInfo, err := os.Stat("/usr/local/raspberry/captures")

	if os.IsNotExist(err) || !fileInfo.IsDir() {
		err = os.MkdirAll("/usr/local/raspberry/captures", 0770)
	}

	if err != nil {
		return err
	}

	i.State = "capture"

	// Create a Capture and an ID
	id := i.Name + time.Now().Format("20060102150405")

	i.process = exec.Command("sudo", "airodump-ng", "-K", "1", "--write", "/usr/local/raspberry/captures/"+id, "--output-format", "cap", "--wps", i.Name)
	err = i.process.Run()

	if err != nil {
		fmt.Println(err)
	}

	i.State = "up"

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
