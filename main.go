package main

import (
	"fmt"

	playerLib "github.com/mlinnem/barter_towns/player"
	worldStateLib "github.com/mlinnem/barter_towns/world_state"
)

const MAX_TIME = 3999

func main() {

	//----setup----

	player := playerLib.Construct()
	worldState := worldStateLib.Construct()

	//----main loop----

	for worldState.Year <= MAX_TIME {
		fmt.Printf("=====Beginning of Year %d=====\n", worldState.Year)

		//Player decisions

		player.MakeDecisions(worldState)

		//World evolves

		worldState.AdvanceTime()

		fmt.Printf("=====End of Year %d=====\n", worldState.Year)
	}
}
