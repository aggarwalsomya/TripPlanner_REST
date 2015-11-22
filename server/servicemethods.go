package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// defines all the REST APIs handlers

// this will vlaidate the response from google with the req asked by user
func ValidateResponseWithRequest(res LocationService, req LocationService) bool {
	return res.Zip == req.Zip || (res.Coordinate.Lat != "" && res.Coordinate.Lng != "")
}

// kind of stubs
func CreateLocation(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	var req LocationService
	err = json.Unmarshal(body, &req)
	if err != nil {
		var err ErrorMsg

		err.ErrorMsg = "Failed to decode the request."
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(err)
		return
	}

	googleresp := getGoogleLocation(req.Address + "+" + req.City + "+" + req.State + "+" + req.Zip)
	fmt.Println("resp is: ", googleresp)

	if !ValidateResponseWithRequest(googleresp, req) {
		// modify the request object itself
		var err ErrorMsg

		err.ErrorMsg = "Invalid Address. No such address exists as per Google service"
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(err)
		return
	} else {
		rand.Seed(time.Now().Unix())
		req.Id = strconv.Itoa(rand.Intn(9999-1) + 1)

		req.Coordinate.Lat = googleresp.Coordinate.Lat
		req.Coordinate.Lng = googleresp.Coordinate.Lng

		// TODO: store the response in the mongo db
		success := setData(string(req.Id), req)
		if !success {
			fmt.Println("Unable to create an entry in the database")
			req.ErrorMsg = "Unable to create an entry in the database"
			w.WriteHeader(500)
		} else {
			w.WriteHeader(201)
		}
	}
	
	json.NewEncoder(w).Encode(req)
}

func GetLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	location_id := vars["location_id"]
	fmt.Println("Id getting is::" + location_id)
	//location_id := "1234"
	fmt.Println("id in get is: " + location_id)
	var res LocationService

	//Get the response from Mongo Db for this Location_Id
	res = getData(location_id)
	if res.ErrorMsg != "" {
		var err ErrorMsg

		err.ErrorMsg = "Location id doesn't exist"
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(err)
		return
	}

	//change this res to the response which needs to be sent back
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(res)
}

func DeleteLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	location_id := vars["location_id"]
	//Delete this location id data from Mongo Db
	success := deleteData(location_id)
	if !success {
		var err ErrorMsg
		fmt.Println("Unable to delete entry from Mongo Db")
		err.ErrorMsg = "Unable to delete data in the database. Probably no such location exist."
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(err)
	} else {
		w.WriteHeader(200)
	}
}

func setNonEmpty(p *string, s1 string, s2 string) {
	if s2 != "" {
		*p = s2
	} else {
		*p = s1
	}
}

func mergeLocations(oldLoc LocationService, newLoc LocationService) LocationService {

	var ret LocationService
	setNonEmpty(&ret.Name, oldLoc.Name, newLoc.Name)
	setNonEmpty(&ret.Address, oldLoc.Address, newLoc.Address)
	setNonEmpty(&ret.City, oldLoc.City, newLoc.City)
	setNonEmpty(&ret.State, oldLoc.State, newLoc.State)
	setNonEmpty(&ret.Zip, oldLoc.Zip, newLoc.Zip)

	setNonEmpty(&ret.Coordinate.Lat, oldLoc.Coordinate.Lat, newLoc.Coordinate.Lat)
	setNonEmpty(&ret.Coordinate.Lng, oldLoc.Coordinate.Lng, newLoc.Coordinate.Lng)

	return ret
}

func PutLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	location_id := vars["location_id"]
	var res LocationService

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		var err ErrorMsg

		err.ErrorMsg = "Failed to decode the request."
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(err)
		return
	}

	// first check that whether an id exists in db or not.
	oldLoc := getData(location_id)
	if oldLoc.ErrorMsg != "" {
		var err ErrorMsg

		err.ErrorMsg = "Location id doesn't exist"
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(err)
		return
	}

	var req LocationService
	err = json.Unmarshal(body, &req)

	googleresp := getGoogleLocation(req.Address + "+" + req.City + "+" + req.State + "+" + req.Zip)
	fmt.Println("resp is: ", googleresp)

	if !ValidateResponseWithRequest(googleresp, req) {
		// modify the request object itself
		res.ErrorMsg = "Invalid Address, cannot update. No such address exists as per Google service"
		w.WriteHeader(400)
	} else {
		googleresp.Name = req.Name
		googleresp.Address = req.Address
		res = mergeLocations(oldLoc, googleresp)
		res.Id = location_id

		// TODO: store the response in the mongo db for this Location ID
		success := updateData(location_id, res)
		if !success {
			var err ErrorMsg

			fmt.Println("Unable to update data in the database")
			err.ErrorMsg = "Unable to update data in the database"
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(err)
			return

		} else {
			w.WriteHeader(201)
		}
	}

	json.NewEncoder(w).Encode(res)
}
