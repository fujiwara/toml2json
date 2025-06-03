# toml2json

A simple command-line tool to convert TOML files to JSON format.

## Features

- Convert TOML to JSON via stdin or file input
- Cross-platform support

## Installation

```bash
go install github.com/fujiwara/toml2json/cmd/toml2json@latest
```

Or build from source:

```bash
git clone https://github.com/fujiwara/toml2json.git
cd toml2json
make
```

## Usage

### From stdin

```bash
echo 'name = "example"' | toml2json
# Output: {"name":"example"}
```

### From file

```bash
toml2json config.toml
# Converts config.toml to JSON and outputs to stdout
```

### Examples

Convert a TOML configuration file:

```bash
# Input: config.toml
[server]
name = "example"
port = 8080

[database]
host = "localhost"
port = 5432

# Command
toml2json config.toml

# Output
{"database":{"host":"localhost","port":5432},"server":{"name":"example","port":8080}}
```

Use with pipes:

```bash
cat config.toml | toml2json | jq .
```

## Error Handling

- Invalid TOML syntax will result in a parsing error
- Non-existent files will result in a file open error
- Signal interruption (Ctrl-C) is handled gracefully

## Development

### Running Tests

```bash
go test -v
```

### Building

```bash
make
```

## LICENSE

MIT

## Author

Fujiwara Shunichiro
