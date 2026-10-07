package restaurants

import (
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/net/html"
)

type PokharaRestaurant struct {
	Restaurant
	dataUrl string
}

var pokharaDays = map[string]int{"mon": 0, "tue": 1, "wed": 2, "thu": 3, "fri": 4, "sat": 5, "sun": 6}

// NewPokharaRestaurant creates the restaurant with url as the public weekly
// menu page and dataUrl as the document the page loads its menu from via JS.
func NewPokharaRestaurant(url string, dataUrl string, name string, id int) *PokharaRestaurant {
	restaurant := new(PokharaRestaurant)
	restaurant.SetDefaultValues()
	restaurant.id = id
	restaurant.url = url
	restaurant.dataUrl = dataUrl
	restaurant.name = name
	return restaurant
}

// pokharaChildText returns the normalized text of the first descendant with
// the given class, or "" when there is none.
func pokharaChildText(node *html.Node, class string) string {
	child, err := findNodeByClass(node, class)
	if err != nil {
		return ""
	}
	text, err := getText(child)
	if err != nil {
		return ""
	}
	return normalizeWhitespace(text)
}

func (restaurant *PokharaRestaurant) parseMeals(node *html.Node, menu *Menu, isSoup bool) {
	if hasKeyValue(node, "class", "meal-card") {
		name := pokharaChildText(node, "meal-name")
		if name == "" {
			return
		}
		price, err := strconv.Atoi(pokharaChildText(node, "price-value"))
		if err != nil {
			price = -1
		}
		menu.Add(isSoup, name, pokharaChildText(node, "meal-description"), price)
		return
	}
	for n := node.FirstChild; n != nil; n = n.NextSibling {
		restaurant.parseMeals(n, menu, isSoup)
	}
}

func (restaurant *PokharaRestaurant) Parse() {
	restaurant.clearMenus()
	restaurant.clearPermanentMenus()
	resp, err := http.Get(restaurant.dataUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// <div class="weekly-panels-container"> holds one
	// <div mc-data="mon" class="weekly-day-panel"> per day, each split into
	// <div class="meal-group"> sections (soups, main courses) of meal cards.
	panels, err := findNodeByClass(doc, "weekly-panels-container")
	if err != nil {
		fmt.Printf("Couldn't find content for restaurant \"%s\"\n", restaurant.name)
		return
	}

	for panel := panels.FirstChild; panel != nil; panel = panel.NextSibling {
		if !hasKeyValue(panel, "class", "weekly-day-panel") {
			continue
		}
		day, err := getAttribute(panel, "mc-data")
		if err != nil {
			continue
		}
		dayIndex, ok := pokharaDays[day]
		if !ok {
			continue
		}

		groups, err := findNodeByClass(panel, "meal-groups-grid")
		if err != nil {
			continue
		}
		for group := groups.FirstChild; group != nil; group = group.NextSibling {
			if !hasKeyValue(group, "class", "meal-group") {
				continue
			}
			isSoup := false
			if title, err := findNodeByClass(group, "group-title"); err == nil {
				isSoup = hasKeyValue(title, "mc-text", "GLOBAL_WEEKLY_MENU_TYPE_SOUP")
			}
			restaurant.parseMeals(group, &restaurant.menus[dayIndex], isSoup)
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
