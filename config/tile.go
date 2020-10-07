package config

import "github.com/paulwrubel/photolum/config/geometry"

type Tile struct {
	Origin geometry.Point  // Top left corner of Tile
	Span   geometry.Vector // Width and Height of Tile
}
