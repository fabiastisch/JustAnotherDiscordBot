package canteenClient

import "encoding/xml"

type Menu struct {
	XMLName xml.Name `xml:"speiseplan"`
	Meal    []Meal   `xml:"essen"`
	Date    Date     `xml:"datum"`
}
