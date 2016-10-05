package jobs

import (
	"time"

	"github.com/fellah/go-helpers/log"
	"github.com/fellah/tcache/db"
	"github.com/fellah/tcache/prefilter"
	"github.com/fellah/tcache/stat"
	"github.com/fellah/tcache/data"
)

var (
	ticker = time.NewTicker(2 * time.Hour)
)

func Start() {
	stat := stat.NewTours()

	for {
		Pipe(stat)
		<-ticker.C
	}
}

func Pipe(stat *stat.Tours) {
	queryOperators()
	queryCities()
	prefilter.PrepareData()

	t, err := makeDownloadTime()
	if err != nil {
		log.Error.Println(err)
		return
	}

	data_channels := [2]chan data.PacketInfo{
		make(chan data.PacketInfo, 100),
		make(chan data.PacketInfo, 100),
	}

	fetchPackets(data_channels, t)

	ends_channels := [2]chan bool{
		make(chan bool),
		make(chan bool),
	}

	fetchTours(data_channels[0], stat, ends_channels[0])
	fetchPartnersData(data_channels[1], stat, ends_channels[1])

	finalize(ends_channels, stat, data_channels)
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

	// Development mode.
	//t := time.Now().In(location).Add(-5 * time.Minute)

	return t.Format(time.RFC3339), nil
}
