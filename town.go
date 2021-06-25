package town

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"os"
	"bufio"
	//"reflect"
	//"time"
)
  

const TOWN_SIZE = 200

const INITIAL_FOOD = 900
const INITIAL_WOOD = 900
const INITIAL_POP = 100

const FOOD_MAINTENANCE_PER_POP = 1
const WOOD_MAINTENANCE_PER_HOUSE = 1

const WOOD_COST_PER_HOUSE = 30
const COLLAPSE_WOOD_RECOVERY = 10

type Town struct {
	Tiles []*Tile
}

type Tile struct {
	TileID   int
	Type     TileType
	Quality  int
	HasHouse bool
}

type TileType int

const (
	Plains TileType = iota
	Forest
)

func (town *Town) getTiles() []*Tile {
	return town.Tiles
}

func (town *Town) collapseRandomHouse() bool {
	var tilesWithHouses = town.getTilesWithHouses()
	var houseCount = len(tilesWithHouses)

	if houseCount == 0 {
		fmt.Printf("SYSTEM WARNING: Tried to collapse a house, but there's no house to collapse!")
		return false
	}
	var r = rand.Float64()
	var rScaledAndFloored = int(math.Floor(r * float64(houseCount)))
	var indexOfHouseToDestroy = rScaledAndFloored

	tilesWithHouses[indexOfHouseToDestroy].HasHouse = false

	return true
}

//TODO: Make this boolean and checked
func (town *Town) buildHouseOn(tileID int) {
	town.Tiles[tileID].HasHouse = true
}

func construct() Town {
	rand.Seed(14)

	var tiles = make([]*Tile, TOWN_SIZE)

	//Set land quality

	for i := 0; i < TOWN_SIZE; i++ {
		r := rand.Float32()
		var t TileType
		var q int
		var hh bool

		hh = false

		//80% plains, 20% forest
		if r < .7 {
			t = Plains
		} else {
			t = Forest
		}

		//quality
		//5% great, 15% good, 30% decent, 40% poor, 10% terrible
		r2 := rand.Float32()

		if r2 > .95 {
			q = 5
		} else if r2 > .80 {
			q = 4
		} else if r2 > .50 {
			q = 3
		} else if r2 > .1 {
			q = 2
		} else if r2 > 0.05 {
			q = 1
		} else {
			q = 0
		}

		tiles[i] = &Tile{Type: t, Quality: q, HasHouse: hh, TileID: i}

		//fmt.Printf("Tile %d is of type %d and of quality %d\n", tiles[i].TileID, tiles[i].Type, tiles[i].Quality)
	}

	//construct town

	town := Town{}

	town.Tiles = tiles

	var plainsTiles []*Tile
	var forestTiles []*Tile

	for i := 0; i < len(tiles); i++ {

		if tiles[i].Type == Plains {
			plainsTiles = append(plainsTiles, tiles[i])
		} else if tiles[i].Type == Forest {
			forestTiles = append(forestTiles, tiles[i])
		}
	}

	return town
}

func advanceTime() {

}