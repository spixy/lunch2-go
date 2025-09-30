package restaurants

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type MealProcessor = func(mealName string) string

type MenickaRestaurant struct {
	Restaurant
	soups         int
	mealProcessor MealProcessor
}

func NewMenickaRestaurant(url string, name string, id int) *MenickaRestaurant {
	restaurant := new(MenickaRestaurant)
	restaurant.SetDefaultValues()
	restaurant.id = id
	restaurant.url = url
	restaurant.name = name
	restaurant.soups = 0
	restaurant.mealProcessor = func(mealName string) string { return mealName }
	return restaurant
}

func NewPadowetzRestaurant(url string, name string, id int) *MenickaRestaurant {
	restaurant := NewMenickaRestaurant(url, name, id)
	restaurant.soups = 2
	return restaurant
}

func NewBogotaRestaurant(url string, name string, id int) *MenickaRestaurant {
	restaurant := NewMenickaRestaurant(url, name, id)
	restaurant.mealProcessor = func(mealName string) string { return strings.ReplaceAll(mealName, "()", "") }
	return restaurant
}

func dayToIndex(day string) (int, error) {
	if day == "Pondělí" {
		return 0, nil
	} else if day == "Úterý" {
		return 1, nil
	} else if day == "Středa" {
		return 2, nil
	} else if day == "Čtvrtek" {
		return 3, nil
	} else if day == "Pátek" {
		return 4, nil
	} else if day == "Sobota" {
		return 5, nil
	} else if day == "Neděle" {
		return 6, nil
	}
	return -1, errors.New("couldn't parse the day")
}

func (restaurant *MenickaRestaurant) Parse() {
	restaurant.clearMenus()
	resp, err := http.Get(restaurant.url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	content, err := findNodeByClass(doc, "obsah")
	if err != nil {
		fmt.Printf("Couldn't find content for restaurant \"%s\"\n", restaurant.name)
		return
	}
	for menu := content.FirstChild; menu != nil; menu = menu.NextSibling {
		if hasKeyValue(menu, "class", "menicka") {
			day, err := findNodeByClass(menu, "nadpis")
			if err != nil {
				continue
			}
			meals, err := findNodeByClass(menu, "popup-gallery")
			if err != nil {
				continue
			}

			dayText, err := getTextDecodeWindows1250(day)
			if err != nil {
				continue
			}
			dayText = strings.Split(dayText, " ")[0]
			dayIndex, err := dayToIndex(dayText)
			if err != nil {
				continue
			}

			var mealIndex = 0
			for meal := meals.FirstChild; meal != nil; meal = meal.NextSibling {
				nameNode, err := findNodeByClass(meal, "polozka")
				if err != nil {
					continue
				}
				name, err := getTextDecodeWindows1250(nameNode)
				if err != nil {
					continue
				}
				price := -1
				priceNode, err := findNodeByClass(meal, "cena")
				if err == nil {
					priceStr, err := getText(priceNode)
					if err != nil {
						continue
					}
					price, err = strconv.Atoi(strings.Split(priceStr, " ")[0])
					if err != nil {
						price = -1
					}
				}

				var isSoup = hasKeyValue(meal, "class", "polevka") || mealIndex < restaurant.soups
				restaurant.menus[dayIndex].Add(isSoup, strings.TrimSpace(restaurant.mealProcessor(name)), "", price)
				mealIndex++
			}
		}
	}
	restaurant.menus[0].SetDay("Monday")
	restaurant.menus[1].SetDay("Tuesday")
	restaurant.menus[2].SetDay("Wednesday")
	restaurant.menus[3].SetDay("Thursday")
	restaurant.menus[4].SetDay("Friday")
	restaurant.menus[5].SetDay("Saturday")
	restaurant.menus[6].SetDay("Sunday")
}
