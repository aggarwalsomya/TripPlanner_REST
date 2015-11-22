package main

import (
    "net/http"
    "log"
		//"fmt"
)

func main() {
	router := NewRouter()
    log.Fatal(http.ListenAndServe(":8082", router))

//    product_id:= getUberProductId("37.3679232","-122.0032597")
//    _ = product_id
//    requestUberForRide("37.3679232","-122.0032597","38.368830","-120.036350",product_id)
}

