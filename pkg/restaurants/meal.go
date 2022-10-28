package restaurants

import "encoding/json"

type Meal struct {
	isSoup bool
	name   string
	desc   string
	price  int
}

func MakeMeal(isSoup bool, name string, desc string, price int) Meal {
	return Meal{isSoup, name, desc, price}
}

func (meal Meal) IsSoup() bool {
	return meal.isSoup
}

func (meal Meal) GetName() string {
	return meal.name
}

func (meal Meal) GetDescription() string {
	return meal.desc
}

func (meal Meal) GetPrice() int {
	return meal.price
}

func (meal *Meal) SetSoup(isSoup bool) {
	meal.isSoup = isSoup
}

func (meal *Meal) SetName(name string) {
	meal.name = name
}

func (meal *Meal) SetDescription(desc string) {
	meal.desc = desc
}

func (meal *Meal) SetPrice(price int) {
	meal.price = price
}

func (meal *Meal) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		IsSoup      bool   `json:"isSoup"`
		Price       int    `json:"price"`
	}{
		Name:        meal.name,
		Description: meal.desc,
		IsSoup:      meal.isSoup,
		Price:       meal.price,
	})
}
