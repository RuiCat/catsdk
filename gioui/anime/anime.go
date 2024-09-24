package anime

/*
import (
	"gioui/anime/canvas"
)

// AnimeObject 动画对象底层
type AnimeObject struct {
	Width, Height int
	ObjectDrawing
	canvas.ContextDrawing
	Value map[string]any
}

// Init 初始化
func (anime *AnimeObject) Init() {
	if anime.Context == nil {
		// 得到绘图大小
		size := anime.Box.Size()
		// 绘制上下文
		anime.ContextDrawing.Context = canvas.NewContext(int(size.X), int(size.Y))
	}
}
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
		cxt.Context = obj.Context // 设置上下文
		cxt.ClearPix()            // 清除场景
		for _, trans := range anime.TransList {
			trans.Transformation(obj, cxt)
		}
		anime.Object.Drawing(cxt) // 对象绘制
		cxt.FillPreserve()        // 绘制填充图像
		cxt.StrokePreserve()      // 绘制线条图像
	}
}

// Play 播放动画
func (anime *Anime) Play(stop uint64, width int, height int, call func(cxt *Context)) {
	// 绘制
	cxt := &Context{}
	for frame := range stop {
		cxt.CurrentFrame = frame // 设置当前帧
		anime.Drawing(cxt)       // 绘制当前帧
		call(cxt)                // 回调
	}
}
*/
