package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check your connection to the Ingenuity Cloud Services API",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if API key exists in configuration
		apiKey := viper.GetString("api_key")
		if apiKey == "" {
			fmt.Println("Not logged in. Please run 'ics-cli auth login' to authenticate.")
			return
		}

		// Make API call to verify the connection
		client := &http.Client{Timeout: 10 * time.Second}
		req, err := http.NewRequest("GET", "https://api.ingenuitycloudservices.com/rest-api/user/details", nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating request:", err)
			return
		}

		req.Header.Add("X-Api-Token", apiKey)

		fmt.Println("Checking connection to Ingenuity Cloud Services API...")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error connecting to API:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			fmt.Println("API key is invalid or expired. Please run 'ics-cli auth login' to authenticate.")
			return
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Fprintf(os.Stderr, "API returned error status code: %d\n", resp.StatusCode)
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
			fmt.Println("Connection successful! (Unable to retrieve username)")
		} else {
			fmt.Printf("Connection successful! Logged in as %s\n", username)
		}
	},
}

func init() {
	authCmd.AddCommand(checkCmd)
}
