package main

import (
    "net/http"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/url"
    "strconv"
)

// This is a bridge class betweeb google rest API and rest of our code. 
// Our code will always use google communicator as an end point. 
// Internally this class would communicate with google using xyz (or anything on this planet) method.

type GoogleResponse struct {
	Results []GoogleResult
}

type GoogleResult struct {

	Address      string               `json:"formatted_address"`
	AddressParts []GoogleAddressPart `json:"address_components"`
	Geometry     Geometry
	Types        []string
}

type GoogleAddressPart struct {

	Name      string `json:"long_name"`
	ShortName string `json:"short_name"`
	Types     []string
}

type Geometry struct {

	Bounds   Bounds
	Location Point
	Type     string
	Viewport Bounds
}
type Bounds struct {
	NorthEast, SouthWest Point
}

type Point struct {
	Lat, Lng float64
}

func valExists(types []string, val string) bool {
	var ret bool = false
	for i := 0; i < len(types); i++ {
		if types[i] == val {
			ret = true
			break
		}
	}

	return ret;
}

func translateGoogleResponseToLocationService(res GoogleResponse) LocationService {
	var ret LocationService 
	
	if len(res.Results) == 0 {
		ret.ErrorMsg = "Empty response from Google. Not a valid address"
	} else {
		// parse the address components
		add := res.Results[0].AddressParts
		for i := 0; i < len(add); i++ {
			if valExists(add[i].Types, "locality") {
				ret.City = add[i].Name
			} else if valExists(add[i].Types, "administrative_area_level_1") {
				ret.State = add[i].Name
			} else if valExists(add[i].Types, "postal_code") {
				ret.Zip = add[i].Name
			}
		}
		
		// now fill lat and longitude
		ret.Coordinate.Lat = strconv.FormatFloat(res.Results[0].Geometry.Location.Lat,'f',7,64)
		ret.Coordinate.Lng = strconv.FormatFloat(res.Results[0].Geometry.Location.Lng,'f',7,64)
	}

	return ret;
}


func getGoogleLocation(Address string) LocationService{

	client := &http.Client{}

	reqURL := "http://maps.google.com/maps/api/geocode/json?address="
	reqURL += url.QueryEscape(Address)
	reqURL += "&sensor=false";
	fmt.Println("URL formed: "+ reqURL)

	req, err := http.NewRequest("GET", reqURL , nil)
	resp, err := client.Do(req)

	var ret LocationService
	if err != nil {
		ret.ErrorMsg = "Got error from google service. Might be an invalid address"
		fmt.Println("error in sending req to google: ", err);	
		return ret
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	
	if err != nil {
		ret.ErrorMsg = "Unable to read response from google service. Might be an invalid address"
		fmt.Println("error in reading response: ", err);	
		return ret
	}
	
	var res GoogleResponse
	err = json.Unmarshal(body, &res)
	
	if err != nil {
		ret.ErrorMsg = "Unable to unmarshall response from google service. Might be an invalid address"
		fmt.Println("error in unmashalling response: ", err);	
		return ret
	}
	
	ret = translateGoogleResponseToLocationService(res)
	return ret;
}