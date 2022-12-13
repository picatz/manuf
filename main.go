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

// writeDiagnostic writes a message to STDERR.
func writeDiagnostic(msg string) {
	os.Stderr.WriteString(msg + "\n")
}

// handleError writes a message to STDERR and exits with a non-zero exit code.
//
// This is a helper function to make error handling more concise, even
// though it is not idiomatic Go.
//
// https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully
func handleError(msg string, err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%s: %v\n", msg, err))
		os.Exit(1)
	}
}

// printRecords prints the records as JSON to STDOUT.
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
	// Get the current user's cache directory, which is where we will store the
	// manufacturer records file (manuf.csv) for future use.
	dir, err := os.UserCacheDir()
	handleError("failed to get current user's cache directory: %w", err)

	// Path to the manufacturer records file (manuf.csv).
	path := filepath.Join(dir, "manuf.csv")

	// Create a context with a 5 minute timeout, which will be used for all HTTP
	// requests made by the index package.
	//
	// This should take a matter of seconds, but we want to be sure that we don't
	// hang indefinitely while retrieving the records, so we set a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create a new HTTP client that will be used for all HTTP requests made by
	// the index package.
	client := cleanhttp.DefaultClient()

	// Check if the cache file exists, and if it does, check if it is older than
	// 30 days. If it is, we will fetch the records again and overwrite the file.
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

	// Get the current time to compare against the cache file's modification time.
	now := time.Now()

	// If the cache file is older than 30 days, we will fetch the records again
	// and overwrite the file.
	if now.Sub(info.ModTime()).Hours() > (24 * 30) {
		writeDiagnostic("cached records are older than 30 days, renewing content")
		records, err := index.HTTPGetAllRecords(ctx, client)
		handleError("failed to fetch all records to populate cache", err)

		err = records.Write(path)
		handleError(fmt.Sprintf("failed to write all fetched records (%d) to populate cache", len(records)), err)

		printRecords(records)
		return
	}

	// If we get here, we have a cache file that is less than 30 days old, so we
	// will read the records from the cache file.
	records, err := index.RecordsFromFile(path)
	handleError("failed to read records from cache", err)

	// Print the records to STDOUT as JSON.
	printRecords(records)
}
