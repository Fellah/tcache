package jobs

import (
	"sync/atomic"

	"github.com/fellah/tcache/log"
)

type statistic struct {
	packets packetsStatistic
	tours   toursStatistic
}

type packetsStatistic struct {
	total uint64
	size  uint64
}

type toursStatistic struct {
	total uint64
	size  uint64
}

func (s *statistic) StorePacketsCount(n uint64) {
	atomic.StoreUint64(&s.packets.total, n)
}

func (s *statistic) StorePacketsSize(n uint64) {
	atomic.StoreUint64(&s.packets.size, n)
}

func (s *statistic) Output() {
	packetsTotal := atomic.LoadUint64(&s.packets.total)
	packetsSize := atomic.LoadUint64(&s.packets.size)

	packetsSize = packetsSize / (1024 * 1024)

	log.Info.Printf("Packets: %d (%d MB)", packetsTotal, packetsSize)
}
