package town_manager

import (
)

func HowMuchWoodForXFood(amount, town) int {
	if (amount <= 0) {
		return 0
	}
	
	var woodDemand = getWoodDemand(town)
	var foodDemand = getFoodDemand(town)

	var exchangeRate = woodDemand / foodDemand

	var offer = math.Floor(amount * exchangeRate)

	if (offer > amount * 20) {
		offer = amount * 20
	}

	offer = math.Max(offer, town.Wood + 1)

	return offer
}

func HowMuchFoodForXWood(amount, town) int {
	if (amount <= 0) {
		return 0
	}
	
	var woodDemand = getWoodDemand(town)
	var foodDemand = getFoodDemand(town)

	var exchangeRate = foodDemand / woodDemand

	var offer = math.Floor(amount * exchangeRate)

	if (offer > amount * 20) {
		offer = amount * 20
	}

	offer = math.Max(offer, town.Food + 1)

	return offer
}

const TIMELINE_PROJECTION_IN_YEARS = 20

func getWoodDemand(town) {

	var additionalHouseCount = TIMELINE_PROJECTION_IN_YEARS / 4 //make a new house every 4 years
	demand_from_maintenance = (town.getHouseCount() + additionalHouseCount / 2) 
	demand_from_building = additionalHouseCount * WOOD_COST_PER_HOUSE
	var woodDemand = (demand_from_maintenance + demand_from_building) / (town.Wood + 1)

	return woodDemand
}

func getFoodDemand(town) {
	var additionalPopulationCount = TIMELINE_PROJECTION_IN_YEARS / 4 //make a new person every 4 years
	var demand_from_maintenance = (town.Population * town.FOOD_MAINTENANCE_PER_POP)
	var foodDemand = (demand_from_maintenance + 1) / (town.Food + 1)

	return foodDemand
}

func main() {
		//---Build buildings

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

		//---Allocate labor

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