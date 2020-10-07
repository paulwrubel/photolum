package config

import "github.com/paulwrubel/photolum/config/geometry/primitive"

type Scene struct {
	Camera  *Camera             // Camera reference
	Objects primitive.Primitive // reference to Objects in the scene
}
