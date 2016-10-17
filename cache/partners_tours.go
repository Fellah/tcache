package cache

import (
	"github.com/fellah/tcache/data"
	"io"
	"crypto/sha1"
	"fmt"
	"strconv"
	"github.com/fellah/tcache/log"
	"github.com/fellah/tcache/db"
	"strings"
	"encoding/base64"
	"sync"
)

const (
	savePdCommitBatchSize = 1000
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

	io.WriteString(h, strconv.Itoa(tour.DptCityId))
	io.WriteString(h, strconv.FormatBool(meal_present))

	hash_sum := h.Sum(nil)
	str := base64.StdEncoding.EncodeToString(hash_sum)

	hash_key := "pt_tg-"+str

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

				"few_econom_tickets_dpt": strconv.Itoa(tour.FewEconomTicketsDpt),
				"few_econom_tickets_rtn": strconv.Itoa(tour.FewEconomTicketsRtn),
				"few_places_in_hotel": strconv.Itoa(tour.FewPlacesInHotel),
				"flags": strconv.FormatInt(tour.Flags, 10),
				"description": tour.Description,
				"tour_url": tour.TourUrl,
				"room_name": tour.RoomName,
				"receiving_party": tour.ReceivingParty,
			})
		}
	} else {
		// Check in DB



		// Save full data of record
		redis_client.HMSet(hash_key, map[string]string{
			"checkin": tour.Checkin,
			"nights": strconv.Itoa(tour.Nights),
			"adults": strconv.Itoa(tour.Adults),
			"kids": strconv.Itoa(tour.Kids),
			"kid1age": strconv.Itoa(kid1age),
			"kid2age": strconv.Itoa(kid2age),
			"kid3age": strconv.Itoa(kid3age),
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

			"few_econom_tickets_dpt": strconv.Itoa(tour.FewEconomTicketsDpt),
			"few_econom_tickets_rtn": strconv.Itoa(tour.FewEconomTicketsRtn),
			"few_places_in_hotel": strconv.Itoa(tour.FewPlacesInHotel),
			"flags": strconv.FormatInt(tour.Flags, 10),
			"description": tour.Description,
			"tour_url": tour.TourUrl,
			"room_name": tour.RoomName,
			"receiving_party": tour.ReceivingParty,
		})

		// Add hash_key to list
		redis_client.RPush("pt_tours_groups_keys", hash_key)
	}
}

func SaveTourGroupsToDB(once_flag *sync.Once) {
	count := redis_client.LLen("pt_tours_groups_keys").Val()
	log.Info.Println("SaveTourGroupsToDB START (", count, ")...")

	batch_size := savePdCommitBatchSize
	transaction, trx_err := db.StartTransaction()
	for row := redis_client.LPop("pt_tours_groups_keys");
		row.Err() == nil && trx_err == nil;
		row = redis_client.LPop("pt_tours_groups_keys") {

		hash_key := row.Val()
		if hash_key == "" {
			continue
		}

		tour, err := redis_client.HGetAll(hash_key).Result()
		if err != nil {
			log.Error.Fatalln(err)
			continue
		}

		key_parts := strings.Split(hash_key, "-")
		group_hash_str := key_parts[1]

		hash_bin, err := base64.StdEncoding.DecodeString(group_hash_str)
		var group_hash_hex string
		if err != nil {
			continue
		} else {
			group_hash_hex = fmt.Sprintf("%x", hash_bin)
		}

		db.SavePartnerTour(group_hash_hex, tour, transaction)
		redis_client.Del(hash_key)

		count--
		if count < 0 {
			break
		}

		batch_size--
		if batch_size <= 0 {
			db.CommitTransaction(transaction)
			transaction, trx_err = db.StartTransaction()
			batch_size = savePdCommitBatchSize
		}
	}
	db.CommitTransaction(transaction)
	log.Info.Println("SaveTourGroupsToDB DONE (", count, ")")

	log.Info.Println("SaveTourGroupsToDB clean partners tours...")
	db.CleanPartnerTours()
	log.Info.Println("SaveTourGroupsToDB clean partners tours DONE...")

	*once_flag = sync.Once{}
}

func ClearTourGroups() {
	for row := redis_client.LPop("pt_tours_groups_keys"); row.Err() == nil;
	    row = redis_client.LPop("pt_tours_groups_keys") {

		hash_key := row.Val()
		if hash_key == "" {
			continue
		}

		redis_client.Del(hash_key)
	}
	redis_client.Del("pt_tours_groups_keys")
}
