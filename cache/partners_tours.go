package cache

import (
	"github.com/fellah/tcache/data"
	"io"
	"crypto/sha1"
	"fmt"
	"strconv"
	"github.com/fellah/tcache/log"
	"github.com/fellah/tcache/db"
)

// Struct for save tours:
// pt_tours_groups_keys - array of tour group keys
// pt_tours_groups:
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

func RegisterTourGroup(tour data.Tour) {
	h := sha1.New()

	meal_present := (tour.MealId != 117 )

	// Add data for digest
	io.WriteString(h, strconv.Itoa(tour.SourceId))
	io.WriteString(h, strconv.Itoa(tour.CountryId))
	io.WriteString(h, strconv.Itoa(tour.TownId))
	io.WriteString(h, strconv.Itoa(tour.Adults))
	io.WriteString(h, tour.Checkin)
	io.WriteString(h, strconv.Itoa(tour.Nights))
	io.WriteString(h, strconv.Itoa(tour.Kids))
	io.WriteString(h, strconv.Itoa(*(tour.Kid1Age)))
	io.WriteString(h, strconv.Itoa(*(tour.Kid2Age)))
	io.WriteString(h, strconv.Itoa(*(tour.Kid3Age)))
	io.WriteString(h, strconv.Itoa(tour.DptCityId))
	io.WriteString(h, strconv.FormatBool(meal_present))

	hash_sum := h.Sum(nil)
	hash_key := fmt.Sprintf("pt_tours_groups-%x", hash_sum)

	if redis_client.Exists(hash_key).Val() {
		old_price_str := redis_client.HGet(hash_key, "price").Val()
		old_price, err := strconv.Atoi(old_price_str)

		if err == nil && old_price > tour.Price {
			// Update tour info data
			redis_client.HMSet(hash_key, map[string]string{
				"price": strconv.Itoa(tour.Price),
				"meal_id": strconv.Itoa(tour.MealId),
				"hotel_id": strconv.Itoa(tour.HotelId),
				"tickets_included": strconv.Itoa(tour.TicketsIncluded),
				"has_econom_tickets_dpt": strconv.Itoa(tour.HasEconomTicketsDpt),
				"has_econom_tickets_rtn": strconv.Itoa(tour.HasEconomTicketsRtn),
				"hotel_is_in_stop": strconv.Itoa(tour.HotelIsInStop),
				"sletat_request_id": strconv.Itoa(tour.RequestId),
				"sletat_offer_id": strconv.FormatInt(tour.OfferId, 10),
			})
		}
	} else {
		// Save full data of record
		redis_client.HMSet(hash_key, map[string]string{
			"checkin": tour.Checkin,
			"nights": strconv.Itoa(tour.Nights),
			"adults": strconv.Itoa(tour.Adults),
			"kids": strconv.Itoa(tour.Kids),
			"kid1age": strconv.Itoa(*(tour.Kid1Age)),
			"kid2age": strconv.Itoa(*(tour.Kid2Age)),
			"kid3age": strconv.Itoa(*(tour.Kid3Age)),
			"dpt_city_id": strconv.Itoa(tour.DptCityId),
			"town_id": strconv.Itoa(tour.TownId),
			"meal_present": strconv.FormatBool(meal_present),

			"price": strconv.Itoa(tour.Price),
			"meal_id": strconv.Itoa(tour.MealId),
			"hotel_id": strconv.Itoa(tour.HotelId),
			"tickets_included": strconv.Itoa(tour.TicketsIncluded),
			"has_econom_tickets_dpt": strconv.Itoa(tour.HasEconomTicketsDpt),
			"has_econom_tickets_rtn": strconv.Itoa(tour.HasEconomTicketsRtn),
			"hotel_is_in_stop": strconv.Itoa(tour.HotelIsInStop),
			"sletat_request_id": strconv.Itoa(tour.RequestId),
			"sletat_offer_id": strconv.FormatInt(tour.OfferId, 10),
		})

		// Add hash_key to list
		redis_client.RPush("pt_tours_groups_keys", hash_key)
	}
}

func SaveTourGroupsToDB() {
	for hash_key := redis_client.LPop("pt_tours_groups_keys").Val();
		hash_key != nil;
		hash_key = redis_client.LPop("pt_tours_groups_keys").Val() {

		tour, err := redis_client.HGetAll(hash_key).Result()
		if err != nil {
			log.Error.Fatalln(err)
			continue
		}

		db.SavePartnerTour(hash_key, tour)
	}
}