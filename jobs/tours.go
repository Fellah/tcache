package jobs

import (
	"log"
	"time"

	"github.com/fellah/tcache/db"
	"github.com/fellah/tcache/sletat"
)

func FetchTours(chPocket <-chan sletat.PacketInfo) {
	for {
		select {
		case pocket, ok := <-chPocket:
			if !ok {
				return
			}
			tours, err := sletat.FetchTours(pocket.Id)
			if err != nil {
				log.Println(err)
			}

			for i := range tours {
				tours[i].CreateDate = pocket.CreateDate
			}

			go db.SaveTours(tours)

		case <-time.After(900 * time.Second):
			log.Println("TIMEOUT")
			return
		}
	}
}
