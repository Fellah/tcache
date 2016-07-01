package sletat

import (
	"compress/gzip"
	"encoding/xml"

	"github.com/fellah/tcache/log"
)

const bulKCacheUrl = "http://bulk.sletat.ru/BulkCacheDownload?packetId="

type Tour struct {
	SourceId   int    `xml:"sourceId,attr"`
	UpdateDate string `xml:"updateDate,attr"`
	Price      int    `xml:"price,attr"`
	CurrencyId int    `xml:"currencyId,attr"`
	Checkin    string `xml:"checkin,attr"`
	Nights     int    `xml:"nights,attr"`
	Adults     int    `xml:"adults,attr"`
	Kids       int    `xml:"kids,attr"`
	Kid1Age    *int   `xml:"kid1age,attr"`
	Kid2Age    *int   `xml:"kid2age,attr"`
	Kid3Age    *int   `xml:"kid3age,attr"`
	HotelId    int    `xml:"hotelId,attr"`
	TownId     int    `xml:"townId,attr"`
	MealId     int    `xml:"mealId,attr"`
	Hash       string `xml:"hash,attr"`

	CreateDate string

	DptCityId int
	CountryId int

	PriceByr int
	PriceEur int
	PriceUsd int
}

func FetchTours(packetId string) (chan Tour, error) {
	url := bulKCacheUrl + packetId
	log.Info.Println("Download:", url)

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}

	tours := make(chan Tour)
	go func() {
		defer resp.Body.Close()
		defer gzipReader.Close()
		defer close(tours)

		decoder := xml.NewDecoder(gzipReader)
		for {
			t, err := decoder.Token()
			if err != nil && err.Error() != "EOF" {
				log.Error.Println(err)
			}

			if t == nil {
				break
			}

			switch se := t.(type) {
			case xml.StartElement:
				if se.Name.Local == "tour" {
					tour := Tour{}
					decoder.DecodeElement(&tour, &se)
					tours <- tour
				}
			}
		}
	}()

	return tours, nil
}
