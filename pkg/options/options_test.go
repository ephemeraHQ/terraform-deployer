package options_test

import (
	"testing"

	"github.com/xmtp-labs/terraform-deployer/pkg/options"
)

func TestCSVUnmarshalFlag(t *testing.T) {
	var c options.CSV
	err := c.UnmarshalFlag("xmtpd, xmtpd-prune")
	if err != nil {
		t.Fatalf("UnmarshalFlag failed: %v", err)
	}

	expected := []string{"xmtpd", "xmtpd-prune"}
	if len(c) != len(expected) {
		t.Fatalf("Expected length %d, got %d", len(expected), len(c))
	}
	for i := range c {
		if c[i] != expected[i] {
			t.Errorf("Index %d: got %s, want %s", i, c[i], expected[i])
		}
	}
}
