package canteenClient

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
)

func Request() string {
	resp, err := http.Get("https://www.swcz.de/bilderspeiseplan/xml.php?plan=1479835489")

	var bodyString string

	if err != nil {
		log.Fatal(err)
	}
	//defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		bodyString = string(bodyBytes)
		fmt.Print(bodyString)

		menu := new(Menu)
		err = xml.Unmarshal(bodyBytes, menu)

		if err != nil {
			log.Fatal(err)
		}

		return menu.Meal[0].Description.Value

	}

	return resp.Status
}
