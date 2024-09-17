package anime

import (
	"gioui/anime/draw"
	"gioui/op/paint"
	"gioui/widget/layout"
	"image"
)

// AnimeManager 动画管理器
type AnimeManager struct {
	StopFrame uint64         // 停止帧
	Renderer  *draw.Renderer // 渲染
	Context   Context        // 动画上下文
	AnimeList []*Anime       // 动画列表
	Image     image.Image
}

// AddAnime 添加动画
func (anime *AnimeManager) Anime(node *Anime) {
	obj := node.Object.getObject()
	obj.Init()
	anime.Renderer.Scene.Add(obj)
	anime.AnimeList = append(anime.AnimeList, node)
}

// Layout 绘制
func (ui *AnimeManager) Update() {
	ui.Context.CurrentFrame++ // 帧自增
	if ui.Context.CurrentFrame > ui.StopFrame {
		ui.Context.CurrentFrame = 0
	}
	// 绘制
	cxt := &Context{}
	for _, anime := range ui.AnimeList {
		cxt.CurrentFrame = ui.Context.CurrentFrame // 设置帧
		anime.Drawing(cxt)                         // 绘制当前对象
	}
	// 渲染
	ui.Image = ui.Renderer.Render()
}

// Layout 绘制

func (ui *AnimeManager) Layout(gtx layout.Context) layout.Dimensions {
	// 渲染
	if ui.Image != nil {
		paint.NewImageOp(ui.Image).Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
	}
	return layout.Dimensions{}
}
