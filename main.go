package main

import (
	playerLib "github.com/mlinnem/barter_towns/player"
)

func main() {

	//----setup----

	player := playerLib.construct()
	worldState := worldState.construct()

	//----main loop----

	for year <= MAX_TIME {
		
		//Player decisions
	
		player.makeDecisions(worldState)

		//World evolves

		worldState.advanceTime()
	}
}
