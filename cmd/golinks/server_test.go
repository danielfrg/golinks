package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSettings(t *testing.T) {
	for _, name := range []string{"GOLINKS_ADDR", "GOLINKS_DATABASE", "GOLINKS_PUBLIC_URL", "GOLINKS_READONLY"} {
		t.Setenv(name, "")
	}
	settings, err := loadSettings()
	if err != nil || settings.Addr != "127.0.0.1:8080" || settings.Database != "./data/golinks.db" {
		t.Fatal(settings, err)
	}
	t.Setenv("GOLINKS_READONLY", "maybe")
	if _, err := loadSettings(); err == nil {
		t.Fatal("accepted invalid boolean")
	}
	t.Setenv("GOLINKS_READONLY", "true")
	t.Setenv("GOLINKS_PUBLIC_URL", "javascript:alert(1)")
	if _, err := loadSettings(); err == nil {
		t.Fatal("accepted invalid public URL")
	}
}
func TestServeShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() {
		done <- serve(ctx, listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) }), time.Second)
	}()
	response, err := http.Get("http://" + listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("shutdown hung")
	}
}
func TestImportCommand(t *testing.T) {
	dir := t.TempDir()
	snapshot := filepath.Join(dir, "links.jsonl")
	if err := os.WriteFile(snapshot, []byte(`{"short":"wiki","long":"https://wiki.test/","created":"2026-10-13T00:00:00Z","lastEdit":"2026-10-13T00:00:00Z","clicks":0}`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := runImport(context.Background(), []string{"--links", snapshot}, filepath.Join(dir, "links.db"), os.Stdout); err != nil {
		t.Fatal(err)
	}
	if err := runImport(context.Background(), nil, filepath.Join(dir, "links.db"), os.Stdout); err == nil {
		t.Fatal("accepted missing snapshot")
	}
}
