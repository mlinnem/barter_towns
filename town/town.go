package town

import (
	tl "github.com/mlinnem/barter_towns/tile"
	"math/rand"
	"sort"
	"fmt"
	"math"
	town_manager "github.com/mlinnem/barter_towns/town_manager"
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
	Tiles []*tl.Tile

	Population int
	Food int
	Wood int
}
//--Constructor
func Construct() *Town {
	rand.Seed(14)

	var tiles = make([]*tl.Tile, TOWN_SIZE)

	//Set land quality

	for i := 0; i < TOWN_SIZE; i++ {
		r := rand.Float32()
		var t tl.TileType
		var q int
		var hh bool

		hh = false

		//80% plains, 20% forest
		if r < .7 {
			t = tl.Plains
		} else {
			t = tl.Forest
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

		tiles[i] = &tl.Tile{Type: t, Quality: q, HasHouse: hh, TileID: i}

		//fmt.Printf("Tile %d is of type %d and of quality %d\n", tiles[i].TileID, tiles[i].Type, tiles[i].Quality)
	}

	//construct town

	town := Town{}

	town.Tiles = tiles

	var plainsTiles []*tl.Tile
	var forestTiles []*tl.Tile

	for i := 0; i < len(tiles); i++ {

		if tiles[i].Type == tl.Plains {
			plainsTiles = append(plainsTiles, tiles[i])
		} else if tiles[i].Type == tl.Forest {
			forestTiles = append(forestTiles, tiles[i])
		}
	}

	return &town
}

//--Utilities

func (town *Town) GetTiles() []*tl.Tile {
	return town.Tiles
}

//TODO: This will be a source of performance issues later?
func (town *Town) GetHouseCount() int {
	return len(tl.WithHouses(town.Tiles))
}

func SortByQualityInPlace(tiles []*tl.Tile) {
	sort.SliceStable(tiles, func(i, j int) bool {
		return tiles[i].Quality > tiles[j].Quality
	})
}

func SetAdjustedTileQuality(tile *tl.Tile, plainsDemand float64, forestDemand float64) float64 {
	if (tile.Type == tl.Plains) {
		return float64(tile.Quality) * plainsDemand;
	} else { //forest TODO: Make else throw real error
		return float64(tile.Quality) * forestDemand;
	} 
}

//--Repercussions

func (town *Town) collapseRandomHouse() bool {
	var tilesWithHouses = tl.WithHouses(town.GetTiles())
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

//--Actions

func (town *Town) BuildHouseOn(tileID int) {
	town.Tiles[tileID].HasHouse = true
}


func (town *Town) AdvanceTime() {
	//food consumption & starvation

	food_demand := town.Population * FOOD_MAINTENANCE_PER_POP

	if food_demand > town.Population {
		food_shortfall := food_demand - town.Food
		fmt.Printf("Starvation! Food shortfall of %d\n", food_shortfall)
		town.Population = town.Population - food_shortfall
		fmt.Printf("%d died\n", food_shortfall)
		food_demand = food_demand - food_shortfall
	}
	town.Food = town.Food - food_demand

	//house maintenance, collapsing, and exposure

	wood_for_maintenance_demand := int(math.Ceil(float64(town.GetHouseCount() * WOOD_MAINTENANCE_PER_HOUSE)))
	if wood_for_maintenance_demand > town.Wood {

		wood_shortfall := wood_for_maintenance_demand - town.Wood
		fmt.Printf("Not enough wood for house maintenance! Wood shortfall of %d\n", wood_shortfall)
		house_collapse_count := int(math.Ceil(float64(wood_shortfall) * (1.0 / WOOD_COST_PER_HOUSE)))
		for n := 0; n < house_collapse_count; n++ {
			town.collapseRandomHouse()
		}
		fmt.Printf("%d houses collapsed\n", house_collapse_count)
		wood_recovered := COLLAPSE_WOOD_RECOVERY * house_collapse_count
		fmt.Printf("%d wood recovered from collapsed houses\n", wood_recovered)
		town.Wood = town.Wood + wood_recovered
	} else {
		town.Wood = town.Wood - wood_for_maintenance_demand
	}

	housing_demand := town.Population

	if housing_demand > town.GetHouseCount() {
		var housing_shortfall = housing_demand - town.GetHouseCount()
		fmt.Printf("Exposure! Housing shortfall of %d\n", housing_shortfall)
		town.Population = town.Population - housing_shortfall
		fmt.Printf("%d died\n", housing_shortfall)
	}

	//reproduction
	
	remaining_food := town.Food
	if town.Population >= 2 && remaining_food/food_demand > 5 && float64(town.GetHouseCount()) > float64((town.Population+1))*1.05 {
		town.Population = town.Population + 1
		fmt.Printf("%d new baby\n", 1)
	}
}