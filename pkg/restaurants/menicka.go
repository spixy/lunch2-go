package restaurants

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type MenickaRestaurant struct {
	Restaurant
}

func MakeMenickaRestaurant(url string, name string) MenickaRestaurant {
	restaurant := MenickaRestaurant{}
	restaurant.url = url
	restaurant.name = name
	return restaurant
}

func NewMenickaRestaurant(url string, name string) *MenickaRestaurant {
	restaurant := new(MenickaRestaurant)
	restaurant.url = url
	restaurant.name = name
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
		if hasClass(menu, "menicka") {
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
				if hasClass(meal, "polevka") {
					restaurant.menus[dayIndex].Add(true, strings.TrimSpace(name), "", price)
				} else {
					restaurant.menus[dayIndex].Add(false, strings.TrimSpace(name), "", price)
				}
			}
		}
	}
}
