package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/Joseph-D1/griphook/internal/vault"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var addCmd = &cobra.Command{
	Use:   "add [service]",
	Short: "Add a new password to the vault.",
	Long:  `This command adds a new password for a given service to the encrypted vault.`,
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

		fmt.Printf("Enter username for %s: ", service)
		byteUsername, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println("Error reading username:", err)
			return
		}
		fmt.Println()

		fmt.Printf("Enter password for %s: ", service)
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println("Error reading password:", err)
			return
		}
		fmt.Println()

		v.Passwords[service] = fmt.Sprintf("%s:%s", string(byteUsername), string(bytePassword))

		if err := vault.SaveVault(v, masterPassword); err != nil {
			fmt.Println("Error saving vault:", err)
			return
		}

		fmt.Printf("Password for %s added successfully!\n", service)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
