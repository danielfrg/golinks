package main

import (
	"fmt"
	"github.com/danielfrg/golinks/internal/mcp"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Settings struct {
	Addr, Database, PublicURL string
	ReadOnly                  bool
	MCPAllowedHosts           []string
}

func env(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
func loadSettings() (Settings, error) {
	settings := Settings{Addr: env("GOLINKS_ADDR", "127.0.0.1:8080"), Database: env("GOLINKS_DATABASE", "./data/golinks.db"), PublicURL: env("GOLINKS_PUBLIC_URL", "http://127.0.0.1:8080")}
	if _, _, err := net.SplitHostPort(settings.Addr); err != nil {
		return settings, fmt.Errorf("GOLINKS_ADDR: %w", err)
	}
	u, err := url.Parse(settings.PublicURL)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil || u.RawQuery != "" || u.Fragment != "" || (u.Path != "" && u.Path != "/") {
		return settings, fmt.Errorf("GOLINKS_PUBLIC_URL must be an HTTP(S) origin without a path")
	}
	settings.PublicURL = strings.TrimSuffix(settings.PublicURL, "/")
	settings.ReadOnly, err = strconv.ParseBool(env("GOLINKS_READONLY", "false"))
	if err != nil {
		return settings, fmt.Errorf("GOLINKS_READONLY: %w", err)
	}
	settings.MCPAllowedHosts, err = mcp.ParseAllowedHosts(os.Getenv("GOLINKS_MCP_ALLOWED_HOSTS"))
	if err != nil {
		return settings, fmt.Errorf("GOLINKS_MCP_ALLOWED_HOSTS: %w", err)
	}
	return settings, nil
}
