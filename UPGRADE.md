# Upgrade Notes

## Interactive Configuration Interface

### Overview
The CLI configuration interface has been enhanced with an interactive prompt-based system to improve user experience. This replaces the previous flag-based configuration approach with a more intuitive step-by-step process.

### New Dependencies
```bash
go get github.com/AlecAivazis/survey/v2
```

### Changes

#### Profile Configuration
- Replaced command-line flags with interactive prompts
- Added visual feedback with emoji icons and colors
- Implemented step-by-step configuration process
- Added input validation for fields like Account ID
- Enhanced error handling and user feedback

#### New Features
1. Interactive Profile Creation
   - Guided profile setup
   - Region selection from predefined list
   - Input validation for critical fields
   - Visual confirmation of successful operations

2. Interactive Profile Removal
   - Profile selection from existing configurations
   - Confirmation prompt before deletion
   - Visual feedback for operation status

### Example Usage

#### Adding a Profile
```bash
gsd config add
```
The command now presents interactive prompts:
- Profile name input
- Region selection menu
- SSO start URL input
- AWS Account ID input with validation
- IAM Role name input

#### Removing a Profile
```bash
gsd config remove
```
The command now shows:
- Interactive profile selection menu
- Deletion confirmation prompt
- Operation status feedback

### Visual Styling
- Added minimal, modern interface design
- Implemented consistent color scheme
- Added emoji indicators for different operations:
  - ⚡ for configuration prompts
  - ✨ for success messages
  - ❌ for errors
  - 🗑 for deletion operations
  - → for selection indicators
