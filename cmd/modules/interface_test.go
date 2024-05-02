package modules

import (
	"testing"
)

// TestGetInterfaces tests the GetInterfaces function
func TestGetInterfaces(t *testing.T) {
	interfaces := GetInterfaces(2, "wlan")

	if interfaces == nil {
		t.Errorf("GetInterfaces() = %v; want nil", interfaces)
	}

	if len(interfaces) != 2 {
		t.Errorf("GetInterfaces() = %v; want 2", len(interfaces))
	}

}
