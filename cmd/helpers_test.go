package cmd

import (
	"os"
	"runtime"
	"testing"
)

func TestGetAPIBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		want     string
	}{
		{
			name:     "default production URL when env not set",
			envValue: "",
			want:     "https://api.ingenuitycloudservices.com",
		},
		{
			name:     "custom URL from environment variable",
			envValue: "https://staging.api.example.com",
			want:     "https://staging.api.example.com",
		},
		{
			name:     "localhost URL for testing",
			envValue: "http://localhost:8080",
			want:     "http://localhost:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env value and restore after test
			originalEnv := os.Getenv("ICS_API_URL")
			defer os.Setenv("ICS_API_URL", originalEnv)

			// Set test environment
			if tt.envValue != "" {
				os.Setenv("ICS_API_URL", tt.envValue)
			} else {
				os.Unsetenv("ICS_API_URL")
			}

			got := getAPIBaseURL()
			if got != tt.want {
				t.Errorf("getAPIBaseURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlueHeading(t *testing.T) {
	text := "Test Heading"
	result := BlueHeading(text)

	// The result should contain the original text
	if result == "" {
		t.Error("BlueHeading() returned empty string")
	}

	// Function should at minimum return the text (may or may not have ANSI codes depending on environment)
	if len(result) < len(text) {
		t.Error("BlueHeading() result is shorter than input text")
	}
}

func TestWhiteText(t *testing.T) {
	text := "Test Text"
	result := WhiteText(text)

	if result == "" {
		t.Error("WhiteText() returned empty string")
	}

	if len(result) < len(text) {
		t.Error("WhiteText() result is shorter than input text")
	}
}

func TestGreenText(t *testing.T) {
	text := "Success"
	result := GreenText(text)

	if result == "" {
		t.Error("GreenText() returned empty string")
	}

	if len(result) < len(text) {
		t.Error("GreenText() result is shorter than input text")
	}
}

func TestRedText(t *testing.T) {
	text := "Error"
	result := RedText(text)

	if result == "" {
		t.Error("RedText() returned empty string")
	}

	if len(result) < len(text) {
		t.Error("RedText() result is shorter than input text")
	}
}

func TestYellowText(t *testing.T) {
	text := "Warning"
	result := YellowText(text)

	if result == "" {
		t.Error("YellowText() returned empty string")
	}

	if len(result) < len(text) {
		t.Error("YellowText() result is shorter than input text")
	}
}

func TestOpenBrowser(t *testing.T) {
	// This test verifies that the function returns expected errors for unsupported platforms
	// We can't actually test browser opening in CI/CD

	tests := []struct {
		name       string
		url        string
		shouldFail bool
	}{
		{
			name:       "valid URL",
			url:        "https://example.com",
			shouldFail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := openBrowser(tt.url)

			// On supported platforms (linux, windows, darwin), the command should be attempted
			// We check that unsupported platforms return an error
			if runtime.GOOS != "linux" && runtime.GOOS != "windows" && runtime.GOOS != "darwin" {
				if err == nil {
					t.Error("openBrowser() should return error on unsupported platform")
				}
			}
		})
	}
}
