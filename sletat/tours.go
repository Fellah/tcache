package sletat

import (
	"compress/gzip"
	"encoding/xml"
	"io"

	"github.com/fellah/tcache/data"
	"github.com/fellah/tcache/log"
)

const bulkCacheUrl = "http://bulk.sletat.ru/BulkCacheDownload?packetId="

func FetchTours(packetId string, tours_channel_count int) (chan data.Tour, error) {
	url := bulkCacheUrl + packetId
	log.Info.Println("Download:", url)

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}

	tours_channel := make(chan data.Tour)
	go func() {
		defer resp.Body.Close()
		defer gzipReader.Close()

		decoder := xml.NewDecoder(gzipReader)
		for {
			t, err := decoder.Token()
			if err != nil && err != io.EOF {
				log.Error.Println(err)
				break
			}

			if err == io.EOF {
				break
			}

			switch se := t.(type) {
			case xml.StartElement:
				if se.Name.Local == "tour" {
					tour := data.Tour{}
					decoder.DecodeElement(&tour, &se)

					tours_channel <- tour
				}
			}
		}

		log.Info.Println("FetchTours FINISH")
		close(tours_channel)
	}()

	return tours_channel, nil
}
