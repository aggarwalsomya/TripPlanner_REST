package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

func getSession() *mgo.Session {
	// Connect to our local mongo
	session, err := mgo.Dial("mongodb://somya:somya@ds041164.mongolab.com:41164/test1")

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session
}

func getData(location_id string) LocationService {
	session := getSession()
	defer session.Close()
	c := session.DB("test1").C("v1")

	result := LocationService{}
	fmt.Println("getting record for : " + location_id);
	err := c.Find(bson.M{"_id":location_id}).One(&result)
    if err != nil {
    	result.ErrorMsg = "invalid location id"
	}

	return result
}

func updateData(location_id string , loc LocationService) bool {
	session := getSession()
	defer session.Close()
	c := session.DB("test1").C("v1")
	
	id := bson.M{"_id": location_id}
	err := c.Update(id, loc)
	if err != nil {
		panic(err)
	}
	return true	
}

func setData(location_id string, loc LocationService) bool {
	session := getSession()
	defer session.Close()
	c := session.DB("test1").C("v1")
	err := c.Insert(&loc)
	if err != nil {
		panic(err)
		return false
	}
	return true
}

func deleteData(location_id string) bool {
	session := getSession()
	defer session.Close()
	c := session.DB("test1").C("v1")
	err:= c.Remove(bson.M{"_id": location_id})
	if err != nil {
		return false
	}
	return true
}

func setUberTripData(trip_id string, tripData UberTripServiceData) bool {
	session := getSession()
	defer session.Close()
	c := session.DB("test1").C("uber_trip_details")
	err := c.Insert(&tripData)
	if err != nil {
		panic(err)
		return false
	}
	return true
}

func getUberTripData(trip_id string) UberTripServiceData {
	session := getSession()
	defer session.Close()
	c := session.DB("test1").C("uber_trip_details")

	result := UberTripServiceData{}
	fmt.Println("Getting record for : " + trip_id);
	err := c.Find(bson.M{"_id":trip_id}).One(&result)
    if err != nil {
    	result.ErrorMsg = "Invalid Trip id"
	}

	return result
}


func updateUberTripData(trip_id string , trip UberTripServiceData) bool {
	session := getSession()
	defer session.Close()
	c := session.DB("test1").C("uber_trip_details")
	
	id := bson.M{"_id": trip_id}
	err := c.Update(id, trip)
	if err != nil {
		panic(err)
		return false
	}
	return true	
}

