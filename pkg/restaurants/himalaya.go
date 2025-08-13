package restaurants

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

type HimalayaRestaurant struct {
	Restaurant
}

func NewHimalayaRestaurant(url string, name string, id int) *HimalayaRestaurant {
	restaurant := new(HimalayaRestaurant)
	restaurant.SetDefaultValues()
	restaurant.id = id
	restaurant.url = url
	restaurant.name = name
	return restaurant
}

func (restaurant *HimalayaRestaurant) Parse() {
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
	daily, err := findNodeById(doc, "menu")
	if err != nil {
		fmt.Printf("Couldn't find content for restaurant \"%s\"\n", restaurant.name)
		return
	}
	for menu := daily.FirstChild; menu != nil; menu = menu.NextSibling {
		nameText, err := getAttribute(menu, "value")
		if err != nil || nameText == "" {
			continue
		}
		restaurant.AddPermanent(false, nameText, "", 0)
	}

	restaurant.menus[0].SetDay("Monday")
	restaurant.menus[1].SetDay("Tuesday")
	restaurant.menus[2].SetDay("Wednesday")
	restaurant.menus[3].SetDay("Thursday")
	restaurant.menus[4].SetDay("Friday")
	restaurant.menus[5].SetDay("Saturday")
	restaurant.menus[6].SetDay("Sunday")
}
