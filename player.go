package main

import (
	"bufio"
)

INITIAL_FOOD_IN_WAREHOUSE := 100
INITIAL_WOOD_IN_WAREHOUSE := 100
struct Player {
	food_in_warehouse int
	wood_in_warehouse int
	scanner Scanner
}

func Construct() {
	return Player{food_in_warehouse: INITIAL_FOOD_IN_WAREHOUSE, wood_in_warehouse: INITIAL_WOOD_IN_WAREHOUSE}
}

func (player *Player) makeDecisions(worldState) {
	
	scanner.Scan()
	result := scanner.Text()
	fmt.Println(scanner.Text())

	if result == "buy 100 food" {
		wood_in_warehouse = wood_in_warehouse - int(math.Round(wood_for_a_food * 100 - .5))
		wood = wood + int(math.Round(wood_for_a_food * 100 - .5))
		food = food - 100
		food_in_warehouse = food_in_warehouse + 100
		fmt.Printf("Bought 100 food\n");
		fmt.Printf("Sold %d wood\n", wood_for_a_food * 100)
	} else if result == "buy 100 wood" {
		food_in_warehouse = food_in_warehouse - int(math.Round(food_for_a_wood * 100 - .5))
		food = food + int(math.Round(food_for_a_wood * 100 - .5))
		wood = wood - 100
		wood_in_warehouse = wood_in_warehouse + 100
		fmt.Printf("Bought 100 wood\n");
		fmt.Printf("Sold %d food\n", food_for_a_wood * 100)
	}

	fmt.Printf("Wood in warehouse: %d\n", wood_in_warehouse)
	fmt.Printf("Food in warehouse: %d\n", food_in_warehouse)
	fmt.Printf("-----------------\n")
}