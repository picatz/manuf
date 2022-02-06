package index

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestHTTPGetAllRecordsAndWriteToFile(t *testing.T) {
	allRecords, err := HTTPGetAllRecords(context.Background(), http.DefaultClient)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("got %d records from HTTP requests", len(allRecords))

	tmpDir := os.TempDir()

	path := filepath.Join(tmpDir, "records.csv")

	err = allRecords.Write(path)
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Cleanup(func() {
		err := os.Remove(path)
		if err != nil {
			t.Logf("failed to cleanup temporary file %q: %v", path, err)
		}
	})

	cachedRecords, err := RecordsFromFile(path)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("got %d records from file", len(cachedRecords))
}
