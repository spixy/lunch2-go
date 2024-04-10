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

func getIndexFromDay(line string) int {
	if len(line) < 3 || line[:3] == "Pol" {
		return -1
	}
	ascii := line[:2]
	unicode := line[:3]
	if ascii == "Po" {
		return 0
	} else if unicode == "Út" {
		return 1
	} else if ascii == "St" {
		return 2
	} else if unicode == "Čt" {
		return 3
	} else if unicode == "Pá" {
		return 4
	}
	return -1
}

func (restaurant *FreshRestaurant) Parse() {
	defer func() {
		recover()
	}()

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
	curIndex := 0
	pricesIndex := -1
	pricesSection := false
	emptyLine := false
	for scanner.Scan() {
		line := scanner.Text()
		dayIndex := getIndexFromDay(line)
		if dayIndex != -1 {
			curIndex = dayIndex
		}
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
			meals[curIndex] = append(meals[curIndex], line[10:])
			prices[curIndex] = append(prices[curIndex], -1)
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
