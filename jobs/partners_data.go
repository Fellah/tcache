package jobs

import (
	"github.com/fellah/tcache/stat"
	"github.com/fellah/tcache/data"
	"sync"
)

const (
	workersNum = 16
)

func fetchPartnersData(packets <-chan data.PacketInfo, stat *stat.Tours, end <-chan bool) {
	wg := new(sync.WaitGroup)
	wg.Add(workersNum)

	var data_channels [workersNum]chan data.PacketInfo;
	for i := range data_channels {
		data_channels[i] = make(chan data.PacketInfo, 32)
	}

	// Run manager for send some data to some channels (based on hash)
	go managerPartnerData(packets, data_channels)

	// Run multiply workers to read from different channels
	for i := 0; i < workersNum; i++ {
		go workerPartnerData(data_channels[i])
	}

	go func() {
		wg.Wait()

		end <- true
		close(end)
	}()
}

// Manager for separate data by workers
func managerPartnerData(loaded_packets <-chan data.PacketInfo, workers []chan data.PacketInfo) {

}

// Worker for partners data
func workerPartnerData(packets <-chan data.PacketInfo) {

}
