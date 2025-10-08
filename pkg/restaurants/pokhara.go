package restaurants

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type PokharaRestaurant struct {
	Restaurant
}

var pokharaDays = [6]string{"PONDELI", "UTERY", "STREDA", "CTVRTEK", "PÁTEK", "SOBOTA"}

func NewPokharaRestaurant(url string, name string, id int) *PokharaRestaurant {
	restaurant := new(PokharaRestaurant)
	restaurant.SetDefaultValues()
	restaurant.id = id
	restaurant.url = url
	restaurant.name = name
	return restaurant
}

func (restaurant *PokharaRestaurant) Parse() {
	restaurant.clearMenus()
	restaurant.clearPermanentMenus()
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

	// first occurence of #menu
	daily, err := findNodeByClass(doc, "col-lg-12")
	if err != nil {
		fmt.Printf("Couldn't find content for restaurant \"%s\"\n", restaurant.name)
		return
	}

	nextDay := 0
	isSoup := false
	for menu := daily.FirstChild; menu != nil; menu = menu.NextSibling {
		nameText, err := getText(menu)
		if err != nil || strings.TrimSpace(nameText) == "" || nameText == "\n" {
			continue
		}
		if nameText == pokharaDays[nextDay] {
			nextDay++
			isSoup = true
			continue
		} else if nextDay == 0 {
			continue
		}

		var textParts = strings.Split(nameText, " ")

		text := nameText
		if len(textParts) > 1 {
			text = strings.Join(textParts[:len(textParts)-1], " ")
		}

		var priceStr = strings.ReplaceAll(strings.ReplaceAll(textParts[len(textParts)-1], "kc", ""), "KC", "")
		price, err := strconv.Atoi(strings.Split(priceStr, " ")[0])
		if err != nil {
			price = -1
		}

		restaurant.menus[nextDay-1].Add(isSoup, text, "", price)
		isSoup = false
	}

	restaurant.menus[0].SetDay("Monday")
	restaurant.menus[1].SetDay("Tuesday")
	restaurant.menus[2].SetDay("Wednesday")
	restaurant.menus[3].SetDay("Thursday")
	restaurant.menus[4].SetDay("Friday")
	restaurant.menus[5].SetDay("Saturday")
	restaurant.menus[6].SetDay("Sunday")
}
