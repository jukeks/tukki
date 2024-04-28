package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureTempDirectory(t *testing.T) {
	dir, cleanup := EnsureTempDirectory()

	if dir == "" {
		t.Fatalf("expected dir to be non-empty")
	}

	if !filepath.IsAbs(dir) {
		t.Fatalf("expected dir to be absolute")
	}

	// ensure dir is in tmp (whatever the OS temp dir is)

	if filepath.Dir(dir) != os.TempDir() {
		t.Errorf("TempDir: %s, dir of tempdir: %s", os.TempDir(), filepath.Dir(os.TempDir()))
		t.Errorf("dir %s, dir of dir: %s", dir, filepath.Dir(dir))
		t.Fatalf("expected dir to be in %s, got %s", os.TempDir(), filepath.Dir(dir))
	}

	// ensure dir exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatalf("expected dir to exist, got %v", err)
	}

	// cleanup
	cleanup()

	// ensure the created tmpdir is gone
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatalf("expected dir to not exist, got %v", err)
	}

	// ensure os.TempDir() is still there
	if _, err := os.Stat(os.TempDir()); os.IsNotExist(err) {
		t.Fatalf("expected os.TempDir() to exist, got %v", err)
	}
}
