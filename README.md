# GSD CLI

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

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd gsd-cli
   ```

2. Build the CLI:
   ```bash
   cd src
   go build -o gsd
   ```

3. Add the binary to your PATH (optional):
   ```bash
   mv gsd /usr/local/bin/
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
gsd switch [profile]
```
If no profile is specified, you will be prompted to select one.

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
gsd open console
```

Open the AWS SSO login page:
```bash
gsd open sso
```

Open a specific AWS service (e.g., Amplify):
```bash
gsd open amplify
```

Specify a profile when opening a service:
```bash
gsd open <service> --profile <profile-name>
```

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
   gsd open console
   ```
   Quickly access the AWS Management Console for the current account.

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

---

## License

This project is licensed under the [MIT License](LICENSE).
