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

func dayToIndex(day string) (int, error) {
	switch day {
	case "Pond\u011bl\u00ed", "Pondeli":
		return 0, nil
	case "\u00dater\u00fd", "Utery":
		return 1, nil
	case "St\u0159eda", "Streda":
		return 2, nil
	case "\u010ctvrtek", "Ctvrtek":
		return 3, nil
	case "P\u00e1tek", "Patek":
		return 4, nil
	case "Sobota":
		return 5, nil
	case "Ned\u011ble", "Nedele":
		return 6, nil
	}
	return -1, errors.New("couldn't parse the day")
}

func getMenickaMealName(node *html.Node) string {
	if node.Type == html.ElementNode {
		if node.Data == "em" || hasKeyValue(node, "class", "poradi") {
			return ""
		}
	}
	if node.Type == html.TextNode {
		return node.Data
	}

	var builder strings.Builder
	for n := node.FirstChild; n != nil; n = n.NextSibling {
		text := getMenickaMealName(n)
		if text == "" {
			continue
		}
		builder.WriteString(text)
		builder.WriteString(" ")
	}
	return builder.String()
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

			dayText, err := getText(day)
			if err != nil {
				continue
			}
			dayText = strings.TrimSpace(strings.Split(dayText, " ")[0])
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
				name := normalizeWhitespace(getMenickaMealName(nameNode))
				price := -1
				priceNode, err := findNodeByClass(meal, "cena")
				if err == nil {
					priceStr, err := getText(priceNode)
					if err != nil {
						continue
					}
					price, err = strconv.Atoi(strings.Split(normalizeWhitespace(priceStr), " ")[0])
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
