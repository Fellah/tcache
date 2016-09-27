package data

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
