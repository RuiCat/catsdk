// Code generated by gioui.org/cpu/cmd/compile DO NOT EDIT.

//go:build linux && (arm64 || arm || amd64)
// +build linux
// +build arm64 arm amd64

package piet

import "gioui/cpu"
import "unsafe"

/*
#cgo LDFLAGS: -lm

#include <stdint.h>
#include <stdlib.h>
#include "abi.h"
#include "runtime.h"
#include "path_coarse_abi.h"
*/
import "C"

var Path_coarseProgramInfo = (*cpu.ProgramInfo)(unsafe.Pointer(&C.path_coarse_program_info))

type Path_coarseDescriptorSetLayout = C.struct_path_coarse_descriptor_set_layout

const Path_coarseHash = "ed67e14c880cf92bdd7a9d520610e8c8b139907ff8b55df20464d353a7f58e79"

func (l *Path_coarseDescriptorSetLayout) Binding0() *cpu.BufferDescriptor {
	return (*cpu.BufferDescriptor)(unsafe.Pointer(&l.binding0))
}

func (l *Path_coarseDescriptorSetLayout) Binding1() *cpu.BufferDescriptor {
	return (*cpu.BufferDescriptor)(unsafe.Pointer(&l.binding1))
}
