package cmd

import (
	"fmt"
	"os"

	"github.com/Joseph-D1/griphook/internal/vault"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm [service]",
	Short: "Remove a password from the vault.",
	Long:  `This command removes a password for a given service from the encrypted vault.`,
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

		if _, ok := v.Passwords[service]; !ok {
			fmt.Printf("Service %s not found in vault.\n", service)
			return
		}

		delete(v.Passwords, service)

		if err := vault.SaveVault(v, masterPassword); err != nil {
			fmt.Println("Error saving vault:", err)
			return
		}

		fmt.Printf("Service %s removed successfully!\n", service)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}