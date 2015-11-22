# cmpe273-assignment3


##Setup

go get github.com/aggarwalsomya/cmpe273-assignment3/server

cd src/github.com/aggarwalsomya/cmpe273-assignment3/server

go run *



## How to use


These are the sample location ids which are already created around my residence 

1. home: 6160
2. costco: 5103
3. pizza: 4404
4. walmart: 8926
5. fresh & easy: 6886


## Planning a trip

This will create a new trip:

curl -v -H "Content-Type: application/json" -X POST -d '{"starting_from_location_id":"6160", "location_ids":["4404", "5103", "6886", "8926"]}' http://localhost:8080/trips

###Response

{"id":"838781","status":"Planning","starting_from_location_id":"6160","best_route_location_ids":["5103","4404","8926","6886"],"total_uber_costs":30,"total_uber_duration":1955,"total_distance":10.13}



## Getting details of the trip

curl -v -H "Content-Type: application/json" GET  http://localhost:8080/trips/838781


###Response

{"id":"838781","status":"Planning","starting_from_location_id":"6160","best_route_location_ids":["5103","4404","8926","6886"],"total_uber_costs":30,"total_uber_duration":1955,"total_distance":10.13}




## Starting the trip


**1st Call:**

curl -v -X  PUT http://localhost:8080/trips/838781/request


**Response**

{"status":"requesting","starting_from_location_id":"6160","best_route_location_ids":["5103","4404","8926","6886"],"total_uber_costs":30,"total_uber_duration":1955,"total_distance":10.13,"next_destination_location_id":"5103","uber_wait_time_eta":12,"id":"838781"}


**2nd Call:**

curl -v -X  PUT http://localhost:8080/trips/838781/request


**Response**

{"id":"838781","status":"requesting","starting_from_location_id":"6160","best_route_location_ids":["5103","4404","8926","6886"],"total_uber_costs":30,"total_uber_duration":1955,"total_distance":10.13,"next_destination_location_id":"4404","uber_wait_time_eta":12}



**3rd Call:**

curl -v -X  PUT http://localhost:8080/trips/838781/request


**Response**

{"id":"838781","status":"requesting","starting_from_location_id":"6160","best_route_location_ids":["5103","4404","8926","6886"],"total_uber_costs":30,"total_uber_duration":1955,"total_distance":10.13,"next_destination_location_id":"8926","uber_wait_time_eta":12}



**4th Call:**

curl -v -X  PUT http://localhost:8080/trips/838781/request


**Response**

{"id":"838781","status":"requesting","starting_from_location_id":"6160","best_route_location_ids":["5103","4404","8926","6886"],"total_uber_costs":30,"total_uber_duration":1955,"total_distance":10.13,"next_destination_location_id":"6886","uber_wait_time_eta":12}



**5th  Call:**

curl -v -X  PUT http://localhost:8080/trips/838781/request


**Response**


{"id":"838781","status":"requesting","starting_from_location_id":"6160","best_route_location_ids":["5103","4404","8926","6886"],"total_uber_costs":30,"total_uber_duration":1955,"total_distance":10.13,"next_destination_location_id":"6160","uber_wait_time_eta":12}


Once the user has reached back to his home, the status  will be changed to *completed*. After that, no more status will be changed. 


**6th call**


curl -v -X  PUT http://localhost:8080/trips/838781/request


**Response**

{"id":"838781","status":"completed","starting_from_location_id":"6160","best_route_location_ids":["5103","4404","8926","6886"],"total_uber_costs":30,"total_uber_duration":1955,"total_distance":10.13,"next_destination_location_id":"","uber_wait_time_eta":0}


