package town_manager

import (
)

func main() {

		//---Build buildings

		//---Allocate labor
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