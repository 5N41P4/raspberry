package modules

import (
	"log"
	"os"
	"os/exec"
)

func (i *Interface) reconStart() error {
	// If the interface is not im monitor mode, try to set it.
	mon := exec.Command("sudo", "airmon-ng", "start", i.Name)
	_ = mon.Run()

	// Delete previous log file
  i.reconClean()

	fileInfo, err := os.Stat("/usr/local/raspberry/recon")

	if os.IsNotExist(err) || !fileInfo.IsDir() {
		err = os.MkdirAll("/usr/local/raspberry/recon", 0770)
	}

  if err != nil {
    log.Println(err)
    return err
  }
  i.State = "recon"

	i.process = exec.Command("sudo", "airodump-ng", "-K", "1", "--write", "/usr/local/raspberry/recon/discovery", "--output-format", "csv", "--wps", i.Name)
	err = i.process.Run()

	i.reconClean()

	i.State = "up"

	return err
}

func (i *Interface) reconClean() {
	os.Remove("/usr/local/raspberry/recon/discovery-01.csv")
}
