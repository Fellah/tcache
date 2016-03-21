package jobs

import (
	"log"
	"time"

	"github.com/fellah/tcache/sletat"
)

func fetchPackets(t time.Time) <-chan sletat.PacketInfo {
	log.Println("Download packets from", t.Format(time.RFC3339))

	chPacket := make(chan sletat.PacketInfo)

	go func(chPacket chan sletat.PacketInfo) {
		err := sletat.FetchPacketsList(t.Format(time.RFC3339), chPacket)
		if err != nil {
			log.Fatal(err)
		}
		close(chPacket)
	}(chPacket)

	return chPacket
}
