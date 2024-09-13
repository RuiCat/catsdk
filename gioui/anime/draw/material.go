package draw

// Material 是用于定义材质的结构体
type Material struct {
	// 色彩值
	Color Color
	// 纹理图像
	Texture Texture
	// 正常纹理图像
	NormalTexture Texture
	// 凸起纹理图像
	BumpTexture Texture
	// 高光纹理图像
	GlossTexture Texture
	// 凸起乘数
	BumpMultiplier float64
	// 发光度
	Emittance float32
	// refractive index 表示折射率
	Index float32
	// 反射锥角（单位：弧度）
	Gloss float32
	// specular 和 refractive 的色彩染色值
	Tint float32
	//金属反射系数
	Reflectivity float64
	// 是否透明
	Transparent bool
}

func DiffuseMaterial(color Color) Material {
	return Material{color, nil, nil, nil, nil, 1, 0, 1, 0, 0, -1, false}
}

func SpecularMaterial(color Color, index float32) Material {
	return Material{color, nil, nil, nil, nil, 1, 0, index, 0, 0, -1, false}
}

func GlossyMaterial(color Color, index, gloss float32) Material {
	return Material{color, nil, nil, nil, nil, 1, 0, index, gloss, 0, -1, false}
}

func ClearMaterial(index float32, gloss float32) Material {
	return Material{Black, nil, nil, nil, nil, 1, 0, index, gloss, 0, -1, true}
}

func TransparentMaterial(color Color, index, gloss, tint float32) Material {
	return Material{color, nil, nil, nil, nil, 1, 0, index, gloss, tint, -1, true}
}

func MetallicMaterial(color Color, gloss, tint float32) Material {
	return Material{color, nil, nil, nil, nil, 1, 0, 1, gloss, tint, 1, false}
}

func LightMaterial(color Color, emittance float32) Material {
	return Material{color, nil, nil, nil, nil, 1, emittance, 1, 0, 0, -1, false}
}

func MaterialAt(shape Shape, point Vector) Material {
	material := shape.MaterialAt(point)
	uv := shape.UV(point)
	if material.Texture != nil {
		material.Color = material.Texture.Sample(uv.X, uv.Y)
	}
	if material.GlossTexture != nil {
		c := material.GlossTexture.Sample(uv.X, uv.Y)
		material.Gloss = (c.R + c.G + c.B) / 3
	}
	return material
}
