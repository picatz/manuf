package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/picatz/manuf/pkg/index"
)

func handleError(msg string, err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("%s: %v\n", msg, err))
		os.Exit(1)
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	records, err := index.HTTPGetAllRecords(ctx, cleanhttp.DefaultClient())
	handleError("failed to get all records", err)

	sort.Slice(records, func(i, j int) bool {
		return records[i].Registry < records[j].Registry
	})

	err = records.Write("manuf.csv")
	handleError("failed to write records to manuf.csv", err)
}
