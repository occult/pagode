package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPatchFeatureMarkersScansCmdButSkipsCmdSetup(t *testing.T) {
	root := t.TempDir()

	webFile := filepath.Join(root, "cmd", "web", "main.go")
	setupFile := filepath.Join(root, "cmd", "setup", "main.go")

	mustWriteMarkerFile(t, webFile, `package main

import "fmt"

func main() {
	// [feature:tasks] start
	fmt.Println("tasks")
	// [feature:tasks] end
	fmt.Println("web")
}
`)

	mustWriteMarkerFile(t, setupFile, `package main

func main() {
	// [feature:tasks] start
	println("setup")
	// [feature:tasks] end
}
`)

	remover := NewRemover(root)
	patched, err := remover.PatchFeatureMarkers([]string{"tasks"})
	if err != nil {
		t.Fatalf("PatchFeatureMarkers() error = %v", err)
	}
	if patched != 1 {
		t.Fatalf("PatchFeatureMarkers() patched %d files, want 1", patched)
	}

	webData, err := os.ReadFile(webFile)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", webFile, err)
	}
	if strings.Contains(string(webData), `[feature:tasks]`) {
		t.Fatalf("expected cmd/web file markers to be removed, got:\n%s", webData)
	}

	setupData, err := os.ReadFile(setupFile)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", setupFile, err)
	}
	if !strings.Contains(string(setupData), `[feature:tasks]`) {
		t.Fatalf("expected cmd/setup file markers to remain, got:\n%s", setupData)
	}
}

func mustWriteMarkerFile(t *testing.T, path, contents string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", path, err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
}
