package jobs

import (
	"time"

	"github.com/fellah/go-helpers/log"
	"github.com/fellah/tcache/db"
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

	//db.RemoveExpiredTours()
	//db.RemoveExistTours(t)

	packets := fetchPackets(t)

	tours := fetchTours(packets)

	end := saveTours(tours)

	finalize(end)
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

	t := time.Now().In(location).Add(-2 * time.Hour)
	t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())

	//t := time.Now().In(location).Add(-5 * time.Minute)

	return t.Format(time.RFC3339), nil
}
