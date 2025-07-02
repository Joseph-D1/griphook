# Griphook - A Secure CLI Password Manager

Griphook is a command-line interface (CLI) password manager designed for security and simplicity. It allows you to securely store, retrieve, list, and remove your passwords directly from your terminal.

## How it Works

Griphook stores your passwords in a single, encrypted file named `vault.enc` located in `~/.griphook/vault.enc`. This file is protected using strong AES-256-GCM encryption. Your master password is never stored directly; instead, it's used to derive the encryption key via the scrypt key derivation function, making it highly resistant to brute-force attacks.

## Commands

Here's a summary of the available commands:

*   **`./griphook init`**
    Initializes the password vault. This is the first command you should run. It will prompt you to set a master password, which will be used to encrypt and decrypt your vault.

*   **`./griphook add <service-name>`**
    Adds a new password entry for a specified service (e.g., `google`, `facebook`). You will be securely prompted to enter the username and password for that service.

*   **`./griphook get <service-name>`**
    Retrieves and displays the stored username and password for a given service. You will need to enter your master password to decrypt the vault.

*   **`./griphook ls`**
    Lists all the service names for which you have stored credentials in your vault.

*   **`./griphook rm <service-name>`**
    Removes a password entry for a specified service from your vault.

## Security

*   **AES-256-GCM Encryption:** Industry-standard symmetric encryption for your vault data.
*   **Scrypt Key Derivation:** Your master password is never stored. Instead, a strong encryption key is derived from it using scrypt, adding a significant layer of security against offline attacks.
*   **Secure Input:** All sensitive inputs (master password, username, password) are handled using `golang.org/x/term` to prevent them from being echoed to your terminal.
