package restaurants

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type GraselRestaurant struct {
	Restaurant
}

var graselDays = map[string]int{
	"pondělí": 0,
	"úterý":   1,
	"středa":  2,
	"čtvrtek": 3,
	"pátek":   4,
	"sobota":  5,
	"neděle":  6,
}

// graselSoupPrice matches the menu soup surcharge in the section header,
// e.g. "Polévka dle denní nabídky – v menu +19 Kč, samostatně 49 Kč."
var graselSoupPrice = regexp.MustCompile(`\+\s*(\d+)\s*Kč`)

func NewGraselRestaurant(url string, name string, id int) *GraselRestaurant {
	restaurant := new(GraselRestaurant)
	restaurant.SetDefaultValues()
	restaurant.id = id
	restaurant.url = url
	restaurant.name = name
	return restaurant
}

// graselPrice extracts the numeric price from strings like "139 Kč".
// Returns -1 when no digits are found.
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

// graselChildText returns the normalized text of the first descendant with
// the given class, or "" when there is none.
func graselChildText(node *html.Node, class string) string {
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

// graselTagText returns the normalized text of the first descendant element
// with the given tag name, or "" when there is none.
func graselTagText(node *html.Node, tag string) string {
	for n := node.FirstChild; n != nil; n = n.NextSibling {
		if n.Type == html.ElementNode && n.Data == tag {
			if text, err := getText(n); err == nil {
				return normalizeWhitespace(text)
			}
			return ""
		}
		if text := graselTagText(n, tag); text != "" {
			return text
		}
	}
	return ""
}

func (restaurant *GraselRestaurant) parseItems(node *html.Node, menu *Menu) {
	if hasKeyValue(node, "class", "menu-list__item") {
		name := graselChildText(node, "menu-list__label")
		if name == "" {
			return
		}
		menu.Add(false, name, graselChildText(node, "menu-list__desc"), graselPrice(graselChildText(node, "menu-list__price")))
		return
	}
	for n := node.FirstChild; n != nil; n = n.NextSibling {
		restaurant.parseItems(n, menu)
	}
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

	// <section id="menu"> holds a <div class="week-day-grid"> with one
	// <article class="menu-card"> per weekday: an <h3> with the Czech day
	// name and a <ul class="menu-list"> of dishes.
	section, err := findNodeById(doc, "menu")
	if err != nil {
		fmt.Printf("Couldn't find content for restaurant \"%s\"\n", restaurant.name)
		return
	}
	grid, err := findNodeByClass(section, "week-day-grid")
	if err != nil {
		fmt.Printf("Couldn't find content for restaurant \"%s\"\n", restaurant.name)
		return
	}

	// The soup is not named on the page ("Polévka dle denní nabídky"), only
	// its surcharge in the section header, so list it as a permanent meal.
	if header, err := findNodeByClass(section, "section__header"); err == nil {
		if text, err := getText(header); err == nil {
			if match := graselSoupPrice.FindStringSubmatch(text); match != nil {
				price, _ := strconv.Atoi(match[1])
				restaurant.AddPermanent(true, "Polévka dle denní nabídky", "", price)
			}
		}
	}

	for card := grid.FirstChild; card != nil; card = card.NextSibling {
		if !hasKeyValue(card, "class", "menu-card") {
			continue
		}
		dayIndex, ok := graselDays[strings.ToLower(graselTagText(card, "h3"))]
		if !ok {
			continue
		}
		restaurant.parseItems(card, &restaurant.menus[dayIndex])
	}

	restaurant.menus[0].SetDay("Monday")
	restaurant.menus[1].SetDay("Tuesday")
	restaurant.menus[2].SetDay("Wednesday")
	restaurant.menus[3].SetDay("Thursday")
	restaurant.menus[4].SetDay("Friday")
	restaurant.menus[5].SetDay("Saturday")
	restaurant.menus[6].SetDay("Sunday")
}
