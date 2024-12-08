package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WeatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

func getWeather(lat, lon, apiKey string) (WeatherResponse, error) {
	var weatherResponse WeatherResponse
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&units=metric&appid=%s", lat, lon, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return weatherResponse, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return weatherResponse, fmt.Errorf("error: status code %d, message: %s", resp.StatusCode, resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&weatherResponse)
	if err != nil {
		return weatherResponse, fmt.Errorf("error decoding weather data: %v", err)
	}

	return weatherResponse, nil
}

func main() {
	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/weather", func(c *gin.Context) {
		lat := c.Query("lat")
		lon := c.Query("lon")
		apiKey := ""

		if lat == "" || lon == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "latitude and longitude are required"})
			return
		}

		weather, err := getWeather(lat, lon, apiKey)
		if err != nil {
			fmt.Println("Error fetching weather:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"location":    weather.Name,
			"temperature": weather.Main.Temp,
			"description": weather.Weather[0].Description,
		})
	})

	router.Run(":8080")
}
