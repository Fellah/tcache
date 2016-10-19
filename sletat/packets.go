package sletat

import (
	"bytes"
	"encoding/xml"
	"github.com/fellah/tcache/data"
	"net/http"
	"github.com/fellah/tcache/log"
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

	log.Info.Println("FetchPacketsList...")

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

	log.Info.Println("FetchPacketsList request for packets data...")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	log.Info.Println("FetchPacketsList request for packets data done")

	log.Info.Println("FetchPacketsList packets data XML decode...")
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
	log.Info.Println("FetchPacketsList packet data:\n", resp.Status)
	if err = xml.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, err
	}
	log.Info.Println("FetchPacketsList packets data XML decode done")

	return envelope.Body.GetPacketListResponse.GetPacketListResult.PacketInfo, nil
}
