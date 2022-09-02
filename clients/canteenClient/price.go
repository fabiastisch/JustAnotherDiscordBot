package canteenClient

import "encoding/xml"

type Price struct {
	XMLName    xml.Name `xml:"pr"`
	Group      string   `xml:"name,attr"`
	GroupPrice string   `xml:",chardata"`
}
