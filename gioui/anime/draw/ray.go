package draw

import (
	"mat/asm/f32"
	"math"
	"math/rand"
)

// Ray 是光线结构体，包含两个成员：
type Ray struct {
	Origin    Vector // 光源位置
	Direction Vector // 光线方向
}

func (r Ray) Position(t float64) Vector {
	return r.Origin.Add(r.Direction.MulScalar(t))
}

func (n Ray) Reflect(i Ray) Ray {
	return Ray{n.Origin, n.Direction.Reflect(i.Direction)}
}

func (n Ray) Refract(i Ray, n1, n2 float32) Ray {
	return Ray{n.Origin, n.Direction.Refract(i.Direction, float64(n1), float64(n2))}
}

func (n Ray) Reflectance(i Ray, n1, n2 float32) float64 {
	return n.Direction.Reflectance(i.Direction, float64(n1), float64(n2))
}

func (r Ray) WeightedBounce(u, v float32, rnd *rand.Rand) Ray {
	radius := f32.Sqrt(u)
	theta := 2 * math.Pi * v
	s := r.Direction.Cross(RandomUnitVector(rnd)).Normalize()
	t := r.Direction.Cross(s)
	d := Vector{}
	d = d.Add(s.MulScalar(float64(radius * f32.Cos(theta))))
	d = d.Add(t.MulScalar(float64(radius * f32.Sin(theta))))
	d = d.Add(r.Direction.MulScalar(float64(f32.Sqrt(1 - u))))
	return Ray{r.Origin, d}
}

func (r Ray) ConeBounce(theta, u, v float32, rnd *rand.Rand) Ray {
	return Ray{r.Origin, Cone(r.Direction, theta, u, v, rnd)}
}

func (i Ray) Bounce(info *HitInfo, u, v float32, bounceType BounceType, rnd *rand.Rand) (Ray, bool, float32) {
	n := info.Ray
	material := info.Material
	n1, n2 := float32(1.0), material.Index
	if info.Inside {
		n1, n2 = n2, n1
	}
	var p float64
	if material.Reflectivity >= 0 {
		p = material.Reflectivity
	} else {
		p = n.Reflectance(i, n1, n2)
	}
	var reflect bool
	switch bounceType {
	case BounceTypeAny:
		reflect = rnd.Float64() < p
	case BounceTypeDiffuse:
		reflect = false
	case BounceTypeSpecular:
		reflect = true
	}
	if reflect {
		reflected := n.Reflect(i)
		return reflected.ConeBounce(material.Gloss, u, v, rnd), true, float32(p)
	} else if material.Transparent {
		refracted := n.Refract(i, n1, n2)
		refracted.Origin = refracted.Origin.Add(refracted.Direction.MulScalar(1e-4))
		return refracted.ConeBounce(material.Gloss, u, v, rnd), true, float32(1 - p)
	} else {
		return n.WeightedBounce(u, v, rnd), false, float32(1 - p)
	}
}
