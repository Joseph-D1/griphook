package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/Joseph-D1/griphook/internal/vault"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the password vault.",
	Long:  `This command initializes the encrypted vault file where your passwords will be stored. You will be prompted to create a master password.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Enter master password: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println("Error reading password:", err)
			return
		}
		fmt.Println()

		fmt.Print("Confirm master password: ")
		byteConfirmPassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println("Error reading confirmation password:", err)
			return
		}
		fmt.Println()

		if string(bytePassword) != string(byteConfirmPassword) {
			fmt.Println("Passwords do not match.")
			return
		}

		v := &vault.Vault{Passwords: make(map[string]string)}
		if err := vault.SaveVault(v, bytePassword); err != nil {
			fmt.Println("Error saving vault:", err)
			return
		}

		fmt.Println("Vault initialized successfully!")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}