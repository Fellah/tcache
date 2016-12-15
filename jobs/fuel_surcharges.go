package jobs

import (
	"github.com/fellah/tcache/data"
	"github.com/fellah/tcache/log"
	"github.com/fellah/tcache/db"
	"time"
)

var fuelSurchargesList []data.FuelSurcharge

func loadFuelSurcharges() {
	var err error

	fuelSurchargesList, err = db.QueryFuelSurchargeOrdered()
	if err != nil {
		log.Error.Fatal(err)
	}
}

func setFuelSurchargesForTour(tour *data.Tour) {
	fsMin := -1
	fsMax := -1
	for _, fsRow := range fuelSurchargesList {
		if isTourAndFuelSurchargeEqual(tour, &fsRow) {
			if fsMin == -1 {
				fsMin = fsRow.Price
			}

			if fsMax == -1 || fsRow.Price > fsMax {
				fsMax = fsRow.Price
			}
		}
	}

	if fsMin == -1 {
		fsMin = 0
	}

	if fsMax == -1 {
		fsMax = 0
	}

	tour.FuelSurchargeMin = fsMin
	tour.FuelSurchargeMax = fsMax
	tour.Price += fsMin
}

func isTourAndFuelSurchargeEqual(tour *data.Tour, fs *data.FuelSurcharge) (bool) {
	templateDate := "2006-01-02"
	templateTime := "2006-01-02T03:04:05Z"

	tourCheckin, err := time.Parse(templateDate, tour.Checkin)
	if err != nil {
		log.Error.Print(tour.Checkin)
		log.Error.Print(err)
	}

	fsStartTime, err := time.Parse(templateTime, fs.StartDate)
	if err != nil {
		log.Error.Print(fs.StartDate)
		log.Error.Print(err)
	}

	fsEndTime, err := time.Parse(templateTime, fs.EndDate)
	if err != nil {
		log.Error.Print(fs.EndDate)
		log.Error.Print(err)
	}

	return (tour.SourceId == fs.SourceId &&
			tourCheckin.Unix() >= fsStartTime.Unix() &&
			tourCheckin.Unix() <= fsEndTime.Unix() &&
			(fs.DptCityId == 0 || fs.DptCityId == tour.DptCityId) &&
			(fs.TownId == 0 || fs.TownId == tour.TownId))
}
