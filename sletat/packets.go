package sletat

import (
	"bytes"
	"encoding/xml"
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

//func FetchPacketsList(date string) (chan PacketInfo, error) {
func FetchPacketsList(date string) ([]PacketInfo, error) {
	var buf bytes.Buffer
	//var packet PacketInfo

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
	defer resp.Body.Close()

	envelope := struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			XMLName               xml.Name `xml:"Body"`
			GetPacketListResponse struct {
				XMLName             xml.Name `xml:"GetPacketListResponse"`
				GetPacketListResult struct {
					XMLName    xml.Name `xml:"GetPacketListResult"`
					PacketInfo []PacketInfo
				}
			}
		}
	}{}
	if err = xml.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, err
	}

	return envelope.Body.GetPacketListResponse.GetPacketListResult.PacketInfo, nil
}
