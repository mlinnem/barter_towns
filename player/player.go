package player

import (
	"bufio"
	"fmt"
	"os"

	broker "github.com/mlinnem/barter_towns/broker"
	worldStateLib "github.com/mlinnem/barter_towns/world_state"
)

const INITIAL_FOOD_IN_WAREHOUSE = 100
const INITIAL_WOOD_IN_WAREHOUSE = 100

type Player struct {
	food_in_warehouse int
	wood_in_warehouse int

	scanner *bufio.Scanner
}

func Construct() *Player {
	return &Player{food_in_warehouse: INITIAL_FOOD_IN_WAREHOUSE, wood_in_warehouse: INITIAL_WOOD_IN_WAREHOUSE, scanner: bufio.NewScanner(os.Stdin)}
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

func (player *Player) HowMuchWoodForXFood(amountInt int) int {
	fmt.Printf("ERROR: Method not yet implemented")
	return 0
}

func (player *Player) HowMuchFoodForXWood(amountInt int) int {
	fmt.Printf("ERROR: Method not yet implemented")
	return 0
}

func (player *Player) MakeDecisions(worldState *worldStateLib.WorldState) {

	for i := range worldState.Towns {

		fmt.Printf("Player decision:")
		player.scanner.Scan()
		result := player.scanner.Text()
		fmt.Println(player.scanner.Text())

		if result == "sell 100 wood" {
			broker.BuyFoodWithXWood(player, 100, worldState.Managers[i])
		} else if result == "sell 100 food" {
			broker.BuyWoodWithXFood(player, 100, worldState.Managers[i])
		}

		fmt.Printf("Wood in warehouse: %d\n", player.wood_in_warehouse)
		fmt.Printf("Food in warehouse: %d\n", player.food_in_warehouse)
		fmt.Printf("-----------------\n")
	}
}
