# Claude Code Instructions

This project uses Claude Code for development assistance.

## Commands

- `make` - Build the project
- `go test -v` - Run tests
- `go mod tidy` - Tidy dependencies

## Project Structure

- `main.go` - Core TOML to JSON conversion logic
- `main_test.go` - Unit tests with comprehensive test cases
- `cmd/toml2json/` - CLI entry point with signal handling
- `testdata/` - Test TOML and JSON files for validation

## Testing

The project includes comprehensive tests covering:
- Standard input processing
- File input processing
- Error handling for invalid TOML
- Context cancellation for signal handling
- Edge cases (empty files, non-existent files)

## Development Notes

- All text files should end with newline characters
- Follow Go conventions for code style
- Signal handling is implemented for graceful interruption during stdin reading
- Tests use testdata files for validation against known TOML/JSON pairs