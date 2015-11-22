# cmpe273-assignment3


##Setup

go get github.com/aggarwalsomya/cmpe273-assignment3/server

cd src/github.com/aggarwalsomya/cmpe273-assignment3/server

go run *



## How to use


These are the sample location ids which were given in quiz 

1. Fairmont Hotel San Francisco : 4096
2. Golden Gate Bridge : 4343
3. Pier 39 : 7030
4. Golden Gate Park : 8624
5. Twin Peaks : 3969


## Planning a trip

This will create a new trip:

curl -v -H "Content-Type: application/json" -X POST -d '{"starting_from_location_id":"4096", "location_ids":["4343", "7030", "8624", "3969"]}' http://localhost:8080/trips

###Response

{"id":"878805","status":"planning","starting_from_location_id":"4096","best_route_location_ids":["3969","4343","8624","7030"],"total_uber_costs":68,"total_uber_duration":4649,"total_distance":24.4}



## Getting details of the trip

curl -v -H "Content-Type: application/json" GET  http://localhost:8080/trips/878805


###Response

{"id":"878805","status":"planning","starting_from_location_id":"4096","best_route_location_ids":["3969","4343","8624","7030"],"total_uber_costs":68,"total_uber_duration":4649,"total_distance":24.4}




## Starting the trip


**1st Call:**

curl -v -X  PUT http://localhost:8080/trips/878805/request


**Response**

{"id":"878805","status":"requesting","starting_from_location_id":"4096","best_route_location_ids":["3969","4343","8624","7030"],"total_uber_costs":68,"total_uber_duration":4649,"total_distance":24.4,"next_destination_location_id":"3969","uber_wait_time_eta":4}


**2nd Call:**

curl -v -X  PUT http://localhost:8080/trips/878805/request


**Response**

{"id":"878805","status":"requesting","starting_from_location_id":"4096","best_route_location_ids":["3969","4343","8624","7030"],"total_uber_costs":68,"total_uber_duration":4649,"total_distance":24.4,"next_destination_location_id":"4343","uber_wait_time_eta":4}



**3rd Call:**

curl -v -X  PUT http://localhost:8080/trips/878805/request


**Response**

{"id":"878805","status":"requesting","starting_from_location_id":"4096","best_route_location_ids":["3969","4343","8624","7030"],"total_uber_costs":68,"total_uber_duration":4649,"total_distance":24.4,"next_destination_location_id":"8624","uber_wait_time_eta":4}



**4th Call:**

curl -v -X  PUT http://localhost:8080/trips/878805/request


**Response**

{"id":"878805","status":"requesting","starting_from_location_id":"4096","best_route_location_ids":["3969","4343","8624","7030"],"total_uber_costs":68,"total_uber_duration":4649,"total_distance":24.4,"next_destination_location_id":"7030","uber_wait_time_eta":4}



**5th  Call:**

curl -v -X  PUT http://localhost:8080/trips/878805/request


**Response**


{"id":"878805","status":"completed","starting_from_location_id":"4096","best_route_location_ids":["3969","4343","8624","7030"],"total_uber_costs":68,"total_uber_duration":4649,"total_distance":24.4,"next_destination_location_id":"4096","uber_wait_time_eta":4}


After this, PUT call will not change the status further and trip remain completed. 

