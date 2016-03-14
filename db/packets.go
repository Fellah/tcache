package db

import (
	"log"
	"time"

	"github.com/lib/pq"

	"github.com/fellah/tcache/sletat"
)

func RemoveExistPackets(t time.Time) {
	_, err := db.Query(`DELETE FROM packets WHERE create_date >= $1`, t)
	if err != nil {
		log.Println(err)
	}
}

func SavePackets(packets []sletat.PacketInfo) {
	txn, err := db.Begin()
	if err != nil {
		log.Println(err)
	}

	stmt, err := txn.Prepare(pq.CopyIn(
		"packets",
		"id",
		"create_date",
		"date_time_from",
		"date_time_to",
		"source_id",
		"country_id",
		"dpt_city_id",
	))
	if err != nil {
		log.Println(err)
	}

	for _, packet := range packets {
		_, err = stmt.Exec(
			packet.Id,
			packet.CreateDate,
			packet.DateTimeFrom,
			packet.DateTimeTo,
			int64(packet.SourceId),
			int64(packet.CountryId),
			int64(packet.DptCityId),
		)
		if err != nil {
			log.Println(err)
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Println(err)
	}

	// TODO: Fix error: `pq: copyin statement has already been closed`
	/*err = stmt.Close()
	if err != nil {
		log.Println(err)
	}*/

	err = txn.Commit()
	if err != nil {
		log.Println(err)
	}
}