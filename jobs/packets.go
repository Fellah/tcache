package jobs

import (
	"log"

	"github.com/fellah/tcache/sletat"
)

const (
	WORKERS_NUM = 10
)

func FetchPacketsList(chPocketId chan<- string) {
	packets, err := sletat.FetchPacketsList("")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < WORKERS_NUM; i++ {
		chPocketId <- packets[i].Id
	}

	close(chPocketId)
}
