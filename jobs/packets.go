package jobs

import (
	"log"
	"time"

	"github.com/fellah/tcache/sletat"
)

func FetchPacketsList(chPocketId chan<- string) {
	t := time.Now().UTC()
	t = t.Add(3 * time.Hour)  // UTC +3h
	t = t.Add(-2 * time.Hour) // 2 hour

	packets, err := sletat.FetchPacketsList(t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}

	for i := range packets {
		chPocketId <- packets[i].Id
	}

	close(chPocketId)
}
