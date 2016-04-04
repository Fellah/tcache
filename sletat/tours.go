package sletat

import (
	"compress/gzip"
	"encoding/xml"
	"log"
)

const (
	BULK_CACHE_URL = "http://bulk.sletat.ru/BulkCacheDownload?packetId="
)

type Tour struct {
	OfferId               int    `xml:"offerId,attr"`
	RequestId             int    `xml:"requestId,attr"`
	SourceId              int    `xml:"sourceId,attr"`
	UpdateDate            string `xml:"updateDate,attr"`
	Price                 int    `xml:"price,attr"`
	CurrencyId            int    `xml:"currencyId,attr"`
	Checkin               string `xml:"checkin,attr"`
	Nights                int    `xml:"nights,attr"`
	Adults                int    `xml:"adults,attr"`
	Kids                  int    `xml:"kids,attr"`
	HotelId               int    `xml:"hotelId,attr"`
	TownId                int    `xml:"townId,attr"`
	StarId                int    `xml:"starId,attr"`
	MealId                int    `xml:"mealId,attr"`
	RoomId                int    `xml:"roomId,attr"`
	RoomName              string `xml:"roomName,attr"`
	HtplaceId             int    `xml:"htplaceId,attr"`
	HotelIsInStop         int    `xml:"hotelIsInStop,attr"`
	TicketsIncluded       int    `xml:"ticketsIncluded,attr"`
	HasEconomTicketsDpt   int    `xml:"hasEconomTicketsDpt,attr"`
	HasEconomTicketsRtn   int    `xml:"hasEconomTicketsRtn,attr"`
	HasBusinessTicketsDpt int    `xml:"hasBusinessTicketsDpt,attr"`
	HasBusinessTicketsRtn int    `xml:"hasBusinessTicketsRtn,attr"`
	TourName              string `xml:"tourName,attr"`
	OriginalPrice         int    `xml:"originalPrice,attr"`
	TourUrl               string `xml:"tourUrl,attr"`
	PriceType             int    `xml:"priceType,attr"`
	Flags                 int    `xml:"flags,attr"`
	Hash                  string `xml:"hash,attr"`

	CreateDate string

	DptCityId int

	PriceByr int
	PriceEur int
	PriceUsd int
}

func FetchTours(packetId string, chRawTour chan<- Tour) error {
	var tour Tour

	url := BULK_CACHE_URL + packetId
	log.Println("Download:", url)

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	decoder := xml.NewDecoder(gzipReader)
	for {
		t, err := decoder.Token()
		if err != nil && err.Error() != "EOF" {
			log.Println(err)
		}

		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "tour" {
				decoder.DecodeElement(&tour, &se)
				chRawTour <- tour
			}
		}
	}

	return nil
}
