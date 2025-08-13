package restaurants

type StaticRestaurant struct {
	Restaurant
}

func NewStaticRestaurant(url string, name string, id int) *StaticRestaurant {
	restaurant := new(StaticRestaurant)
	restaurant.SetDefaultValues()
	restaurant.id = id
	restaurant.url = url
	restaurant.name = name
	return restaurant
}

func (restaurant *StaticRestaurant) Parse() {
	restaurant.clearMenus()
	restaurant.clearPermanentMenus()

	restaurant.AddPermanent(false, restaurant.url, "", 0)

	restaurant.menus[0].SetDay("Monday")
	restaurant.menus[1].SetDay("Tuesday")
	restaurant.menus[2].SetDay("Wednesday")
	restaurant.menus[3].SetDay("Thursday")
	restaurant.menus[4].SetDay("Friday")
	restaurant.menus[5].SetDay("Saturday")
	restaurant.menus[6].SetDay("Sunday")
}
