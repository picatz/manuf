package index

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/sync/errgroup"
)

// HTTPGetAllRecords fetches all records all public listing URLs concurrently.
func HTTPGetAllRecords(ctx context.Context, client *http.Client) (Records, error) {
	var allRecords sync.Map

	eg, gtx := errgroup.WithContext(ctx)

	// Fetch all records concurrently.
	for i := range AllPublicListingURLs {
		url := AllPublicListingURLs[i]
		eg.Go(func() error {
			records, err := HTTPGetRecords(gtx, client, string(url))
			if err != nil {
				return fmt.Errorf("failed to get records for %q: %w", url, err)
			}

			allRecords.Store(url, records)
			return nil
		})
	}

	// Wait for all records to be fetched.
	err := eg.Wait()
	if err != nil {
		return nil, err
	}

	// Merge all records into a single slice.
	records := Records{}

	allRecords.Range(func(key, value interface{}) bool {
		records = append(records, value.(Records)...)
		return true
	})

	// Return all records.
	return records, nil
}

func HTTPGetRecords(ctx context.Context, client *http.Client, recordsURL string) (Records, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, recordsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request to %q: %w", recordsURL, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got non-200 HTTP response %d from %q: %v", resp.StatusCode, recordsURL, http.StatusText(resp.StatusCode))
	}

	defer resp.Body.Close()

	records, err := RecordsFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV records from HTTP response body from %q: %w", recordsURL, err)
	}

	return records, nil
}
