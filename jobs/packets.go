package jobs

import (
	"log"

	"github.com/fellah/tcache/sletat"
)

const (
	WORKERS_NUM = 10
)

func FetchPacketsList(chPocketId chan<- string) {
	packets, err := sletat.FetchPacketsList("2016-03-10T20:00:00Z")
	if err != nil {
		log.Fatal(err)
	}

	for i := range packets {
		log.Println(packets[i].CreateDate)
	}

	/*for i := 0; i < WORKERS_NUM; i++ {
		chPocketId <- packets[i].Id
	}*/

	close(chPocketId)
}
