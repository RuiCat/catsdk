// Code generated by build.go. DO NOT EDIT.

package gio

import (
	_ "embed"
	"runtime"

	"gioui/shader"
)

var (
	Shader_blit_frag = [...]shader.Sources{
		{
			Name:   "blit.frag",
			Inputs: []shader.InputLocation{{Name: "vUV", Location: 0, Semantic: "TEXCOORD", SemanticIndex: 0, Type: 0x0, Size: 2}, {Name: "opacity", Location: 1, Semantic: "TEXCOORD", SemanticIndex: 1, Type: 0x0, Size: 1}},
			Uniforms: shader.UniformsReflection{
				Locations: []shader.UniformLocation{{Name: "_color.color", Type: 0x0, Size: 4, Offset: 112}},
				Size:      16,
			},
		},
		{
			Name:   "blit.frag",
			Inputs: []shader.InputLocation{{Name: "vUV", Location: 0, Semantic: "TEXCOORD", SemanticIndex: 0, Type: 0x0, Size: 2}, {Name: "opacity", Location: 1, Semantic: "TEXCOORD", SemanticIndex: 1, Type: 0x0, Size: 1}},
			Uniforms: shader.UniformsReflection{
				Locations: []shader.UniformLocation{{Name: "_gradient.color1", Type: 0x0, Size: 4, Offset: 96}, {Name: "_gradient.color2", Type: 0x0, Size: 4, Offset: 112}},
				Size:      32,
			},
		},
		{
			Name:     "blit.frag",
			Inputs:   []shader.InputLocation{{Name: "vUV", Location: 0, Semantic: "TEXCOORD", SemanticIndex: 0, Type: 0x0, Size: 2}, {Name: "opacity", Location: 1, Semantic: "TEXCOORD", SemanticIndex: 1, Type: 0x0, Size: 1}},
			Textures: []shader.TextureBinding{{Name: "tex", Binding: 0}},
		},
	}
	//go:embed zblit.frag.0.spirv
	zblit_frag_0_spirv string
	//go:embed zblit.frag.0.glsl100es
	zblit_frag_0_glsl100es string
	//go:embed zblit.frag.0.glsl150
	zblit_frag_0_glsl150 string
	//go:embed zblit.frag.0.dxbc
	zblit_frag_0_dxbc string
	//go:embed zblit.frag.1.spirv
	zblit_frag_1_spirv string
	//go:embed zblit.frag.1.glsl100es
	zblit_frag_1_glsl100es string
	//go:embed zblit.frag.1.glsl150
	zblit_frag_1_glsl150 string
	//go:embed zblit.frag.1.dxbc
	zblit_frag_1_dxbc string
	//go:embed zblit.frag.2.spirv
	zblit_frag_2_spirv string
	//go:embed zblit.frag.2.glsl100es
	zblit_frag_2_glsl100es string
	//go:embed zblit.frag.2.glsl150
	zblit_frag_2_glsl150 string
	//go:embed zblit.frag.2.dxbc
	zblit_frag_2_dxbc string
	Shader_blit_vert  = shader.Sources{
		Name:   "blit.vert",
		Inputs: []shader.InputLocation{{Name: "pos", Location: 0, Semantic: "TEXCOORD", SemanticIndex: 0, Type: 0x0, Size: 2}, {Name: "uv", Location: 1, Semantic: "TEXCOORD", SemanticIndex: 1, Type: 0x0, Size: 2}},
		Uniforms: shader.UniformsReflection{
			Locations: []shader.UniformLocation{{Name: "_block.transform", Type: 0x0, Size: 4, Offset: 0}, {Name: "_block.uvTransformR1", Type: 0x0, Size: 4, Offset: 16}, {Name: "_block.uvTransformR2", Type: 0x0, Size: 4, Offset: 32}, {Name: "_block.opacity", Type: 0x0, Size: 1, Offset: 48}, {Name: "_block.fbo", Type: 0x0, Size: 1, Offset: 52}},
			Size:      56,
		},
	}
	//go:embed zblit.vert.0.spirv
	zblit_vert_0_spirv string
	//go:embed zblit.vert.0.glsl100es
	zblit_vert_0_glsl100es string
	//go:embed zblit.vert.0.glsl150
	zblit_vert_0_glsl150 string
	//go:embed zblit.vert.0.dxbc
	zblit_vert_0_dxbc string
	Shader_copy_frag  = shader.Sources{
		Name:     "copy.frag",
		Inputs:   []shader.InputLocation{{Name: "vUV", Location: 0, Semantic: "TEXCOORD", SemanticIndex: 0, Type: 0x0, Size: 2}},
		Textures: []shader.TextureBinding{{Name: "tex", Binding: 0}},
	}
	//go:embed zcopy.frag.0.spirv
	zcopy_frag_0_spirv string
	//go:embed zcopy.frag.0.glsl100es
	zcopy_frag_0_glsl100es string
	//go:embed zcopy.frag.0.glsl150
	zcopy_frag_0_glsl150 string
	//go:embed zcopy.frag.0.dxbc
	zcopy_frag_0_dxbc string
	Shader_copy_vert  = shader.Sources{
		Name:   "copy.vert",
		Inputs: []shader.InputLocation{{Name: "pos", Location: 0, Semantic: "TEXCOORD", SemanticIndex: 0, Type: 0x0, Size: 2}, {Name: "uv", Location: 1, Semantic: "TEXCOORD", SemanticIndex: 1, Type: 0x0, Size: 2}},
		Uniforms: shader.UniformsReflection{
			Locations: []shader.UniformLocation{{Name: "_block.scale", Type: 0x0, Size: 2, Offset: 0}, {Name: "_block.pos", Type: 0x0, Size: 2, Offset: 8}, {Name: "_block.uvScale", Type: 0x0, Size: 2, Offset: 16}},
			Size:      24,
		},
	}
	//go:embed zcopy.vert.0.spirv
	zcopy_vert_0_spirv string
	//go:embed zcopy.vert.0.glsl100es
	zcopy_vert_0_glsl100es string
	//go:embed zcopy.vert.0.glsl150
	zcopy_vert_0_glsl150 string
	//go:embed zcopy.vert.0.dxbc
	zcopy_vert_0_dxbc string
	Shader_cover_frag = [...]shader.Sources{
		{
			Name:   "cover.frag",
			Inputs: []shader.InputLocation{{Name: "vCoverUV", Location: 0, Semantic: "TEXCOORD", SemanticIndex: 0, Type: 0x0, Size: 2}, {Name: "vUV", Location: 1, Semantic: "TEXCOORD", SemanticIndex: 1, Type: 0x0, Size: 2}},
			Uniforms: shader.UniformsReflection{
				Locations: []shader.UniformLocation{{Name: "_color.color", Type: 0x0, Size: 4, Offset: 112}},
				Size:      16,
			},
			Textures: []shader.TextureBinding{{Name: "cover", Binding: 1}},
		},
		{
			Name:   "cover.frag",
			Inputs: []shader.InputLocation{{Name: "vCoverUV", Location: 0, Semantic: "TEXCOORD", SemanticIndex: 0, Type: 0x0, Size: 2}, {Name: "vUV", Location: 1, Semantic: "TEXCOORD", SemanticIndex: 1, Type: 0x0, Size: 2}},
			Uniforms: shader.UniformsReflection{
				Locations: []shader.UniformLocation{{Name: "_gradient.color1", Type: 0x0, Size: 4, Offset: 96}, {Name: "_gradient.color2", Type: 0x0, Size: 4, Offset: 112}},
				Size:      32,
			},
			Textures: []shader.TextureBinding{{Name: "cover", Binding: 1}},
		},
		{
			Name:     "cover.frag",
			Inputs:   []shader.InputLocation{{Name: "vCoverUV", Location: 0, Semantic: "TEXCOORD", SemanticIndex: 0, Type: 0x0, Size: 2}, {Name: "vUV", Location: 1, Semantic: "TEXCOORD", SemanticIndex: 1, Type: 0x0, Size: 2}},
			Textures: []shader.TextureBinding{{Name: "tex", Binding: 0}, {Name: "cover", Binding: 1}},
		},
	}
	//go:embed zcover.frag.0.spirv
	zcover_frag_0_spirv string
	//go:embed zcover.frag.0.glsl100es
	zcover_frag_0_glsl100es string
	//go:embed zcover.frag.0.glsl150
	zcover_frag_0_glsl150 string
	//go:embed zcover.frag.0.dxbc
	zcover_frag_0_dxbc string
	//go:embed zcover.frag.1.spirv
	zcover_frag_1_spirv string
	//go:embed zcover.frag.1.glsl100es
	zcover_frag_1_glsl100es string
	//go:embed zcover.frag.1.glsl150
	zcover_frag_1_glsl150 string
	//go:embed zcover.frag.1.dxbc
	zcover_frag_1_dxbc string
	//go:embed zcover.frag.2.spirv
	zcover_frag_2_spirv string
	//go:embed zcover.frag.2.glsl100es
	zcover_frag_2_glsl100es string
	//go:embed zcover.frag.2.glsl150
	zcover_frag_2_glsl150 string
	//go:embed zcover.frag.2.dxbc
	zcover_frag_2_dxbc string
	Shader_cover_vert  = shader.Sources{
		Name:   "cover.vert",
		Inputs: []shader.InputLocation{{Name: "pos", Location: 0, Semantic: "TEXCOORD", SemanticIndex: 0, Type: 0x0, Size: 2}, {Name: "uv", Location: 1, Semantic: "TEXCOORD", SemanticIndex: 1, Type: 0x0, Size: 2}},
		Uniforms: shader.UniformsReflection{
			Locations: []shader.UniformLocation{{Name: "_block.transform", Type: 0x0, Size: 4, Offset: 0}, {Name: "_block.uvCoverTransform", Type: 0x0, Size: 4, Offset: 16}, {Name: "_block.uvTransformR1", Type: 0x0, Size: 4, Offset: 32}, {Name: "_block.uvTransformR2", Type: 0x0, Size: 4, Offset: 48}, {Name: "_block.fbo", Type: 0x0, Size: 1, Offset: 64}},
			Size:      68,
		},
	}
	//go:embed zcover.vert.0.spirv
	zcover_vert_0_spirv string
	//go:embed zcover.vert.0.glsl100es
	zcover_vert_0_glsl100es string
	//go:embed zcover.vert.0.glsl150
	zcover_vert_0_glsl150 string
	//go:embed zcover.vert.0.dxbc
	zcover_vert_0_dxbc string
	Shader_input_vert  = shader.Sources{
		Name:   "input.vert",
		Inputs: []shader.InputLocation{{Name: "inPos", Location: 0, Semantic: "TEXCOORD", SemanticIndex: 0, Type: 0x0, Size: 3}, {Name: "inColor", Location: 1, Semantic: "TEXCOORD", SemanticIndex: 1, Type: 0x0, Size: 3}, {Name: "inUV", Location: 2, Semantic: "TEXCOORD", SemanticIndex: 2, Type: 0x0, Size: 2}, {Name: "innormal", Location: 3, Semantic: "TEXCOORD", SemanticIndex: 3, Type: 0x0, Size: 3}},
	}
	//go:embed zinput.vert.0.spirv
	zinput_vert_0_spirv string
	//go:embed zinput.vert.0.glsl100es
	zinput_vert_0_glsl100es string
	//go:embed zinput.vert.0.glsl150
	zinput_vert_0_glsl150 string
	//go:embed zinput.vert.0.dxbc
	zinput_vert_0_dxbc    string
	Shader_intersect_frag = shader.Sources{
		Name:     "intersect.frag",
		Inputs:   []shader.InputLocation{{Name: "vUV", Location: 0, Semantic: "TEXCOORD", SemanticIndex: 0, Type: 0x0, Size: 2}},
		Textures: []shader.TextureBinding{{Name: "cover", Binding: 0}},
	}
	//go:embed zintersect.frag.0.spirv
	zintersect_frag_0_spirv string
	//go:embed zintersect.frag.0.glsl100es
	zintersect_frag_0_glsl100es string
	//go:embed zintersect.frag.0.glsl150
	zintersect_frag_0_glsl150 string
	//go:embed zintersect.frag.0.dxbc
	zintersect_frag_0_dxbc string
	Shader_intersect_vert  = shader.Sources{
		Name:   "intersect.vert",
		Inputs: []shader.InputLocation{{Name: "pos", Location: 0, Semantic: "TEXCOORD", SemanticIndex: 0, Type: 0x0, Size: 2}, {Name: "uv", Location: 1, Semantic: "TEXCOORD", SemanticIndex: 1, Type: 0x0, Size: 2}},
		Uniforms: shader.UniformsReflection{
			Locations: []shader.UniformLocation{{Name: "_block.uvTransform", Type: 0x0, Size: 4, Offset: 0}, {Name: "_block.subUVTransform", Type: 0x0, Size: 4, Offset: 16}},
			Size:      32,
		},
	}
	//go:embed zintersect.vert.0.spirv
	zintersect_vert_0_spirv string
	//go:embed zintersect.vert.0.glsl100es
	zintersect_vert_0_glsl100es string
	//go:embed zintersect.vert.0.glsl150
	zintersect_vert_0_glsl150 string
	//go:embed zintersect.vert.0.dxbc
	zintersect_vert_0_dxbc string
	Shader_material_frag   = shader.Sources{
		Name:   "material.frag",
		Inputs: []shader.InputLocation{{Name: "vUV", Location: 0, Semantic: "TEXCOORD", SemanticIndex: 0, Type: 0x0, Size: 2}},
		Uniforms: shader.UniformsReflection{
			Locations: []shader.UniformLocation{{Name: "_color.emulateSRGB", Type: 0x0, Size: 1, Offset: 16}},
			Size:      4,
		},
		Textures: []shader.TextureBinding{{Name: "tex", Binding: 0}},
	}
	//go:embed zmaterial.frag.0.spirv
	zmaterial_frag_0_spirv string
	//go:embed zmaterial.frag.0.glsl100es
	zmaterial_frag_0_glsl100es string
	//go:embed zmaterial.frag.0.glsl150
	zmaterial_frag_0_glsl150 string
	//go:embed zmaterial.frag.0.dxbc
	zmaterial_frag_0_dxbc string
	Shader_material_vert  = shader.Sources{
		Name:   "material.vert",
		Inputs: []shader.InputLocation{{Name: "pos", Location: 0, Semantic: "TEXCOORD", SemanticIndex: 0, Type: 0x0, Size: 2}, {Name: "uv", Location: 1, Semantic: "TEXCOORD", SemanticIndex: 1, Type: 0x0, Size: 2}},
		Uniforms: shader.UniformsReflection{
			Locations: []shader.UniformLocation{{Name: "_block.scale", Type: 0x0, Size: 2, Offset: 0}, {Name: "_block.pos", Type: 0x0, Size: 2, Offset: 8}},
			Size:      16,
		},
	}
	//go:embed zmaterial.vert.0.spirv
	zmaterial_vert_0_spirv string
	//go:embed zmaterial.vert.0.glsl100es
	zmaterial_vert_0_glsl100es string
	//go:embed zmaterial.vert.0.glsl150
	zmaterial_vert_0_glsl150 string
	//go:embed zmaterial.vert.0.dxbc
	zmaterial_vert_0_dxbc string
	Shader_simple_frag    = shader.Sources{
		Name:     "simple.frag",
		Inputs:   []shader.InputLocation{{Name: "vUV", Location: 0, Semantic: "TEXCOORD", SemanticIndex: 0, Type: 0x0, Size: 2}, {Name: "vertexColor", Location: 1, Semantic: "TEXCOORD", SemanticIndex: 1, Type: 0x0, Size: 3}, {Name: "vNormal", Location: 2, Semantic: "TEXCOORD", SemanticIndex: 2, Type: 0x0, Size: 3}},
		Textures: []shader.TextureBinding{{Name: "vTexture", Binding: 0}},
	}
	//go:embed zsimple.frag.0.spirv
	zsimple_frag_0_spirv string
	//go:embed zsimple.frag.0.glsl100es
	zsimple_frag_0_glsl100es string
	//go:embed zsimple.frag.0.glsl150
	zsimple_frag_0_glsl150 string
	//go:embed zsimple.frag.0.dxbc
	zsimple_frag_0_dxbc string
	Shader_stencil_frag = shader.Sources{
		Name:   "stencil.frag",
		Inputs: []shader.InputLocation{{Name: "vFrom", Location: 0, Semantic: "TEXCOORD", SemanticIndex: 0, Type: 0x0, Size: 2}, {Name: "vCtrl", Location: 1, Semantic: "TEXCOORD", SemanticIndex: 1, Type: 0x0, Size: 2}, {Name: "vTo", Location: 2, Semantic: "TEXCOORD", SemanticIndex: 2, Type: 0x0, Size: 2}},
	}
	//go:embed zstencil.frag.0.spirv
	zstencil_frag_0_spirv string
	//go:embed zstencil.frag.0.glsl100es
	zstencil_frag_0_glsl100es string
	//go:embed zstencil.frag.0.glsl150
	zstencil_frag_0_glsl150 string
	//go:embed zstencil.frag.0.dxbc
	zstencil_frag_0_dxbc string
	Shader_stencil_vert  = shader.Sources{
		Name:   "stencil.vert",
		Inputs: []shader.InputLocation{{Name: "corner", Location: 0, Semantic: "TEXCOORD", SemanticIndex: 0, Type: 0x0, Size: 1}, {Name: "maxy", Location: 1, Semantic: "TEXCOORD", SemanticIndex: 1, Type: 0x0, Size: 1}, {Name: "from", Location: 2, Semantic: "TEXCOORD", SemanticIndex: 2, Type: 0x0, Size: 2}, {Name: "ctrl", Location: 3, Semantic: "TEXCOORD", SemanticIndex: 3, Type: 0x0, Size: 2}, {Name: "to", Location: 4, Semantic: "TEXCOORD", SemanticIndex: 4, Type: 0x0, Size: 2}},
		Uniforms: shader.UniformsReflection{
			Locations: []shader.UniformLocation{{Name: "_block.transform", Type: 0x0, Size: 4, Offset: 0}, {Name: "_block.pathOffset", Type: 0x0, Size: 2, Offset: 16}},
			Size:      24,
		},
	}
	//go:embed zstencil.vert.0.spirv
	zstencil_vert_0_spirv string
	//go:embed zstencil.vert.0.glsl100es
	zstencil_vert_0_glsl100es string
	//go:embed zstencil.vert.0.glsl150
	zstencil_vert_0_glsl150 string
	//go:embed zstencil.vert.0.dxbc
	zstencil_vert_0_dxbc string
)

func init() {
	const (
		opengles = runtime.GOOS == "linux" || runtime.GOOS == "freebsd" || runtime.GOOS == "openbsd" || runtime.GOOS == "windows" || runtime.GOOS == "js" || runtime.GOOS == "android" || runtime.GOOS == "darwin" || runtime.GOOS == "ios"
		opengl   = runtime.GOOS == "darwin"
		d3d11    = runtime.GOOS == "windows"
		vulkan   = runtime.GOOS == "linux" || runtime.GOOS == "android"
	)
	if vulkan {
		Shader_blit_frag[0].SPIRV = zblit_frag_0_spirv
	}
	if opengles {
		Shader_blit_frag[0].GLSL100ES = zblit_frag_0_glsl100es
	}
	if opengl {
		Shader_blit_frag[0].GLSL150 = zblit_frag_0_glsl150
	}
	if d3d11 {
		Shader_blit_frag[0].DXBC = zblit_frag_0_dxbc
	}
	if runtime.GOOS == "darwin" {
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
		} else {
		}
	}
	if vulkan {
		Shader_blit_frag[1].SPIRV = zblit_frag_1_spirv
	}
	if opengles {
		Shader_blit_frag[1].GLSL100ES = zblit_frag_1_glsl100es
	}
	if opengl {
		Shader_blit_frag[1].GLSL150 = zblit_frag_1_glsl150
	}
	if d3d11 {
		Shader_blit_frag[1].DXBC = zblit_frag_1_dxbc
	}
	if runtime.GOOS == "darwin" {
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
		} else {
		}
	}
	if vulkan {
		Shader_blit_frag[2].SPIRV = zblit_frag_2_spirv
	}
	if opengles {
		Shader_blit_frag[2].GLSL100ES = zblit_frag_2_glsl100es
	}
	if opengl {
		Shader_blit_frag[2].GLSL150 = zblit_frag_2_glsl150
	}
	if d3d11 {
		Shader_blit_frag[2].DXBC = zblit_frag_2_dxbc
	}
	if runtime.GOOS == "darwin" {
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
		} else {
		}
	}
	if vulkan {
		Shader_blit_vert.SPIRV = zblit_vert_0_spirv
	}
	if opengles {
		Shader_blit_vert.GLSL100ES = zblit_vert_0_glsl100es
	}
	if opengl {
		Shader_blit_vert.GLSL150 = zblit_vert_0_glsl150
	}
	if d3d11 {
		Shader_blit_vert.DXBC = zblit_vert_0_dxbc
	}
	if runtime.GOOS == "darwin" {
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
		} else {
		}
	}
	if vulkan {
		Shader_copy_frag.SPIRV = zcopy_frag_0_spirv
	}
	if opengles {
		Shader_copy_frag.GLSL100ES = zcopy_frag_0_glsl100es
	}
	if opengl {
		Shader_copy_frag.GLSL150 = zcopy_frag_0_glsl150
	}
	if d3d11 {
		Shader_copy_frag.DXBC = zcopy_frag_0_dxbc
	}
	if runtime.GOOS == "darwin" {
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
		} else {
		}
	}
	if vulkan {
		Shader_copy_vert.SPIRV = zcopy_vert_0_spirv
	}
	if opengles {
		Shader_copy_vert.GLSL100ES = zcopy_vert_0_glsl100es
	}
	if opengl {
		Shader_copy_vert.GLSL150 = zcopy_vert_0_glsl150
	}
	if d3d11 {
		Shader_copy_vert.DXBC = zcopy_vert_0_dxbc
	}
	if runtime.GOOS == "darwin" {
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
		} else {
		}
	}
	if vulkan {
		Shader_cover_frag[0].SPIRV = zcover_frag_0_spirv
	}
	if opengles {
		Shader_cover_frag[0].GLSL100ES = zcover_frag_0_glsl100es
	}
	if opengl {
		Shader_cover_frag[0].GLSL150 = zcover_frag_0_glsl150
	}
	if d3d11 {
		Shader_cover_frag[0].DXBC = zcover_frag_0_dxbc
	}
	if runtime.GOOS == "darwin" {
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
		} else {
		}
	}
	if vulkan {
		Shader_cover_frag[1].SPIRV = zcover_frag_1_spirv
	}
	if opengles {
		Shader_cover_frag[1].GLSL100ES = zcover_frag_1_glsl100es
	}
	if opengl {
		Shader_cover_frag[1].GLSL150 = zcover_frag_1_glsl150
	}
	if d3d11 {
		Shader_cover_frag[1].DXBC = zcover_frag_1_dxbc
	}
	if runtime.GOOS == "darwin" {
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
		} else {
		}
	}
	if vulkan {
		Shader_cover_frag[2].SPIRV = zcover_frag_2_spirv
	}
	if opengles {
		Shader_cover_frag[2].GLSL100ES = zcover_frag_2_glsl100es
	}
	if opengl {
		Shader_cover_frag[2].GLSL150 = zcover_frag_2_glsl150
	}
	if d3d11 {
		Shader_cover_frag[2].DXBC = zcover_frag_2_dxbc
	}
	if runtime.GOOS == "darwin" {
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
		} else {
		}
	}
	if vulkan {
		Shader_cover_vert.SPIRV = zcover_vert_0_spirv
	}
	if opengles {
		Shader_cover_vert.GLSL100ES = zcover_vert_0_glsl100es
	}
	if opengl {
		Shader_cover_vert.GLSL150 = zcover_vert_0_glsl150
	}
	if d3d11 {
		Shader_cover_vert.DXBC = zcover_vert_0_dxbc
	}
	if runtime.GOOS == "darwin" {
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
		} else {
		}
	}
	if vulkan {
		Shader_input_vert.SPIRV = zinput_vert_0_spirv
	}
	if opengles {
		Shader_input_vert.GLSL100ES = zinput_vert_0_glsl100es
	}
	if opengl {
		Shader_input_vert.GLSL150 = zinput_vert_0_glsl150
	}
	if d3d11 {
		Shader_input_vert.DXBC = zinput_vert_0_dxbc
	}
	if runtime.GOOS == "darwin" {
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
		} else {
		}
	}
	if vulkan {
		Shader_intersect_frag.SPIRV = zintersect_frag_0_spirv
	}
	if opengles {
		Shader_intersect_frag.GLSL100ES = zintersect_frag_0_glsl100es
	}
	if opengl {
		Shader_intersect_frag.GLSL150 = zintersect_frag_0_glsl150
	}
	if d3d11 {
		Shader_intersect_frag.DXBC = zintersect_frag_0_dxbc
	}
	if runtime.GOOS == "darwin" {
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
		} else {
		}
	}
	if vulkan {
		Shader_intersect_vert.SPIRV = zintersect_vert_0_spirv
	}
	if opengles {
		Shader_intersect_vert.GLSL100ES = zintersect_vert_0_glsl100es
	}
	if opengl {
		Shader_intersect_vert.GLSL150 = zintersect_vert_0_glsl150
	}
	if d3d11 {
		Shader_intersect_vert.DXBC = zintersect_vert_0_dxbc
	}
	if runtime.GOOS == "darwin" {
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
		} else {
		}
	}
	if vulkan {
		Shader_material_frag.SPIRV = zmaterial_frag_0_spirv
	}
	if opengles {
		Shader_material_frag.GLSL100ES = zmaterial_frag_0_glsl100es
	}
	if opengl {
		Shader_material_frag.GLSL150 = zmaterial_frag_0_glsl150
	}
	if d3d11 {
		Shader_material_frag.DXBC = zmaterial_frag_0_dxbc
	}
	if runtime.GOOS == "darwin" {
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
		} else {
		}
	}
	if vulkan {
		Shader_material_vert.SPIRV = zmaterial_vert_0_spirv
	}
	if opengles {
		Shader_material_vert.GLSL100ES = zmaterial_vert_0_glsl100es
	}
	if opengl {
		Shader_material_vert.GLSL150 = zmaterial_vert_0_glsl150
	}
	if d3d11 {
		Shader_material_vert.DXBC = zmaterial_vert_0_dxbc
	}
	if runtime.GOOS == "darwin" {
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
		} else {
		}
	}
	if vulkan {
		Shader_simple_frag.SPIRV = zsimple_frag_0_spirv
	}
	if opengles {
		Shader_simple_frag.GLSL100ES = zsimple_frag_0_glsl100es
	}
	if opengl {
		Shader_simple_frag.GLSL150 = zsimple_frag_0_glsl150
	}
	if d3d11 {
		Shader_simple_frag.DXBC = zsimple_frag_0_dxbc
	}
	if runtime.GOOS == "darwin" {
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
		} else {
		}
	}
	if vulkan {
		Shader_stencil_frag.SPIRV = zstencil_frag_0_spirv
	}
	if opengles {
		Shader_stencil_frag.GLSL100ES = zstencil_frag_0_glsl100es
	}
	if opengl {
		Shader_stencil_frag.GLSL150 = zstencil_frag_0_glsl150
	}
	if d3d11 {
		Shader_stencil_frag.DXBC = zstencil_frag_0_dxbc
	}
	if runtime.GOOS == "darwin" {
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
		} else {
		}
	}
	if vulkan {
		Shader_stencil_vert.SPIRV = zstencil_vert_0_spirv
	}
	if opengles {
		Shader_stencil_vert.GLSL100ES = zstencil_vert_0_glsl100es
	}
	if opengl {
		Shader_stencil_vert.GLSL150 = zstencil_vert_0_glsl150
	}
	if d3d11 {
		Shader_stencil_vert.DXBC = zstencil_vert_0_dxbc
	}
	if runtime.GOOS == "darwin" {
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
		} else {
		}
	}
}
