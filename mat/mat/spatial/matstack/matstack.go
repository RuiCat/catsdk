// This file is generated from mgl32/matstack/matstack.go; DO NOT EDIT

package matstack

import (
	"errors"
	"mat/mat/spatial/f32"
 
)

// A MatStack is an OpenGL-style matrix stack,
// usually used for things like scenegraphs. This allows you
// to easily maintain matrix state per call level.
type MatStack  []f32.Mat4 

func NewMatStack  () *MatStack  {
	return &MatStack {f32.Ident4 ()}
}

// Push copies the top element and pushes it on the stack.
func (ms *MatStack ) Push() {
	(*ms) = append(*ms, (*ms)[len(*ms)-1])
}

// Pop removes the first element of the matrix from the stack, if there is only
// one element left there is an error.
func (ms *MatStack ) Pop() error {
	if len(*ms) == 1 {
		return errors.New("Cannot pop from mat stack, at minimum stack length of 1")
	}
	(*ms) = (*ms)[:len(*ms)-1]

	return nil
}

// RightMul multiplies the current top of the matrix by the argument.
func (ms *MatStack ) RightMul(m f32.Mat4 ) {
	(*ms)[len(*ms)-1] = (*ms)[len(*ms)-1].Mul4(m)
}

// LeftMul multiplies the current top of the matrix by the argument.
func (ms *MatStack ) LeftMul(m f32.Mat4 ) {
	(*ms)[len(*ms)-1] = m.Mul4((*ms)[len(*ms)-1])
}

// Peek returns the top element.
func (ms *MatStack ) Peek() f32.Mat4  {
	return (*ms)[len(*ms)-1]
}

// Load rewrites the top element of the stack with m
func (ms *MatStack ) Load(m f32.Mat4 ) {
	(*ms)[len(*ms)-1] = m
}

// LoadIdent is a shortcut for Load(mgl.Ident4())
func (ms *MatStack ) LoadIdent() {
	(*ms)[len(*ms)-1] = f32.Ident4 ()
}
