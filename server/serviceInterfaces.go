package main

import (
)

type LocationService struct {
	Id string `json:"id" bson:"_id"`
	Name string `json:"name"`
	Address string `json:"address"`
	City string `json:"city"`
	State string `json:"state"`
	Zip string `json:"zip"`
	Coordinate struct {
		Lat string `json:"lat"`
		Lng string `json:"lng"`
	}  `json:"coordinate"`
}


type TripServiceReq struct {
	Starting_from_location_id string `json:"starting_from_location_id"`
	Location_ids []string `json:"location_ids"`
}


type UberEstimatesResponse struct {
	Price float32
	Distance float32
	Duration float32
	ErrorMsg string
}

type ErrorMsg struct {
	ErrorMsg string `json:"error"`
}

type TripServiceData struct {
	Id string `json:"id"`
	Status string `json:"status"`
	Starting_from_location_id string `json:"starting_from_location_id"`
	Best_route_location_ids []string `json:"best_route_location_ids"`
	Total_uber_costs float32 `json:"total_uber_costs"`
	Total_uber_duration float32 `json:"total_uber_duration"`
	Total_distance float32 `json:"total_distance"`
}

type UberTripServiceData struct {
	IdGlobal string `json:"id" bson:"_id"`
	TripServiceData				// anonymous field
	Next_destination_location_id string `json:"next_destination_location_id"`
	Uber_wait_time_eta int `json:"uber_wait_time_eta"`
}
