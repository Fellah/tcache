package jobs

import (
	"github.com/fellah/tcache/data"
	"sync"
	"github.com/fellah/tcache/cache"
	"github.com/fellah/tcache/prefilter"
)

const (
	pWorkersNum = 16
)

func fetchPartnersData(tours <-chan data.Tour, end chan<- bool) {
	wg := new(sync.WaitGroup)
	wg.Add(pWorkersNum)

	var data_channels []chan data.Tour = make([]chan data.Tour, pWorkersNum)
	for i := range data_channels {
		data_channels[i] = make(chan data.Tour, 32)
	}

	// Run manager for send same data to same channels (based on hash)
	manager_stop := make(chan bool)
	go managerPartnerData(tours, data_channels, manager_stop)

	// Run multiply workers to read from different channels
	for i := 0; i < pWorkersNum; i++ {
		go workerPartnerData(data_channels[i], wg)
	}

	go func() {
		wg.Wait()

		manager_stop <- true
		end <- true
		close(end)
		close(manager_stop)
	}()
}

// Manager for separate data by workers
func managerPartnerData(tours <-chan data.Tour, workers []chan data.Tour, stop_channel <-chan bool) {
	for {
		select {
		case tour := <-tours:
			// Filters
			if prefilter.ForPartnersTours(tour) {
				cache.RegisterTourGroup(tour)
			}
		case <-stop_channel:
			return
		}
	}
}

// Worker for partners data
func workerPartnerData(packets <-chan data.Tour, wg *sync.WaitGroup) {






	wg.Done()
}
