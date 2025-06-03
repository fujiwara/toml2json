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