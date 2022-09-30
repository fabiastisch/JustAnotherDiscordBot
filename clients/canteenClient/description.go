package canteenClient

import "encoding/xml"

type Descritption struct {
	XMLName xml.Name `xml:"deutsch"`
	Value   string   `xml:",chardata"`
}
