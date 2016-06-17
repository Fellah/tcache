package jobs

import (
	"github.com/fellah/tcache/log"
	"github.com/fellah/tcache/sletat"
)

func fetchPackets(t string, stat *statistic) chan sletat.PacketInfo {
	packets := make(chan sletat.PacketInfo)

	go func() {
		log.Info.Println("Download packets from", t)
		packetsList, size, err := sletat.FetchPacketsList(t)
		if err != nil {
			log.Error.Println(err)
		}

		for _, packet := range packetsList {
			packets <- packet
		}

		stat.StorePacketsCount(uint64(len(packetsList)))
		stat.StorePacketsSize(size)

		close(packets)
	}()

	return packets
}
