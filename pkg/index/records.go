package index

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

// Assignment is either a MAC addresses (most common?), Subnetwork Access Protocol identifiers,
// World Wide Names for Fibre Channel devices or vendor blocks in EDID.
type Assignment string

type (
	MACAddress                         Assignment // https://en.wikipedia.org/wiki/MAC_address
	SubnetworkAccessProtocolIdentifier Assignment // https://en.wikipedia.org/wiki/Subnetwork_Access_Protocol
	WorldWideName                      Assignment // https://en.wikipedia.org/wiki/World_Wide_Name
	FibreChannelDevice                 Assignment // https://en.wikipedia.org/wiki/Fibre_Channel
	EDIDVendorBlock                    Assignment // https://en.wikipedia.org/wiki/Extended_Display_Identification_Data
)

// Registry is an IEEE public OUI listing.
type Registry string

const (
	RegistryMA_L Registry = "MA-L" // http://standards-oui.ieee.org/oui/oui.csv
	RegistryMA_M Registry = "MA-M" // http://standards-oui.ieee.org/oui28/mam.csv
	RegistryMA_S Registry = "MA-S" // http://standards-oui.ieee.org/oui36/oui36.csv
	RegistryIAB  Registry = "IAB"  // http://standards-oui.ieee.org/iab/iab.csv
	RegistryCID  Registry = "CID"  // http://standards-oui.ieee.org/cid/cid.csv
)

// Record contains information for a organizationally unique identifier (OUI),
// a 24-bit number that uniquely identifies a vendor, manufacturer, or other
// organization.
//
// Only assignment from "MA-L" registry assigns new OUI.
type Record struct {
	Registry            Registry
	Assignment          Assignment
	OrganizationName    string
	OrganizationAddress string
}

// CSV returns a slice of strings representing the record in CSV format.
func (r *Record) CSV() []string {
	return []string{
		string(r.Registry),
		string(r.Assignment),
		r.OrganizationName,
		r.OrganizationAddress,
	}
}

// func (a Assignment) MACPrefix() string {
// 	parts := strings.Split(string(a), "")
//
// 	return strings.Join(parts[0:2], "") + ":" + strings.Join(parts[2:4], "") + ":" + strings.Join(parts[4:], "")
// }

// Records is a slice of individual record entries.
type Records []*Record

// Write writes the records to a CSV file.
func (r Records) Write(filepath string) error {
	// Open file for writing.
	fh, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640)
	if err != nil {
		return fmt.Errorf("failed to open file for writing %q: %w", filepath, err)
	}
	defer fh.Close()

	// Write records to file.
	w := csv.NewWriter(fh)
	err = w.Write(strings.Split("Registry,Assignment,Organization Name,Organization Address", ","))
	if err != nil {
		return fmt.Errorf("failed to write header to file %q: %w", filepath, err)
	}
	for _, record := range r {
		err = w.Write(record.CSV())
		if err != nil {
			return fmt.Errorf("failed to write record to file %q: %w", filepath, err)
		}
	}
	w.Flush()
	err = w.Error()
	if err != nil {
		return fmt.Errorf("failed to flush records to file %q: %w", filepath, err)
	}
	return nil
}

// RecordsFromReader reads records from a CSV file reader.
func RecordsFromReader(rdr io.Reader) (Records, error) {
	records := Records{}

	r := csv.NewReader(rdr)

	// Read records, skipping header.
	for {
		parts, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to read CSV records: %w", err)
		}
		if parts[0] == "Registry" {
			continue // skip header
		}
		record := &Record{
			Registry:            Registry(parts[0]),
			Assignment:          Assignment(parts[1]),
			OrganizationName:    parts[2],
			OrganizationAddress: strings.TrimSpace(parts[3]),
		}
		records = append(records, record)
	}

	return records, nil
}

// RecordsFromFile reads records from a CSV file path.
func RecordsFromFile(filepath string) (Records, error) {
	fh, err := os.OpenFile(filepath, os.O_RDONLY, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading %q: %w", filepath, err)
	}
	defer fh.Close()
	return RecordsFromReader(fh)
}
