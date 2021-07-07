package tile

import (
	"sort"
)

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

func WithHouses(in_tiles []*Tile) []*Tile {
	var out_tiles []*Tile
	for _, tile := range in_tiles {
		if tile.HasHouse {
			out_tiles = append(out_tiles, tile)
		}
	}

	return out_tiles
}

func WithoutHouses(in_tiles []*Tile) []*Tile {
	var out_tiles []*Tile
	for _, tile := range in_tiles {
		if !tile.HasHouse {
			out_tiles = append(out_tiles, tile)
		}
	}

	return out_tiles
}

func ThatAreForest(in_tiles []*Tile) []*Tile {
	var out_tiles []*Tile
	for _, tile := range in_tiles {
		if tile.Type == Forest {
			out_tiles = append(out_tiles, tile)
		}
	}

	return out_tiles
}

func ThatArePlains(in_tiles []*Tile) []*Tile {
	var out_tiles []*Tile
	for _, tile := range in_tiles {
		if tile.Type == Plains {
			out_tiles = append(out_tiles, tile)
		}
	}

	return out_tiles
}

func SortByQualityInPlace(tiles []*Tile) {
	sort.SliceStable(tiles, func(i, j int) bool {
		return tiles[i].Quality > tiles[j].Quality
	})
}

func GetAdjustedTileQuality(tile *Tile, plainsDemand float64, forestDemand float64) float64 {
	if tile.Type == Plains {
		return float64(tile.Quality) * plainsDemand
	} else { //forest TODO: Make else throw real error
		return float64(tile.Quality) * forestDemand
	}
}

func SortByDemandAdjustedQualityInPlace(tiles []*Tile, plainsDemand float64, forestDemand float64) {
	sort.SliceStable(tiles, func(i, j int) bool {
		return GetAdjustedTileQuality(tiles[i], plainsDemand, forestDemand) > GetAdjustedTileQuality(tiles[j], plainsDemand, forestDemand)
	})
}
