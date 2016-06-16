package jobs

import (
	"time"

	"github.com/fellah/tcache/db"
	"github.com/avialeta/api/log"
)

var (
	ticker = time.NewTicker(2 * time.Hour)
)

func Start() {
	for {
		Pipe()
		<-ticker.C
	}
}

func Pipe() {
	queryOperators()
	queryCities()

	t, err := makeDownloadTime()
	if err != nil {
		log.Error.Println(err)
		return
	}

	stat := new(statistic)

	//db.RemoveExpiredTours()
	//db.RemoveExistTours(t)

	packets := fetchPackets(t, stat)

	tours := fetchTours(packets, stat)

	end := saveTours(tours, stat)

	finalize(end, stat)
}

func End() {
	db.Close()
	ticker.Stop()
}

func makeDownloadTime() (string, error) {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return "", err
	}

	t := time.Now().In(location)
	t = t.Add(-2 * time.Hour)  // UTC +3h (Moscow time)

	t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())

	// Test code.
	//t = t.Add(15 * time.Minute)  // UTC +3h (Moscow time)

	return t.Format(time.RFC3339), nil
}
