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
	id := i.Name + "_" + time.Now().Format("02.01.2006_15:04")
	path := "/usr/local/raspberry/captures/" + id + "/"

	err = os.MkdirAll(path, 0770)
	if err != nil {
		fmt.Println("Couldn't create dir!")
		return err
	}

	if i.TargetBssid == "" {
		i.process = exec.Command("sudo", "airodump-ng", "-K", "1", "--write", path, "--output-format", "cap,csv", "--wps", i.Name)
	}
	if i.TargetBssid != "" {
		fmt.Println(i.TargetBssid + ", " + i.TargetChannel + ", " + i.TargetStation)
		i.process = exec.Command("sudo", "airodump-ng", "-K", "1", "-c", i.TargetChannel, "--bssid", i.TargetBssid, "--write", path, "--output-format", "cap,csv", "--wps", i.Name)
	}

	err = i.process.Run()

	if err != nil {
		fmt.Println(err)
	}

	i.State = "up"

	return err
}

func CaptureDelete(capName string) {
	os.RemoveAll("/usr/local/raspberry/captures/" + capName)
}
