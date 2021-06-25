package main

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
  

const MAX_TIME = 2000

const TOWN_SIZE = 200

const INITIAL_FOOD = 900
const INITIAL_WOOD = 900

const WOOD_COST_PER_HOUSE = 30

const WOOD_MAINTENANCE_PER_HOUSE = 1
const COLLAPSE_WOOD_RECOVERY = 10

//not current used everywhere, June 23
const FOOD_MAINTENANCE_PER_POP = 1

const INITIAL_POP = 100

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

func (town *Town) getPlainsTiles() []*Tile {
	var plainsTiles []*Tile
	for i := 0; i < len(town.Tiles); i++ {
		if town.Tiles[i].Type == Plains {
			plainsTiles = append(plainsTiles, town.Tiles[i])

		}
	}

	return plainsTiles
}

func (town *Town) getPlainsTilesWithHouses() []*Tile {
	var plainsTiles = town.getPlainsTiles()

	var plainsTilesWithHouses []*Tile

	for i := 0; i < len(plainsTiles); i++ {
		var plainsTile = plainsTiles[i]

		if plainsTile.HasHouse == true {

			plainsTilesWithHouses = append(plainsTilesWithHouses, plainsTile)

		}
	}

	return plainsTilesWithHouses

}

func (town *Town) getTilesWithHouses() []*Tile {
	var tilesWithHouses []*Tile
	for i := 0; i < len(town.Tiles); i++ {
		if town.Tiles[i].HasHouse == true {
			tilesWithHouses = append(tilesWithHouses, town.Tiles[i])

		}
	}

	return tilesWithHouses
}

func (town *Town) getForestTiles() []*Tile {
	var forestTiles []*Tile
	for i := 0; i < len(town.Tiles); i++ {
		if town.Tiles[i].Type == Forest {
			forestTiles = append(forestTiles, town.Tiles[i])
		}
	}

	return forestTiles
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

func withHouses(in_tiles []*Tile) []*Tile {
	var out_tiles []*Tile
	for _, tile := range in_tiles {
		if tile.HasHouse {
			out_tiles = append(out_tiles, tile)
		}
	}

	return out_tiles
}

func withoutHouses(in_tiles []*Tile) []*Tile {
	var out_tiles []*Tile
	for _, tile := range in_tiles {
		if !tile.HasHouse {
			out_tiles = append(out_tiles, tile)
		}
	}

	return out_tiles
}

func thatAreForest(in_tiles []*Tile) []*Tile {
	var out_tiles []*Tile
	for _, tile := range in_tiles {
		if tile.Type == Forest {
			out_tiles = append(out_tiles, tile)
		}
	}

	return out_tiles
}

func thatArePlains(in_tiles []*Tile) []*Tile {
	var out_tiles []*Tile
	for _, tile := range in_tiles {
		if tile.Type == Plains {
			out_tiles = append(out_tiles, tile)
		}
	}

	return out_tiles
}

//TODO: Make houses and pop global?
func unoccupiedHouses(houses int, pop int) int {
	return int(math.Max(float64(houses - pop), 0))
}

func generateTown(seed int64) Town {
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

func sortByQualityInPlace(tiles []*Tile) {
	sort.SliceStable(tiles, func(i, j int) bool {
		return tiles[i].Quality > tiles[j].Quality
	})
}

func sortByDemandAdjustedQualityInPlace(tiles []*Tile, plainsDemand float64, forestDemand float64) {
	sort.SliceStable(tiles, func(i, j int) bool {
		return getAdjustedTileQuality(tiles[i], plainsDemand, forestDemand) > getAdjustedTileQuality(tiles[j], plainsDemand, forestDemand)
	})
}

func getAdjustedTileQuality(tile *Tile, plainsDemand float64, forestDemand float64) float64 {
	if (tile.Type == Plains) {
		return float64(tile.Quality) * plainsDemand;
	} else { //forest TODO: Make else throw real error
		return float64(tile.Quality) * forestDemand;
	} 
}

func main() {


	food_in_warehouse := 100
	wood_in_warehouse := 100

	scanner := bufio.NewScanner(os.Stdin)

	//setup

	year := 0
	pop := INITIAL_POP
	food := INITIAL_FOOD
	wood := INITIAL_WOOD

	town := generateTown(5)

	for year <= MAX_TIME && pop > 0 {

		
		

		//Status

		fmt.Printf("-----year %d-----\n", year)
		fmt.Printf("population: %d\n", pop)
		fmt.Printf("food: %d\n", food)
		fmt.Printf("wood: %d\n", wood)

		//Building Strategy:

		houseCount := len(town.getTilesWithHouses())
		fmt.Printf("House count: %d\n", houseCount);
		wood_needs_for_next_30_years := houseCount * 30 * WOOD_MAINTENANCE_PER_HOUSE

		var counter = 0

		for wood > wood_needs_for_next_30_years && houseCount < pop*2 && wood >= WOOD_COST_PER_HOUSE {
			counter = counter + 1
			fmt.Printf("loop %d in house building \n", counter)
			houseCount = len(town.getTilesWithHouses())
			wood_needs_for_next_30_years = houseCount * 30 * WOOD_MAINTENANCE_PER_HOUSE

			//calculate food demand & wood demand for next 100 years, assuming 1.5x population

			hundred_year_food_need := int(float64(pop) * 1.25 * 100)
			hundred_year_wood_need := int(float64(pop)*1.25*100 + (float64(pop) * .5 * WOOD_COST_PER_HOUSE))

			lt_wood_demand := hundred_year_wood_need / (wood + 10)
			lt_food_demand := hundred_year_food_need / (food + 10)

			fmt.Printf("Wood demand level: %d\n", lt_wood_demand)
			fmt.Printf("Food demand level: %d\n", lt_food_demand)
			//build a house on best land for that demand.

			var bestTilesToBuildOn []*Tile

			// if lt_wood_demand > lt_food_demand {
			// 	//var allTiles = town.getTiles()
			// 	//fmt.Printf("All tiles count: %n", len(allTiles));
			// 	//fmt.Printf("With houses count: %n", len(allTiles));
			// 	bestTilesToBuildOn = thatAreForest(withoutHouses(town.getTiles()))
			// } else {
			// 	bestTilesToBuildOn = thatArePlains(withoutHouses(town.getTiles()))
			// }

			bestTilesToBuildOn = withoutHouses(town.getTiles())

			sortByDemandAdjustedQualityInPlace(bestTilesToBuildOn, float64(lt_food_demand), float64(lt_wood_demand))

			//Build house
			if len(bestTilesToBuildOn) > 0 {
				fmt.Printf("Good land to build on \n");
				var bestTileToBuildOn = bestTilesToBuildOn[0]
				wood = wood - 30
				town.buildHouseOn(bestTileToBuildOn.TileID)
				fmt.Printf("built a house on tile %d, type of (%d) with quality %d\n", bestTileToBuildOn.TileID, bestTileToBuildOn.Type, bestTileToBuildOn.Quality)
			} else {
				fmt.Printf("No land to build on!\n");
				break
			}

		}

		//Allocate labor


		

		food_cost := 0.0
		wood_cost := 0.0

		existing_wood_maintain := houseCount * WOOD_MAINTENANCE_PER_HOUSE
		existing_food_maintain := 20 * pop * FOOD_MAINTENANCE_PER_POP

		fmt.Printf("Existing wood maintain, %d\n", existing_wood_maintain);
		fmt.Printf("Existing food maintain, %d\n", existing_food_maintain);
	
		new_wood_maintain := (20 / 2) * WOOD_MAINTENANCE_PER_HOUSE
		new_wood_build := int(math.Max((float64(20 - unoccupiedHouses(houseCount, pop)) * WOOD_COST_PER_HOUSE), 0.0))
		st_wood_demand := existing_wood_maintain + new_wood_maintain + new_wood_build
		wood_cost = float64(st_wood_demand) / float64(wood + 10)

			
		new_food_maintain := (20 / 2) * FOOD_MAINTENANCE_PER_POP
		st_food_demand := existing_food_maintain + new_food_maintain
		food_cost = float64(st_food_demand) / float64(food + 10)
	
	

		//Prices
		food_for_a_wood := float64(wood_cost) / float64(food_cost) 
		wood_for_a_food := float64(food_cost) / float64(wood_cost)

		pop_unallocated := pop

		var allHouses = withHouses(town.getTiles())
		sortByDemandAdjustedQualityInPlace(allHouses, food_cost, wood_cost)
		var houseIndex = 0

		for pop_unallocated > 0 && houseIndex < len(allHouses) {
			var topHouse = allHouses[houseIndex]
			pop_unallocated = pop_unallocated - 1
			houseIndex = houseIndex + 1
			if topHouse.Type == Plains {
				food = food + topHouse.Quality;
				fmt.Printf("produced %d food from tile %d\n", topHouse.Quality, topHouse.TileID)
			} else {
				wood = wood + topHouse.Quality;
				fmt.Printf("produced %d wood from tile %d\n", topHouse.Quality, topHouse.TileID)
			}
		}
		

		//Update

		year += 1

		//food consumption & starvation

		food_demand := pop
		if food_demand > food {

			food_shortfall := food_demand - food
			fmt.Printf("Starvation! Food shortfall of %d\n", food_shortfall)
			pop = pop - food_shortfall
			fmt.Printf("%d died\n", food_shortfall)
			food_demand = food_demand - food_shortfall
		}
		food = food - food_demand

		//house maintenance & collapsing

		var houses = len(town.getTilesWithHouses())

		wood_for_maintenance_demand := int(math.Ceil(float64(houses * WOOD_MAINTENANCE_PER_HOUSE)))
		if wood_for_maintenance_demand > wood {

			wood_shortfall := wood_for_maintenance_demand - wood
			fmt.Printf("Not enough wood for house maintenance! Wood shortfall of %d\n", wood_shortfall)
			house_collapse_count := int(math.Ceil(float64(wood_shortfall) * (1.0 / WOOD_COST_PER_HOUSE)))
			for n := 0; n < house_collapse_count; n++ {
				town.collapseRandomHouse()
			}
			fmt.Printf("%d houses collapsed\n", house_collapse_count)
			wood_recovered := COLLAPSE_WOOD_RECOVERY * house_collapse_count
			fmt.Printf("%d wood recovered from collapsed houses\n", wood_recovered)
			wood = wood_recovered

		} else {
			wood = wood - wood_for_maintenance_demand
		}
		//exposure check

		housing_demand := pop

		if housing_demand > houses {
			var housing_shortfall = housing_demand - houses
			fmt.Printf("Exposure! Housing shortfall of %d\n", housing_shortfall)
			pop = pop - housing_shortfall
			fmt.Printf("%d died\n", housing_shortfall)
		}

		//reproduction

		

		var houses_now = len(town.getTilesWithHouses())
		remaining_food := food
		if pop >= 2 && remaining_food/food_demand > 5 && float64(houses_now) > float64((pop+1))*1.05 {
			pop = pop + 1
			fmt.Printf("%d new baby\n", 1)
		}

		//trade
		fmt.Printf("You can buy 100 wood for %d food\n", int(math.Max(food_for_a_wood * 100, 1)))
		fmt.Printf("You can buy 100 food for %d wood\n", int(math.Max(wood_for_a_food * 100, 1)))

		
		scanner.Scan()
		result := scanner.Text()
		fmt.Println(scanner.Text())

		if result == "buy 100 food" {
			wood_in_warehouse = wood_in_warehouse - int(math.Round(wood_for_a_food * 100 - .5))
			wood = wood + int(math.Round(wood_for_a_food * 100 - .5))
			food = food - 100
			food_in_warehouse = food_in_warehouse + 100
			fmt.Printf("Bought 100 food\n");
			fmt.Printf("Sold %d wood\n", wood_for_a_food * 100)
		} else if result == "buy 100 wood" {
			food_in_warehouse = food_in_warehouse - int(math.Round(food_for_a_wood * 100 - .5))
			food = food + int(math.Round(food_for_a_wood * 100 - .5))
			wood = wood - 100
			wood_in_warehouse = wood_in_warehouse + 100
			fmt.Printf("Bought 100 wood\n");
			fmt.Printf("Sold %d food\n", food_for_a_wood * 100)
		}

		fmt.Printf("Wood in warehouse: %d\n", wood_in_warehouse)
		fmt.Printf("Food in warehouse: %d\n", food_in_warehouse)
		fmt.Printf("-----------------\n")

		
	}
}
