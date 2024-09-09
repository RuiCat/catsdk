package draw

import "math"

type Cube struct {
	Min      Vector
	Max      Vector
	Material Material
	Box      Box
}

func NewCube(min, max Vector, material Material) *Cube {
	box := Box{min, max}
	return &Cube{min, max, material, box}
}

func (c *Cube) Compile() {
}

func (c *Cube) BoundingBox() Box {
	return c.Box
}

func (c *Cube) Intersect(r Ray) Hit {
	n := c.Min.Sub(r.Origin).Div(r.Direction)
	f := c.Max.Sub(r.Origin).Div(r.Direction)
	n, f = n.Min(f), n.Max(f)
	t0 := math.Max(math.Max(n.X, n.Y), n.Z)
	t1 := math.Min(math.Min(f.X, f.Y), f.Z)
	if t0 > 0 && t0 < t1 {
		return Hit{c, t0, nil}
	}
	return NoHit
}

func (c *Cube) UV(p Vector) Vector {
	p = p.Sub(c.Min).Div(c.Max.Sub(c.Min))
	return Vector{X: p.X, Y: p.Z, Z: 0}
}

func (c *Cube) MaterialAt(p Vector) Material {
	return c.Material
}

func (c *Cube) NormalAt(p Vector) Vector {
	switch {
	case p.X < c.Min.X+EPS:
		return Vector{X: -1, Y: 0, Z: 0}
	case p.X > c.Max.X-EPS:
		return Vector{X: 1, Y: 0, Z: 0}
	case p.Y < c.Min.Y+EPS:
		return Vector{X: 0, Y: -1, Z: 0}
	case p.Y > c.Max.Y-EPS:
		return Vector{X: 0, Y: 1, Z: 0}
	case p.Z < c.Min.Z+EPS:
		return Vector{X: 0, Y: 0, Z: -1}
	case p.Z > c.Max.Z-EPS:
		return Vector{X: 0, Y: 0, Z: 1}
	}
	return Vector{X: 0, Y: 1, Z: 0}
}

func (c *Cube) Mesh() *Mesh {
	a := c.Min
	b := c.Max
	z := Vector{}
	m := c.Material
	v000 := Vector{X: a.X, Y: a.Y, Z: a.Z}
	v001 := Vector{X: a.X, Y: a.Y, Z: b.Z}
	v010 := Vector{X: a.X, Y: b.Y, Z: a.Z}
	v011 := Vector{X: a.X, Y: b.Y, Z: b.Z}
	v100 := Vector{X: b.X, Y: a.Y, Z: a.Z}
	v101 := Vector{X: b.X, Y: a.Y, Z: b.Z}
	v110 := Vector{X: b.X, Y: b.Y, Z: a.Z}
	v111 := Vector{X: b.X, Y: b.Y, Z: b.Z}
	triangles := []*Triangle{
		NewTriangle(v000, v100, v110, z, z, z, m),
		NewTriangle(v000, v110, v010, z, z, z, m),
		NewTriangle(v001, v101, v111, z, z, z, m),
		NewTriangle(v001, v111, v011, z, z, z, m),
		NewTriangle(v000, v100, v101, z, z, z, m),
		NewTriangle(v000, v101, v001, z, z, z, m),
		NewTriangle(v010, v110, v111, z, z, z, m),
		NewTriangle(v010, v111, v011, z, z, z, m),
		NewTriangle(v000, v010, v011, z, z, z, m),
		NewTriangle(v000, v011, v001, z, z, z, m),
		NewTriangle(v100, v110, v111, z, z, z, m),
		NewTriangle(v100, v111, v101, z, z, z, m),
	}
	return NewMesh(triangles)
}
