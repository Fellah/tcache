package jobs

import (
	"github.com/fellah/tcache/data"
	"sync"
	"github.com/fellah/tcache/cache"
	"github.com/fellah/tcache/prefilter"
)

func fetchPartnersData(tours <-chan data.Tour, end chan<- bool) {
	wg := new(sync.WaitGroup)
	wg.Add(1)

	// Run manager for send same data to same channels (based on hash)
	collector_stop := make(chan bool)
	go collectPartnerData(tours, collector_stop, wg)

	go func() {
		wg.Wait()

		cache.SaveTourGroupsToDB()

		collector_stop <- true
		end <- true
		close(end)
	}()
}

// Manager for separate data by workers
func collectPartnerData(tours <-chan data.Tour, stop_channel <-chan bool, wg *sync.WaitGroup) {
	for {
		select {
		case tour := <-tours:
			// Filters
			if prefilter.ForPartnersTours(&tour) {
				cache.RegisterTourGroup(tour)
			}
		case <-stop_channel:
			wg.Done()
			return
		}
	}
}
