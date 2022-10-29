package main

import (
	"fmt"
	"net/http"

	r "git.zvon.tech/zv0n/lunch-go/pkg/restaurants"
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
	fmt.Println("Hello, World!")
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/2752-u-drevaka-beergrill.html", "U Dřeváka"))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/4116-padagali.html", "Padagali"))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/5448-light-of-india.html", "Light of India"))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/2609-pizzeria-al-capone.html", "Al Capone"))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/3830-suzies-steak-pub.html", "Suzie's"))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/6468-diva-bara.html", "Divá Bára"))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/6695-u-karla.html", "U Karla"))
	restaurants = append(restaurants, r.NewFreshRestaurant("http://www.fresh-menu.cz/", "Fresh"))

	refreshInternal()
	fmt.Println("Initial parsing finished")

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/get_all", getAll)
	router.GET("/get", get)
	router.GET("/refresh", refresh)
	router.Run("localhost:8080")
	/*
	   test := restaurants.MakeMenickaRestaurant("https://www.menicka.cz/2752-u-drevaka-beergrill.html", "U Dřeváka")
	   test.Parse()

	   	for _, menu := range test.GetMenus() {
	   		for _, meal := range menu.GetMeals() {
	   			if meal.IsSoup() {
	   				fmt.Print("Soup: ")
	   			} else {
	   				fmt.Print("Meal: ")
	   			}
	   			fmt.Println(meal.GetName(), "; ", meal.GetPrice(), "Kč")
	   		}
	   	}
	*/
}
