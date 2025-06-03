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

	// Create a context-aware reader for stdin
	if len(args) == 0 {
		input = newContextReader(ctx, stdin)
	}

	var data interface{}
	if _, err := toml.NewDecoder(input).Decode(&data); err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("failed to parse TOML: %w", err)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	fmt.Fprintln(stdout, string(jsonData))
	return nil
}

type contextReader struct {
	ctx context.Context
	r   io.Reader
}

func newContextReader(ctx context.Context, r io.Reader) io.Reader {
	return &contextReader{ctx: ctx, r: r}
}

func (cr *contextReader) Read(p []byte) (n int, err error) {
	if cr.ctx.Err() != nil {
		return 0, cr.ctx.Err()
	}

	// Use a goroutine to read from the underlying reader
	type result struct {
		n   int
		err error
	}
	ch := make(chan result, 1)
	
	go func() {
		n, err := cr.r.Read(p)
		ch <- result{n, err}
	}()

	select {
	case <-cr.ctx.Done():
		return 0, cr.ctx.Err()
	case res := <-ch:
		return res.n, res.err
	}
}
