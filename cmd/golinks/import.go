package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/danielfrg/golinks/internal/store"
)

func runImport(ctx context.Context, args []string, databasePath string, stdout io.Writer) error {
	flags := flag.NewFlagSet("import", flag.ContinueOnError)
	flags.SetOutput(stdout)
	linksFile := flags.String("links", "", "JSONL link snapshot")
	statsFile := flags.String("stats", "", "optional CSV click snapshot")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *linksFile == "" || flags.NArg() != 0 {
		return errors.New("usage: golinks import --links <snapshot.jsonl> [--stats <stats.csv>]")
	}
	source, err := os.Open(*linksFile)
	if err != nil {
		return err
	}
	defer source.Close()
	var stats io.Reader
	if *statsFile != "" {
		file, err := os.Open(*statsFile)
		if err != nil {
			return err
		}
		defer file.Close()
		stats = file
	}
	database, err := store.Open(databasePath)
	if err != nil {
		return err
	}
	defer database.Close()
	result, err := database.Import(ctx, source, stats)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(stdout, "Imported %d links; skipped %d existing links.\n", result.Imported, result.Skipped)
	if err == nil && result.SkippedStats > 0 {
		_, err = fmt.Fprintf(stdout, "Skipped %d stats rows for links absent from the snapshot.\n", result.SkippedStats)
	}
	return err
}
