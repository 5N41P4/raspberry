package modules

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/5N41P4/raspberry/internal/data"
)

var CaptureBasePath = "/usr/local/raspberry/captures"

func (i *Interface) Capture(t *data.Target) error {
	// If the interface is not im monitor mode, try to set it.
	mon := exec.Command("sudo", "airmon-ng", "start", i.Name)
	_ = mon.Run()

	fileInfo, err := os.Stat(CaptureBasePath)

	if os.IsNotExist(err) || !fileInfo.IsDir() {
		err = os.MkdirAll(CaptureBasePath, 0770)
	}

	if err != nil {
		return err
	}

	i.State = InterfaceStates[Capture]

	// Create a Capture and an ID
	id := i.Name + "_" + time.Now().Format("02.01.2006_15:04")
	path := CaptureBasePath + "/" + id + "/"

	err = os.MkdirAll(path, 0770)
	if err != nil {
		fmt.Println("Couldn't create dir!")
		return err
	}

	if t.Bssid == "" {
		i.process = exec.Command("sudo", "airodump-ng", "-K", "1", "--write", path, "--output-format", "cap,csv", "--wps", i.Name)
	}
	if t.Bssid != "" {
		fmt.Println(i.Target.Bssid + ", " + i.Target.Channel + ", " + i.Target.Station)
		i.process = exec.Command("sudo", "airodump-ng", "-K", "1", "-c", i.Target.Channel, "--bssid", i.Target.Bssid, "--write", path, "--output-format", "cap,csv", "--wps", i.Name)
	}

	err = i.process.Run()

	if err != nil {
		fmt.Println(err)
	}

	i.State = InterfaceStates[Up]

	return err
}
