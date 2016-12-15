package db

import (
	"github.com/fellah/tcache/data"
	"github.com/fellah/tcache/log"
)

func QueryFuelSurchargeOrdered() ([]data.FuelSurcharge, error) {
	rows, err := db.Query("SELECT COALESCE(source_id, 0), COALESCE(dpt_city_id, 0), " +
		"COALESCE(country_id, 0), COALESCE(town_id, 0), COALESCE(airport_id, 0), " +
		"COALESCE(aircompany_id, 0), COALESCE(flight_number, ''), COALESCE(host_id, 0)," +
		"start_date, end_date, COALESCE(currency_id, 0), COALESCE(price, 0) " +
		"FROM sletat_fuel_surcharges " +
		"ORDER BY source_id, start_date, " +
		"end_date, dpt_city_id, country_id, town_id")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	fsRows := make([]data.FuelSurcharge, 0)

	for rows.Next() {
		var row data.FuelSurcharge

		err = rows.Scan(&row.SourceId, &row.DptCityId, &row.CountryId, &row.TownId, &row.AirportId,
				&row.AirCompanyId, &row.FlightNumber, &row.HostId, &row.StartDate,
				&row.EndDate, &row.CurrencyId, &row.Price)
		if err != nil {
			log.Error.Println(err)
		}

		fsRows = append(fsRows, row)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return fsRows, nil
}
