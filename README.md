# Squirrel CLI - A Local Password Manager

Squirrel is a command-line tool for securely managing passwords and sensitive information. It is designed for tech-savvy users and individuals who want full control over their data without relying on cloud services. **All data is encrypted locally using AES encryption with a key derived from the user's master password.** This means that your passwords are only as secure as your master password, and **if the master password is lost, all data will be irretrievable**. Squirrel ensures maximum security by keeping everything stored and encrypted on your local machine.

## Features

- AES encryption to protect sensitive information
- Key derivation using PBKDF2 or scrypt to ensure secure password storage
- Works across macOS, Linux, and Windows
- Simple CLI interface for ease of use

## Installation

### Prerequisites

Ensure you have Go installed. If not, you can download and install it from [Go's official website](https://golang.org/dl/).

### Install Squirrel

1. Clone the repository:

   ```bash
   git clone https://github.com/ehsun7b/squirrel.git
   ```

2. Navigate into the project directory:

   ```bash
   cd squirrel
   ```

3. Build the project:

   ```bash
   go build -o squirrel
   ```

4. Move the binary to a directory in your `PATH` (optional):

   ```bash
   sudo mv squirrel /usr/local/bin/
   ```

## Usage

To start using Squirrel, run the following command in your terminal:

```bash
squirrel
```

You will be prompted to create a master password. **Ensure you remember this password as it will be required to encrypt and decrypt your stored data. If lost, your data cannot be recovered.**

### Adding an Entry

To add a new password or entry:

TODO: fix command for adding entry

```bash
squirrel add
```

Follow the prompts to input the necessary information.

### Retrieving an Entry

To retrieve a stored password:

TODO: fix command for retrieving an entry

```bash
squirrel get
```

You will be prompted for your master password to decrypt the stored information.

### Listing All Entries

To list all stored entries:

```bash
squirrel list
```

### Deleting an Entry

To delete an entry:

```bash
squirrel delete
```

### Help

For additional commands and usage information:

```bash
squirrel --help
```

## Security Considerations

Squirrel stores all encrypted data on your local machine and never sends your data over the internet. The encryption key is derived from your master password using PBKDF2 or scrypt to ensure strong protection against brute-force attacks.

**Warning:** If you forget your master password, there is no way to recover your data. The encryption is designed to be secure, so there are no backdoors.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
