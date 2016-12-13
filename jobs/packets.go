package jobs

import (
	"github.com/fellah/tcache/data"
	"github.com/fellah/tcache/log"
	"github.com/fellah/tcache/sletat"
	"github.com/fellah/tcache/prefilter"
)

func fetchPackets(t string) (channel chan data.PacketInfo) {
	packets := make(chan data.PacketInfo)

	go func() {
		log.Info.Println("Download packets from", t)
		packetsList, err := sletat.FetchPacketsList(t)
		if err != nil {
			log.Error.Println(err)
		}

		log.Info.Println("fetchPackets list...")
		for _, packet := range packetsList {
			if skipPacket(&packet) {
				log.Info.Println("fetchPackets packet skip...")
				continue
			}

			log.Info.Println("fetchPackets packet to work")
			channel <- packet
		}

		close(packets)
		log.Info.Println("fetchPackets done")
	}()

	return packets
}

func skipPacket(packet *data.PacketInfo) bool {
	if !prefilter.IsDepartCityActive(packet.DptCityId) &&
	   !prefilter.IsPartnersDepartCityActive(packet.DptCityId) {
		return true
	}

	return false
}
