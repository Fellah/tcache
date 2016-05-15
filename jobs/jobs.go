package jobs

import (
	"time"

	"github.com/fellah/tcache/db"
)

var ticker = time.NewTicker(2 * time.Hour)

func Start() {
	for {
		Pipe()
		<-ticker.C
	}
}

func Pipe() {
	queryOperators()

	t := time.Now().UTC()
	t = t.Add(3 * time.Hour)  // UTC +3h
	t = t.Add(-2 * time.Hour) // 2 hour

	// Set time to the hour begin.
	t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())

	db.RemoveExpiredTours()
	db.RemoveExistTours(t)

	packets := fetchPackets(t)

	tours := fetchTours(packets)

	end := saveTours(tours)

	finalize(end)
}

func End() {
	db.Close()
	ticker.Stop()
}
