package logger

import (
	"io"
	"os"
	"testing"
)

func TestSetQuietSuppressesOutput(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("create pipe: %v", err)
	}
	originalStdout := os.Stdout
	os.Stdout = writer
	defer func() {
		os.Stdout = originalStdout
		SetQuiet(false)
	}()

	SetQuiet(true)
	Info("hidden")
	writer.Close()
	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	reader.Close()
	if len(output) != 0 {
		t.Fatalf("quiet logger wrote %q", output)
	}
}
