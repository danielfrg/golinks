package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielfrg/golinks/internal/httpapi"
	"github.com/danielfrg/golinks/internal/store"
	"github.com/danielfrg/golinks/internal/web"
)

func run(ctx context.Context, args []string) error {
	settings, err := loadSettings()
	if err != nil {
		return err
	}
	if len(args) > 0 {
		if args[0] == "import" {
			return runImport(ctx, args[1:], settings.Database, os.Stdout)
		}
		return errors.New("usage: golinks [import --links <file> [--stats <file>]]")
	}
	database, err := store.Open(settings.Database)
	if err != nil {
		return err
	}
	defer database.Close()
	assets, err := web.Assets()
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp", settings.Addr)
	if err != nil {
		return err
	}
	log.Printf("Golinks listening at http://%s (database: %s)", listener.Addr(), settings.Database)
	return serve(ctx, listener, httpapi.NewHandler(assets, database, httpapi.Options{PublicURL: settings.PublicURL, ReadOnly: settings.ReadOnly, MCPAllowedHosts: settings.MCPAllowedHosts}), 5*time.Second)
}
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, os.Args[1:]); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
