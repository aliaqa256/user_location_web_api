package controllers

import (
	"fmt"
	"os"
	"time"

	"github.com/aliaqa256/user_location_web_api/cache"
	"github.com/aliaqa256/user_location_web_api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)




func RootHandler(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"ok":      true,
		"message": "server is up !!!!",
	})
}

func CreateUserHandler(c *fiber.Ctx) error {
	modelApp := models.NewApplication()


	var payload struct {
		Id int `json:"id ,omitempty"`
		Name string `json:"name,omitempty"`
	}
	err := c.BodyParser(&payload)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"ok":      false,
			"message": "error decoding request body",
		})
	}
	createUserSql := `INSERT INTO users (name) VALUES ($1) RETURNING id`
	err = modelApp.DB.QueryRow(createUserSql, payload.Name).Scan(&payload.Id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"ok":      false,
			"message": "error creating user",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"ok":      true,
		"message": "user created successfully",
		"id": payload.Id,
	})

}


func AddUserInfoHandler(c *fiber.Ctx) error {
	modelApp := models.NewApplication()
	var payload  struct {
			Id int `json:"id ,omitempty"`
	UserId int `json:"user_id"`
	Longitude float64 `json:"longitude"`
	Latitude float64 `json:"latitude"`
	Speed float64 `json:"speed"`
	Created_at string `json:"created_at,omitempty"`
	}
	err := c.BodyParser(&payload)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"ok":      false,
			"message": "error decoding request body",
		})
	}
	createUserInfoSql := `INSERT INTO user_info (user_id,longitude,latitude,speed) VALUES ($1,$2,$3,$4) RETURNING id`
	err = modelApp.DB.QueryRow(createUserInfoSql, payload.UserId,payload.Longitude,payload.Latitude,payload.Speed).Scan(&payload.Id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"ok":      false,
			"message": "error creating user info",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"ok":      true,
		"message": "user info created successfully",
		"id": payload.Id,
	})

}

func GetLastLocationHandler(c *fiber.Ctx) error {
	modelApp := models.NewApplication()
	var payload  struct {
		Id int `json:"id ,omitempty"`
		UserId int `json:"user_id"`
		Longitude float64 `json:"longitude"`
		Latitude float64 `json:"latitude"`
		Speed float64 `json:"speed"`
		Created_at string `json:"created_at,omitempty"`
	}
	id:=c.Params("id")
	// chck if data is in cache if yes return it if no get it from db and store it in cache
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	redisHost := os.Getenv("REDIS_URI")
	cacheApp:=cache.NewRedisCache(redisHost, 0, 10*time.Second)
	data,isany:=cacheApp.Get(id)
	if isany {
		fmt.Println("data is in cache")
		return c.Status(200).JSON(fiber.Map{
			"ok":      true,
			"message": "user info retrieved successfully",
			"data": data,
		})
	}else{
		fmt.Println("data is not in cache")
		getLastLocationSql := `SELECT * FROM user_info WHERE user_id = $1 ORDER BY id DESC LIMIT 1`
		err := modelApp.DB.QueryRow(getLastLocationSql, id).Scan(&payload.Id,&payload.UserId,&payload.Longitude,&payload.Latitude,&payload.Speed,&payload.Created_at)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"ok":      false,
				"message": "error getting last location",
			})
		}
		// store data in cache
		cacheApp.Set(id,payload)
		return c.Status(200).JSON(fiber.Map{
			"ok":      true,
			"message": "last location retrieved successfully",
			"payload": payload,
		})
	}


	
}

func GetPastLocationsHandler(c *fiber.Ctx) error {
	modelApp := models.NewApplication()
	var payload  struct {
	UserId int `json:"user_id"`
	StartTime string `json:"start_time"`
	EndTime string `json:"end_time"`
	}

	type Info struct {
		Id int `json:"id ,omitempty"`
		UserId int `json:"user_id"`
		Longitude float64 `json:"longitude"`
		Latitude float64 `json:"latitude"`
		Speed float64 `json:"speed"`
		Created_at string `json:"created_at,omitempty"`
	}

	var info = Info{}
	var infoArray [] Info
	err := c.BodyParser(&payload)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"ok":      false,
			"message": "error decoding request body",
		})
	}
	err = godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	redisHost := os.Getenv("REDIS_URI")
	cacheApp:=cache.NewRedisCache(redisHost, 0, 10*time.Second)
	data,isany:=cacheApp.Get( fmt.Sprintf("%v_past",payload.UserId))
	if isany {
		fmt.Println("data is in cache")
		return c.Status(200).JSON(fiber.Map{
			"ok":      true,
			"message": "user info retrieved successfully",
			"payload": data,
		})
	}else{
		fmt.Println("data is not in cache")
	getPastLocationsSql := `SELECT * FROM user_info WHERE user_id = $1 AND created_at BETWEEN $2 AND $3`
	rows, err := modelApp.DB.Query(getPastLocationsSql, payload.UserId,payload.StartTime,payload.EndTime)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"ok":      false,
			"payload": "error getting past locations",
		})
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&info.Id,&info.UserId,&info.Longitude,&info.Latitude,&info.Speed,&info.Created_at)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"ok":      false,
				"message": "error getting past locations",
			})
		}
		infoArray = append(infoArray, info)
	}

	// store data in cache
	cacheApp.Set(fmt.Sprintf("%v_past",payload.UserId),infoArray)
	return c.Status(200).JSON(fiber.Map{
		"ok":      true,
		"message": "past locations retrieved successfully",
		"payload": infoArray,
	})
}
}