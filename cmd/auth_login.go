package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

// UserResponse represents the structure of the API response
type UserResponse struct {
	Data struct {
		UserProfile struct {
			Username string `json:"username"`
		} `json:"userProfile"`
	} `json:"data"`
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to your Ingenuity Cloud Services Account (requires an API Key)",
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, err := promptForAPIKey(cmd)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading API key:", err)
			return
		}

		// Verify API key by making a test API call
		client := &http.Client{Timeout: 10 * time.Second}
		req, err := http.NewRequest("GET", "https://api.ingenuitycloudservices.com/rest-api/user/details", nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating request:", err)
			return
		}

		req.Header.Add("X-Api-Token", apiKey)

		fmt.Println("Verifying API key...")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error connecting to API:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			fmt.Fprintln(os.Stderr, "Error: Invalid API key. Authentication failed.")
			return
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Fprintf(os.Stderr, "Error: API returned status code %d\n", resp.StatusCode)
			return
		}

		// Read and parse the response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading API response:", err)
			return
		}

		var userResp UserResponse
		if err := json.Unmarshal(body, &userResp); err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing API response:", err)
			return
		}

		username := userResp.Data.UserProfile.Username
		if username == "" {
			fmt.Fprintln(os.Stderr, "Warning: Could not get username from API response")
		}

		// Store API key in Viper configuration
		viper.Set("api_key", apiKey)

		// Save the config file
		if err := viper.WriteConfig(); err != nil {
			// If the config file doesn't exist, create it
			if err := viper.SafeWriteConfig(); err != nil {
				fmt.Fprintln(os.Stderr, "Error saving API key to config:", err)
				return
			}
		}

		if username != "" {
			fmt.Printf("Successfully logged in as %s\n", username)
		} else {
			fmt.Println("Successfully logged in to Ingenuity Cloud Services!")
		}
	},
}

// promptForAPIKey asks the user to input their API key
func promptForAPIKey(cmd *cobra.Command) (string, error) {
	// Check if API key was provided via flag
	if apiKey, _ := cmd.Flags().GetString("key"); apiKey != "" {
		return apiKey, nil
	}

	// Otherwise, prompt user for API key
	fmt.Print("Enter your Ingenuity Cloud Services API key: ")

	// Try to read API key securely (hidden input)
	apiKeyBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		// Fall back to regular input if secure input fails
		reader := bufio.NewReader(os.Stdin)
		apiKey, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		fmt.Println() // Add newline after input
		return strings.TrimSpace(apiKey), nil
	}

	fmt.Println() // Add newline after hidden input
	return strings.TrimSpace(string(apiKeyBytes)), nil
}

func init() {
	authCmd.AddCommand(loginCmd)

	// Add flag for providing API key directly via command line
	loginCmd.Flags().StringP("key", "k", "", "API key for Ingenuity Cloud Services")
}
