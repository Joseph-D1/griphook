package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var rootCmd = &cobra.Command{
	Use:   "griphook",
	Short: "A simple, secure password manager CLI.",
	Long:  `Griphook is a CLI tool for managing your passwords securely. It uses AES-256 encryption to protect your data.`,
}

func readMasterPassword() ([]byte, error) {
	fmt.Print("Enter master password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, fmt.Errorf("error reading password: %w", err)
	}
	fmt.Println()
	return bytePassword, nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI: %s", err)
		os.Exit(1)
	}
}
