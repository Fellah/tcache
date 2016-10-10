package sletat

import (
	"compress/gzip"
	"encoding/xml"
	"io"

	"github.com/fellah/tcache/data"
	"github.com/fellah/tcache/log"
	"github.com/fellah/tcache/prefilter"
)

const bulkCacheUrl = "http://bulk.sletat.ru/BulkCacheDownload?packetId="

func FetchTours(packetId string, tours_channel_count int) ([]chan data.Tour, error) {
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

	var tours_channels []chan data.Tour = make([]chan data.Tour, tours_channel_count)
	for i := 0; i < tours_channel_count; i++ {
		tours_channels[i] = make(chan data.Tour)
	}
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
					for index,tours := range tours_channels {
						if prefilter.ForHotel(&tour, index) {
							tours <- tour
						}
					}
				}
			}
		}

		log.Info.Println("FetchTours FINISH")
		for _,channel := range tours_channels {
			close(channel)
		}
		log.Info.Println("FetchTours tours channels closed")
	}()

	return tours_channels, nil
}
