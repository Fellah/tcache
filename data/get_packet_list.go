package data

import "encoding/xml"

type GetPacketList struct {
	XMLName         xml.Name `xml:"urn:SletatRu:Contracts:Bulk:Soap11Gate:v1 GetPacketList"`
	CreateDatePoint string   `xml:"createDatePoint"`
}
