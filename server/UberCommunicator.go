package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"bytes"
)

type UberResponse struct {
	Prices []UberResult
}

type UberResult struct {
	Localized_display_name string `json:"localized_display_name"`
	Duration               float32 `json:"duration"`
	Distance               float32 `json:"distance"`
	Low_estimate           float32 `json:"low_estimate"`
}

type ErrorDetail struct {
	Title string `json:title"`
}

type UberETA struct {
	ETA float32 `json:eta"`
	Errors []ErrorDetail `json:errors"`
}

type UberProducts struct {
	Products []UberProduct
}

type UberProduct struct {
	Display_name string `json:"display_name"`
	Product_id string `json:"product_id"`
}

type UberRideRequest struct {
	ProductId string 	`json:"product_id"`
	Start_Lat string 	`json:"start_latitude"`
	Start_Long	string `json:"start_longitude"`
	End_Lat		string `json:"end_latitude"`
	End_Long	string	`json:"end_longitude"`
}

func translateUberResponse(res UberResponse) UberEstimatesResponse {
	var ret UberEstimatesResponse
	
	if len(res.Prices) == 0 {
		ret.ErrorMsg = "Empty response from Uber. Not a valid address"
	} else {
		add := res.Prices
		for i := 0; i < len(add); i++ {
			if add[i].Localized_display_name == "uberX" {
				ret.Duration = add[i].Duration
				ret.Distance = add[i].Distance
				ret.Price = add[i].Low_estimate
			}
		}
	}
	return ret
}

//this function will give the details like price, distance and time for the 2 locations passed.
func getUberEstimates(StartLat, StartLong, EndLat, EndLong string) UberEstimatesResponse {

	client := &http.Client{}

	reqURL := "https://sandbox-api.uber.com/v1/estimates/price?start_latitude="
	reqURL += url.QueryEscape(StartLat)

	reqURL += "&start_longitude="
	reqURL += url.QueryEscape(StartLong)

	reqURL += "&end_latitude="
	reqURL += url.QueryEscape(EndLat)

	reqURL += "&end_longitude="
	reqURL += url.QueryEscape(EndLong)

	reqURL += "&server_token=5xVAKNR9mZj2BuXh7-_gKTWtvScwEJiel5-1iD5r"

	fmt.Println("URL formed: " + reqURL)

	req, err := http.NewRequest("GET", reqURL, nil)
	resp, err := client.Do(req)

	var ret UberEstimatesResponse
	if err != nil {
		ret.ErrorMsg = "Got error from Uber service. Might be invalid coordinates"
		fmt.Println("error in sending req to Uber: ", err)
		return ret
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		ret.ErrorMsg = "Got error from Uber service. Might be invalid coordinates"
		fmt.Println("error in reading response: ", err)
		return ret
	}

	var res UberResponse
	err = json.Unmarshal(body, &res)
	fmt.Println("resp is " , res);
	

	if err != nil {
		ret.ErrorMsg = "Unable to unmarshall response from Uber service. Might be invalid coordinates"
		fmt.Println("error in unmashalling response: ", err)
		return ret
	}

	ret = translateUberResponse(res)
	fmt.Println("Return struct is:",ret);
	return ret
}

// it will reutrn the eta for the uber. -1 for any error
func requestUberForRide(startLat string, startLon string, destLat string, destLong string, product_id string) (eta float32, errMsg string) {

	client := &http.Client{}

	eta = -1
	
	var rideReq UberRideRequest
	rideReq.ProductId = product_id
	rideReq.Start_Lat = startLat
	rideReq.Start_Long = startLon
	rideReq.End_Lat = destLat
	rideReq.End_Long = destLong
	b, err := json.Marshal(rideReq)
	if err != nil {
		errMsg = "error in json marshalling:"
		return  eta, errMsg
	}
	
	reqURL := "https://sandbox-api.uber.com/v1/requests?"
	reqURL += "access_token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsicmVxdWVzdCJdLCJzdWIiOiJiMWQxODcyYS05NGE0LTQ1MDQtYWU5Yi00OGM1Y2FhZGMyNDciLCJpc3MiOiJ1YmVyLXVzMSIsImp0aSI6ImYxNTUzNDc1LWEyZDEtNDU2Yy04ZDU0LWYzNWJiMzQzNzg4NCIsImV4cCI6MTQ1MDI0ODQ4NCwiaWF0IjoxNDQ3NjU2NDgzLCJ1YWN0IjoiVERndmNsRmtTOGVndzNZVENqZks1eFZ1b253aVVVIiwibmJmIjoxNDQ3NjU2MzkzLCJhdWQiOiJvdHhiV3ZaM1J5OHNaalgtMzFCZE42cjF2ZXA1MTIxUyJ9.Y4EvbGQEDGqgOJm4jWWQu6jJAER80d9nz7yA0q0_VuB8poJD-_Of_Xaim0RGw6mtolqJe2zh-tKpHyDx861UK-oFSeocB1wpuQ7yRdhXEqH0AtepGIRTVl4gXPmKd3_bv9ERWJ0BymGY4xM64iT0jJjI2kAcY0YXixhgblabJzUGWTn9GhHR7PFvqxq1sUR3yDWJEFYyvJ2bbbHRWoOX-uSUmeZ53RKshTnPYwK_NVIl30LKx9kUxHVqxamsYCw4XhWkKTwLDhDMOzb70X8o4lgC1ceZMyAAvAevScIxrWJsQA1eorqojsEOnOeRxJdMxNy9Yx0942z1Dm3qdlMWZg"

	fmt.Println("URL formed: " + reqURL)
	fmt.Println("body is" + string(b))

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(b))
    req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	var res UberETA
	if err != nil {
		fmt.Println("error in sending ride req to Uber: ", err)
		return eta, errMsg
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("error in reading response: Might be invalid coordinates:", err)
		return eta, errMsg
	}

	err = json.Unmarshal(body, &res)
	fmt.Println("resp is " , string(body));

	if err != nil {
		fmt.Println("error in unmashalling response: ", err)
		fmt.Println("resp is " , res);
		return eta, errMsg
	}

	if len(res.Errors) > 0 {
		errMsg = res.Errors[0].Title
	}

	eta = res.ETA
	return eta, errMsg
}

func getUberProductId(lat, long string) string {

	client := &http.Client{}

	reqURL := "https://sandbox-api.uber.com/v1/products?latitude="
	reqURL += url.QueryEscape(lat)

	reqURL += "&longitude="
	reqURL += url.QueryEscape(long)

	reqURL += "&access_token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsicmVxdWVzdCJdLCJzdWIiOiJiMWQxODcyYS05NGE0LTQ1MDQtYWU5Yi00OGM1Y2FhZGMyNDciLCJpc3MiOiJ1YmVyLXVzMSIsImp0aSI6ImYxNTUzNDc1LWEyZDEtNDU2Yy04ZDU0LWYzNWJiMzQzNzg4NCIsImV4cCI6MTQ1MDI0ODQ4NCwiaWF0IjoxNDQ3NjU2NDgzLCJ1YWN0IjoiVERndmNsRmtTOGVndzNZVENqZks1eFZ1b253aVVVIiwibmJmIjoxNDQ3NjU2MzkzLCJhdWQiOiJvdHhiV3ZaM1J5OHNaalgtMzFCZE42cjF2ZXA1MTIxUyJ9.Y4EvbGQEDGqgOJm4jWWQu6jJAER80d9nz7yA0q0_VuB8poJD-_Of_Xaim0RGw6mtolqJe2zh-tKpHyDx861UK-oFSeocB1wpuQ7yRdhXEqH0AtepGIRTVl4gXPmKd3_bv9ERWJ0BymGY4xM64iT0jJjI2kAcY0YXixhgblabJzUGWTn9GhHR7PFvqxq1sUR3yDWJEFYyvJ2bbbHRWoOX-uSUmeZ53RKshTnPYwK_NVIl30LKx9kUxHVqxamsYCw4XhWkKTwLDhDMOzb70X8o4lgC1ceZMyAAvAevScIxrWJsQA1eorqojsEOnOeRxJdMxNy9Yx0942z1Dm3qdlMWZg"

	fmt.Println("URL formed: " + reqURL)

	req, err := http.NewRequest("GET", reqURL, nil)
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("error in sending ride req to Uber: ", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("error in reading response: ", err)
		return ""
	}

	var res UberProducts
	err = json.Unmarshal(body, &res)
//	fmt.Println("resp is " , res);

	if err != nil {
		fmt.Println("error in unmashalling response: ", err)
		return ""
	}

	var product_id string
	if len(res.Products) == 0 {
		fmt.Println("Empty response from Uber. Not a valid address");
	} else {
		add := res.Products
		for i := 0; i < len(add); i++ {
			if add[i].Display_name == "uberX" {
				product_id = add[i].Product_id
			}
		}
	}
	
	fmt.Println("product_id:",product_id)
	return product_id
}








