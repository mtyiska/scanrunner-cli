# Contributing to SCANRUNNER-CLI

Thank you for considering contributing to **SCANRUNNER-CLI**! We welcome contributions from everyone, whether you're fixing a bug, proposing new features, improving documentation, or suggesting ideas. This guide will help you get started.

---

## Table of Contents
- [How to Contribute](#how-to-contribute)
- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Reporting Issues](#reporting-issues)
- [Submitting Pull Requests](#submitting-pull-requests)
- [Coding Guidelines](#coding-guidelines)
- [Development Workflow](#development-workflow)

---

## How to Contribute
1. Fork the repository to your GitHub account.
2. Clone your fork to your local machine.
3. Create a new branch for your changes.
4. Commit your changes and push them to your fork.
5. Submit a pull request (PR) to the main repository.

---

## Code of Conduct
We expect contributors to adhere to our [Code of Conduct](CODE_OF_CONDUCT.md). Be respectful and considerate in all interactions.

---

## Getting Started
### Prerequisites
- **Go** version 1.20 or later installed.
- Familiarity with Git, GitHub, and command-line tools.
- Recommended tools: `make` (for automation), `docker` (for testing environments).

### Setting Up the Project
1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/scanrunner-cli.git
   cd scanrunner-cli
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the project:
   ```bash
   go build -o scanrunner main.go
   ```

4. Run tests to ensure everything works:
   ```bash
   go test ./...
   ```

---

## Reporting Issues
If you find a bug or have a feature request, please file an issue in the GitHub repository.

### Steps for Filing Issues
1. Check if the issue already exists in the [issues list](https://github.com/your-username/scanrunner-cli/issues).
2. If not, open a new issue and provide:
   - A clear title.
   - Steps to reproduce (if applicable).
   - Expected and actual behavior.
   - Any relevant logs, screenshots, or code snippets.

---

## Submitting Pull Requests
1. **Fork and Clone**:
   - Fork the repository to your GitHub account.
   - Clone your fork locally:
     ```bash
     git clone https://github.com/your-username/scanrunner-cli.git
     ```

2. **Create a Branch**:
   - Use descriptive branch names like `fix-bug-xyz` or `add-new-feature`.
     ```bash
     git checkout -b branch-name
     ```

3. **Make Changes**:
   - Follow the coding guidelines outlined below.

4. **Commit Changes**:
   - Write clear commit messages:
     ```bash
     git commit -m "Fix: Correct YAML parsing error in scan command"
     ```

5. **Push to Your Fork**:
   ```bash
   git push origin branch-name
   ```

6. **Submit a Pull Request**:
   - Open a PR against the `main` branch.
   - Include a description of the changes, a reference to the issue (if applicable), and testing steps.

---

## Coding Guidelines
- **Code Style**:
  - Follow the Go standard coding style (use `gofmt`).
  - Write clean, self-documenting code with clear variable and function names.

- **Testing**:
  - Write tests for any new functionality or bug fixes.
  - Place tests in the appropriate `test/` subdirectory (e.g., `test/fileparser_test.go`).

- **Documentation**:
  - Update relevant documentation (e.g., `README.md`, `docs/`) for any new features or changes.
  - Ensure CLI commands include clear descriptions and usage examples.

- **Commit Messages**:
  - Use descriptive messages:
    ```
    Type: Short description (50 chars max)

    Optional detailed explanation (if needed, wrapped at 72 chars).
    ```

    Examples:
    - `Feat: Add AI-based suggestion generator`
    - `Fix: Resolve JSON parsing error in validate command`

---

## Development Workflow
### Running the CLI
Run the built CLI tool locally:
```bash
./scanrunner scan --path ./example-files
```

### Running Tests
Run all tests:
```bash
go test ./...
```

Run specific tests:
```bash
go test ./internal/compliance -v
```

### Building the CLI
Build the CLI binary:
```bash
go build -o scanrunner main.go
```

### Using `make` (Optional)
A `Makefile` is included for automation:
```bash
make build       # Build the CLI
make test        # Run all tests
make lint        # Run linters
```

---

## Need Help?
- Open an issue for questions or concerns.
- Check out the [discussions](https://github.com/your-username/scanrunner-cli/discussions) for ongoing conversations.

Thank you for contributing to SCANRUNNER-CLI! 🚀

