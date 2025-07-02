package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Joseph-D1/griphook/internal/vault"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [service]",
	Short: "Retrieve a password from the vault.",
	Long:  `This command retrieves and displays the password for a given service from the encrypted vault.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]

		masterPassword, err := readMasterPassword()
		if err != nil {
			fmt.Println(err)
			return
		}

		v, err := vault.LoadVault(masterPassword)
		if err != nil {
			fmt.Println("Error loading vault:", err)
			return
		}

		credentials, ok := v.Passwords[service]
		if !ok {
			fmt.Printf("No password found for service: %s\n", service)
			return
		}

		parts := strings.SplitN(credentials, ":", 2)
		if len(parts) != 2 {
			fmt.Printf("Invalid credentials format for %s\n", service)
			return
		}

		username := parts[0]
		password := parts[1]

		fmt.Printf("Username: %s\n", username)
		fmt.Printf("Password: %s\n", password)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
