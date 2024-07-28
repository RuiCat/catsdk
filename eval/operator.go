package eval

import (
	"fmt"
	"reflect"
)

var actionsOperator = map[action]string{
	aAdd:     "Add",     // +
	aMul:     "Mul",     // *
	aQuo:     "Quo",     // /
	aSub:     "Sub",     // -
	aAnd:     "And",     // &
	aXor:     "Xor",     // ^
	aGreater: "Greater", // >
	aLower:   "Lower",   // <
	aShl:     "Shl",     // <<
	aShr:     "Shr",     // >>
}

func (check typecheck) operator(n *node) bool {
	t := n.child[0].typ.TypeOf()
	switch t.Kind() {
	case reflect.Array, reflect.Struct, reflect.Interface:
		if name, ok := actionsOperator[n.action]; ok {
			fnext := n.child[0].typ.getMethod(name)
			if fnext == nil {
				return false
			}
			n.fnext = fnext
			n.gen = operator
			return true
		}
	}
	return false
}

func operator(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.concrete().TypeOf()
	dest := genValueOutput(n, typ)
	vb := genValue(n.child[1])
	fn := &node{}
	fn.child = []*node{n.child[0], n.child[0].typ.getMethod(actionsOperator[n.action])}
	fn.typ = n.child[1].typ
	fn.recv = n.child[1].recv
	fn.val = n.child[0]
	fn.gen = nop
	fn.action = aGetSym
	if err := matchSelectorMethod(n.scope, fn); err == nil {
		vfn := genFunctionWrapper(fn)
		(*n).exec = func(f *frame) bltn {
			dest(f).Set(vfn(f).Call([]reflect.Value{vb(f)})[0])
			return next
		}
	} else {
		panic(fmt.Sprintf("operator %s %s", actionsOperator[n.action], err))
	}
}
