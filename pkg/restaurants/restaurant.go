package restaurants

import "encoding/json"

type RestaurantInterface interface {
	// public
	GetMenus() [7]Menu
	Parse()
	AddPermanentMeal(meal Meal)
	MarshalJSON() ([]byte, error)
	GetSpecificDayObject(days []int) RestaurantJSON

	// private
	clearMenus()
}

type Restaurant struct {
	RestaurantInterface
	id        int
	url       string
	name      string
	menus     [7]Menu
	permanent []Meal
}

type RestaurantJSON struct {
	Id             int    `json:"id"`
	Restaurant     string `json:"restaurant"`
	DailyMenus     []Menu `json:"dailymenus"`
	PermanentMeals []Meal `json:"permanentmeals"`
}

func (restaurant *Restaurant) SetDefaultValues() {
	restaurant.url = ""
	restaurant.name = ""
	restaurant.menus = [7]Menu{}
	restaurant.permanent = []Meal{}
}

func (restaurant *Restaurant) AddPermanent(isSoup bool, name string, desc string, price int) {
	restaurant.AddPermanentMeal(MakeMeal(isSoup, name, desc, price))
}

func (restaurant *Restaurant) AddPermanentMeal(meal Meal) {
	restaurant.permanent = append(restaurant.permanent, meal)
}

func (restaurant Restaurant) GetMenus() [7]Menu {
	return restaurant.menus
}

func (restaurant *Restaurant) clearMenus() {
	for i := 0; i < 7; i++ {
		restaurant.menus[i] = MakeMenuDefault()
	}
}

func (restaurant *Restaurant) MarshalJSON() ([]byte, error) {
	return json.Marshal(&RestaurantJSON{
		Id:             restaurant.id,
		Restaurant:     restaurant.name,
		DailyMenus:     restaurant.menus[:],
		PermanentMeals: restaurant.permanent,
	})
}

func (restaurant *Restaurant) GetSpecificDayObject(days []int) RestaurantJSON {
	obj := RestaurantJSON{
		Id:             restaurant.id,
		Restaurant:     restaurant.name,
		PermanentMeals: restaurant.permanent,
	}
	for _, index := range days {
		obj.DailyMenus = append(obj.DailyMenus, restaurant.menus[index])
	}
	return obj
}
