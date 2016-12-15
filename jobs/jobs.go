package jobs

import (
	"time"

	"github.com/fellah/go-helpers/log"
	"github.com/fellah/tcache/db"
	"github.com/fellah/tcache/prefilter"
	"github.com/fellah/tcache/stat"
	"github.com/fellah/tcache/cache"
	"sync"
)

var (
	ticker_save_data = time.NewTicker(60 * time.Minute)
	save_wait_group = sync.WaitGroup{}
)

func Start() {
	log.Info.Println("START....")

	stat := stat.NewTours()

	CronSaveTourGroupsToDB();

	for {
		Pipe(stat)
		stat.Idle <- 1
		time.Sleep(time.Minute)
	}
}

func Pipe(stat *stat.Tours) {
	queryOperators()
	loadFuelSurcharges()
	prefilter.PrepareData()
	cache.Init()

	t, err := makeDownloadTime()
	if err != nil {
		log.Error.Println(err)
		return
	}

	packets_channel := fetchPackets(t)

	end := make(chan bool)

	fetchTours(packets_channel, stat, end)

	finalize(end, stat)
}


func finalize(end chan bool, stat *stat.Tours) {
	// wait end signal from all channels
	<-end
	close(end)

	stat.Output()

	log.Info.Println("END")
}


func End() {
	db.Close()
	ticker_save_data.Stop()
}

func makeDownloadTime() (string, error) {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return "", err
	}

	t := time.Now().In(location).Add(-2 * time.Hour)
	t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())

	// Development mode.
	//t := time.Now().In(location).Add(-5 * time.Minute)

	return t.Format(time.RFC3339), nil
}
