package index

type PublicListingURL string

const (
	OUIURL   PublicListingURL = "http://standards-oui.ieee.org/oui/oui.csv"
	CIDURL   PublicListingURL = "http://standards-oui.ieee.org/cid/cid.csv"
	IABURL   PublicListingURL = "http://standards-oui.ieee.org/iab/iab.csv"
	MAMURL   PublicListingURL = "http://standards-oui.ieee.org/oui28/mam.csv"
	OUI36URL PublicListingURL = "http://standards-oui.ieee.org/oui36/oui36.csv"
)
