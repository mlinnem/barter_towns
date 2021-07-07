package world_state

import (
	"fmt"

	t "github.com/mlinnem/barter_towns/town"
	tm "github.com/mlinnem/barter_towns/town_manager"
)

type WorldState struct {
	Towns    []*t.Town
	Managers []*tm.TownManager
	Year     int
}

const NUM_TOWNS = 1

func Construct() *WorldState {
	towns := make([]*t.Town, NUM_TOWNS)
	managers := make([]*tm.TownManager, NUM_TOWNS)

	for i := range towns {
		fmt.Printf("Making town %d...\n", i)
		towns[i] = t.Construct()
		managers[i] = tm.Construct(towns[i])
	}
	return &WorldState{Towns: towns, Managers: managers, Year: 0}
}

func (worldState *WorldState) AdvanceTime() {
	for i, town := range worldState.Towns {
		var town_manager = worldState.Managers[i]
		town_manager.TakeActions()
		town.AdvanceTime()
	}

	worldState.Year = worldState.Year + 1
}
