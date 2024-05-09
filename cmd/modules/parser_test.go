package modules

import (
	"fmt"
	"testing"
)

// Parse a testfile to check if the CSV parser works
func TestParseCSV(t *testing.T) {
	// csvData which is now located in the testfile folder

	// csvData := `BSSID, First time seen, Last time seen, channel, Speed, Privacy, Cipher, Authentication, Power, # beacons, # IV, LAN IP, ID-length, ESSID, Key
	// A4:CE:DA:87:4D:50, 2024-04-11 21:15:17, 2024-04-11 21:15:31, 56, 1733, WPA3 WPA2, CCMP, SAE PSK, -87,      109,       37,   0.  0.  0.  0,   9, xep-27452,

	// Station MAC, First time seen, Last time seen, Power, # packets, BSSID, Probed ESSIDs
	// DC:A6:32:51:5E:72, 2024-04-11 21:15:18, 2024-04-11 21:17:01, -35,       79, 00:00:00:00:00:00,`

	aps, cls, err := ParseCSV("/users/aurel/Source/raspberry/testfiles/try-02.csv")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(aps) <= 1 {
		fmt.Printf("APs: %v\n", aps)
		t.Fatalf("Expected >= 1 AP, got %d", len(aps))
	}

	ap := aps[0]
	fmt.Printf("ap: %v\n", ap)
	if ap.Bssid != "A4:CE:DA:87:4D:50" {
		t.Errorf("Expected BSSID to be 'A4:CE:DA:87:4D:50', got '%s'", ap.Bssid)
	}

	if ap.Auth != "SAE PSK" {
		t.Errorf("Expected Auth to be 'SAE PSK', got '%s'", ap.Auth)
	}

	if ap.Essid != "xep-27452" {
		t.Errorf("Expected ESSID to be 'xep-27452', got '%s'", ap.Essid)
	}

	if len(cls) <= 1 {
		t.Fatalf("Expected >= 1 client, got %d", len(cls))
	}

	cl := cls[0]
	fmt.Printf("cl: %v\n", cl)
	if cl.Station != "DC:A6:32:51:5E:72" {
		t.Errorf("Expected Station to be 'DC:A6:32:51:5E:72', got '%s'", cl.Station)
	}
	if cl.Packets != 79 {
		t.Errorf("Expected Packets to be 79, got %d", cl.Packets)
	}
}
