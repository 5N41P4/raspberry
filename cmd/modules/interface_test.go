package modules

import (
	"testing"
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
