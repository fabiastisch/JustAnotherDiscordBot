package canteenClient

import (
	"encoding/xml"
	"log"
	"net/http"
	"strconv"
	_ "strconv"
	"strings"
)

var data = `
<speiseplan xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="speiseplan.xsd">
<datum tag="2" monat="9" jahr="2022"/>
<essen id="8363" kategorie="Ofen" bewertung="0" img="true" img_small="https://www.swcz.de/bilderspeiseplan/bilder_190/8363.png" img_big="https://www.swcz.de/bilderspeiseplan/bilder_350/8363.png" schwein="false" rind="false" vegetarisch="true" alkohol="false">
<deutsch>
Bei uns wird noch selber zerrupft!: Hausgemachter Kaiserschmarrn mit Rosinen und Apfelmus (3,15,19,81)
</deutsch>
<pr gruppe="S">2.40</pr>
<pr gruppe="M">4.30</pr>
<pr gruppe="G">5.40</pr>
</essen>
<essen id="11774" kategorie="Schneller Teller 1" bewertung="0" img="false" schwein="true" rind="false" vegetarisch="false" alkohol="false">
<deutsch>
Hirtenrolle mit Frischkäsefüllung an hausgemachtem Zaziki dazu Pommes frites (15,18,19,49,51,81) und kleiner Weißkraut-Paprika-Salat
</deutsch>
<pr gruppe="S">3.00</pr>
<pr gruppe="M">5.00</pr>
<pr gruppe="G">6.10</pr>
</essen>
<essen id="11618" kategorie="Schneller Teller 2" bewertung="0" img="true" img_small="https://www.swcz.de/bilderspeiseplan/bilder_190/11618.png" img_big="https://www.swcz.de/bilderspeiseplan/bilder_350/11618.png" schwein="false" rind="false" vegetarisch="false" alkohol="false">
<deutsch>
Statt Waffel mal Falafel!: 6 Kichererbsenbällchen "Falafel" (21,81) mit Avocado-Hummus-Dip (3,23,49) dazu Pommes frites und kleiner Weißkraut-Paprikasalat
</deutsch>
<pr gruppe="S">2.90</pr>
<pr gruppe="M">4.90</pr>
<pr gruppe="G">6.00</pr>
</essen>
<essen id="11743" kategorie="Wok" bewertung="0" img="true" img_small="https://www.swcz.de/bilderspeiseplan/bilder_190/11743.png" img_big="https://www.swcz.de/bilderspeiseplan/bilder_350/11743.png" schwein="false" rind="false" vegetarisch="true" alkohol="false">
<deutsch>
Gebratene Gnocchi in Gemüse-Gorgonzola-Sahne-Soße mit Schnittlauch (15,19)
</deutsch>
<pr gruppe="S">2.50</pr>
<pr gruppe="M">4.40</pr>
<pr gruppe="G">5.50</pr>
</essen>
<essen id="11382" kategorie="xXx cafete⁵⁵ xXx" bewertung="0" img="true" img_small="https://www.swcz.de/bilderspeiseplan/bilder_190/11382.png" img_big="https://www.swcz.de/bilderspeiseplan/bilder_350/11382.png" schwein="false" rind="false" vegetarisch="true" alkohol="false">
<deutsch>
"Kalter Kaffee macht schön!" Hausgemachter Eiskaffee mit 1er Kugel Vanilleeis und frischer Schlagsahne (19), (auch vegan)
</deutsch>
<pr gruppe="S">2.20</pr>
<pr gruppe="M">2.40</pr>
<pr gruppe="G">2.40</pr>
</essen>
<essen id="11404" kategorie="xXx cafete⁵⁵ xXx" bewertung="0" img="true" img_small="https://www.swcz.de/bilderspeiseplan/bilder_190/11404.png" img_big="https://www.swcz.de/bilderspeiseplan/bilder_350/11404.png" schwein="false" rind="false" vegetarisch="true" alkohol="false">
<deutsch>
Ist Dir zu heiß? In unserer Cafeteria gibt's lecker lokal und natürlich produziertes Kugeleis! Probiert jetzt unsere neuen Sorten: Gurke, Joghurt-Sanddorn, Salted Caramel und Mango!
</deutsch>
<pr gruppe="S">1.50</pr>
<pr gruppe="M">1.60</pr>
<pr gruppe="G">1.60</pr>
</essen>
<essen id="11406" kategorie="xXx cafete⁵⁵ xXx" bewertung="0" img="true" img_small="https://www.swcz.de/bilderspeiseplan/bilder_190/11406.png" img_big="https://www.swcz.de/bilderspeiseplan/bilder_350/11406.png" schwein="false" rind="false" vegetarisch="true" alkohol="false">
<deutsch>
Du hast Bock auf ein Stück Kuchen danach? Dann hols Dir aus der Torten- und Kuchenvitrine in der cafete⁵⁵! Hier findet jeder das passende Stück.
</deutsch>
<pr gruppe="S">0.00</pr>
<pr gruppe="M">0.00</pr>
<pr gruppe="G">0.00</pr>
</essen>
</speiseplan>
`

func Request() string {

	resp, err := http.Get("https://www.swcz.de/bilderspeiseplan/xml.php?plan=1479835489")

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == http.StatusOK {

		//		var bodyString string
		//		bodyBytes, err := io.ReadAll(resp.Body)
		//		if err != nil {
		//			log.Fatal(err)
		//		}

		//		bodyString = string(bodyBytes)

		/*
			menu := new(Menu)
			err = xml.Unmarshal(bodyBytes, menu)

		*/

		menu := new(Menu)
		err = xml.Unmarshal([]byte(data), menu)

		if err != nil {
			log.Fatal(err)
		}

		if len(menu.Meal) == 0 {
			return "Heute ist kein Speiseplan verfügbar"
		}

		return buildResponse(menu)

	}

	return resp.Status
}

func buildResponse(menu *Menu) string {
	var stringBuilder strings.Builder
	for i, meal := range menu.Meal {
		strconv.Itoa(i)
		stringBuilder.WriteString("\n" + meal.Description.Value)
		stringBuilder.WriteString("\n" + meal.ImgSmall)
		stringBuilder.WriteString("\n")
		vegeterian, _ := strconv.ParseBool(meal.Vegetarian)
		beef, _ := strconv.ParseBool(meal.Beef)
		alcohol, _ := strconv.ParseBool(meal.Alcohol)
		pig, _ := strconv.ParseBool(meal.Pig)

		stringBuilder.WriteString("\n\nHinweis:")
		if vegeterian {
			stringBuilder.WriteString("\n\t\t Vegetarisch")
		}
		if beef {
			stringBuilder.WriteString("\n\t\t Rind")
		}
		if alcohol {
			stringBuilder.WriteString("\n\t\t Alkohol")
		}
		if pig {
			stringBuilder.WriteString("\n\t\t Schwein")
		}

		stringBuilder.WriteString("\n\nPreise:")
		for _, price := range meal.Price {
			if price.Group == "S" {
				stringBuilder.WriteString("\n\t\t Student: ")
			}
			if price.Group == "M" {
				stringBuilder.WriteString("\n\t\t Mitarbeiter: ")
			}
			if price.Group == "G" {
				stringBuilder.WriteString("\n\t\t Gast: ")
			}
			stringBuilder.WriteString(price.GroupPrice)
		}
		stringBuilder.WriteString("\n\n\n\n")
	}
	return stringBuilder.String()
}
