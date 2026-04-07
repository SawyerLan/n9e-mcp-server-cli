package client

import (
	"testing"
	"time"
)

func TestNewClientWithOptionsUsesDefaultTimeout(t *testing.T) {
	c, err := NewClientWithOptions("token", "http://localhost:17000", "test-agent", Options{})
	if err != nil {
		t.Fatalf("NewClientWithOptions() error = %v", err)
	}

	if c.httpClient.Timeout != DefaultTimeout {
		t.Fatalf("http timeout = %s, want %s", c.httpClient.Timeout, DefaultTimeout)
	}
}

func TestNewClientWithOptionsUsesConfiguredTimeout(t *testing.T) {
	c, err := NewClientWithOptions("token", "http://localhost:17000", "test-agent", Options{
		Timeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewClientWithOptions() error = %v", err)
	}

	if c.httpClient.Timeout != 5*time.Second {
		t.Fatalf("http timeout = %s, want 5s", c.httpClient.Timeout)
	}
}
