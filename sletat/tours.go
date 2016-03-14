package sletat

import (
	"compress/gzip"
	"encoding/xml"
	//"log"
)

const (
	BULK_CACHE_URL = "http://bulk.sletat.ru/BulkCacheDownload?packetId="
)

type Tours struct {
	XMLName xml.Name `xml:"tours"`
	Tour    []Tour   `xml:"tour"`
}

type Tour struct {
	XMLName xml.Name `xml:"tour"`

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
	HtplaceId             int    `xml:"htplaceId,attr"`
	HotelIsInStop         int    `xml:"hotelIsInStop,attr"`
	TicketsIncluded       int    `xml:"ticketsIncluded,attr"`
	HasEconomTicketsDpt   int    `xml:"hasEconomTicketsDpt,attr"`
	HasEconomTicketsRtn   int    `xml:"hasEconomTicketsRtn,attr"`
	HasBusinessTicketsDpt int    `xml:"hasBusinessTicketsDpt,attr"`
	HasBusinessTicketsRtn int    `xml:"hasBusinessTicketsRtn,attr"`
	TourName              string `xml:"tourName,attr"`
	HotelName             string `xml:"hotelName,attr"`
	TownName              string `xml:"townName,attr"`
	StarName              string `xml:"starName,attr"`
	MealName              string `xml:"mealName,attr"`
	RoomName              string `xml:"roomName,attr"`
	HtplaceName           string `xml:"htplaceName,attr"`
	OriginalHotelName     string `xml:"originalHotelName,attr"`
	OriginalTownName      string `xml:"originalTownName,attr"`
	OriginalStarName      string `xml:"originalStarName,attr"`
	OriginalMealName      string `xml:"originalMealName,attr"`
	OriginalRoomName      string `xml:"originalRoomName,attr"`
	OriginalHtplaceName   string `xml:"originalHtplaceName,attr"`
	OriginalCountryName   string `xml:"originalCountryName,attr"`
	OriginalDptCityName   string `xml:"originalDptCityName,attr"`
	OriginalCurrencyId    int    `xml:"originalCurrencyId,attr"`
	OriginalCurrencyName  string `xml:"originalCurrencyName,attr"`
	OriginalPrice         int    `xml:"originalPrice,attr"`
	TourUrl               string `xml:"tourUrl,attr"`
	PriceType             int    `xml:"priceType,attr"`
	Flags                 int    `xml:"flags,attr"`
	Hash                  string `xml:"hash,attr"`
	CreateDate            string
}

func FetchTours(packetId string) ([]Tour, error) {
	url := BULK_CACHE_URL + packetId
	// log.Println("Download:", url)

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	gzipReader, err := gzip.NewReader(resp.Body)
	defer gzipReader.Close()

	tours := Tours{}
	if err = xml.NewDecoder(gzipReader).Decode(&tours); err != nil {
		return nil, err
	}

	return tours.Tour, nil
}
