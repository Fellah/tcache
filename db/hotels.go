package db

import (
	"github.com/fellah/tcache/log"
)

func QueryHotels() ([]int, error) {
	rows, err := db.Query("SELECT sletat_hotel_id FROM sletat_hotels WHERE active = true AND images_count > 0 AND name IS NOT NULL AND name != ''")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var hotelId int
	hotelsIds := make([]int, 0)

	for rows.Next() {
		err = rows.Scan(&hotelId)
		if err != nil {
			log.Error.Println(err)
		}

		hotelsIds = append(hotelsIds, hotelId)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return hotelsIds, nil
}
