package jobs

import (
	"log"
	"time"

	"github.com/fellah/tcache/db"
	"github.com/fellah/tcache/sletat"
)

func FetchPackets(chSavePocket chan<- sletat.PacketInfo) {
	t := time.Now().UTC()
	t = t.Add(3 * time.Hour)  // UTC +3h
	t = t.Add(-2 * time.Hour) // 2 hour

	t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())

	log.Println("Download packets from", t.Format(time.RFC3339))

	packets, err := sletat.FetchPacketsList(t.Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}

	packets = packets

	for i := range packets {
		chSavePocket <- packets[i]
	}

	db.RemoveExistPackets(t)
	db.RemoveExistTours(t)
	go db.SavePackets(packets)

	close(chSavePocket)
}
