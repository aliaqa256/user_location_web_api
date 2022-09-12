# user_location_web_api

This is a simple web api that saves and returns the location of a user .

## apis
/ - root handler - returns "server is running"
/user - gets name and creates new user returns id
/user/info - gets user_id,lontitude,latitude,speed and saves it to db
/user/lastlocation/{id} - gets user id and returns last location
/user/pastlocations - gets user id,start time , endtime and returns all locations that user has been between 2 dates

postman collection json file url :http://alilotfidev.ir/files/user%20location.postman_collection.json