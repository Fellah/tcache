package jobs

import (
	"os"
	"net/http"
	"log"
	"fmt"

	"io/ioutil"
	"strings"
)

const (
	BULK_SLETAT_RU_URL = "http://bulk.sletat.ru/"
	BULK_SLETAT_RU_LOGIN = "BULK_SLETAT_RU_LOGIN"
	BULK_SLETAT_RU_PASSWORD = "BULK_SLETAT_RU_PASSWORD"

)

var (
	BulkSletatRuLogin string
	BulkSletatRuPassword string
)

var reqGetPacketList = `<soap:Envelope xmlns:soap=http://schemas.xmlsoap.org/soap/envelope/
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
  <soap:Header>
    <AuthInfo xmlns="urn:SletatRu:DataTypes:AuthData:v1">
      <Login>******</Login>
      <Password>******</Password>
    </AuthInfo>
  </soap:Header>
  <soap:Body>
  </soap:Body>
</soap:Envelope>`

func init() {
	BulkSletatRuLogin = os.Getenv(BULK_SLETAT_RU_LOGIN)
	BulkSletatRuPassword = os.Getenv(BULK_SLETAT_RU_PASSWORD)
}

func GetPacketList() {
	client := new(http.Client)

	url := BULK_SLETAT_RU_URL + "GetPacketList"

	resp, err := client.Post(url, "text/xml; charset=utf-8", strings.NewReader(reqGetPacketList))
	if err != nil {
		log.Println(err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("cool", string(b))
}
