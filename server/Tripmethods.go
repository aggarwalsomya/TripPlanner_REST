package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

var _ = fmt.Println

func CreateTripPlan(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	var req TripServiceReq
	var res ErrorableTripServiceData
	err = json.Unmarshal(body, &req)
	if err != nil {
		var err ErrorMsg
		err.ErrorMsg = "Failed to decode the request."
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(err)
		return
	}

	fmt.Println("got the req: ", req)
	// get the best path
	ret_g = nil
	generatePerm(req.Starting_from_location_id, req.Location_ids)
	//	// use it and then clear again
	res = bestRouteFinder(ret_g)
	res.IdGlobal = res.Id
	res.Next_destination_location_id = ""
	res.Uber_wait_time_eta = -1

	ret_g = nil

	if res.ErrorMsg != "" {
		var err ErrorMsg

		err.ErrorMsg = "Failed to find the best path. Must be some error in input"
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(err)
		return
	}

	// set in mongo also
	success := setUberTripData(res.Id, res)
	if !success {
		var err ErrorMsg
		fmt.Println("Unable to create an entry for trip data in the database")
		err.ErrorMsg = "Unable to create an entry in the database"
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(res)

	} else {
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(res.UberTripServiceData.TripServiceData)
	}
}

func GetTripPlan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trip_id := vars["trip_id"]
	var res ErrorableTripServiceData

	//Get the response from Mongo Db for this Location_Id
	res = getUberTripData(trip_id)
	if res.ErrorMsg != "" {
		var err ErrorMsg

		err.ErrorMsg = "trip_id doesn't exist"
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(err)
		return
	}

	//change this res to the response which needs to be sent back
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(res.UberTripServiceData.TripServiceData)
}

func RequestNextRide(w http.ResponseWriter, res ErrorableTripServiceData, startLoc string, endLoc string) {
	res.Next_destination_location_id = endLoc
	startLocData := getData(startLoc)
	if res.ErrorMsg != "" {
		var err ErrorMsg

		err.ErrorMsg = "start id doesn't exist"
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(err)
		return
	}

	endLocData := getData(endLoc)
	if res.ErrorMsg != "" {
		var err ErrorMsg

		err.ErrorMsg = "end loc id doesn't exist"
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(err)
		return
	}

	// get the product id

	prod_id := getUberProductId(startLocData.Coordinate.Lat, startLocData.Coordinate.Lng)
	if prod_id == "" {
		var err ErrorMsg

		err.ErrorMsg = "Failed to get any available uber products at this time. May be nothing is available"
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
		return
	}

	uber_eta, errMsg := requestUberForRide(startLocData.Coordinate.Lat, startLocData.Coordinate.Lng,
		endLocData.Coordinate.Lat, endLocData.Coordinate.Lng, prod_id)

	if uber_eta < 0 || errMsg != "" {
		var err ErrorMsg

		err.ErrorMsg = "Failed to request the ride from uber: " + errMsg
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
		return
	}

	res.Uber_wait_time_eta = int(uber_eta)
	
	//save in mongo db
	success := updateUberTripData(res.Id, res)
	if !success {
		var err ErrorMsg

		fmt.Println("Unable to update the trip status in the database")
		err.ErrorMsg = "Unable to update the trip status in the database"
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(res.UberTripServiceData)
	}
}

func PutTripPlan(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in the put method")
	vars := mux.Vars(r)
	trip_id := vars["trip_id"]

	//Get the response from Mongo Db for this Location_Id
	res := getUberTripData(trip_id)
	if res.ErrorMsg != "" {
		var err ErrorMsg

		err.ErrorMsg = "trip_id doesn't exist"
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(err)
		return
	}

	if res.Status == "completed" || res.Next_destination_location_id == res.Starting_from_location_id {
		res.Status = "completed"
		res.Uber_wait_time_eta = 0
		res.Next_destination_location_id = ""
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(res.UberTripServiceData)
		return
	}

	// find the next destination
	var startLoc string
	var endLoc string

	if res.Next_destination_location_id == "" {
		startLoc = res.Starting_from_location_id
		endLoc = res.Best_route_location_ids[0]
		res.Status = "requesting"
	} else {
		res.Status = "requesting"
		for i := 0; i < len(res.Best_route_location_ids); i++ {
			if res.Best_route_location_ids[i] == res.Next_destination_location_id {
				startLoc = res.Next_destination_location_id
				if i == len(res.Best_route_location_ids)-1 {
					// set home as end point
					endLoc = res.Starting_from_location_id
				} else {
					endLoc = res.Best_route_location_ids[i+1]
				}
			}
		}
	}

	if startLoc == "" || endLoc == "" {
		// some error
		var err ErrorMsg
		fmt.Println("wrong state is set in db for this trip: ", trip_id)
		err.ErrorMsg = "wrong state is set in db for this trip"
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
		return
	}

	// get the coordinates and set the request
	RequestNextRide(w, res, startLoc, endLoc)
}
