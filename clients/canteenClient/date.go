package canteenClient

import "encoding/xml"

type Date struct {
	XMLName xml.Name `xml:"datum"`
	Day     string   `xml:"tag,attr"`
	Month   string   `xml:"monat,attr"`
	Year    string   `xml:"jahr,attr"`
}
