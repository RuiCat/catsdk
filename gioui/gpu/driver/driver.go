// SPDX-License-Identifier: Unlicense OR MIT

package driver

import (
	"errors"
	"image"
	"time"

	"gioui/shader"
	"mat/mat/spatial/f32color"
)

// Device 代表GPU驱动程序的抽象接口，提供了多种GPU API（如OpenGL、Direct3D）用于Gio渲染操作。
type Device interface {
	// 开始一个新的绘制帧，并清空背景，如果clear为true则清空背景。
	BeginFrame(target RenderTarget, clear bool, viewport image.Point) Texture // 开始新一帧的渲染，指定目标、是否清空背景和视口。
	// 结束当前帧的渲染操作。
	EndFrame()
	// 获取当前设备的支持能力。
	Caps() Caps
	// 创建一个新的计时器，用于计时。
	NewTimer() Timer
	// 是否所有计时器测量都有效。
	IsTimeContinuous() bool // 判断是否所有计时器测量都有效。
	// 创建纹理
	NewTexture(format TextureFormat, width, height int, minFilter, magFilter TextureFilter, bindings BufferBinding) (Texture, error)
	// 创建一个不可改变的缓存区（Immutable Buffer），指定类型和数据。
	NewImmutableBuffer(typ BufferBinding, data []byte) (Buffer, error)
	// 创建一个可变化的缓存区（Mutable Buffer），指定类型和大小。
	NewBuffer(typ BufferBinding, size int) (Buffer, error)
	// 创建一个计算程序（Compute Program），指定程序源代码。
	NewComputeProgram(shader shader.Sources) (Program, error)
	// 创建一个顶点着色器（Vertex Shader），指定程序源代码。
	NewVertexShader(src shader.Sources) (VertexShader, error)
	// 创建一个片段着色器（Fragment Shader），指定程序源代码。
	NewFragmentShader(src shader.Sources) (FragmentShader, error)
	// 创建一个管线描述符，指定管线参数。
	NewPipeline(desc PipelineDesc) (Pipeline, error)
	// 设置视口，指定x、y、宽度和高度。
	Viewport(x, y, width, height int)
	// 绘制数组操作，指定偏移和数量。
	DrawArrays(off, count int)
	// 绘制元素操作，指定偏移和数量。
	DrawElements(off, count int)
	// 开始渲染操作，指定目标和加载描述符。
	BeginRenderPass(t Texture, desc LoadDesc)
	// 结束渲染操作。
	EndRenderPass()
	// 准备一个纹理，指定目标纹理。
	PrepareTexture(t Texture)
	// 绑定程序，指定管线。
	BindProgram(p Program)
	// 绑定管线，指定管线。
	BindPipeline(p Pipeline)
	// 绑定纹理，指定单位和纹理。
	BindTexture(unit int, t Texture)
	// 绑定缓存区，指定偏移。
	BindVertexBuffer(b Buffer, offset int)
	// 绑定索引缓存区。
	BindIndexBuffer(b Buffer)
	// 绑定图像纹理，指定单位和纹理。
	BindImageTexture(unit int, texture Texture)
	// 绑定通用缓存区，指定缓存区和偏移。
	BindUniforms(buf Buffer)
	// 绑定存储缓存区，指定绑定点和缓存区。
	BindStorageBuffer(binding int, buf Buffer)
	// 开始计算操作。
	BeginCompute()
	// 结束计算操作。
	EndCompute()
	// 复制纹理，指定目标位置、源纹理和源区域。
	CopyTexture(dst Texture, dstOrigin image.Point, src Texture, srcRect image.Rectangle)
	// 分发计算，指定x、y、z坐标。
	DispatchCompute(x, y, z int)
	// 释放设备资源。
	Release()
}

var ErrDeviceLost = errors.New("GPU device lost")

type LoadDesc struct {
	Action     LoadAction
	ClearColor f32color.RGBA
}

type Pipeline interface {
	Release()
}

type PipelineDesc struct {
	VertexShader   VertexShader
	FragmentShader FragmentShader
	VertexLayout   VertexLayout
	BlendDesc      BlendDesc
	PixelFormat    TextureFormat
	Topology       Topology
}

type VertexLayout struct {
	Inputs []InputDesc
	Stride int
}

// InputDesc describes a vertex attribute as laid out in a Buffer.
type InputDesc struct {
	Type shader.DataType
	Size int

	Offset int
}

type BlendDesc struct {
	Enable               bool
	SrcFactor, DstFactor BlendFactor
}

type BlendFactor uint8

type Topology uint8

type TextureFilter uint8
type TextureFormat uint8

type BufferBinding uint8

type LoadAction uint8

type Features uint

type Caps struct {
	// BottomLeftOrigin is true if the driver has the origin in the lower left
	// corner. The OpenGL driver returns true.
	BottomLeftOrigin bool
	Features         Features
	MaxTextureSize   int
}

type VertexShader interface {
	Release()
}

type FragmentShader interface {
	Release()
}

type Program interface {
	Release()
}

type Buffer interface {
	Release()
	Upload(data []byte)
	Download(data []byte) error
}

type Timer interface {
	Begin()
	End()
	Duration() (time.Duration, bool)
	Release()
}

type Texture interface {
	RenderTarget
	Upload(offset, size image.Point, pixels []byte, stride int)
	ReadPixels(src image.Rectangle, pixels []byte, stride int) error
	Release()
}

const (
	BufferBindingIndices BufferBinding = 1 << iota
	BufferBindingVertices
	BufferBindingUniforms
	BufferBindingTexture
	BufferBindingFramebuffer
	BufferBindingShaderStorageRead
	BufferBindingShaderStorageWrite
)

const (
	TextureFormatSRGBA TextureFormat = iota
	TextureFormatFloat
	TextureFormatRGBA8
	// TextureFormatOutput denotes the format used by the output framebuffer.
	TextureFormatOutput
)

const (
	FilterNearest TextureFilter = iota
	FilterLinear
	FilterLinearMipmapLinear
)

const (
	FeatureTimers Features = 1 << iota
	FeatureFloatRenderTargets
	FeatureCompute
	FeatureSRGB
)

const (
	TopologyTriangleStrip Topology = iota
	TopologyTriangles
)

const (
	BlendFactorOne BlendFactor = iota
	BlendFactorOneMinusSrcAlpha
	BlendFactorZero
	BlendFactorDstColor
)

const (
	LoadActionKeep LoadAction = iota
	LoadActionClear
	LoadActionInvalidate
)

var ErrContentLost = errors.New("buffer content lost")

func (f Features) Has(feats Features) bool {
	return f&feats == feats
}

func DownloadImage(d Device, t Texture, img *image.RGBA) error {
	r := img.Bounds()
	if err := t.ReadPixels(r, img.Pix, img.Stride); err != nil {
		return err
	}
	if d.Caps().BottomLeftOrigin {
		// OpenGL origin is in the lower-left corner. Flip the image to
		// match.
		flipImageY(r.Dx()*4, r.Dy(), img.Pix)
	}
	return nil
}

func flipImageY(stride, height int, pixels []byte) {
	// Flip image in y-direction. OpenGL's origin is in the lower
	// left corner.
	row := make([]uint8, stride)
	for y := 0; y < height/2; y++ {
		y1 := height - y - 1
		dest := y1 * stride
		src := y * stride
		copy(row, pixels[dest:])
		copy(pixels[dest:], pixels[src:src+len(row)])
		copy(pixels[src:], row)
	}
}

func UploadImage(t Texture, offset image.Point, img *image.RGBA) {
	var pixels []byte
	size := img.Bounds().Size()
	min := img.Rect.Min
	start := img.PixOffset(min.X, min.Y)
	end := img.PixOffset(min.X+size.X, min.Y+size.Y-1)
	pixels = img.Pix[start:end]
	t.Upload(offset, size, pixels, img.Stride)
}
