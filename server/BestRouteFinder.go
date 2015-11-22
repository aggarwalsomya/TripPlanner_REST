package main

import (
	"fmt"
	"math/rand"
		"time"
		"strconv"
)

var _ = fmt.Println;
type Coordinate struct {
		Lat string 
		Lng string 
}

var idToCoordMap = make(map[string]Coordinate)
var locationToUberRespMap = make(map[string]UberEstimatesResponse)

func clearCache() {
			idToCoordMap = make(map[string]Coordinate)
			locationToUberRespMap = make(map[string]UberEstimatesResponse)
}

func bestRouteFinder(perm []perms) ErrorableTripServiceData {
	clearCache()
	var tripdata ErrorableTripServiceData
	var finalPath perms			// this will have the final paths
	var finaluberresp UberEstimatesResponse //this will store the most effective total time, tot duration and tot distance for the final path
	
	//generate the map for the location Id to coordinates - one time task only
	populateIdToCoorMap(perm[0].perm)
	
	// find the best route from the perm
	for i:=0; i < len(perm); i++ {
		uberresp := findPathCost(perm[i])
		
		if (uberresp.ErrorMsg != "") {
			tripdata.ErrorMsg = uberresp.ErrorMsg 
			return tripdata
		}

		if(i == 0) {
			finaluberresp = uberresp
			finalPath = perm[i]
		} else {
			//compare the price which is extract from uber response for this path with the lowest price stored till now.
			if(uberresp.Price < finaluberresp.Price) {
				finaluberresp = uberresp
				finalPath = perm[i]
			} else if (uberresp.Price == finaluberresp.Price) {
				if (uberresp.Distance < finaluberresp.Distance) {
					finaluberresp = uberresp
					finalPath = perm[i]
				}
			}
		}
	}
	
	//Now when you have the best uber response and best path, generate the Trip response strcuture and return it back.
    rand.Seed(time.Now().Unix())
	tripdata.Id = strconv.Itoa(rand.Intn(999999 - 1) + 1)
	tripdata.Status ="planning"
	tripdata.Starting_from_location_id = finalPath.perm[0]
	tripdata.Best_route_location_ids = finalPath.perm[1:len(finalPath.perm) - 1]
	tripdata.Total_uber_costs = finaluberresp.Price
	tripdata.Total_uber_duration =  finaluberresp.Duration
	tripdata.Total_distance = finaluberresp.Distance
	tripdata.ErrorMsg = ""
	
	// form a ret struct to return the best path. So simple
	fmt.Println("best path is : " , tripdata);
	return tripdata
}

//this will find the total cost, distance and duration for one of the n possible paths
func findPathCost(in_perm perms)UberEstimatesResponse {

	var ret UberEstimatesResponse 
	ret.Price = 0
	ret.Distance = 0
	ret.Duration = 0
	
	// it will find the cost for this path
	permData := in_perm.perm
	for i:=0; i < len(permData) - 1; i++{
		var location_1 = permData[i]		
		var location_2 = permData[i+1]
		
		//get the cost estimate from the uber api here.
		//using Dynamic programming here. SO that uber is not called again and again
		uberesp := getTripEstimate(location_1, location_2)
		
		ret.Price = ret.Price + uberesp.Price
		ret.Distance = ret.Distance + uberesp.Distance
		ret.Duration = ret.Duration + uberesp.Duration
		ret.ErrorMsg += uberesp.ErrorMsg
	}
	
//	fmt.Println("cost of path: ", in_perm.perm, " is : ", ret);
	return ret
}

//gets the uber estimtes for 2 locations from uber only if they have not been fetched previously.
//using Dynamic programming here. SO that uber is not called again and again
func getTripEstimate(location_1, location_2 string)UberEstimatesResponse {
	_,ok := locationToUberRespMap[location_1+"-"+location_2]
	var uberresp UberEstimatesResponse

	//if value is already stored in the map, dont fetch again.
	if ok {
		uberresp = locationToUberRespMap[location_1+"-"+location_2]
	} else {
		lat_1:= idToCoordMap[location_1].Lat
		lng_1:= idToCoordMap[location_1].Lng
		
		lat_2:= idToCoordMap[location_2].Lat
		lng_2:= idToCoordMap[location_2].Lng
		fmt.Println("getting uber res for: " , location_1+"-"+location_2);
		uberresp = getUberEstimates(lat_1, lng_1, lat_2, lng_2);
		locationToUberRespMap[location_1+"-"+location_2] = uberresp
	}
	return uberresp
}

//this will populate the location id to coordinates map.
//It will not fetch again from mongo DB if it is already present.
func populateIdToCoorMap(permData []string) {
	for i:=0; i < len(permData); i++{
		var locationId = permData[i]
		_, ok := idToCoordMap[locationId]
		if ok {
		//this entry is already present in the map. DO nothing.
		} else {
		
			locationservice := getData(locationId)
			var coord Coordinate
		
			//extract the lat and long from mongo db
			lat := locationservice.Coordinate.Lat
			lng := locationservice.Coordinate.Lng
			
			//fill them in a loca struct
			coord.Lat = lat
			coord.Lng = lng
			
			idToCoordMap[locationId] = coord
		}	
	}
}
