# GSD CLI

[![Go Version](https://img.shields.io/github/go-mod/go-version/aphexlog/gsd)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/aphexlog/gsd)](https://github.com/aphexlog/gsd/releases/latest)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/github/actions/workflow/status/aphexlog/gsd/release.yml)](https://github.com/aphexlog/gsd/actions)

**GSD CLI** (Get Shit Done) is a user-friendly command-line tool designed to simplify authentication and credential management for AWS services. It provides an intuitive interface for managing AWS profiles, authenticating via SSO, and accessing AWS services, so you can focus on getting things done.

---

## Features

- **Interactive Configuration**: User-friendly interactive prompts for managing AWS profiles
- **Simplified Authentication**: Authenticate with AWS using SSO or static credentials
- **Profile Management**: Easily switch between AWS profiles
- **Credential Validation**: Verify the validity of your current credentials
- **Quick Access to AWS Services**: Open the AWS Management Console, SSO login page, or specific AWS services directly from the CLI
- **User-Friendly Commands**: Intuitive commands with clear error messages and helpful output

---

## Installation

### Option 1: Download Binary (Recommended)

Visit the [Releases](https://github.com/aphexlog/gsd/releases) page and download the latest version for your platform:

- For macOS:
  ```bash
  # For Apple Silicon (M1/M2)
  curl -LO https://github.com/aphexlog/gsd/releases/latest/download/gsd_Darwin_arm64.tar.gz
  tar xzf gsd_Darwin_arm64.tar.gz
  sudo mv gsd /usr/local/bin/

  # For Intel Macs
  curl -LO https://github.com/aphexlog/gsd/releases/latest/download/gsd_Darwin_x86_64.tar.gz
  tar xzf gsd_Darwin_x86_64.tar.gz
  sudo mv gsd /usr/local/bin/
  ```

- For Linux:
  ```bash
  # For AMD64
  curl -LO https://github.com/aphexlog/gsd/releases/latest/download/gsd_Linux_x86_64.tar.gz
  tar xzf gsd_Linux_x86_64.tar.gz
  sudo mv gsd /usr/local/bin/

  # For ARM64
  curl -LO https://github.com/aphexlog/gsd/releases/latest/download/gsd_Linux_arm64.tar.gz
  tar xzf gsd_Linux_arm64.tar.gz
  sudo mv gsd /usr/local/bin/
  ```

- For Windows:
  Download the appropriate ZIP file from the releases page and add it to your PATH.

### Option 2: Build from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/aphexlog/gsd.git
   cd gsd
   ```

2. Build the CLI:
   ```bash
   go build -o gsd
   ```

3. Add the binary to your PATH (optional):
   ```bash
   sudo mv gsd /usr/local/bin/
   ```

---

## Usage

### Authentication

Authenticate with AWS using the `login` command:
```bash
gsd login [profile]
```
If no profile is specified, the default profile is used. This command automatically detects your AWS SSO profile or prompts you to authenticate.

### Profile Management

Switch to a different profile:
```bash
gsd switch
```
You will be presented with an interactive menu to select your profile.

### Configuration Management

List all configured profiles:
```bash
gsd config ls
```

Add a new profile interactively:
```bash
gsd config add
```
The CLI will guide you through:
- Profile name selection
- AWS region selection
- Authentication method choice (SSO or Access Keys)
- Required configuration details

Remove an existing profile:
```bash
gsd config remove
```
Provides an interactive menu to:
- Select the profile to remove
- Confirm deletion

Edit an existing profile:
```bash
gsd config edit
```
Interactive interface to:
- Select the profile to edit
- Choose what to modify (Region, SSO Configuration, or Access Keys)
- Update the selected configuration

### Open AWS Services

Open the AWS Management Console for the current account:
```bash
gsd open
```
Select "Console (Main)" from the interactive menu.

Open the AWS SSO login page:
```bash
gsd open
```
Select "SSO" from the interactive menu.

Open a specific AWS service:
```bash
gsd open
```
Select the desired service from the interactive menu.

### Credential Validation

Check the currently authenticated profile and credentials:
```bash
gsd whoami
```

---

## Example Workflow

1. **Add a New Profile**:
   ```bash
   gsd config add
   ```
   Follow the interactive prompts to configure your profile.

2. **Login**:
   ```bash
   gsd login my-profile
   ```
   Authenticate using your configured profile.

3. **Switch Profile**:
   ```bash
   gsd switch
   ```
   Select a different AWS profile from the interactive menu.

4. **Open AWS Console**:
   ```bash
   gsd open
   ```
   Select "Console (Main)" to access the AWS Management Console.

5. **Check Current Profile**:
   ```bash
   gsd whoami
   ```
   View the currently authenticated profile and credentials.

---

## Roadmap

- Add support for additional AWS services
- Improve output formatting (e.g., JSON, tables)
- Integrate secure credential storage
- Add profile import/export functionality
- Enhance interactive configuration options

---

## Contributing

Contributions are welcome! Please submit a pull request or open an issue to suggest improvements or report bugs.

Created by [aphexlog](https://github.com/aphexlog)

---

## License

This project is licensed under the [MIT License](LICENSE).
