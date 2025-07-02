# Griphook - Technical Documentation for Language Models

## 1. Project Overview

Griphook is a command-line interface (CLI) password manager developed in Go. Its primary purpose is to provide a secure, local solution for users to store, retrieve, list, and manage their sensitive credentials (usernames and passwords) in an encrypted vault file. The design prioritizes security, ease of use via the command line, and self-containment.

## 2. Core Functionality

Griphook exposes several subcommands, each handling a specific aspect of password management:

*   **`griphook init`**
    *   **Purpose:** Initializes the encrypted vault file (`vault.enc`). This is the first command a user must execute.
    *   **Process:** Prompts the user to create and confirm a master password. This master password is then used to derive an encryption key, which encrypts an empty vault structure and saves it to `~/.griphook/vault.enc`.

*   **`griphook add <service-name>`**
    *   **Purpose:** Adds a new username and password entry for a specified service (e.g., `github`, `mybank`).
    *   **Process:** Requires the user to provide the master password to decrypt the vault. It then securely prompts for the username and password for the given service. These are stored as a `username:password` string within the vault's internal map. The updated vault is then re-encrypted and saved.

*   **`griphook get <service-name>`**
    *   **Purpose:** Retrieves and displays the stored username and password for a specific service.
    *   **Process:** Requires the master password to decrypt the vault. It then looks up the `username:password` string for the requested service, parses it, and prints the username and password to the console.

*   **`griphook ls`**
    *   **Purpose:** Lists all the service names currently stored in the vault.
    *   **Process:** Requires the master password to decrypt the vault. It then iterates through the keys (service names) of the vault's internal password map and prints each service name.

*   **`griphook rm <service-name>`**
    *   **Purpose:** Removes a password entry for a specified service from the vault.
    *   **Process:** Requires the master password to decrypt the vault. It deletes the entry corresponding to the service name from the vault's internal map. The modified vault is then re-encrypted and saved.

## 3. Technical Architecture and Implementation Details

### 3.1. Language and Frameworks
*   **Go (Golang):** The entire application is written in Go, chosen for its performance, strong type system, concurrency features, and ability to compile into a single, self-contained binary.
*   **Cobra:** Used for building the command-line interface. Cobra provides a robust framework for defining commands, subcommands, flags, and handling command execution flow.
    *   `main.go`: Entry point, calls `cmd.Execute()`.
    *   `cmd/root.go`: Defines the base `griphook` command and sets up global flags/functions (like `readMasterPassword`).
    *   `cmd/*.go`: Each subcommand (`init`, `add`, `get`, `ls`, `rm`) is defined in its own file within the `cmd` directory, registered with the `rootCmd`.

### 3.2. Security Mechanisms
*   **Encryption Standard (AES-256-GCM):** All vault data is encrypted using the Advanced Encryption Standard (AES) with a 256-bit key in Galois/Counter Mode (GCM). AES-GCM provides both confidentiality and authenticity (integrity) of the data.
    *   **`crypto/aes`:** Go's standard library package for AES block cipher.
    *   **`crypto/cipher`:** Provides interfaces for stream ciphers and authenticated encryption (like GCM).
*   **Key Derivation Function (scrypt):** The master password provided by the user is never stored directly. Instead, a cryptographically strong key is derived from it using `scrypt` (`golang.org/x/crypto/scrypt`). Scrypt is a password-based key derivation function designed to be computationally intensive, making brute-force attacks significantly harder.
    *   A unique, randomly generated salt is used for each key derivation, stored alongside the encrypted data.
*   **Secure Input (`golang.org/x/term`):** When the user enters sensitive information (master password, username, password), the `term.ReadPassword` function is used. This function prevents the input from being echoed to the terminal, protecting against shoulder-surfing and logging.

### 3.3. Vault Structure and Storage
*   **Vault File:** The encrypted data is stored in a single file: `~/.griphook/vault.enc`.
*   **Internal Data Structure:** Within the Go application, the vault's plaintext content is represented by a `Vault` struct:
    ```go
    type Vault struct {
        Passwords map[string]string `json:"passwords"`
    }
    ```
    This `Passwords` map stores service names (string keys) and their corresponding `username:password` (string values).
*   **Serialization:** The `Vault` struct is marshaled to/from JSON format before encryption/decryption using Go's `encoding/json` package.
*   **File Format:** The `vault.enc` file contains the randomly generated salt prepended to the AES-GCM ciphertext. This allows the application to retrieve the salt needed for key derivation during decryption.

### 3.4. Data Flow (Encryption/Decryption)
1.  **Encryption (`SaveVault`):**
    *   The `Vault` struct is marshaled into a JSON byte array (plaintext).
    *   A random `salt` is generated.
    *   An encryption `key` is derived from the user's master password and the `salt` using `scrypt`.
    *   The plaintext is encrypted using AES-256-GCM with the derived `key`.
    *   The `salt` is prepended to the resulting ciphertext.
    *   The combined `salt + ciphertext` is written to `~/.griphook/vault.enc`.
2.  **Decryption (`LoadVault`):**
    *   The `vault.enc` file is read.
    *   The `salt` is extracted from the beginning of the file.
    *   An encryption `key` is derived from the user's master password and the extracted `salt` using `scrypt`.
    *   The remaining ciphertext is decrypted using AES-256-GCM with the derived `key`.
    *   The decrypted plaintext (JSON) is unmarshaled back into the `Vault` struct.

## 4. Project Directory Structure

```
griphook/
├── main.go             # Application entry point
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── README.md           # User-facing README
├── docs.md             # This detailed documentation
├── .gitignore          # Git ignore rules (e.g., for 'griphook' binary, .griphook/)
├── cmd/                # Contains Cobra command definitions
│   ├── root.go         # Root command and shared functions (e.g., readMasterPassword)
│   ├── add.go          # 'add' subcommand implementation
│   ├── get.go          # 'get' subcommand implementation
│   ├── init.go         # 'init' subcommand implementation
│   ├── ls.go           # 'ls' subcommand implementation
│   └── rm.go           # 'rm' subcommand implementation
└── internal/           # Internal packages not exposed externally
    └── vault/          # Core encryption and vault management logic
        └── vault.go    # Defines Vault struct, encryption/decryption, key derivation
```
