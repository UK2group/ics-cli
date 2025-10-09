package cmd

import (
	"testing"
)

func TestVersionVariables(t *testing.T) {
	// Test that version variables are properly initialized
	if Version == "" {
		t.Error("Version should not be empty")
	}

	// BuildTime and Commit may be "unknown" at test time, which is acceptable
	if BuildTime != "unknown" && BuildTime == "" {
		t.Error("BuildTime should either be 'unknown' or have a value")
	}

	if Commit != "unknown" && Commit == "" {
		t.Error("Commit should either be 'unknown' or have a value")
	}
}

func TestVersionFormat(t *testing.T) {
	// Version should follow semantic versioning pattern (at least X.Y.Z)
	if len(Version) < 5 {
		t.Errorf("Version format appears invalid: %s (should be at least X.Y.Z format)", Version)
	}

	// Check that it contains periods (basic semver check)
	foundPeriod := false
	for _, char := range Version {
		if char == '.' {
			foundPeriod = true
			break
		}
	}

	if !foundPeriod {
		t.Errorf("Version %s does not appear to follow semantic versioning format", Version)
	}
}
