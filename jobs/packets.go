package jobs

import (
	"github.com/fellah/tcache/data"
	"github.com/fellah/tcache/log"
	"github.com/fellah/tcache/sletat"
)

func fetchPackets(t string) chan data.PacketInfo {
	packets := make(chan data.PacketInfo)

	go func() {
		log.Info.Println("Download packets from", t)
		packetsList, err := sletat.FetchPacketsList(t)
		if err != nil {
			log.Error.Println(err)
		}

		for _, packet := range packetsList {
			if skipPacket(&packet) {
				continue
			}

			if !isOperatorActive(packet.SourceId) {
				continue
			}

			packets <- packet
		}

		close(packets)
	}()

	return packets
}

func skipPacket(packet *data.PacketInfo) bool {
	if !isDepartCityActive(packet.DptCityId) {
		return true
	}

	return false
}
