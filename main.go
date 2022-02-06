package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/picatz/manuf/pkg/index"
)

func writeDiagnostic(msg string) {
	os.Stderr.WriteString(msg + "\n")
}

func handleError(msg string, err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%s: %v\n", msg, err))
		os.Exit(1)
	}
}

func printRecords(records index.Records) error {
	enc := json.NewEncoder(os.Stdout)

	for _, record := range records {
		err := enc.Encode(record)
		if err != nil {
			return fmt.Errorf("failed to write record as JSON to STDOUT: %w", err)
		}
	}

	return nil
}

func main() {
	dir, err := os.UserCacheDir()
	handleError("failed to get current user's cache directory: %w", err)

	path := filepath.Join(dir, "manuf.csv")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	client := cleanhttp.DefaultClient()

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			writeDiagnostic(fmt.Sprintf("no cache %q found, creating", path))
			records, err := index.HTTPGetRecords(ctx, client, string(index.RawGitHubURL))
			handleError("failed to fetch all records to populate cache", err)

			err = records.Write(path)
			handleError(fmt.Sprintf("failed to write all fetched records (%d) to populate cache", len(records)), err)

			printRecords(records)
			return
		} else {
			handleError(fmt.Sprintf("failed to check previous cache %q", path), err)
			return
		}
	}

	now := time.Now()

	if now.Sub(info.ModTime()).Hours() > (24 * 30) {
		writeDiagnostic("cached records are older than 30 days, renewing content")
		records, err := index.HTTPGetAllRecords(ctx, client)
		handleError("failed to fetch all records to populate cache", err)

		err = records.Write(path)
		handleError(fmt.Sprintf("failed to write all fetched records (%d) to populate cache", len(records)), err)

		printRecords(records)
		return
	}

	records, err := index.RecordsFromFile(path)
	handleError("failed to read records from cache", err)

	printRecords(records)
}
