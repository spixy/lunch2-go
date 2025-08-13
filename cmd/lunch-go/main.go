package main

import (
	"fmt"
	"net/http"

	r "git.zvon.tech/zv0n/lunch-go/pkg/restaurants"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var restaurants []r.RestaurantInterface

func getAll(c *gin.Context) {
	c.JSON(http.StatusOK, restaurants)
}

func dayToIndex(day string) int {
	if day == "monday" {
		return 0
	} else if day == "tuesday" {
		return 1
	} else if day == "wednesday" {
		return 2
	} else if day == "thursday" {
		return 3
	} else if day == "friday" {
		return 4
	} else if day == "saturday" {
		return 5
	}
	return 6
}

func get(c *gin.Context) {
	day := dayToIndex(c.Query("day"))
	var responseObjects []r.RestaurantJSON
	for _, restaurant := range restaurants {
		obj := restaurant.GetSpecificDayObject([]int{day})
		responseObjects = append(responseObjects, obj)
	}
	c.JSON(http.StatusOK, responseObjects)
}

func refreshInternal() {
	for _, restaurant := range restaurants {
		restaurant.Parse()
	}
}

func refresh(c *gin.Context) {
	refreshInternal()
	c.Status(http.StatusOK)
}

func main() {
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/2787-tusto-titanium.html", "Tusto", 1))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/2743-restaurant-padowetz.html", "Padowetz", 2))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/7220-gourmet-u-vankovky.html", "Gourmet", 3))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/3870-pizzerie-basilico.html", "Basilico", 4))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/5876-restaurace-u-emila.html", "U Emila", 5))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/5369-restaurace-bogota.html", "Bogota", 6))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/2687-kometa-arena-pub.html", "Kometa Arena Pub", 7))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/3160-restaurace-na-tahu-.html", "Na-tahu", 8))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/3165-potrefena-husa-zelny-trh.html", "Potrefená Husa", 9))
	restaurants = append(restaurants, r.NewStaticRestaurant("https://himalayarestaurace.cz/denni-menu", "Himalaya", 10))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/4116-padagali.html", "Padagali", 11))

	refreshInternal()
	fmt.Println("Initial parsing finished")

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/get_all", getAll)
	router.GET("/get", get)
	router.GET("/refresh", refresh)
	router.Run("0.0.0.0:8080")
}
