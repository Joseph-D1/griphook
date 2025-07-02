package cmd

import (
	"fmt"
	"os"

	"github.com/Joseph-D1/griphook/internal/vault"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all services in the vault.",
	Long:  `This command lists all the service names for which passwords are stored in the encrypted vault.`,
	Run: func(cmd *cobra.Command, args []string) {
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

		if len(v.Passwords) == 0 {
			fmt.Println("No services found in the vault.")
			return
		}

		fmt.Println("Services in vault:")
		for service := range v.Passwords {
			fmt.Printf("- %s\n", service)
		}
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
