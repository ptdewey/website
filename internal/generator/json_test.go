package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteJSONFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.json")

	data := map[string]any{
		"name": "test",
		"count": 42,
	}

	if err := WriteJSONFile(data, path); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	s := string(content)
	if s == "" {
		t.Fatal("output file is empty")
	}
	// Should be indented
	if len(s) > 0 && s[0] != '{' {
		t.Errorf("expected JSON object, got %q", s[:10])
	}
}

func TestWriteJSONFileUnmarshalable(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")

	// Channels can't be marshaled to JSON
	err := WriteJSONFile(make(chan int), path)
	if err == nil {
		t.Fatal("expected error for unmarshalable data")
	}
}
