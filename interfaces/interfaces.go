package interfaces

import (
)

type ITrader interface {

	HowMuchWoodForXFood(int) int
	HowMuchFoodForXWood(int) int

	Food() int
	Wood() int

	SetFood(int)
	SetWood(int)

}
