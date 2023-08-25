package restaurants

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type TaoRestaurant struct {
	Restaurant
}

func NewTaoRestaurant(url string, name string, id int) *TaoRestaurant {
	restaurant := new(TaoRestaurant)
	restaurant.SetDefaultValues()
	restaurant.id = id
	restaurant.url = url
	restaurant.name = name
	return restaurant
}

func (restaurant *TaoRestaurant) Parse() {
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

	// first occurence of tydenni-menu-div
	daily, err := findNodeByClass(doc, "ct-div-block tydenni-menu-div")
	if err != nil {
		fmt.Printf("Couldn't find content for restaurant \"%s\"\n", restaurant.name)
		return
	}
	for menu := daily; menu != nil; menu = menu.NextSibling {
		nameNode, err := findNodeByClass(menu, "ct-span")
		if err != nil {
			continue
		}
		nameText, err := getText(nameNode)
		if err != nil {
			continue
		}

		textElements := strings.Split(nameText, "..")
		meal := textElements[0]
		priceText := textElements[len(textElements)-1]
		priceNum := -1
		soup := false
		if len(priceText) < 2 {
			soup = true
		} else {
			priceNum, err = strconv.Atoi(strings.TrimLeft(strings.TrimSpace(strings.Split(priceText, "k")[0]), "."))
			if err != nil {
				priceNum = -1
			}
		}
		restaurant.AddPermanent(soup, strings.TrimSpace(meal), "", priceNum)
	}

	special := daily.Parent.NextSibling.NextSibling
	for i := 0; i < 5; i++ {
		nameNode, err := findNodeByClass(special, "ct-span")
		if err != nil {
			continue
		}
		nameText, err := getText(nameNode)
		if err != nil {
			continue
		}

		textElements := strings.Split(nameText, "..")
		meal := textElements[0]
		priceText := textElements[len(textElements)-1]
		priceNum := -1
		soup := false
		if len(priceText) == 0 {
			soup = true
		} else {
			priceNum, err = strconv.Atoi(strings.TrimLeft(strings.TrimSpace(strings.Split(priceText, "k")[0]), "."))
			if err != nil {
				priceNum = -1
			}
		}
		restaurant.menus[i].Add(soup, strings.TrimSpace(meal), "", priceNum)
		special = special.NextSibling
	}
	restaurant.menus[0].SetDay("Monday")
	restaurant.menus[1].SetDay("Tuesday")
	restaurant.menus[2].SetDay("Wednesday")
	restaurant.menus[3].SetDay("Thursday")
	restaurant.menus[4].SetDay("Friday")
	restaurant.menus[5].SetDay("Saturday")
	restaurant.menus[6].SetDay("Sunday")
}
