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
	//restaurants = append(restaurants, r.NewFreshRestaurant("http://www.fresh-menu.cz/", "Fresh", 1))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/2752-u-drevaka-beergrill.html", "U Dřeváka", 2))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/4116-padagali.html", "Padagali", 3))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/5448-light-of-india.html", "Light of India", 4))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/2609-pizzeria-al-capone.html", "Al Capone", 5))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/3830-suzies-steak-pub.html", "Suzie's", 6))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/6468-diva-bara.html", "Divá Bára", 7))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/6695-u-karla.html", "U Karla", 8))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/7470-the-immigrant-.html", "The Immigrant", 9))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/2721-selepka.html", "Šelepka", 10))
	restaurants = append(restaurants, r.NewTaoRestaurant("https://www.taorestaurant.cz/tydenni_menu/nabidka/", "Tao", 11))
	restaurants = append(restaurants, r.NewMenickaRestaurant("https://www.menicka.cz/3854-na-ruzku.html", "Na Růžku", 12))

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
