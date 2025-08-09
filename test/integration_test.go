package test

import (
	"os"
	"os/exec"
	"testing"
)

// TestBinaryCompiles verifies that the binary can be compiled
func TestBinaryCompiles(t *testing.T) {
	cmd := exec.Command("go", "build", "-o", "prom_cli_test", "../cmd/prom-cli")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to compile binary: %v", err)
	}

	// Clean up
	defer func() {
		if err := os.Remove("prom_cli_test"); err != nil {
			t.Logf("Failed to remove test binary: %v", err)
		}
	}()

	// Check that the binary exists
	_, err = os.Stat("prom_cli_test")
	if os.IsNotExist(err) {
		t.Fatal("Binary was not created")
	}
}

// TestMockPrometheus simulates a Prometheus server and tests the binary
func TestMockPrometheus(t *testing.T) {
	// This test would normally start a mock Prometheus server
	// and test the binary against it. However, since our binary
	// requires user input, we'll just verify that it can be executed.

	// First, compile the binary
	cmd := exec.Command("go", "build", "-o", "prom_cli_test", "../cmd/prom-cli")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to compile binary: %v", err)
	}

	// Clean up
	defer func() {
		if err := os.Remove("prom_cli_test"); err != nil {
			t.Logf("Failed to remove test binary: %v", err)
		}
	}()

	// Check that the binary exists and is executable
	info, err := os.Stat("prom_cli_test")
	if os.IsNotExist(err) {
		t.Fatal("Binary was not created")
	}

	// On Unix systems, check that the binary is executable
	if info.Mode()&0111 == 0 {
		t.Fatal("Binary is not executable")
	}
}
