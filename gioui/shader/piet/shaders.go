// Code generated by build.go. DO NOT EDIT.

package piet

import (
	_ "embed"
	"runtime"

	"gioui/shader"
)

var (
	Shader_backdrop_comp = shader.Sources{
		Name:           "backdrop.comp",
		StorageBuffers: []shader.BufferBinding{{Name: "Memory", Binding: 0}, {Name: "ConfigBuf", Binding: 1}},
		WorkgroupSize:  [3]int{128, 1, 1},
	}
	//go:embed zbackdrop.comp.0.spirv
	zbackdrop_comp_0_spirv string
	//go:embed zbackdrop.comp.0.dxbc
	zbackdrop_comp_0_dxbc string
	//go:embed zbackdrop.comp.0.metallibmacos
	zbackdrop_comp_0_metallibmacos string
	//go:embed zbackdrop.comp.0.metallibios
	zbackdrop_comp_0_metallibios string
	//go:embed zbackdrop.comp.0.metallibiossimulator
	zbackdrop_comp_0_metallibiossimulator string
	Shader_binning_comp                   = shader.Sources{
		Name:           "binning.comp",
		StorageBuffers: []shader.BufferBinding{{Name: "Memory", Binding: 0}, {Name: "ConfigBuf", Binding: 1}},
		WorkgroupSize:  [3]int{128, 1, 1},
	}
	//go:embed zbinning.comp.0.spirv
	zbinning_comp_0_spirv string
	//go:embed zbinning.comp.0.dxbc
	zbinning_comp_0_dxbc string
	//go:embed zbinning.comp.0.metallibmacos
	zbinning_comp_0_metallibmacos string
	//go:embed zbinning.comp.0.metallibios
	zbinning_comp_0_metallibios string
	//go:embed zbinning.comp.0.metallibiossimulator
	zbinning_comp_0_metallibiossimulator string
	Shader_coarse_comp                   = shader.Sources{
		Name:           "coarse.comp",
		StorageBuffers: []shader.BufferBinding{{Name: "Memory", Binding: 0}, {Name: "ConfigBuf", Binding: 1}},
		WorkgroupSize:  [3]int{128, 1, 1},
	}
	//go:embed zcoarse.comp.0.spirv
	zcoarse_comp_0_spirv string
	//go:embed zcoarse.comp.0.dxbc
	zcoarse_comp_0_dxbc string
	//go:embed zcoarse.comp.0.metallibmacos
	zcoarse_comp_0_metallibmacos string
	//go:embed zcoarse.comp.0.metallibios
	zcoarse_comp_0_metallibios string
	//go:embed zcoarse.comp.0.metallibiossimulator
	zcoarse_comp_0_metallibiossimulator string
	Shader_elements_comp                = shader.Sources{
		Name:           "elements.comp",
		StorageBuffers: []shader.BufferBinding{{Name: "Memory", Binding: 0}, {Name: "SceneBuf", Binding: 2}, {Name: "StateBuf", Binding: 3}, {Name: "ConfigBuf", Binding: 1}},
		WorkgroupSize:  [3]int{32, 1, 1},
	}
	//go:embed zelements.comp.0.spirv
	zelements_comp_0_spirv string
	//go:embed zelements.comp.0.dxbc
	zelements_comp_0_dxbc string
	//go:embed zelements.comp.0.metallibmacos
	zelements_comp_0_metallibmacos string
	//go:embed zelements.comp.0.metallibios
	zelements_comp_0_metallibios string
	//go:embed zelements.comp.0.metallibiossimulator
	zelements_comp_0_metallibiossimulator string
	Shader_kernel4_comp                   = shader.Sources{
		Name:           "kernel4.comp",
		Images:         []shader.ImageBinding{{Name: "images", Binding: 3}, {Name: "image", Binding: 2}},
		StorageBuffers: []shader.BufferBinding{{Name: "Memory", Binding: 0}, {Name: "ConfigBuf", Binding: 1}},
		WorkgroupSize:  [3]int{16, 8, 1},
	}
	//go:embed zkernel4.comp.0.spirv
	zkernel4_comp_0_spirv string
	//go:embed zkernel4.comp.0.dxbc
	zkernel4_comp_0_dxbc string
	//go:embed zkernel4.comp.0.metallibmacos
	zkernel4_comp_0_metallibmacos string
	//go:embed zkernel4.comp.0.metallibios
	zkernel4_comp_0_metallibios string
	//go:embed zkernel4.comp.0.metallibiossimulator
	zkernel4_comp_0_metallibiossimulator string
	Shader_path_coarse_comp              = shader.Sources{
		Name:           "path_coarse.comp",
		StorageBuffers: []shader.BufferBinding{{Name: "Memory", Binding: 0}, {Name: "ConfigBuf", Binding: 1}},
		WorkgroupSize:  [3]int{32, 1, 1},
	}
	//go:embed zpath_coarse.comp.0.spirv
	zpath_coarse_comp_0_spirv string
	//go:embed zpath_coarse.comp.0.dxbc
	zpath_coarse_comp_0_dxbc string
	//go:embed zpath_coarse.comp.0.metallibmacos
	zpath_coarse_comp_0_metallibmacos string
	//go:embed zpath_coarse.comp.0.metallibios
	zpath_coarse_comp_0_metallibios string
	//go:embed zpath_coarse.comp.0.metallibiossimulator
	zpath_coarse_comp_0_metallibiossimulator string
	Shader_tile_alloc_comp                   = shader.Sources{
		Name:           "tile_alloc.comp",
		StorageBuffers: []shader.BufferBinding{{Name: "Memory", Binding: 0}, {Name: "ConfigBuf", Binding: 1}},
		WorkgroupSize:  [3]int{128, 1, 1},
	}
	//go:embed ztile_alloc.comp.0.spirv
	ztile_alloc_comp_0_spirv string
	//go:embed ztile_alloc.comp.0.dxbc
	ztile_alloc_comp_0_dxbc string
	//go:embed ztile_alloc.comp.0.metallibmacos
	ztile_alloc_comp_0_metallibmacos string
	//go:embed ztile_alloc.comp.0.metallibios
	ztile_alloc_comp_0_metallibios string
	//go:embed ztile_alloc.comp.0.metallibiossimulator
	ztile_alloc_comp_0_metallibiossimulator string
)

func init() {
	const (
		opengles = runtime.GOOS == "linux" || runtime.GOOS == "freebsd" || runtime.GOOS == "openbsd" || runtime.GOOS == "windows" || runtime.GOOS == "js" || runtime.GOOS == "android" || runtime.GOOS == "darwin" || runtime.GOOS == "ios"
		opengl   = runtime.GOOS == "darwin"
		d3d11    = runtime.GOOS == "windows"
		vulkan   = runtime.GOOS == "linux" || runtime.GOOS == "android"
	)
	if vulkan {
		Shader_backdrop_comp.SPIRV = zbackdrop_comp_0_spirv
	}
	if opengles {
	}
	if opengl {
	}
	if d3d11 {
		Shader_backdrop_comp.DXBC = zbackdrop_comp_0_dxbc
	}
	if runtime.GOOS == "darwin" {
		Shader_backdrop_comp.MetalLib = zbackdrop_comp_0_metallibmacos
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
			Shader_backdrop_comp.MetalLib = zbackdrop_comp_0_metallibiossimulator
		} else {
			Shader_backdrop_comp.MetalLib = zbackdrop_comp_0_metallibios
		}
	}
	if vulkan {
		Shader_binning_comp.SPIRV = zbinning_comp_0_spirv
	}
	if opengles {
	}
	if opengl {
	}
	if d3d11 {
		Shader_binning_comp.DXBC = zbinning_comp_0_dxbc
	}
	if runtime.GOOS == "darwin" {
		Shader_binning_comp.MetalLib = zbinning_comp_0_metallibmacos
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
			Shader_binning_comp.MetalLib = zbinning_comp_0_metallibiossimulator
		} else {
			Shader_binning_comp.MetalLib = zbinning_comp_0_metallibios
		}
	}
	if vulkan {
		Shader_coarse_comp.SPIRV = zcoarse_comp_0_spirv
	}
	if opengles {
	}
	if opengl {
	}
	if d3d11 {
		Shader_coarse_comp.DXBC = zcoarse_comp_0_dxbc
	}
	if runtime.GOOS == "darwin" {
		Shader_coarse_comp.MetalLib = zcoarse_comp_0_metallibmacos
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
			Shader_coarse_comp.MetalLib = zcoarse_comp_0_metallibiossimulator
		} else {
			Shader_coarse_comp.MetalLib = zcoarse_comp_0_metallibios
		}
	}
	if vulkan {
		Shader_elements_comp.SPIRV = zelements_comp_0_spirv
	}
	if opengles {
	}
	if opengl {
	}
	if d3d11 {
		Shader_elements_comp.DXBC = zelements_comp_0_dxbc
	}
	if runtime.GOOS == "darwin" {
		Shader_elements_comp.MetalLib = zelements_comp_0_metallibmacos
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
			Shader_elements_comp.MetalLib = zelements_comp_0_metallibiossimulator
		} else {
			Shader_elements_comp.MetalLib = zelements_comp_0_metallibios
		}
	}
	if vulkan {
		Shader_kernel4_comp.SPIRV = zkernel4_comp_0_spirv
	}
	if opengles {
	}
	if opengl {
	}
	if d3d11 {
		Shader_kernel4_comp.DXBC = zkernel4_comp_0_dxbc
	}
	if runtime.GOOS == "darwin" {
		Shader_kernel4_comp.MetalLib = zkernel4_comp_0_metallibmacos
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
			Shader_kernel4_comp.MetalLib = zkernel4_comp_0_metallibiossimulator
		} else {
			Shader_kernel4_comp.MetalLib = zkernel4_comp_0_metallibios
		}
	}
	if vulkan {
		Shader_path_coarse_comp.SPIRV = zpath_coarse_comp_0_spirv
	}
	if opengles {
	}
	if opengl {
	}
	if d3d11 {
		Shader_path_coarse_comp.DXBC = zpath_coarse_comp_0_dxbc
	}
	if runtime.GOOS == "darwin" {
		Shader_path_coarse_comp.MetalLib = zpath_coarse_comp_0_metallibmacos
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
			Shader_path_coarse_comp.MetalLib = zpath_coarse_comp_0_metallibiossimulator
		} else {
			Shader_path_coarse_comp.MetalLib = zpath_coarse_comp_0_metallibios
		}
	}
	if vulkan {
		Shader_tile_alloc_comp.SPIRV = ztile_alloc_comp_0_spirv
	}
	if opengles {
	}
	if opengl {
	}
	if d3d11 {
		Shader_tile_alloc_comp.DXBC = ztile_alloc_comp_0_dxbc
	}
	if runtime.GOOS == "darwin" {
		Shader_tile_alloc_comp.MetalLib = ztile_alloc_comp_0_metallibmacos
	}
	if runtime.GOOS == "ios" {
		if runtime.GOARCH == "amd64" {
			Shader_tile_alloc_comp.MetalLib = ztile_alloc_comp_0_metallibiossimulator
		} else {
			Shader_tile_alloc_comp.MetalLib = ztile_alloc_comp_0_metallibios
		}
	}
}
