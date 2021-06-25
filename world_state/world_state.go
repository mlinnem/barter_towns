package world_state

import (
)

struct WorldState {
	Towns []*Town
	Year int
}

func construct() WorldState {
	return WorldState{Towns: zzzzz, Year: 0}
}

func (worldState *WorldState) advanceTime() {
	for town range := town {
		town_manager.makeDecisions(town)
		town.advanceTime()
	}
}