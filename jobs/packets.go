package jobs

import (
	"time"

	"github.com/fellah/tcache/log"
	"github.com/fellah/tcache/sletat"
)

func fetchPackets(t time.Time) chan sletat.PacketInfo {
	packets := make(chan sletat.PacketInfo)

	go func() {
		log.Info.Println("Download packets from", t.Format(time.RFC3339))
		packetsList, err := sletat.FetchPacketsList(t.Format(time.RFC3339))
		if err != nil {
			log.Error.Println(err)
		}

		for _, packet := range packetsList {
			packets <- packet
		}

		close(packets)
	}()

	return packets
}
