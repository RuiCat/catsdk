package anime

// Confing 配置结构体
type Confing func(obj *Anime)

// New 创建对象
func New(obj ObjectDrawing, confing ...Confing) *Anime {
	anime := &Anime{Object: obj}
	for _, c := range confing {
		c(anime)
	}
	obj.Init()
	return anime
}

// ConfingPosition 配置绘图位置
func ConfingPosition(x, y float64) func(obj *Anime) {
	return func(obj *Anime) {
		o := obj.Object.getObject()
		o.Min.X = x
		o.Min.Y = y
		o.Min.Z = 0
	}
}

// ConfingFrame 起始帧
func ConfingFrame(start, len uint64) func(obj *Anime) {
	return func(obj *Anime) {
		obj.StartFrame = start
		obj.StopFrame = start + len
	}
}
