package cloudinary

import (
	"context"
	"testing"
)

func TestUnconfiguredClient(t *testing.T) {
	c, err := New("", "", "")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if c.IsConfigured() {
		t.Fatalf("empty credentials must not configure the client")
	}
	if _, err := c.SignUpload("u", "tr"); err == nil {
		t.Fatalf("expected not-configured error")
	}
	if _, err := c.DeliveryURL("x"); err == nil {
		t.Fatalf("expected not-configured error")
	}
	if err := c.Destroy(context.Background(), "x"); err == nil {
		t.Fatalf("expected not-configured error")
	}
}
