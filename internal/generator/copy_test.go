package generator

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCopyDirIncremental(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	// Create source files
	os.WriteFile(filepath.Join(src, "a.txt"), []byte("hello"), 0644)
	os.MkdirAll(filepath.Join(src, "sub"), 0755)
	os.WriteFile(filepath.Join(src, "sub", "b.txt"), []byte("world"), 0644)

	if err := CopyDirIncremental(src, dst); err != nil {
		t.Fatal(err)
	}

	// Verify files were copied
	data, err := os.ReadFile(filepath.Join(dst, "a.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Errorf("a.txt = %q, want %q", string(data), "hello")
	}

	data, err = os.ReadFile(filepath.Join(dst, "sub", "b.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "world" {
		t.Errorf("sub/b.txt = %q, want %q", string(data), "world")
	}
}

func TestCopyDirIncrementalSkipsUnchanged(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	srcFile := filepath.Join(src, "file.txt")
	dstFile := filepath.Join(dst, "file.txt")

	os.WriteFile(srcFile, []byte("original"), 0644)

	// First copy
	if err := CopyDirIncremental(src, dst); err != nil {
		t.Fatal(err)
	}

	// Make dst file newer
	futureTime := time.Now().Add(time.Hour)
	os.Chtimes(dstFile, futureTime, futureTime)

	// Overwrite dst with different content
	os.WriteFile(dstFile, []byte("modified"), 0644)
	os.Chtimes(dstFile, futureTime, futureTime)

	// Second copy should skip since dst is newer
	if err := CopyDirIncremental(src, dst); err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(dstFile)
	if string(data) != "modified" {
		t.Errorf("file should not have been overwritten, got %q", string(data))
	}
}

func TestCopyDirIncrementalUpdatesOlder(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	srcFile := filepath.Join(src, "file.txt")
	dstFile := filepath.Join(dst, "file.txt")

	// Create dst file first with old content
	os.WriteFile(dstFile, []byte("old"), 0644)
	pastTime := time.Now().Add(-time.Hour)
	os.Chtimes(dstFile, pastTime, pastTime)

	// Create src file (newer)
	os.WriteFile(srcFile, []byte("new"), 0644)

	if err := CopyDirIncremental(src, dst); err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(dstFile)
	if string(data) != "new" {
		t.Errorf("file should have been updated, got %q", string(data))
	}
}

func TestCopyDirIncrementalNonexistentSrc(t *testing.T) {
	dst := t.TempDir()
	err := CopyDirIncremental("/nonexistent/path", dst)
	if err == nil {
		t.Fatal("expected error for nonexistent source")
	}
}
