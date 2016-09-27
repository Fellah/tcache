package sletat

import (
	"io"
	"compress/gzip"
	"encoding/xml"

	"github.com/fellah/tcache/log"
	"github.com/fellah/tcache/data"
	"github.com/fellah/tcache/prefilter"
)

const bulkCacheUrl = "http://bulk.sletat.ru/BulkCacheDownload?packetId="

func FetchTours(packetId string) (chan data.Tour, error) {
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

	tours := make(chan data.Tour)
	go func() {
		defer resp.Body.Close()
		defer gzipReader.Close()
		defer close(tours)

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

					if prefilter.TourEnable(&tour) {
						tours <- tour
					}
				}
			}
		}
	}()

	return tours, nil
}
