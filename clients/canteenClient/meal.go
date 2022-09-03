package canteenClient

import "encoding/xml"

type Meal struct {
	XMLName     xml.Name     `xml:"essen"`
	Id          string       `xml:"id,attr"`
	Category    string       `xml:"kategorie,attr"`
	Rating      string       `xml:"bewertung,attr"`
	Img         string       `xml:"img,attr"`
	ImgSmall    string       `xml:"img_small,attr"`
	ImgBig      string       `xml:"img_big,attr"`
	Pig         string       `xml:"schwein,attr"`
	Beef        string       `xml:"rind,attr"`
	Vegetarian  string       `xml:"vegetarisch,attr"`
	Alcohol     string       `xml:"alkohol,attr"`
	Description Descritption `xml:"deutsch"`
	Price       []Price      `xml:"pr"`
}
