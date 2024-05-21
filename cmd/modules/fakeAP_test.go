package modules

import (
	"testing"

	"github.com/5N41P4/raspberry/internal/data"
)

func TestGetFakeAP(t *testing.T) {
	i := &Interface{Name: "en0"}
	f := &data.Target{
		Bssid:   "A4:CE:DA:87:4D:50",
		Channel: "6",
		Essid:   "FakeAP",
		Cipher:  "WPA2",
	}

	fakeAp, err := NewFakeAP(i.Name, f)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	i.FakeAP = fakeAp

	if i.FakeAP == nil {
		t.Fatal("Expected FakeAP to be initialized")
	}

	if i.FakeAP.Target.Bssid != f.Bssid {
		t.Errorf("Expected BSSID to be '%s', got '%s'", f.Bssid, i.FakeAP.Target.Bssid)
	}

	if i.FakeAP.Target.Channel != f.Channel {
		t.Errorf("Expected Channel to be %s, got %s", f.Channel, i.FakeAP.Target.Channel)
	}

	if i.FakeAP.Target.Essid != f.Essid {
		t.Errorf("Expected ESSID to be '%s', got '%s'", f.Essid, i.FakeAP.Target.Essid)
	}

	if i.FakeAP.Target.Cipher != f.Cipher {
		t.Errorf("Expected Cipher to be '%s', got '%s'", f.Cipher, i.FakeAP.Target.Cipher)
	}
}
