package data

type PacketInfo struct {
	CountryId    int    `xml:"CountryId"`
	CreateDate   string `xml:"CreateDate"`
	DateTimeFrom string `xml:"DateTimeFrom"`
	DateTimeTo   string `xml:"DateTimeTo"`
	DptCityId    int    `xml:"DptCityId"`
	Id           string `xml:"Id"`
	SourceId     int    `xml:"SourceId"`
}
