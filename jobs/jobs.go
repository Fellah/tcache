package jobs

import (
	"time"

	"github.com/fellah/tcache/sletat"
)

const WORKERS_NUM = 10

var ticker = time.NewTicker(2 * time.Hour)

func Start() {
	for {
		Pipe()
		<-ticker.C
	}
}

func Stop() {
	ticker.Stop()
}

func Pipe() {
	QueryOperators()

	chPocket := make(chan sletat.PacketInfo)

	go FetchPackets(chPocket)

	for i := 0; i < WORKERS_NUM; i++ {
		go FetchTours(chPocket)
	}
}
