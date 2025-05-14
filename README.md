
# GSD CLI

**GSD CLI** (Get Shit Done) is a user-friendly command-line tool designed to simplify authentication and credential management for AWS services. It provides an intuitive interface for managing AWS profiles, authenticating via SSO, and opening the AWS Console or SSO login page, so you can focus on getting things done.

---

## Features

- **Simplified Authentication**: Authenticate with AWS using SSO or static credentials.
- **Profile Management**: Easily switch between AWS profiles.
- **Credential Validation**: Verify the validity of your current credentials.
- **Quick Access to AWS Console**: Open the AWS Management Console or SSO login page directly from the CLI.
- **User-Friendly Commands**: Intuitive commands with clear error messages and helpful output.

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

Authenticate with AWS using the `auth login` command:
```bash
gsd auth login
```
This command automatically detects your AWS SSO profile or prompts you to authenticate.

### Profile Management

List available AWS profiles:
```bash
gsd auth list
```

Switch to a different profile:
```bash
gsd auth switch <profile-name>
```

### Credential Validation

Validate your current credentials:
```bash
gsd auth validate
```

### Open AWS Console or SSO Login Page

Open the AWS Management Console for the current account:
```bash
gsd open
```

Open the AWS SSO login page:
```bash
gsd open --sso
```

---

## Example Workflow

1. **Login**:
   ```bash
   gsd auth login
   ```
   Authenticate using your SSO profile or credentials.

2. **Switch Profile**:
   ```bash
   gsd auth switch my-profile
   ```
   Switch to a different AWS profile.

3. **Validate Credentials**:
   ```bash
   gsd auth validate
   ```
   Ensure your credentials are valid and ready to use.

4. **Open AWS Console**:
   ```bash
   gsd open
   ```
   Quickly access the AWS Management Console for the current account.

---

## Roadmap

- Add support for additional AWS services (e.g., S3, EC2).
- Improve output formatting (e.g., JSON, tables).
- Integrate secure credential storage.

---

## Contributing

Contributions are welcome! Please submit a pull request or open an issue to suggest improvements or report bugs.

---

## License

This project is licensed under the [MIT License](LICENSE).
