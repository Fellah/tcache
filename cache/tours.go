package cache

import (
	"github.com/fellah/tcache/data"
	"io"
	"crypto/sha1"
	"strconv"
	"github.com/fellah/tcache/log"
	"github.com/fellah/tcache/db"
	"encoding/base64"
	"sync"
)

const (
	saveMapCommitBatchSize = 1000
	redisMapToursListKey = "mt_tours_groups_keys"
)

// Struct for save tours:
// mt_tours_groups_keys - array of tour group keys
// mt_tg_X:
//   key - sha1 from tour group key data (group key)
//   value - hash data:
//     // group data:
//     checkin, nights, adults, kids, kid1age, kid2age, kid3age, dpt_city_id, town_id, meal_present, operator_id
//     // info data:
//     meal_id
//     price - minimal price of tour
//     hotel_id
//     tickets_included
//     has_econom_tickets_dpt
//     has_econom_tickets_rtn
//     hotel_is_in_stop
//     request_id
//     offer_id

func RegisterMapTourGroup(tour data.Tour) {
	h := sha1.New()

	// Add data for digest
	// hotel_id, checkin, dpt_city_id, nights, adults, meal_id, kids, kid1age, kid2age, kid3age
	io.WriteString(h, strconv.Itoa(tour.HotelId))
	io.WriteString(h, tour.Checkin)
	io.WriteString(h, strconv.Itoa(tour.DptCityId))
	io.WriteString(h, strconv.Itoa(tour.Nights))
	io.WriteString(h, strconv.Itoa(tour.Adults))
	io.WriteString(h, strconv.Itoa(tour.MealId))
	io.WriteString(h, strconv.Itoa(tour.Kids))

	kid1age := -1
	if tour.Kid1Age != nil {
		kid1age = *(tour.Kid1Age)
	}

	io.WriteString(h, strconv.Itoa(kid1age))

	kid2age := -1
	if tour.Kid2Age != nil {
		kid2age = *(tour.Kid2Age)
	}

	io.WriteString(h, strconv.Itoa(kid2age))

	kid3age := -1
	if tour.Kid3Age != nil {
		kid3age = *(tour.Kid3Age)
	}

	io.WriteString(h, strconv.Itoa(kid3age))

	hash_sum := h.Sum(nil)
	str := base64.StdEncoding.EncodeToString(hash_sum)

	hash_key := "mt_tg-"+str

	if redis_client.Exists(hash_key).Val() {
		old_price_str := redis_client.HGet(hash_key, "price").Val()
		old_price, err := strconv.Atoi(old_price_str)

		if err == nil && old_price > tour.Price {
			// Update tour info data
			redis_client.HMSet(hash_key, map[string]string{
				"price": strconv.Itoa(tour.Price),
				"town_id": strconv.Itoa(tour.TownId),
				"tickets_included": strconv.Itoa(tour.TicketsIncluded),
				"has_econom_tickets_dpt": strconv.Itoa(tour.HasEconomTicketsDpt),
				"has_econom_tickets_rtn": strconv.Itoa(tour.HasEconomTicketsRtn),
				"hotel_is_in_stop": strconv.Itoa(tour.HotelIsInStop),

				"currency_id": strconv.Itoa(tour.CurrencyId),
				"create_date": tour.CreateDate,
				"update_date": tour.UpdateDate,
				"price_byr": strconv.Itoa(tour.PriceByr),
				"price_eur": strconv.Itoa(tour.PriceEur),
				"price_usd": strconv.Itoa(tour.PriceUsd),
			})
		}
	} else {
		// Save full data of record
		redis_client.HMSet(hash_key, map[string]string{
			"hotel_id": strconv.Itoa(tour.HotelId),
			"checkin": tour.Checkin,
			"nights": strconv.Itoa(tour.Nights),
			"adults": strconv.Itoa(tour.Adults),
			"kids": strconv.Itoa(tour.Kids),
			"kid1age": strconv.Itoa(kid1age),
			"kid2age": strconv.Itoa(kid2age),
			"kid3age": strconv.Itoa(kid3age),
			"dpt_city_id": strconv.Itoa(tour.DptCityId),
			"meal_id": strconv.Itoa(tour.MealId),

			"price": strconv.Itoa(tour.Price),
			"town_id": strconv.Itoa(tour.TownId),
			"tickets_included": strconv.Itoa(tour.TicketsIncluded),
			"has_econom_tickets_dpt": strconv.Itoa(tour.HasEconomTicketsDpt),
			"has_econom_tickets_rtn": strconv.Itoa(tour.HasEconomTicketsRtn),
			"hotel_is_in_stop": strconv.Itoa(tour.HotelIsInStop),

			"currency_id": strconv.Itoa(tour.CurrencyId),
			"create_date": tour.CreateDate,
			"update_date": tour.UpdateDate,
			"price_byr": strconv.Itoa(tour.PriceByr),
			"price_eur": strconv.Itoa(tour.PriceEur),
			"price_usd": strconv.Itoa(tour.PriceUsd),
		})

		// Add hash_key to list
		redis_client.RPush(redisMapToursListKey, hash_key)
	}
}

func SaveMapTourGroupsToDB(once_flag *sync.Once) {
	count := redis_client.LLen(redisMapToursListKey).Val()
	log.Info.Println("SaveMapTourGroupsToDB START (", count, ")...")

	batch_size := saveMapCommitBatchSize
	transaction, trx_err := db.StartTransaction()
	for row := redis_client.LPop(redisMapToursListKey);
		row.Err() == nil && trx_err == nil;
		row = redis_client.LPop(redisMapToursListKey) {

		hash_key := row.Val()
		if hash_key == "" {
			continue
		}

		tour, err := redis_client.HGetAll(hash_key).Result()
		if err != nil {
			log.Error.Fatalln(err)
			continue
		}

		db.SaveMapTour(tour, transaction)
		redis_client.Del(hash_key)

		count--
		if count < 0 {
			break
		}

		batch_size--
		if batch_size <= 0 {
			db.CommitTransaction(transaction)
			transaction, trx_err = db.StartTransaction()
			batch_size = saveMapCommitBatchSize
		}
	}
	db.CommitTransaction(transaction)
	log.Info.Println("SaveMapTourGroupsToDB DONE (", count, ")")

	log.Info.Println("SaveMapTourGroupsToDB clean partners tours...")
	db.CleanMapTours()
	log.Info.Println("SaveMapTourGroupsToDB clean partners tours DONE...")

	*once_flag = sync.Once{}
}

