package broker

import (
	"fmt"

	inter "github.com/mlinnem/barter_towns/interfaces"
)

func BuyWoodWithXFood(player inter.ITrader, amount int, town inter.ITrader) bool {
	if player.Food() < amount {
		fmt.Printf("ERROR: player tried to spend %d food but only had %d\n", amount, player.Food())
		return false
	}

	inReturn := town.HowMuchWoodForXFood(amount)

	fmt.Printf("Town offers %d wood for your %d food\n", inReturn, amount)

	if town.Wood() < inReturn {
		fmt.Printf("ERROR: Town offered to provide %d wood, but only has %d wood left\n", inReturn, town.Wood)
		return false
	}

	player.SetFood(player.Food() - amount)
	town.SetFood(town.Food() + amount)

	player.SetWood(player.Wood() + inReturn)
	town.SetWood(town.Wood() - inReturn)

	return true
}

func BuyFoodWithXWood(player inter.ITrader, amount int, town inter.ITrader) bool {
	if player.Wood() < amount {
		fmt.Printf("ERROR: player tried to spend %d wood but only had %d\n", amount, player.Wood())
		return false
	}

	inReturn := town.HowMuchFoodForXWood(amount)

	fmt.Printf("Town offers %d food for your %d wood\n", inReturn, amount)

	if town.Food() < inReturn {
		fmt.Printf("ERROR: Town offered to provide %d food, but only has %d food left\n", inReturn, town.Food())
		return false
	}

	player.SetWood(player.Wood() - amount)
	town.SetWood(town.Wood() + amount)

	player.SetFood(player.Food() + inReturn)
	town.SetFood(town.Food() - inReturn)

	return true
}
