package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCleanStaleEntFilesRemovesGeneratedAdminFiles(t *testing.T) {
	root := t.TempDir()

	mustWriteFile(t, filepath.Join(root, "ent", "generate.go"), "package ent\n")
	mustWriteFile(t, filepath.Join(root, "ent", "entc.go"), "package ent\n")
	mustWriteFile(t, filepath.Join(root, "ent", "client.go"), "package ent\n")
	mustWriteFile(t, filepath.Join(root, "ent", "ent.go"), "package ent\n")
	mustWriteFile(t, filepath.Join(root, "ent", "mutation.go"), "package ent\n")
	mustWriteFile(t, filepath.Join(root, "ent", "runtime.go"), "package ent\n")
	mustWriteFile(t, filepath.Join(root, "ent", "tx.go"), "package ent\n")
	mustWriteFile(t, filepath.Join(root, "ent", "schema", "user.go"), "package schema\n")
	mustWriteFile(t, filepath.Join(root, "ent", "schema", "passwordtoken.go"), "package schema\n")
	mustWriteFile(t, filepath.Join(root, "ent", "admin", "extension.go"), "package admin\n")
	mustWriteFile(t, filepath.Join(root, "ent", "admin", "handler.go"), "package admin\n")
	mustWriteFile(t, filepath.Join(root, "ent", "admin", "types.go"), "package admin\n")
	mustWriteFile(t, filepath.Join(root, "ent", "admin", "templates", "handler.tmpl"), "{{ define \"x\" }}{{ end }}\n")
	mustWriteFile(t, filepath.Join(root, "ent", "hook", "hook.go"), "package hook\n")
	mustWriteFile(t, filepath.Join(root, "ent", "runtime", "runtime.go"), "package runtime\n")
	mustWriteFile(t, filepath.Join(root, "ent", "migrate", "migrate.go"), "package migrate\n")
	mustWriteFile(t, filepath.Join(root, "ent", "predicate", "predicate.go"), "package predicate\n")
	mustWriteFile(t, filepath.Join(root, "ent", "enttest", "enttest.go"), "package enttest\n")
	mustWriteFile(t, filepath.Join(root, "ent", "user.go"), "package ent\n")
	mustWriteFile(t, filepath.Join(root, "ent", "user_create.go"), "package ent\n")
	mustWriteFile(t, filepath.Join(root, "ent", "user_delete.go"), "package ent\n")
	mustWriteFile(t, filepath.Join(root, "ent", "user_query.go"), "package ent\n")
	mustWriteFile(t, filepath.Join(root, "ent", "user_update.go"), "package ent\n")
	mustWriteFile(t, filepath.Join(root, "ent", "user", "where.go"), "package user\n")
	mustWriteFile(t, filepath.Join(root, "ent", "passwordtoken.go"), "package ent\n")
	mustWriteFile(t, filepath.Join(root, "ent", "passwordtoken_create.go"), "package ent\n")
	mustWriteFile(t, filepath.Join(root, "ent", "passwordtoken_delete.go"), "package ent\n")
	mustWriteFile(t, filepath.Join(root, "ent", "passwordtoken_query.go"), "package ent\n")
	mustWriteFile(t, filepath.Join(root, "ent", "passwordtoken_update.go"), "package ent\n")
	mustWriteFile(t, filepath.Join(root, "ent", "passwordtoken", "where.go"), "package passwordtoken\n")
	mustWriteFile(t, filepath.Join(root, "ent", "paymentcustomer.go"), "package ent\n")
	mustWriteFile(t, filepath.Join(root, "ent", "paymentcustomer_create.go"), "package ent\n")
	mustWriteFile(t, filepath.Join(root, "ent", "paymentcustomer_delete.go"), "package ent\n")
	mustWriteFile(t, filepath.Join(root, "ent", "paymentcustomer_query.go"), "package ent\n")
	mustWriteFile(t, filepath.Join(root, "ent", "paymentcustomer_update.go"), "package ent\n")
	mustWriteFile(t, filepath.Join(root, "ent", "paymentcustomer", "where.go"), "package paymentcustomer\n")

	cleanup := NewCleanup(root)
	if err := cleanup.CleanStaleEntFiles(); err != nil {
		t.Fatalf("CleanStaleEntFiles() error = %v", err)
	}

	assertExists(t, filepath.Join(root, "ent", "generate.go"))
	assertExists(t, filepath.Join(root, "ent", "entc.go"))
	assertExists(t, filepath.Join(root, "ent", "client.go"))
	assertExists(t, filepath.Join(root, "ent", "ent.go"))
	assertExists(t, filepath.Join(root, "ent", "mutation.go"))
	assertExists(t, filepath.Join(root, "ent", "runtime.go"))
	assertExists(t, filepath.Join(root, "ent", "tx.go"))
	assertExists(t, filepath.Join(root, "ent", "schema", "user.go"))
	assertExists(t, filepath.Join(root, "ent", "schema", "passwordtoken.go"))
	assertExists(t, filepath.Join(root, "ent", "admin", "extension.go"))
	assertExists(t, filepath.Join(root, "ent", "admin", "handler.go"))
	assertExists(t, filepath.Join(root, "ent", "admin", "types.go"))
	assertExists(t, filepath.Join(root, "ent", "admin", "templates", "handler.tmpl"))
	assertExists(t, filepath.Join(root, "ent", "hook", "hook.go"))
	assertExists(t, filepath.Join(root, "ent", "runtime", "runtime.go"))
	assertExists(t, filepath.Join(root, "ent", "migrate", "migrate.go"))
	assertExists(t, filepath.Join(root, "ent", "predicate", "predicate.go"))
	assertExists(t, filepath.Join(root, "ent", "enttest", "enttest.go"))
	assertExists(t, filepath.Join(root, "ent", "user.go"))
	assertExists(t, filepath.Join(root, "ent", "user_create.go"))
	assertExists(t, filepath.Join(root, "ent", "user_delete.go"))
	assertExists(t, filepath.Join(root, "ent", "user_query.go"))
	assertExists(t, filepath.Join(root, "ent", "user_update.go"))
	assertExists(t, filepath.Join(root, "ent", "user", "where.go"))
	assertExists(t, filepath.Join(root, "ent", "passwordtoken.go"))
	assertExists(t, filepath.Join(root, "ent", "passwordtoken_create.go"))
	assertExists(t, filepath.Join(root, "ent", "passwordtoken_delete.go"))
	assertExists(t, filepath.Join(root, "ent", "passwordtoken_query.go"))
	assertExists(t, filepath.Join(root, "ent", "passwordtoken_update.go"))
	assertExists(t, filepath.Join(root, "ent", "passwordtoken", "where.go"))

	assertNotExists(t, filepath.Join(root, "ent", "paymentcustomer.go"))
	assertNotExists(t, filepath.Join(root, "ent", "paymentcustomer_create.go"))
	assertNotExists(t, filepath.Join(root, "ent", "paymentcustomer_delete.go"))
	assertNotExists(t, filepath.Join(root, "ent", "paymentcustomer_query.go"))
	assertNotExists(t, filepath.Join(root, "ent", "paymentcustomer_update.go"))
	assertNotExists(t, filepath.Join(root, "ent", "paymentcustomer", "where.go"))
}

func TestPrepareEntForRegenerationRemovesGeneratedAdminFiles(t *testing.T) {
	root := t.TempDir()

	mustWriteFile(t, filepath.Join(root, "ent", "admin", "extension.go"), "package admin\n")
	mustWriteFile(t, filepath.Join(root, "ent", "admin", "handler.go"), "package admin\n")
	mustWriteFile(t, filepath.Join(root, "ent", "admin", "types.go"), "package admin\n")
	mustWriteFile(t, filepath.Join(root, "ent", "admin", "templates", "handler.tmpl"), "{{ define \"x\" }}{{ end }}\n")

	cleanup := NewCleanup(root)
	if err := cleanup.prepareEntForRegeneration(); err != nil {
		t.Fatalf("prepareEntForRegeneration() error = %v", err)
	}

	assertExists(t, filepath.Join(root, "ent", "admin", "extension.go"))
	assertExists(t, filepath.Join(root, "ent", "admin", "templates", "handler.tmpl"))
	assertNotExists(t, filepath.Join(root, "ent", "admin", "handler.go"))
	assertNotExists(t, filepath.Join(root, "ent", "admin", "types.go"))
}

func TestVerifyFrontendBuildRunsTypecheckThenViteBuild(t *testing.T) {
	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, "node_modules", ".gitkeep"), "")

	var calls []struct {
		dir  string
		name string
		args []string
	}

	cleanup := NewCleanup(root)
	cleanup.runCommand = func(dir, name string, args ...string) error {
		calls = append(calls, struct {
			dir  string
			name string
			args []string
		}{
			dir:  dir,
			name: name,
			args: append([]string(nil), args...),
		})
		return nil
	}

	if err := cleanup.VerifyFrontendBuild(); err != nil {
		t.Fatalf("VerifyFrontendBuild() error = %v", err)
	}

	want := []struct {
		dir  string
		name string
		args []string
	}{
		{dir: root, name: "npx", args: []string{"tsc", "--noEmit"}},
		{dir: root, name: "npx", args: []string{"vite", "build"}},
	}

	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("VerifyFrontendBuild() calls = %#v, want %#v", calls, want)
	}
}

func mustWriteFile(t *testing.T, path, contents string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", path, err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
}

func assertExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected %q to exist, got error: %v", path, err)
	}
}

func assertNotExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected %q to be removed, got err=%v", path, err)
	}
}
