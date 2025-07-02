# Griphook - Development Roadmap

This roadmap outlines potential future features and considerations for the Griphook CLI password manager, along with implementation ideas and deployment strategies.

## Phase 1: Core Enhancements

This phase focuses on adding highly requested and valuable features that enhance the immediate usability of Griphook.

### 1. Password Generation

**Goal:** Allow users to generate strong, random passwords directly from the CLI.

**Implementation Ideas:**
*   **Command:** Introduce a new command, e.g., `griphook generate`.
*   **Options:** Support flags for specifying password length, inclusion of numbers, symbols, uppercase/lowercase letters (e.g., `--length 16 --numbers --symbols`).
*   **Logic:** Implement a function that generates a cryptographically secure random string based on the specified criteria.
*   **Integration:** The generated password could optionally be added directly to the vault using the `add` command, or simply printed to stdout.

### 2. Copy to Clipboard

**Goal:** Enable users to copy retrieved passwords or generated passwords directly to their system clipboard for convenience and to avoid displaying sensitive information on screen.

**Implementation Ideas:**
*   **Cross-Platform Library:** Utilize a Go library for cross-platform clipboard access (e.g., `github.com/atotto/clipboard`).
*   **Integration with `get`:** Add a flag to the `get` command, e.g., `griphook get <service> --copy`. If this flag is present, the password is copied to the clipboard instead of being printed to stdout.
*   **Integration with `generate`:** The `generate` command could also have a `--copy` flag to copy the newly generated password.
*   **Security Consideration:** Implement a mechanism to clear the clipboard after a short delay (e.g., 30-60 seconds) for enhanced security, if the chosen library supports it.

## Phase 2: Advanced Features

This phase explores more complex features that enhance the overall management and security capabilities of Griphook.

### 1. Password Strength Indicator

**Goal:** Provide feedback on the strength of user-provided passwords or generated passwords.

**Implementation Ideas:**
*   **Library:** Integrate a Go library for password strength assessment (e.g., `github.com/nbutton/password-strength`).
*   **Integration:** Display strength feedback during `init` (for master password), `add` (for service passwords), and `generate` commands.

### 2. Import/Export Functionality

**Goal:** Allow users to import passwords from other formats (e.g., CSV) or export their Griphook vault data.

**Implementation Ideas:**
*   **Commands:** `griphook import <file>` and `griphook export <file>`.
*   **Format:** Define a secure, structured format for export (e.g., encrypted JSON or a custom binary format).
*   **Security:** Ensure imported/exported data is handled securely, especially during decryption/encryption processes.

### 3. Password History

**Goal:** Keep a history of old passwords for a service, allowing users to revert or review past credentials.

**Implementation Ideas:**
*   **Vault Structure:** Modify the `Vault` struct to store a slice of passwords for each service, along with timestamps.
*   **Commands:** Add commands like `griphook history <service>` or `griphook revert <service> <timestamp>`.

### 4. Two-Factor Authentication (2FA) Integration

**Goal:** Add support for storing and generating Time-based One-Time Passwords (TOTP) for 2FA.

**Implementation Ideas:**
*   **Library:** Use a Go library for TOTP generation (e.g., `github.com/pquerna/otp`).
*   **Vault Structure:** Extend the vault to store TOTP secrets.
*   **Commands:** `griphook totp add <service> <secret>` and `griphook totp get <service>`.

## Phase 3: Deployment & Distribution

This phase focuses on making Griphook easily accessible and usable across different platforms.

### 1. Cross-Platform Binaries

**Goal:** Provide pre-compiled binaries for common operating systems (Linux, macOS, Windows).

**Implementation Ideas:**
*   **Go's Cross-Compilation:** Leverage Go's built-in cross-compilation capabilities (`GOOS`, `GOARCH` environment variables).
*   **Build Script:** Create a simple shell script or a `Makefile` to automate the build process for different targets.

### 2. Release Management

**Goal:** Automate the release process, including tagging, binary creation, and GitHub Releases.

**Implementation Ideas:**
*   **GitHub Actions:** Set up a CI/CD pipeline using GitHub Actions to:
    *   Run tests on push.
    *   Build cross-platform binaries on new tag creation.
    *   Create a GitHub Release with attached binaries.

### 3. Package Managers (Future Consideration)

**Goal:** Distribute Griphook via popular package managers for easier installation.

**Implementation Ideas:**
*   **Homebrew (macOS):** Create a Homebrew formula.
*   **APT/YUM (Linux):** Explore creating `.deb` or `.rpm` packages, or providing a simple installation script.
*   **Scoop/Chocolatey (Windows):** Investigate Windows package managers.

### 4. Documentation Hosting

**Goal:** Host comprehensive documentation online.

**Implementation Ideas:**
*   **GitHub Pages:** Use GitHub Pages to host the `docs.md` content (converted to HTML) or a dedicated documentation site generator.
