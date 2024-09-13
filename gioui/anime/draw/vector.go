package draw

import (
	"mat/mat/spatial/r3"
	"math/rand"
)

type Vector = r3.Vec

func V(x, y, z float64) Vector {
	return Vector{X: x, Y: y, Z: z}
}

func RandomUnitVector(rnd *rand.Rand) Vector {
	for {
		var x, y, z float64
		if rnd == nil {
			x = rand.Float64()*2 - 1
			y = rand.Float64()*2 - 1
			z = rand.Float64()*2 - 1
		} else {
			x = rnd.Float64()*2 - 1
			y = rnd.Float64()*2 - 1
			z = rnd.Float64()*2 - 1
		}
		if x*x+y*y+z*z > 1 {
			continue
		}
		return Vector{X: x, Y: y, Z: z}.Normalize()
	}
}
