package main

import (
	"encoding/json"
	"net/http"
		"github.com/go-chi/chi/v5"
)



type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content any `json:"content,omitempty"`
	ID      int    `json:"id,omitempty"`
}

type User struct {
	Id int `json:"id ,omitempty"`
	Name string `json:"name"`
	Created_at string `json:"created_at,omitempty"`
}

type UserInfo struct {
	Id int `json:"id ,omitempty"`
	UserId int `json:"user_id"`
	Longitude float64 `json:"longitude"`
	Latitude float64 `json:"latitude"`
	Speed float64 `json:"speed"`
	Created_at string `json:"created_at,omitempty"`
}

type PastlocationsBody struct {
	UserId int `json:"user_id"`
	StartTime string `json:"start_time"`
	EndTime string `json:"end_time"`
}


func (app *application) RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var response jsonResponse

	response.OK = true
	response.Message = "server is up !!!!"

	json.NewEncoder(w).Encode(response)
}

func (app *application) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var response jsonResponse
	// get request body
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response.OK = false
		response.Message = "error decoding request body"
		json.NewEncoder(w).Encode(response)
		return
	}
	createUserSql := `INSERT INTO users (name) VALUES ($1) RETURNING id`
	err = app.db.QueryRow(createUserSql, user.Name).Scan(&user.Id)
	if err != nil {
		response.OK = false
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	response.OK = true
	response.Message = "user created successfully"
	response.ID = user.Id
	json.NewEncoder(w).Encode(response)


	
}

func (app *application) AddUserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var response jsonResponse
	// get request body
	var userInfo UserInfo
	err := json.NewDecoder(r.Body).Decode(&userInfo)
	if err != nil {
		response.OK = false
		response.Message = "error decoding request body"
		json.NewEncoder(w).Encode(response)
		return
	}
	createUserInfoSql := `INSERT INTO user_info (user_id,longitude,latitude,speed) VALUES ($1,$2,$3,$4) RETURNING id`
	err = app.db.QueryRow(createUserInfoSql, userInfo.UserId,userInfo.Longitude,userInfo.Latitude,userInfo.Speed).Scan(&userInfo.Id)
	if err != nil {
		response.OK = false
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	response.OK = true
	response.Message = "user info added successfully"
	response.ID = userInfo.Id
	json.NewEncoder(w).Encode(response)
}

// GetLastLocation
func (app *application) GetLastLocation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// get id from url
	id := chi.URLParam(r, "id")
	var response jsonResponse

	var userInfo UserInfo
	getLastLocationSql := `SELECT * FROM user_info WHERE user_id = $1 ORDER BY id DESC LIMIT 1`
	err := app.db.QueryRow(getLastLocationSql, id).Scan(&userInfo.Id,&userInfo.UserId,&userInfo.Longitude,&userInfo.Latitude,&userInfo.Speed,&userInfo.Created_at)
	if err != nil {
		response.OK = false
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	response.OK = true
	response.Message = "user info fetched successfully"
	response.Content = userInfo
	json.NewEncoder(w).Encode(response)


	
}

// GetLastLocation between two dates
func (app *application) GetPastLocations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var pastlocationsBody PastlocationsBody
	// get id request body
	json.NewDecoder(r.Body).Decode(&pastlocationsBody)
	var response jsonResponse

	var userInfos []UserInfo
	getLastLocationSql := `SELECT * FROM user_info WHERE user_id = $1 AND created_at BETWEEN $2 AND $3`
	rows, err := app.db.Query(getLastLocationSql, pastlocationsBody.UserId, pastlocationsBody.StartTime, pastlocationsBody.EndTime)
	if err != nil {
		response.OK = false
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var userInfo UserInfo
		err := rows.Scan(&userInfo.Id,&userInfo.UserId,&userInfo.Longitude,&userInfo.Latitude,&userInfo.Speed,&userInfo.Created_at)
		if err != nil {
			response.OK = false
			response.Message = err.Error()
			json.NewEncoder(w).Encode(response)
			return
		}
		userInfos = append(userInfos, userInfo)
	}
	response.OK = true
	response.Message = "user info fetched successfully"
	response.Content = userInfos
	json.NewEncoder(w).Encode(response)




}