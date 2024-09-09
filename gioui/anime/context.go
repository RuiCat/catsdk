package anime

import "gioui/anime/draw"

// Context 对象上下文
type Context struct {
	*draw.Context        // 绘图画布
	CurrentFrame  uint64 // 当前帧
}

// NewContext 创建绘图平面
func (cxt *Context) NewContext(axis draw.Axis, width int, height int) *draw.Context {
	return nil
}

// ObjectTransformation 底层接口用于对动画对象进行变换
type ObjectTransformation interface {
	Transformation(*AnimeObject, *Context)
}

// ObjectDrawing 底层接口用于绘制对象
type ObjectDrawing interface {
	Init()                   // 初始化
	Drawing(*Context)        // 绘制
	getObject() *AnimeObject // 底层接口
}

/*
// AnimeNode 动画节点
type AnimeNode struct {
	ObjectList               []ObjectDrawing        // 动画列表
	ObjectTransformationList []ObjectTransformation // 对每一个动画对象执行变换
}

// AddAnime 添加动画
func (anime *AnimeNode) AddAnime(object ...ObjectDrawing) {
	anime.ObjectList = append(anime.ObjectList, object...)
}

// AnimeManager 动画管理器
type AnimeManager struct {
	Context   Context      // 动画上下文
	AnimeList []*AnimeNode // 动画列表
}

// AddAnime 添加动画
func (anime *AnimeManager) AddAnime() *AnimeNode {
	node := &AnimeNode{}
	anime.AnimeList = append(anime.AnimeList, node)
	return node
}
*/
