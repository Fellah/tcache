package jobs

import (
	"github.com/fellah/tcache/data"
	"sync"
	"github.com/fellah/tcache/cache"
	"github.com/fellah/tcache/prefilter"
	"github.com/fellah/tcache/log"
)

func fetchPartnersData(tours <-chan data.Tour, wg *sync.WaitGroup) {
	log.Info.Println("fetchPartnersData...")
	wg.Add(1)
	go func() {
		for tour := range tours {
			// Filters
			if prefilter.ForPartnersTours(&tour) {
				cache.RegisterTourGroup(tour)
			}
		}
		log.Info.Println("fetchPartnersData done")
		wg.Done()
	}()
}
