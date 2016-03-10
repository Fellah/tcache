package sletat

import (
	"encoding/xml"
	"net/http"
	"os"
)

const (
	URL          = "http://bulk.sletat.ru/Main.svc"
	ENV_LOGIN    = "SLETAT_LOGIN"
	ENV_PASSWORD = "SLETAT_PASSWORD"
)

var (
	login    = os.Getenv(ENV_LOGIN)
	password = os.Getenv(ENV_PASSWORD)
	client   = http.Client{}
)

type Request struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Header  RequestHeader
	Body    RequestBody
}

type RequestHeader struct {
	XMLName  xml.Name `xml:"Header"`
	AuthInfo AuthInfo
}

type AuthInfo struct {
	XMLName  xml.Name `xml:"urn:SletatRu:DataTypes:AuthData:v1 AuthInfo"`
	Login    string   `xml:"Login"`
	Password string   `xml:"Password"`
}

type RequestBody struct {
	XMLName    xml.Name `xml:"Body"`
	SOAPAction interface{}
}
