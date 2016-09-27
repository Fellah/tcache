package sletat

import (
	"bytes"
	"encoding/xml"
	"github.com/fellah/tcache/data"
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

func FetchPacketsList(date string) ([]data.PacketInfo, error) {
	var buf bytes.Buffer

	request.Body.SOAPAction = data.GetPacketList{
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
					PacketInfo []data.PacketInfo
				}
			}
		}
	}{}
	if err = xml.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, err
	}

	return envelope.Body.GetPacketListResponse.GetPacketListResult.PacketInfo, nil
}
