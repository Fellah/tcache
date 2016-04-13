package sletat

import (
	"bytes"
	"encoding/xml"
	"log"
	"net/http"
)

var request = Request{
	Header: RequestHeader{
		AuthInfo: AuthInfo{
			Login:    login,
			Password: password,
		},
	},
}

// SOAPAction
type GetPacketList struct {
	XMLName         xml.Name `xml:"urn:SletatRu:Contracts:Bulk:Soap11Gate:v1 GetPacketList"`
	CreateDatePoint string   `xml:"createDatePoint"`
}

type PacketInfo struct {
	CountryId    int    `xml:"CountryId"`
	CreateDate   string `xml:"CreateDate"`
	DateTimeFrom string `xml:"DateTimeFrom"`
	DateTimeTo   string `xml:"DateTimeTo"`
	DptCityId    int    `xml:"DptCityId"`
	Id           string `xml:"Id"`
	SourceId     int    `xml:"SourceId"`
}

func FetchPacketsList(date string) (chan PacketInfo, error) {
	var buf bytes.Buffer
	var packet PacketInfo

	request.Body.SOAPAction = GetPacketList{
		CreateDatePoint: date,
	}

	enc := xml.NewEncoder(&buf)
	if err := enc.Encode(request); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, URL, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "text/xml;charset=UTF-8")
	req.Header.Add("SOAPAction", "urn:SletatRu:Contracts:Bulk:Soap11Gate:v1/Soap11Gate/GetPacketList")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}


	packets := make(chan PacketInfo)
	go func() {
		defer resp.Body.Close()
		defer close(packets)

		decoder := xml.NewDecoder(resp.Body)
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
				if se.Name.Local == "PacketInfo" {
					decoder.DecodeElement(&packet, &se)
					packets <- packet
				}
			}
		}
	}()

	return packets, nil
}
