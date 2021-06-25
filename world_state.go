package main

import (
	"player"
)

struct WorldState {
	Towns []*Town
	Year int
}

func construct() WorldState {
	return WorldState{Towns: zzzzz, Year: 0}
}