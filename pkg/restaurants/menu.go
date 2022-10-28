package restaurants

import "encoding/json"

type Menu struct {
	meals []Meal
	valid bool
	day   string
}

func MakeMenu(meals []Meal, day string) Menu {
	return Menu{meals, true, day}
}

func (menu *Menu) Add(isSoup bool, name string, desc string, price int) {
	menu.AddMeal(MakeMeal(isSoup, name, desc, price))
}

func (menu *Menu) AddMeal(meal Meal) {
	menu.meals = append(menu.meals, meal)
}

func (menu Menu) GetMeals() []Meal {
	return menu.meals
}

func (menu *Menu) SetInvalidMenu(invalid bool) {
	menu.valid = !invalid
}

func (menu *Menu) SetValidMenu(valid bool) {
	menu.valid = valid
}

func (menu Menu) IsValid() bool {
	return menu.valid
}

func (menu *Menu) SetDay(day string) {
	menu.day = day
}

func (menu Menu) GetDay() string {
	return menu.day
}

func (menu *Menu) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Meals []Meal `json:"meals"`
		Day   string `json:"day"`
	}{
		Meals: menu.meals,
		Day:   menu.day,
	})
}
