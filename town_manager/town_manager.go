package town_manager

import (
	"fmt"
	"math"

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

	fmt.Printf("Wood demand is...%8.2f\n", woodDemand)
	fmt.Printf("Food demand is...%8.2f\n", foodDemand)

	var exchangeRate = foodDemand / woodDemand

	var offer = math.Floor(amount * exchangeRate)

	if offer > amount*20 {
		offer = amount * 20
	}

	//never give more than half your stockpile
	if offer > float64(manager.Town.Wood)/2 {
		offer = float64(manager.Town.Wood) / 2
	}

	//Never offer less than 0
	offer = math.Max(offer, 0)

	return int(offer)
}

func (manager *TownManager) HowMuchFoodForXWood(amountInt int) int {
	var amount = float64(amountInt)
	if amount <= 0 {
		return 0
	}

	var woodDemand = manager.getWoodDemand()
	var foodDemand = manager.getFoodDemand()

	fmt.Printf("Wood demand is...%8.2f\n", woodDemand)
	fmt.Printf("Food demand is...%8.2f\n", foodDemand)

	var exchangeRate = woodDemand / foodDemand

	var offer = math.Floor(amount * exchangeRate)

	if offer > amount*20 {
		offer = amount * 20
	}

	//never give more than half your stockpile
	if offer > float64(manager.Town.Food)/2 {
		offer = float64(manager.Town.Food) / 2
	}

	//Never offer less than 0
	offer = math.Max(offer, 0)

	return int(offer)
}

func (manager *TownManager) Food() int {
	return manager.Town.Food
}

func (manager *TownManager) SetFood(amount int) {
	manager.Town.Food = amount
}

func (manager *TownManager) Wood() int {
	return manager.Town.Wood
}

func (manager *TownManager) SetWood(amount int) {
	manager.Town.Wood = amount
}

func (manager *TownManager) woodProductionCapacity() int {
	var producingTiles = ti.ThatAreForest(ti.WithHouses(manager.Town.GetTiles()))
	var sumProductionCapacity = 0
	for _, tile := range producingTiles {
		sumProductionCapacity = sumProductionCapacity + tile.Quality
	}

	return sumProductionCapacity
}

func (manager *TownManager) foodProductionCapacity() int {
	var producingTiles = ti.ThatArePlains(ti.WithHouses(manager.Town.GetTiles()))
	var sumProductionCapacity = 0
	for _, tile := range producingTiles {
		sumProductionCapacity = sumProductionCapacity + tile.Quality
	}

	return sumProductionCapacity
}

func (manager *TownManager) rateOfWoodConsumption() int {
	return manager.Town.GetHouseCount() * to.WOOD_MAINTENANCE_PER_HOUSE
}

func (manager *TownManager) rateOfFoodConsumption() int {
	return manager.Town.Population * to.FOOD_MAINTENANCE_PER_POP
}

const TIMELINE_PROJECTION_IN_YEARS = 20

func (manager *TownManager) getWoodDemand() float64 {

	var additionalHouseCount = int(TIMELINE_PROJECTION_IN_YEARS / 4) //make a new house every 4 years
	var demand_from_maintenance = (manager.Town.GetHouseCount() + int(additionalHouseCount/2)) * 20
	fmt.Printf("demand from maintenance:%d\n", demand_from_maintenance)
	var demand_from_building = additionalHouseCount * to.WOOD_COST_PER_HOUSE
	fmt.Printf("demand from building:%d\n", demand_from_building)
	var denominator = float64(manager.Town.Wood + 1.0)
	fmt.Printf("denominator:%d\n", denominator)
	fmt.Printf("manager.Town.Wood:%d\n", manager.Town.Wood)
	var woodDemand = (float64(demand_from_maintenance) + float64(demand_from_building)) / denominator
	fmt.Printf("demand for wood%d\n", woodDemand)
	return float64(woodDemand)
}

func (manager *TownManager) getFoodDemand() float64 {
	var additionalPopulationCount = TIMELINE_PROJECTION_IN_YEARS / 4 //make a new person every 4 years
	var demand_from_maintenance = ((manager.Town.Population + additionalPopulationCount) * to.FOOD_MAINTENANCE_PER_POP) * TIMELINE_PROJECTION_IN_YEARS
	fmt.Printf("Food demand from maintenance: %d\n", demand_from_maintenance)
	var foodDemand = (float64(demand_from_maintenance) + 1.0) / (float64(manager.Town.Food) + 1.0)

	return float64(foodDemand)
}

func (manager *TownManager) TakeActions() {

	//---Build buildings

	houseCount := manager.Town.GetHouseCount()
	//wood_needs_for_next_30_years := houseCount * 30 * to.WOOD_MAINTENANCE_PER_HOUSE

	var counter = 0

	//Ok to build houses if...
	for houseCount < manager.Town.Population*2 && manager.Town.Wood >= to.WOOD_COST_PER_HOUSE {
		counter = counter + 1
		houseCount = manager.Town.GetHouseCount()
		//wood_needs_for_next_30_years = houseCount * 30 * to.WOOD_MAINTENANCE_PER_HOUSE

		//calculate food demand & wood demand for next 100 years, assuming 1.5x population

		hundred_year_food_need := int(float64(manager.Town.Population) * 1.25 * 100)
		hundred_year_wood_need := int(float64(manager.Town.Population)*1.25*100 + (float64(manager.Town.Population) * .5 * to.WOOD_COST_PER_HOUSE))

		lt_wood_demand := hundred_year_wood_need / (manager.Town.Wood + 10)
		lt_food_demand := hundred_year_food_need / (manager.Town.Food + 10)

		//build a house on best land for that demand.

		var bestTilesToBuildOn []*ti.Tile

		bestTilesToBuildOn = ti.WithoutHouses(manager.Town.GetTiles())

		ti.SortByDemandAdjustedQualityInPlace(bestTilesToBuildOn, float64(lt_food_demand), float64(lt_wood_demand))

		//Build house
		if len(bestTilesToBuildOn) > 0 {
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

	new_wood_maintain := (20 / 2) * to.WOOD_MAINTENANCE_PER_HOUSE
	new_wood_build := int(math.Max((float64(20-unoccupiedHouses(houseCount, manager.Town.Population)) * to.WOOD_COST_PER_HOUSE), 0.0))
	st_wood_demand := existing_wood_maintain + new_wood_maintain + new_wood_build
	wood_cost = float64(st_wood_demand) / float64(manager.Town.Wood+10)

	new_food_maintain := (20 / 2) * to.FOOD_MAINTENANCE_PER_POP
	st_food_demand := existing_food_maintain + new_food_maintain
	food_cost = float64(st_food_demand) / float64(manager.Town.Food+10)

	pop_unallocated := manager.Town.Population

	var allHouses = ti.WithHouses(manager.Town.GetTiles())
	ti.SortByDemandAdjustedQualityInPlace(allHouses, food_cost, wood_cost)
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

//TODO: Make houses and pop global?
func unoccupiedHouses(houses int, pop int) int {
	return int(math.Max(float64(houses-pop), 0))
}
