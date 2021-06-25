package player

import (
	"bufio"
	"os"
	broker "github.com/mlinnem/barter_towns/broker"
	world_state "github.com/mlinnem/barter_towns/world_state"
)

const INITIAL_FOOD_IN_WAREHOUSE = 100
const INITIAL_WOOD_IN_WAREHOUSE = 100

type Player struct {
	food_in_warehouse int
	wood_in_warehouse int

	scanner *bufio.Scanner
}

func Construct() *Player {
	return &Player{food_in_warehouse: INITIAL_FOOD_IN_WAREHOUSE, wood_in_warehouse: INITIAL_WOOD_IN_WAREHOUSE, scanner : bufio.NewScanner(os.Stdin)}
}

func (player *Player) Food() int {
	return player.food_in_warehouse
}

func (player *Player) SetFood(value int) {
	player.food_in_warehouse = value
}

func (player *Player) Wood() int {
	return player.wood_in_warehouse
}

func (player *Player) SetWood(value int) {
	player.wood_in_warehouse = value
}

func (player *Player) makeDecisions(worldState *world_state.WorldState) {

	scanner.Scan()
	result := scanner.Text()
	fmt.Println(scanner.Text())

	if result == "buy 100 food" {
		broker.BuyWoodWithFood(player, 100, town)
	} else if result == "buy 100 wood" {
		broker.BuyFoodWithWood(player, 100, town)
	}

	fmt.Printf("Wood in warehouse: %d\n", wood_in_warehouse)
	fmt.Printf("Food in warehouse: %d\n", food_in_warehouse)
	fmt.Printf("-----------------\n")
}