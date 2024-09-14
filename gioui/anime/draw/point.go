package draw

import "mat/mat/spatial/r2"

type Point = r2.Point

// Point3D 3D点
type Point3D struct {
	// 包围盒
	box Box
	// 点位置
	Point Vector
	// 法线向量
	Normal Vector
	// 材质属性
	Material Material
}

// NewPoint3D 创建一个新的 Point3D 对象。
// point 点位置
// normal 法线向量
// material 材质属性
func NewPoint3D(point, normal Vector, material Material) *Point3D {
	box := Box{Min: point.AddScalar(-0.05), Max: point.AddScalar(0.05)}
	normal = normal.Normalize()
	return &Point3D{box, point, normal, material}
}

func (p *Point3D) Compile() {}

// BoundingBox 获取包围盒
func (p *Point3D) BoundingBox() Box {
	return p.box
}

// Intersect 检查射线是否与点相交
func (p *Point3D) Intersect(ray Ray) Hit {
	to := ray.Origin.Sub(p.Point)
	b := to.Dot(ray.Direction)
	c := b*b - to.Dot(to)
	if uint(c) == 0 {
		return Hit{p, c, nil}
	}
	return NoHit
}

// 计算 UV 坐标
func (p *Point3D) UV(a Vector) Vector {
	return Vector{}
}

// MaterialAt 获取材质属性
func (p *Point3D) MaterialAt(a Vector) Material {
	return p.Material
}

// NormalAt 获取法线向量
func (p *Point3D) NormalAt(a Vector) Vector {
	return p.Normal
}
