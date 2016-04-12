package jobs

import (
	"log"
	"time"

	"github.com/fellah/tcache/sletat"
)

func fetchPackets(t time.Time) chan sletat.PacketInfo {
	log.Println("Download packets from", t.Format(time.RFC3339))

	packets, err := sletat.FetchPacketsList(t.Format(time.RFC3339))
	if err != nil {
		log.Println(err)
		return nil
	}

	return packets
}
