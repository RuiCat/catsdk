package draw

import (
	"math"
	"math/rand"
)

// LightMode 表示光照模式，枚举值有：
type LightMode int

const (
	LightModeRandom = iota // 随机光照
	LightModeAll           // 全体光照
)

// SpecularMode 表示specular模式，枚举值有：
type SpecularMode int

const (
	SpecularModeNaive = iota // 简单镜面反射
	SpecularModeFirst        // 首次镜面反射
	SpecularModeAll          // 全体镜面反射
)

// BounceType 表示反弹类型，枚举值有：
type BounceType int

const (
	BounceTypeAny      = iota // 任意反弹
	BounceTypeDiffuse         // 散射反弹
	BounceTypeSpecular        // 镜面反弹
)

type Sampler interface {
	Sample(scene *Scene, ray Ray, rnd *rand.Rand) Color
}

func NewSampler(firstHitSamples, maxBounces int) *DefaultSampler {
	return &DefaultSampler{firstHitSamples, maxBounces, true, true, LightModeRandom, SpecularModeNaive}
}

func NewDirectSampler() *DefaultSampler {
	return &DefaultSampler{1, 0, true, false, LightModeAll, SpecularModeAll}
}

// DefaultSampler 是一个结构体，用于存储采样器的参数
type DefaultSampler struct {
	FirstHitSamples int          // 第一次采样数
	MaxBounces      int          // 最大反弹次数
	DirectLighting  bool         // 直接光照模式
	SoftShadows     bool         // 软阴影模式
	LightMode       LightMode    // 光照模式
	SpecularMode    SpecularMode // specular模式
}

func (s *DefaultSampler) Sample(scene *Scene, ray Ray, rnd *rand.Rand) Color {
	return s.sample(scene, ray, true, s.FirstHitSamples, 0, rnd)
}

func (s *DefaultSampler) sample(scene *Scene, ray Ray, emission bool, samples, depth int, rnd *rand.Rand) Color {
	// 如果当前层数超过最大反弹次数，返回材质色
	if depth > s.MaxBounces {
		return Black
	}
	// 检查ray与场景是否相交，如果没有则返回环境光
	hit, list := scene.Intersect(ray)
	if !hit.Ok() {
		return s.sampleEnvironment(scene, ray)
	}
	// 获取物体的信息和材质
	info := hit.Info(ray)
	material := info.Material
	// 如果物体有自发光，直接返回结果
	result := material.Color
	if material.Emittance > 0 {
		// 如果启用直射光并且不是自发光，返回结果
		if s.DirectLighting && !emission {
			return Black
		}
		// 计算自发光贡献
		result = result.Add(material.Color.MulScalar(material.Emittance * float32(samples)))
	}
	// 确定反弹模式
	n := int(math.Sqrt(float64(samples)))
	var ma, mb BounceType
	if s.SpecularMode == SpecularModeAll || (depth == 0 && s.SpecularMode == SpecularModeFirst) {
		ma = BounceTypeDiffuse  // 反弹为漫射
		mb = BounceTypeSpecular // 反弹为镜面
	} else {
		ma = BounceTypeAny // 反弹为任何类型
		mb = BounceTypeAny // 反弹为任何类型
	}
	// 循环计算反弹次数
	for u := 0; u < n; u++ {
		for v := 0; v < n; v++ {
			for mode := ma; mode <= mb; mode++ {
				// 计算随机采样点
				fu := (float32(u) + rnd.Float32()) / float32(n)
				fv := (float32(v) + rnd.Float32()) / float32(n)
				// 计算新的射线和反弹信息
				newRay, reflected, p := ray.Bounce(&info, fu, fv, mode, rnd)
				if mode == BounceTypeAny {
					p = 1 // 如果是任意类型，直接设置权值为1
				}
				if p > 0 {
					indirect := s.sample(scene, newRay, reflected, 1, depth+1, rnd)
					if reflected {
						// specular
						// 如果是镜面反弹，计算间接光和材质混合
						tinted := indirect.Mix(material.Color.Mul(indirect), material.Tint)
						result = result.Add(tinted.MulScalar(p))
					}
					if !reflected {
						// diffuse
						// 如果是漫射反弹，计算间接光和材质混合
						direct := Black
						if s.DirectLighting {
							// 如果启用直射光，计算直接光
							direct = s.sampleLights(scene, info.Ray, rnd)
						}
						result = result.Add(material.Color.Mul(direct.Add(indirect)).MulScalar(p))
					}
				}
			}
		}
	}
	// 处理透明像素
	color := Black
	if result.A != 1 {
		for _, hit := range list {
			if hit.Ok() {
				position := ray.Position(hit.T)
				material := MaterialAt(hit.Shape, position)
				// 混合颜色
				color = material.Color.AddColor(color)
				// 检查是否继续
				if material.Color.A == 1 {
					break
				}
			} else {
				color = color.AddColor(s.sampleEnvironment(scene, ray))
				break
			}
		}
	}
	return result.Add(color).DivScalar(float32(n * n)) // 最终结果除以采样次数
}

func (s *DefaultSampler) sampleEnvironment(scene *Scene, ray Ray) Color {
	if scene.Texture != nil {
		d := ray.Direction
		u := math.Atan2(d.Z, d.X) + scene.TextureAngle
		v := math.Atan2(d.Y, Vector{X: d.X, Y: 0, Z: d.Z}.Length())
		u = (u + math.Pi) / (2 * math.Pi)
		v = (v + math.Pi/2) / math.Pi
		return scene.Texture.Sample(u, v)
	}
	return scene.Color
}

func (s *DefaultSampler) sampleLights(scene *Scene, n Ray, rnd *rand.Rand) Color {
	nLights := len(scene.Lights)
	if nLights == 0 {
		return Black
	}

	if s.LightMode == LightModeAll {
		var result Color
		for _, light := range scene.Lights {
			result = result.Add(s.sampleLight(scene, n, rnd, light))
		}
		return result
	} else {
		// pick a random light
		light := scene.Lights[rand.Intn(nLights)]
		return s.sampleLight(scene, n, rnd, light).MulScalar(float32(nLights))
	}
}

func (s *DefaultSampler) sampleLight(scene *Scene, n Ray, rnd *rand.Rand, light Shape) Color {
	// get bounding sphere center and radius
	var center Vector
	var radius float64
	switch t := light.(type) {
	case *Sphere:
		radius = t.Radius
		center = t.Center
	default:
		// get bounding sphere from bounding box
		box := t.BoundingBox()
		radius = box.OuterRadius()
		center = box.Center()
	}

	// get random point in disk
	point := center
	if s.SoftShadows {
		for {
			x := rnd.Float64()*2 - 1
			y := rnd.Float64()*2 - 1
			if x*x+y*y <= 1 {
				l := center.Sub(n.Origin).Normalize()
				u := l.Cross(RandomUnitVector(rnd)).Normalize()
				v := l.Cross(u)
				point = Vector{}
				point = point.Add(u.MulScalar(x * radius))
				point = point.Add(v.MulScalar(y * radius))
				point = point.Add(center)
				break
			}
		}
	}

	// construct ray toward light point
	ray := Ray{n.Origin, point.Sub(n.Origin).Normalize()}

	// get cosine term
	diffuse := ray.Direction.Dot(n.Direction)
	if diffuse <= 0 {
		return Black
	}

	// check for light visibility
	hit, _ := scene.Intersect(ray)
	if !hit.Ok() || hit.Shape != light {
		return Black
	}

	// compute solid angle (hemisphere coverage)
	hyp := center.Sub(n.Origin).Length()
	opp := radius
	theta := math.Asin(opp / hyp)
	adj := opp / math.Tan(theta)
	d := math.Cos(theta) * adj
	r := math.Sin(theta) * adj
	coverage := (r * r) / (d * d)

	// TODO: fix issue where hyp < opp (point inside sphere)
	if hyp < opp {
		coverage = 1
	}
	coverage = math.Min(coverage, 1)

	// get material properties from light
	material := MaterialAt(light, point)

	// combine factors
	m := material.Emittance * float32(diffuse) * float32(coverage)
	return material.Color.MulScalar(m)
}
