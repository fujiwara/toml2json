package toml2json

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
)

func TestRunWithArgs_FromStdin(t *testing.T) {
	tomlInput := `[server]
name = "example"
port = 8080

[database]
host = "localhost"
port = 5432`

	stdin := strings.NewReader(tomlInput)
	var stdout bytes.Buffer

	err := RunWithArgs(context.Background(), []string{}, stdin, &stdout)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := strings.TrimSpace(stdout.String())
	expected := `{"database":{"host":"localhost","port":5432},"server":{"name":"example","port":8080}}`

	if output != expected {
		t.Errorf("expected %s, got %s", expected, output)
	}
}

func TestRunWithArgs_FromFile(t *testing.T) {
	testCases := []struct {
		name     string
		tomlFile string
		jsonFile string
	}{
		{"simple", "testdata/simple.toml", "testdata/simple.json"},
		{"complex", "testdata/complex.toml", "testdata/complex.json"},
		{"empty", "testdata/empty.toml", "testdata/empty.json"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expectedBytes, err := os.ReadFile(tc.jsonFile)
			if err != nil {
				t.Fatalf("failed to read expected JSON file: %v", err)
			}
			expected := strings.TrimSpace(string(expectedBytes))

			var stdout bytes.Buffer
			err = RunWithArgs(context.Background(), []string{tc.tomlFile}, nil, &stdout)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			output := strings.TrimSpace(stdout.String())
			if output != expected {
				t.Errorf("expected %s, got %s", expected, output)
			}
		})
	}
}

func TestRunWithArgs_InvalidTOML(t *testing.T) {
	invalidToml := `[server
name = "broken"`

	stdin := strings.NewReader(invalidToml)
	var stdout bytes.Buffer

	err := RunWithArgs(context.Background(), []string{}, stdin, &stdout)
	if err == nil {
		t.Fatal("expected error for invalid TOML")
	}

	if !strings.Contains(err.Error(), "failed to parse TOML") {
		t.Errorf("expected parse error, got: %v", err)
	}
}

func TestRunWithArgs_NonExistentFile(t *testing.T) {
	var stdout bytes.Buffer

	err := RunWithArgs(context.Background(), []string{"nonexistent.toml"}, nil, &stdout)
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}

	if !strings.Contains(err.Error(), "failed to open file") {
		t.Errorf("expected file open error, got: %v", err)
	}
}

func TestRunWithArgs_ContextCancellation(t *testing.T) {
	// Create a context that will be canceled
	ctx, cancel := context.WithCancel(context.Background())
	
	// Create a slow reader that will block
	slowReader := &slowReader{delay: 100}
	var stdout bytes.Buffer

	// Start the conversion in a goroutine
	errCh := make(chan error, 1)
	go func() {
		err := RunWithArgs(ctx, []string{}, slowReader, &stdout)
		errCh <- err
	}()

	// Cancel the context after a short delay
	cancel()

	// Wait for the function to return
	err := <-errCh
	if err == nil {
		t.Fatal("expected error due to context cancellation")
	}

	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got: %v", err)
	}
}

// slowReader simulates a slow input source
type slowReader struct {
	delay int
	pos   int
}

func (sr *slowReader) Read(p []byte) (n int, err error) {
	if sr.pos == 0 {
		// First read returns some data
		data := []byte("key = \"value\"\n")
		copy(p, data)
		sr.pos += len(data)
		return len(data), nil
	}
	// Subsequent reads would block indefinitely in real stdin
	select {}
}
