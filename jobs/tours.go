package jobs

import (
	"log"
	"time"

	"github.com/fellah/tcache/sletat"
)

func FetchBulkCacheDownload(chPocketId <-chan string) {
	for {
		select {
		case pocketId, ok := <-chPocketId:
			if !ok {
				return
			}
			tours, err := sletat.FetchBulkCacheDownload(pocketId)
			if err != nil {
				log.Println(err)
			}

			log.Println(len(tours))

		case <-time.After(900 * time.Second):
			log.Println("TIMEOUT")
			return
		}
	}

}
