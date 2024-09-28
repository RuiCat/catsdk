package gio

import "gioui/shader"

func init() {
	Shader_input_vert = shader.Sources{
		Name: "input.vert",
		Inputs: []shader.InputLocation{
			{Name: "inPos", Location: 0, Semantic: "TEXCOORD", SemanticIndex: 0, Type: 0x0, Size: 3},
			{Name: "inColor", Location: 1, Semantic: "TEXCOORD", SemanticIndex: 1, Type: 0x0, Size: 3},
			{Name: "inUV", Location: 2, Semantic: "TEXCOORD", SemanticIndex: 2, Type: 0x0, Size: 2}},
		Uniforms: shader.UniformsReflection{
			Locations: []shader.UniformLocation{{Name: "_block.Matrix", Type: 0x0, Size: 16, Offset: 0}},
			Size:      64,
		},
	}
	Shader_simple_frag = shader.Sources{
		Name:     "simple.frag",
		Inputs:   []shader.InputLocation{{Name: "vUV", Location: 0, Semantic: "TEXCOORD", SemanticIndex: 0, Type: 0x0, Size: 2}, {Name: "vertexColor", Location: 1, Semantic: "TEXCOORD", SemanticIndex: 1, Type: 0x0, Size: 4}},
		Textures: []shader.TextureBinding{{Name: "vTexture", Binding: 0}},
	}
}
