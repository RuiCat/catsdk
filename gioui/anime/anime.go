package anime

import (
	"gioui/anime/draw"
	"mat/mat/spatial/r3"
)

// AnimeObject 动画对象底层
type AnimeObject struct {
	ObjectDrawing
	draw.Axis // 绘图平面
	r3.Box    // 任何动画对象均为一个基础 Box 包围盒为基础.
	r3.Mat    // 变换矩阵
	Value     map[string]any
}

func (anime *AnimeObject) Init()                    {}
func (anime *AnimeObject) GetValue(v string) any    { return anime.Value[v] }
func (anime *AnimeObject) AddValue(v string, k any) { anime.Value[v] = k }

// getObject 动画底层
func (anime *AnimeObject) getObject() *AnimeObject { return anime }

// Anime 动画
type Anime struct {
	StartFrame uint64                 // 起始帧
	StopFrame  uint64                 // 停止帧
	Object     ObjectDrawing          // 动画对象
	TransList  []ObjectTransformation // 动画变换
}

// Drawing 绘制动画
func (anime *Anime) Drawing(cxt *Context) {
	obj := anime.Object.getObject()
	if (cxt.CurrentFrame >= anime.StartFrame) && (cxt.CurrentFrame <= anime.StopFrame || anime.StopFrame == anime.StartFrame) {
		for _, trans := range anime.TransList {
			trans.Transformation(obj, cxt)
		}
		anime.Object.Drawing(cxt) // 对象绘制
	}
}

// Play 播放动画
func (anime *Anime) Play(stop uint64, width int, height int, call func(cxt *draw.Context)) {
	cxt := &Context{Context: draw.NewContext(width, height)}
	for frame := range stop {
		cxt.CurrentFrame = frame // 设置当前帧
		cxt.ClearPix()           // 清除场景
		anime.Drawing(cxt)       // 绘制当前帧
		cxt.FillPreserve()       // 绘制填充图像
		cxt.StrokePreserve()     // 绘制线条图像
		call(cxt.Context)        // 回调
	}
}
