package index

// PublicListingURL is a URL to a public listing of MAC address
// prefixes and their associated manufacturer.
type PublicListingURL string

const (
	// Official IEEE OUI public listing URLs. For some reason,
	// this is only available over HTTP.
	OUIURL   PublicListingURL = "http://standards-oui.ieee.org/oui/oui.csv"
	CIDURL   PublicListingURL = "http://standards-oui.ieee.org/cid/cid.csv"
	IABURL   PublicListingURL = "http://standards-oui.ieee.org/iab/iab.csv"
	MAMURL   PublicListingURL = "http://standards-oui.ieee.org/oui28/mam.csv"
	OUI36URL PublicListingURL = "http://standards-oui.ieee.org/oui36/oui36.csv"

	// Special public listing URL for the picatz/manuf project to serve
	// the contents of the manuf.csv file over HTTPS.
	RawGitHubURL PublicListingURL = "https://raw.githubusercontent.com/picatz/manuf/main/manuf.csv"
)

// AllPublicListingURLs is a slice of all public listing URLs, excluding
// the RawGitHubURL which is compiled from those other URLs.
var AllPublicListingURLs = []PublicListingURL{
	OUIURL,
	CIDURL,
	IABURL,
	MAMURL,
	OUI36URL,
}
