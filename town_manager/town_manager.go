package town_manager

import (
	"fmt"
	"math"
	"sort"

	ti "github.com/mlinnem/barter_towns/tile"
	to "github.com/mlinnem/barter_towns/town"
)

type TownManager struct {
	Town *to.Town
}

func Construct(town *to.Town) *TownManager {
	var manager = &TownManager{Town: town}
	return manager
}

func (manager *TownManager) HowMuchWoodForXFood(amountInt int) int {
	var amount = float64(amountInt)
	if amount <= 0 {
		return 0
	}

	woodDemand := manager.getWoodDemand()
	foodDemand := manager.getFoodDemand()

	var exchangeRate = woodDemand / foodDemand

	var offer = math.Floor(amount * exchangeRate)

	if offer > amount*20 {
		offer = amount * 20
	}

	offer = math.Min(offer, float64(manager.Town.Wood+1))

	return int(offer)
}

func (manager *TownManager) HowMuchFoodForXWood(amountInt int) int {
	var amount = float64(amountInt)
	if amount <= 0 {
		return 0
	}

	var woodDemand = manager.getWoodDemand()
	var foodDemand = manager.getFoodDemand()

	var exchangeRate = foodDemand / woodDemand

	var offer = math.Floor(amount * exchangeRate)

	if offer > amount*20 {
		offer = amount * 20
	}

	offer = math.Max(offer, float64(manager.Town.Food+1))

	return int(offer)
}

func (manager *TownManager) Food() int {
	return manager.Town.Food
}

func (manager *TownManager) SetFood(amount int) {
	manager.Town.Food = amount
}

func (manager *TownManager) Wood() int {
	return manager.Town.Food
}

func (manager *TownManager) SetWood(amount int) {
	manager.Town.Wood = amount
}

const TIMELINE_PROJECTION_IN_YEARS = 20

func (manager *TownManager) getWoodDemand() float64 {

	var additionalHouseCount = int(TIMELINE_PROJECTION_IN_YEARS / 4) //make a new house every 4 years
	var demand_from_maintenance = manager.Town.GetHouseCount() + int(additionalHouseCount/2)
	var demand_from_building = additionalHouseCount * to.WOOD_COST_PER_HOUSE
	var woodDemand = (demand_from_maintenance + demand_from_building) / (manager.Town.Wood + 1)

	return float64(woodDemand)
}

func (manager *TownManager) getFoodDemand() float64 {
	var additionalPopulationCount = TIMELINE_PROJECTION_IN_YEARS / 4 //make a new person every 4 years
	var demand_from_maintenance = ((manager.Town.Population + additionalPopulationCount) * to.FOOD_MAINTENANCE_PER_POP)
	var foodDemand = (demand_from_maintenance + 1) / (manager.Town.Food + 1)

	return float64(foodDemand)
}

func (manager *TownManager) TakeActions() {

	//---Build buildings

	houseCount := manager.Town.GetHouseCount()
	fmt.Printf("House count: %d\n", houseCount)
	wood_needs_for_next_30_years := houseCount * 30 * to.WOOD_MAINTENANCE_PER_HOUSE

	var counter = 0

	for manager.Town.Wood > wood_needs_for_next_30_years && houseCount < manager.Town.Population*2 && manager.Town.Wood >= to.WOOD_COST_PER_HOUSE {
		counter = counter + 1
		fmt.Printf("loop %d in house building \n", counter)
		houseCount = manager.Town.GetHouseCount()
		wood_needs_for_next_30_years = houseCount * 30 * to.WOOD_MAINTENANCE_PER_HOUSE

		//calculate food demand & wood demand for next 100 years, assuming 1.5x population

		hundred_year_food_need := int(float64(manager.Town.Population) * 1.25 * 100)
		hundred_year_wood_need := int(float64(manager.Town.Population)*1.25*100 + (float64(manager.Town.Population) * .5 * to.WOOD_COST_PER_HOUSE))

		lt_wood_demand := hundred_year_wood_need / (manager.Town.Wood + 10)
		lt_food_demand := hundred_year_food_need / (manager.Town.Food + 10)

		fmt.Printf("Wood demand level: %d\n", lt_wood_demand)
		fmt.Printf("Food demand level: %d\n", lt_food_demand)
		//build a house on best land for that demand.

		var bestTilesToBuildOn []*ti.Tile

		// if lt_wood_demand > lt_food_demand {
		// 	//var allTiles = town.getTiles()
		// 	//fmt.Printf("All tiles count: %n", len(allTiles));
		// 	//fmt.Printf("With houses count: %n", len(allTiles));
		// 	bestTilesToBuildOn = thatAreForest(withoutHouses(town.getTiles()))
		// } else {
		// 	bestTilesToBuildOn = thatArePlains(withoutHouses(town.getTiles()))
		// }

		bestTilesToBuildOn = ti.WithoutHouses(manager.Town.GetTiles())

		sortByDemandAdjustedQualityInPlace(bestTilesToBuildOn, float64(lt_food_demand), float64(lt_wood_demand))

		//Build house
		if len(bestTilesToBuildOn) > 0 {
			fmt.Printf("Good land to build on \n")
			var bestTileToBuildOn = bestTilesToBuildOn[0]
			manager.Town.Wood = manager.Town.Wood - 30
			manager.Town.BuildHouseOn(bestTileToBuildOn.TileID)
			fmt.Printf("built a house on tile %d, type of (%d) with quality %d\n", bestTileToBuildOn.TileID, bestTileToBuildOn.Type, bestTileToBuildOn.Quality)
		} else {
			fmt.Printf("No land to build on!\n")
			break
		}

	}

	//---Allocate labor

	food_cost := 0.0
	wood_cost := 0.0

	existing_wood_maintain := houseCount * to.WOOD_MAINTENANCE_PER_HOUSE
	existing_food_maintain := 20 * manager.Town.Population * to.FOOD_MAINTENANCE_PER_POP

	fmt.Printf("Existing wood maintain, %d\n", existing_wood_maintain)
	fmt.Printf("Existing food maintain, %d\n", existing_food_maintain)

	new_wood_maintain := (20 / 2) * to.WOOD_MAINTENANCE_PER_HOUSE
	new_wood_build := int(math.Max((float64(20-unoccupiedHouses(houseCount, manager.Town.Population)) * to.WOOD_COST_PER_HOUSE), 0.0))
	st_wood_demand := existing_wood_maintain + new_wood_maintain + new_wood_build
	wood_cost = float64(st_wood_demand) / float64(manager.Town.Wood+10)

	new_food_maintain := (20 / 2) * to.FOOD_MAINTENANCE_PER_POP
	st_food_demand := existing_food_maintain + new_food_maintain
	food_cost = float64(st_food_demand) / float64(manager.Town.Food+10)

	pop_unallocated := manager.Town.Population

	var allHouses = ti.WithHouses(manager.Town.GetTiles())
	sortByDemandAdjustedQualityInPlace(allHouses, food_cost, wood_cost)
	var houseIndex = 0

	for pop_unallocated > 0 && houseIndex < len(allHouses) {
		var topHouse = allHouses[houseIndex]
		pop_unallocated = pop_unallocated - 1
		houseIndex = houseIndex + 1
		if topHouse.Type == ti.Plains {
			manager.Town.Food = manager.Town.Food + topHouse.Quality
			fmt.Printf("produced %d food from tile %d\n", topHouse.Quality, topHouse.TileID)
		} else { //TODO: Will  break when more tiles added
			manager.Town.Wood = manager.Town.Wood + topHouse.Quality
			fmt.Printf("produced %d wood from tile %d\n", topHouse.Quality, topHouse.TileID)
		}
	}

}

func withHouses(in_tiles []*ti.Tile) []*ti.Tile {
	var out_tiles []*ti.Tile
	for _, tile := range in_tiles {
		if tile.HasHouse {
			out_tiles = append(out_tiles, tile)
		}
	}

	return out_tiles
}

func withoutHouses(in_tiles []*ti.Tile) []*ti.Tile {
	var out_tiles []*ti.Tile
	for _, tile := range in_tiles {
		if !tile.HasHouse {
			out_tiles = append(out_tiles, tile)
		}
	}

	return out_tiles
}

func thatAreForest(in_tiles []*ti.Tile) []*ti.Tile {
	var out_tiles []*ti.Tile
	for _, tile := range in_tiles {
		if tile.Type == ti.Forest {
			out_tiles = append(out_tiles, tile)
		}
	}

	return out_tiles
}

func thatArePlains(in_tiles []*ti.Tile) []*ti.Tile {
	var out_tiles []*ti.Tile
	for _, tile := range in_tiles {
		if tile.Type == ti.Plains {
			out_tiles = append(out_tiles, tile)
		}
	}

	return out_tiles
}

func getAdjustedTileQuality(tile *ti.Tile, plainsDemand float64, forestDemand float64) float64 {
	if tile.Type == ti.Plains {
		return float64(tile.Quality) * plainsDemand
	} else { //forest TODO: Make else throw real error
		return float64(tile.Quality) * forestDemand
	}
}

func sortByDemandAdjustedQualityInPlace(tiles []*ti.Tile, plainsDemand float64, forestDemand float64) {
	sort.SliceStable(tiles, func(i, j int) bool {
		return getAdjustedTileQuality(tiles[i], plainsDemand, forestDemand) > getAdjustedTileQuality(tiles[j], plainsDemand, forestDemand)
	})
}

//TODO: Make houses and pop global?
func unoccupiedHouses(houses int, pop int) int {
	return int(math.Max(float64(houses-pop), 0))
}
