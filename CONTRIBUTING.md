# Contributing to hann-cli

Thank you for your interest in contributing to `hann-cli`!

## Getting Started

1. Fork the repository on GitHub.
2. Clone your fork locally:
   ```bash
   git clone https://github.com/<your-username>/hann-cli.git
   cd hann-cli
   ```
3. Create a feature branch:
   ```bash
   git checkout -b feature/my-new-feature
   ```

## Development Guidelines

- **Go Version**: Ensure you have Go 1.22+ installed.
- **Local-First & Privacy**: Do not add telemetry, external tracking calls, or credentials reading.
- **Error Handling**: Use clean, human-friendly error messages. Avoid raw panics.
- **Formatting**: Always format your code with `go fmt ./...`.
- **Testing**: Add unit tests for any new parsers or commands.

## Running Tests

Before committing, ensure all tests pass:

```bash
make fmt
make lint
make test
```

## Submitting Pull Requests

1. Commit your changes with clear, descriptive commit messages.
2. Push your branch to GitHub.
3. Open a Pull Request against the `main` or `master` branch.
4. Describe what your PR accomplishes and verify all CI checks pass.
