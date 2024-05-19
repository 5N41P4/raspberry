package modules

import (
	"testing"
	"time"
)

// TestGetInterfaces tests the GetInterfaces function
func TestGetInterfaces(t *testing.T) {
	configWlan := []string{"wlan0", "wlan1"}
	configEn := []string{"en0", "en1"}
	interfacesWlan := GetInterfaces(configWlan)
	interfacesEn := GetInterfaces(configEn)

	if interfacesWlan == nil && interfacesEn == nil {
		t.Errorf("GetInterfaces returned nothing")
	}

	if len(interfacesEn) <= 1 && len(interfacesWlan) <= 1 {
		t.Errorf("GetInterfaces(en) returned: %d\n GetInterfaces(wlan) returned: %d", len(interfacesEn), len(interfacesWlan))
	}

}
func TestStopAfter(t *testing.T) {
	interfaceName := "wlan0"
	delay := 1 // minutes

	// Create a mock Interface
	mockInterface := &Interface{
		Name: interfaceName,
	}

	// Start the StopAfter method with the specified delay
	mockInterface.StopAfter(delay)

	// Wait for the specified delay + 1 second
	time.Sleep((time.Duration(delay) + 1) * time.Minute)

	// Verify that the Stop method was called
	if mockInterface.process != nil && mockInterface.process.Process != nil {
		t.Errorf("Stop method was not called")
	}

	// Verify that the Deauth channel was closed
	if mockInterface.Deauth != nil && mockInterface.Deauth.DeauthCan != nil {
		select {
		case _, ok := <-mockInterface.Deauth.DeauthCan:
			if ok {
				t.Errorf("Deauth channel was not closed")
			}
		default:
			t.Errorf("Deauth channel was not closed")
		}
	}
}
