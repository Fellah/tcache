package jobs

import (
	"time"

	"github.com/fellah/tcache/db"
	"github.com/fellah/tcache/sletat"
)

const WORKERS_NUM = 16

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

	t := time.Now().UTC()
	t = t.Add(3 * time.Hour)  // UTC +3h
	t = t.Add(-2 * time.Hour) // 2 hour

	// Set time to the hour begin.
	t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())

	chPacket := fetchPackets(t)
	chTour := make(chan sletat.Tour)
	for i := 0; i < WORKERS_NUM; i++ {
		go fetchTours(chPacket, chTour)
	}
	db.RemoveExistTours(t)
	go saveTours(chTour)
}
