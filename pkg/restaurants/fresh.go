package restaurants

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"code.sajari.com/docconv"
	"golang.org/x/net/html"
)

type FreshRestaurant struct {
	Restaurant
}

func NewFreshRestaurant(url string, name string, id int) *FreshRestaurant {
	restaurant := new(FreshRestaurant)
	restaurant.SetDefaultValues()
	restaurant.id = id
	restaurant.url = url
	restaurant.name = name
	return restaurant
}

func (restaurant *FreshRestaurant) Parse() {
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

	link, err := findNodeById(doc, "menu-item-519")
	if err != nil {
		fmt.Printf("Couldn't find menu node for restaurant \"%s\"\n", restaurant.name)
		return
	}
	linkAddress, err := getAttribute(link.FirstChild, "href")
	if err != nil {
		fmt.Printf("Couldn't get PDF link for restaurant \"%s\"\n", restaurant.name)
		return
	}
	pdf, err := http.Get(linkAddress)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pdf.Body.Close()
	pdftxt, err := docconv.Convert(pdf.Body, "application/pdf", true)
	if err != nil {
		fmt.Println(err)
		return
	}

	meals := [5][]string{}
	prices := [5][]int{}
	scanner := bufio.NewScanner(strings.NewReader(pdftxt.Body))
	curIndex := -1
	pricesIndex := -1
	pricesSection := false
	emptyLine := false
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			emptyLine = true
			pricesSection = false
			continue
		}
		if strings.Contains(line, "Kč") && emptyLine {
			pricesIndex++
			pricesSection = true
		}
		emptyLine = false
		// 10 because soup name starts after 10 and vacations have no soup provided
		if !pricesSection && len(line) > 10 && line[:3] == "Pol" {
			curIndex++
			meals[curIndex] = append(meals[curIndex], line[10:])
			prices[curIndex] = append(prices[curIndex], -1)
		}
		if curIndex < 0 {
			curIndex++
		}
		if !pricesSection && line[1] == '.' && len(line) > 2 {
			meals[curIndex] = append(meals[curIndex], line[3:])
		}
		if pricesSection && len(line) != 0 {
			priceInt, err := strconv.Atoi(strings.Split(line, " ")[0])
			if err != nil {
				priceInt = -1
			}
			prices[pricesIndex] = append(prices[pricesIndex], priceInt)
		}
	}

	for i := 0; i < 5; i++ {
		if len(meals[i]) == 0 || len(prices[i]) == 0 {
			continue
		}
		for ind, meal := range meals[i] {
			if ind >= len(prices[i]) {
				break
			}
			restaurant.menus[i].Add(ind == 0, meal, "", prices[i][ind])
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
