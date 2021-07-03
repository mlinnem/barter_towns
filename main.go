package main

import (
	playerLib "github.com/mlinnem/barter_towns/player"
	worldStateLib "github.com/mlinnem/barter_towns/world_state"
)

const MAX_TIME = 200

func main() {

	//----setup----

	player := playerLib.Construct()
	worldState := worldStateLib.Construct()

	//----main loop----

	for worldState.Year <= MAX_TIME {

		//Player decisions

		player.MakeDecisions(worldState)

		//World evolves

		worldState.AdvanceTime()
	}
}
