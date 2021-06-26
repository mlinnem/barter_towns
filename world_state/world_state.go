package world_state

import (
	t "github.com/mlinnem/barter_towns/town"
	town_manager "github.com/mlinnem/barter_towns/town_manager"
)

type WorldState struct {
	Towns []*t.Town
	Year int
}

func construct() WorldState {
	towns := make([]*t.Town, 1)
	for i := range towns {
		towns[i] = t.Construct()
	}
	return WorldState{Towns: towns, Year: 0}
}

func (worldState *WorldState) advanceTime() {
	for _, town := range worldState.Towns {
		town_manager.makeDecisions(town)
		town.advanceTime()
	}
}