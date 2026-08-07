package payment

import (
	"testing"
)

func TestProviderRegistry(t *testing.T) {
	providers := []string{"cashea", "spidi", "wayupay", "biopago"}
	for _, name := range providers {
		p := GetProvider(name)
		if p == nil {
			t.Errorf("GetProvider(%q) returned nil", name)
		}
	}

	if p := GetProvider("unknown"); p != nil {
		t.Error("GetProvider(\"unknown\") should return nil")
	}
}
