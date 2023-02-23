package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type NominatimResponse struct {
	Address struct {
		City     string `json:"city"`
		Road     string `json:"road"`
		Postcode string `json:"postcode"`
		State    string `json:"state"`
		Country  string `json:"country"`
	} `json:"address"`
}

func main() {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, error=${error}\n",
	}))
	// Endpoint to retrieve the user's location and insert it into the database
	e.GET("/location", func(c echo.Context) error {
		// Retrieve the user's location
		latitude, _ := strconv.ParseFloat(c.QueryParam("latitude"), 64)
		longitude, _ := strconv.ParseFloat(c.QueryParam("longitude"), 64)
		// city := c.QueryParam("city")

		// Insert the location into the database
		// location := Location{City: city, Latitude: latitude, Longitude: longitude}

		url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f&zoom=18&addressdetails=1", latitude, longitude)
		response, err := http.Get(url)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "error"})
		}

		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "error"})
		}

		var nominatimResponse NominatimResponse
		err = json.Unmarshal(body, &nominatimResponse)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "error"})
		}

		city := nominatimResponse.Address.City
		street := nominatimResponse.Address.Road
		postcode := nominatimResponse.Address.Postcode
		state := nominatimResponse.Address.State
		country := nominatimResponse.Address.Country
		urlLocation := fmt.Sprintf("https://www.openstreetmap.org/#map=19/%f/%f", latitude, longitude)
		// Return the inserted location
		log.Println("Road : ", street)
		log.Println("City : ", city)
		log.Println("Postal Code:", postcode)
		log.Println("State:", state)
		log.Println("Country:", country)
		log.Println("Latitude: ", latitude)
		log.Println("Longitude: ", longitude)
		log.Println("location : ", urlLocation)
		// Return the inserted location
		return c.JSON(http.StatusOK, map[string]interface{}{
			"Road":        street,
			"City":        city,
			"Postal Code": postcode,
			"State":       state,
			"Country":     country,
			"Latitude":    latitude,
			"Longitude":   longitude,
			"message":     "Succesfull Show Location",
		})
	})

	// Start the server
	if err := e.Start(":8000"); err != nil {
		log.Println(err.Error())
	}
}
