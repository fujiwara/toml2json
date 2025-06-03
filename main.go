package toml2json

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

func Run(ctx context.Context) error {
	return RunWithArgs(ctx, os.Args[1:], os.Stdin, os.Stdout)
}

func RunWithArgs(ctx context.Context, args []string, stdin io.Reader, stdout io.Writer) error {
	var input io.Reader = stdin
	
	if len(args) > 0 {
		file, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", args[0], err)
		}
		defer file.Close()
		input = file
	}

	var data interface{}
	if _, err := toml.NewDecoder(input).Decode(&data); err != nil {
		return fmt.Errorf("failed to parse TOML: %w", err)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	fmt.Fprintln(stdout, string(jsonData))
	return nil
}
