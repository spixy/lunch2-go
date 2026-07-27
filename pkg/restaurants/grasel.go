package restaurants

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type GraselRestaurant struct {
	Restaurant
}

func NewGraselRestaurant(url string, name string, id int) *GraselRestaurant {
	restaurant := new(GraselRestaurant)
	restaurant.SetDefaultValues()
	restaurant.id = id
	restaurant.url = url
	restaurant.name = name
	return restaurant
}

// graselPrice extracts the numeric price from strings like "139 Kč",
// "+19 Kč" or "Samostatně 49 Kč". Returns -1 when no digits are found.
func graselPrice(text string) int {
	var builder strings.Builder
	for _, r := range text {
		if r >= '0' && r <= '9' {
			builder.WriteRune(r)
		}
	}
	price, err := strconv.Atoi(builder.String())
	if err != nil {
		return -1
	}
	return price
}

func (restaurant *GraselRestaurant) Parse() {
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

	// The menu is a fixed weekday offering: the same twelve dishes are
	// listed as <article class="dish"> inside <div class="dish-grid">.
	grid, err := findNodeByClass(doc, "dish-grid")
	if err != nil {
		fmt.Printf("Couldn't find content for restaurant \"%s\"\n", restaurant.name)
		return
	}

	// Optional daily soup, priced as "+19 Kč" (the plus is stripped).
	if soupNode, err := findNodeByClass(doc, "menu-list__price"); err == nil {
		if text, err := getText(soupNode); err == nil {
			restaurant.AddPermanent(true, "Denní polévka", "", graselPrice(text))
		}
	}

	for dish := grid.FirstChild; dish != nil; dish = dish.NextSibling {
		if !hasKeyValue(dish, "class", "dish") {
			continue
		}

		nameNode, err := findNodeByClass(dish, "dish__name")
		if err != nil {
			continue
		}
		name, err := getText(nameNode)
		if err != nil {
			continue
		}
		name = normalizeWhitespace(name)

		desc := ""
		if descNode, err := findNodeByClass(dish, "dish__desc"); err == nil {
			if text, err := getText(descNode); err == nil {
				desc = normalizeWhitespace(text)
			}
		}

		price := -1
		if priceNode, err := findNodeByClass(dish, "dish__price"); err == nil {
			if text, err := getText(priceNode); err == nil {
				price = graselPrice(text)
			}
		}

		restaurant.AddPermanent(false, name, desc, price)
	}

	restaurant.menus[0].SetDay("Monday")
	restaurant.menus[1].SetDay("Tuesday")
	restaurant.menus[2].SetDay("Wednesday")
	restaurant.menus[3].SetDay("Thursday")
	restaurant.menus[4].SetDay("Friday")
	restaurant.menus[5].SetDay("Saturday")
	restaurant.menus[6].SetDay("Sunday")
}
