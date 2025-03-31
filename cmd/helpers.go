package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

// makeAPIRequest is a generic function to handle API calls with proper error handling
func makeAPIRequest(method string, timeout time.Duration, url string, body io.Reader, result interface{}) error {

	client := &http.Client{Timeout: timeout * time.Second}

	// Check if API key exists in configuration
	apiKey := viper.GetString("api_key")
	if apiKey == "" {
		return fmt.Errorf("not logged in. Please run 'ics-cli auth login' to authenticate")
	}

	// Make API call
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("X-Api-Token", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error connecting to API: %w", err)
	}
	defer resp.Body.Close()

	// Handle status codes
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("API key is invalid or expired")
	case http.StatusNotFound:
		return fmt.Errorf("resource not found: %s", url)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned error status code: %d", resp.StatusCode)
	}

	// Read and parse the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading API response: %w", err)
	}

	if err := json.Unmarshal(respBody, result); err != nil {
		return fmt.Errorf("error parsing API response: %w", err)
	}

	return nil
}

// openBrowser opens a URL in the default browser
func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}

// BlueHeading returns a string formatted as a blue heading
func BlueHeading(text string) string {
	return color.New(color.FgBlue, color.Bold).Sprint(text)
}

// WhiteText returns a string formatted as white text
func WhiteText(text string) string {
	return color.New(color.FgWhite).Sprint(text)
}

// GreenText returns a string formatted as green text
func GreenText(text string) string {
	return color.New(color.FgGreen).Sprint(text)
}

// RedText returns a string formatted as red text
func RedText(text string) string {
	return color.New(color.FgRed).Sprint(text)
}

// YellowText returns a string formatted as yellow text
func YellowText(text string) string {
	return color.New(color.FgYellow).Sprint(text)
}
