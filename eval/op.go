package eval

import (
	"go/constant"
	"go/token"
	"reflect"
)

// Arithmetic operators

func add(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.concrete().TypeOf()
	isInterface := n.typ.TypeOf().Kind() == reflect.Interface
	dest := genValueOutput(n, typ)
	c0, c1 := n.child[0], n.child[1]

	switch typ.Kind() {
	case reflect.String:
		switch {
		case isInterface:
			v0 := genValue(c0)
			v1 := genValue(c1)
			n.exec = func(f *frame) bltn {
				dest(f).Set(reflect.ValueOf(v0(f).String() + v1(f).String()).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			s0 := vString(c0.rval)
			v1 := genValue(c1)
			n.exec = func(f *frame) bltn {
				dest(f).SetString(s0 + v1(f).String())
				return next
			}
		case c1.rval.IsValid():
			v0 := genValue(c0)
			s1 := vString(c1.rval)
			n.exec = func(f *frame) bltn {
				dest(f).SetString(v0(f).String() + s1)
				return next
			}
		default:
			v0 := genValue(c0)
			v1 := genValue(c1)
			n.exec = func(f *frame) bltn {
				dest(f).SetString(v0(f).String() + v1(f).String())
				return next
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch {
		case isInterface:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i + j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vInt(c0.rval)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetInt(i + j)
				return next
			}
		case c1.rval.IsValid():
			v0 := genValueInt(c0)
			j := vInt(c1.rval)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetInt(i + j)
				return next
			}
		default:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetInt(i + j)
				return next
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch {
		case isInterface:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i + j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vUint(c0.rval)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetUint(i + j)
				return next
			}
		case c1.rval.IsValid():
			j := vUint(c1.rval)
			v0 := genValueUint(c0)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetUint(i + j)
				return next
			}
		default:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetUint(i + j)
				return next
			}
		}
	case reflect.Float32, reflect.Float64:
		switch {
		case isInterface:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i + j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vFloat(c0.rval)
			v1 := genValueFloat(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetFloat(i + j)
				return next
			}
		case c1.rval.IsValid():
			j := vFloat(c1.rval)
			v0 := genValueFloat(c0)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetFloat(i + j)
				return next
			}
		default:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetFloat(i + j)
				return next
			}
		}
	case reflect.Complex64, reflect.Complex128:
		switch {
		case isInterface:
			v0 := genComplex(c0)
			v1 := genComplex(c1)
			n.exec = func(f *frame) bltn {
				dest(f).Set(reflect.ValueOf(v0(f) + v1(f)).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			r0 := vComplex(c0.rval)
			v1 := genComplex(c1)
			n.exec = func(f *frame) bltn {
				dest(f).SetComplex(r0 + v1(f))
				return next
			}
		case c1.rval.IsValid():
			r1 := vComplex(c1.rval)
			v0 := genComplex(c0)
			n.exec = func(f *frame) bltn {
				dest(f).SetComplex(v0(f) + r1)
				return next
			}
		default:
			v0 := genComplex(c0)
			v1 := genComplex(c1)
			n.exec = func(f *frame) bltn {
				dest(f).SetComplex(v0(f) + v1(f))
				return next
			}
		}
	}
}

func addConst(n *node) {
	v0, v1 := n.child[0].rval, n.child[1].rval
	isConst := (v0.IsValid() && isConstantValue(v0.Type())) && (v1.IsValid() && isConstantValue(v1.Type()))
	t := n.typ.rtype
	if isConst {
		t = constVal
	}
	n.rval = reflect.New(t).Elem()
	switch {
	case isConst:
		v := constant.BinaryOp(vConstantValue(v0), token.ADD, vConstantValue(v1))
		n.rval.Set(reflect.ValueOf(v))
	case isString(t):
		n.rval.SetString(vString(v0) + vString(v1))
	case isComplex(t):
		n.rval.SetComplex(vComplex(v0) + vComplex(v1))
	case isFloat(t):
		n.rval.SetFloat(vFloat(v0) + vFloat(v1))
	case isUint(t):
		n.rval.SetUint(vUint(v0) + vUint(v1))
	case isInt(t):
		n.rval.SetInt(vInt(v0) + vInt(v1))
	}
}

func and(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.concrete().TypeOf()
	isInterface := n.typ.TypeOf().Kind() == reflect.Interface
	dest := genValueOutput(n, typ)
	c0, c1 := n.child[0], n.child[1]

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch {
		case isInterface:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i & j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vInt(c0.rval)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetInt(i & j)
				return next
			}
		case c1.rval.IsValid():
			v0 := genValueInt(c0)
			j := vInt(c1.rval)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetInt(i & j)
				return next
			}
		default:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetInt(i & j)
				return next
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch {
		case isInterface:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i & j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vUint(c0.rval)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetUint(i & j)
				return next
			}
		case c1.rval.IsValid():
			j := vUint(c1.rval)
			v0 := genValueUint(c0)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetUint(i & j)
				return next
			}
		default:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetUint(i & j)
				return next
			}
		}
	}
}

func andConst(n *node) {
	v0, v1 := n.child[0].rval, n.child[1].rval
	isConst := (v0.IsValid() && isConstantValue(v0.Type())) && (v1.IsValid() && isConstantValue(v1.Type()))
	t := n.typ.rtype
	if isConst {
		t = constVal
	}
	n.rval = reflect.New(t).Elem()
	switch {
	case isConst:
		v := constant.BinaryOp(constant.ToInt(vConstantValue(v0)), token.AND, constant.ToInt(vConstantValue(v1)))
		n.rval.Set(reflect.ValueOf(v))
	case isUint(t):
		n.rval.SetUint(vUint(v0) & vUint(v1))
	case isInt(t):
		n.rval.SetInt(vInt(v0) & vInt(v1))
	}
}

func andNot(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.concrete().TypeOf()
	isInterface := n.typ.TypeOf().Kind() == reflect.Interface
	dest := genValueOutput(n, typ)
	c0, c1 := n.child[0], n.child[1]

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch {
		case isInterface:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i &^ j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vInt(c0.rval)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetInt(i &^ j)
				return next
			}
		case c1.rval.IsValid():
			v0 := genValueInt(c0)
			j := vInt(c1.rval)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetInt(i &^ j)
				return next
			}
		default:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetInt(i &^ j)
				return next
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch {
		case isInterface:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i &^ j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vUint(c0.rval)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetUint(i &^ j)
				return next
			}
		case c1.rval.IsValid():
			j := vUint(c1.rval)
			v0 := genValueUint(c0)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetUint(i &^ j)
				return next
			}
		default:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetUint(i &^ j)
				return next
			}
		}
	}
}

func andNotConst(n *node) {
	v0, v1 := n.child[0].rval, n.child[1].rval
	isConst := (v0.IsValid() && isConstantValue(v0.Type())) && (v1.IsValid() && isConstantValue(v1.Type()))
	t := n.typ.rtype
	if isConst {
		t = constVal
	}
	n.rval = reflect.New(t).Elem()
	switch {
	case isConst:
		v := constant.BinaryOp(constant.ToInt(vConstantValue(v0)), token.AND_NOT, constant.ToInt(vConstantValue(v1)))
		n.rval.Set(reflect.ValueOf(v))
	case isUint(t):
		n.rval.SetUint(vUint(v0) &^ vUint(v1))
	case isInt(t):
		n.rval.SetInt(vInt(v0) &^ vInt(v1))
	}
}

func mul(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.concrete().TypeOf()
	isInterface := n.typ.TypeOf().Kind() == reflect.Interface
	dest := genValueOutput(n, typ)
	c0, c1 := n.child[0], n.child[1]

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch {
		case isInterface:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i * j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vInt(c0.rval)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetInt(i * j)
				return next
			}
		case c1.rval.IsValid():
			v0 := genValueInt(c0)
			j := vInt(c1.rval)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetInt(i * j)
				return next
			}
		default:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetInt(i * j)
				return next
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch {
		case isInterface:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i * j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vUint(c0.rval)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetUint(i * j)
				return next
			}
		case c1.rval.IsValid():
			j := vUint(c1.rval)
			v0 := genValueUint(c0)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetUint(i * j)
				return next
			}
		default:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetUint(i * j)
				return next
			}
		}
	case reflect.Float32, reflect.Float64:
		switch {
		case isInterface:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i * j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vFloat(c0.rval)
			v1 := genValueFloat(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetFloat(i * j)
				return next
			}
		case c1.rval.IsValid():
			j := vFloat(c1.rval)
			v0 := genValueFloat(c0)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetFloat(i * j)
				return next
			}
		default:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetFloat(i * j)
				return next
			}
		}
	case reflect.Complex64, reflect.Complex128:
		switch {
		case isInterface:
			v0 := genComplex(c0)
			v1 := genComplex(c1)
			n.exec = func(f *frame) bltn {
				dest(f).Set(reflect.ValueOf(v0(f) * v1(f)).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			r0 := vComplex(c0.rval)
			v1 := genComplex(c1)
			n.exec = func(f *frame) bltn {
				dest(f).SetComplex(r0 * v1(f))
				return next
			}
		case c1.rval.IsValid():
			r1 := vComplex(c1.rval)
			v0 := genComplex(c0)
			n.exec = func(f *frame) bltn {
				dest(f).SetComplex(v0(f) * r1)
				return next
			}
		default:
			v0 := genComplex(c0)
			v1 := genComplex(c1)
			n.exec = func(f *frame) bltn {
				dest(f).SetComplex(v0(f) * v1(f))
				return next
			}
		}
	}
}

func mulConst(n *node) {
	v0, v1 := n.child[0].rval, n.child[1].rval
	isConst := (v0.IsValid() && isConstantValue(v0.Type())) && (v1.IsValid() && isConstantValue(v1.Type()))
	t := n.typ.rtype
	if isConst {
		t = constVal
	}
	n.rval = reflect.New(t).Elem()
	switch {
	case isConst:
		v := constant.BinaryOp(vConstantValue(v0), token.MUL, vConstantValue(v1))
		n.rval.Set(reflect.ValueOf(v))
	case isComplex(t):
		n.rval.SetComplex(vComplex(v0) * vComplex(v1))
	case isFloat(t):
		n.rval.SetFloat(vFloat(v0) * vFloat(v1))
	case isUint(t):
		n.rval.SetUint(vUint(v0) * vUint(v1))
	case isInt(t):
		n.rval.SetInt(vInt(v0) * vInt(v1))
	}
}

func or(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.concrete().TypeOf()
	isInterface := n.typ.TypeOf().Kind() == reflect.Interface
	dest := genValueOutput(n, typ)
	c0, c1 := n.child[0], n.child[1]

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch {
		case isInterface:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i | j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vInt(c0.rval)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetInt(i | j)
				return next
			}
		case c1.rval.IsValid():
			v0 := genValueInt(c0)
			j := vInt(c1.rval)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetInt(i | j)
				return next
			}
		default:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetInt(i | j)
				return next
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch {
		case isInterface:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i | j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vUint(c0.rval)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetUint(i | j)
				return next
			}
		case c1.rval.IsValid():
			j := vUint(c1.rval)
			v0 := genValueUint(c0)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetUint(i | j)
				return next
			}
		default:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetUint(i | j)
				return next
			}
		}
	}
}

func orConst(n *node) {
	v0, v1 := n.child[0].rval, n.child[1].rval
	isConst := (v0.IsValid() && isConstantValue(v0.Type())) && (v1.IsValid() && isConstantValue(v1.Type()))
	t := n.typ.rtype
	if isConst {
		t = constVal
	}
	n.rval = reflect.New(t).Elem()
	switch {
	case isConst:
		v := constant.BinaryOp(constant.ToInt(vConstantValue(v0)), token.OR, constant.ToInt(vConstantValue(v1)))
		n.rval.Set(reflect.ValueOf(v))
	case isUint(t):
		n.rval.SetUint(vUint(v0) | vUint(v1))
	case isInt(t):
		n.rval.SetInt(vInt(v0) | vInt(v1))
	}
}

func quo(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.concrete().TypeOf()
	isInterface := n.typ.TypeOf().Kind() == reflect.Interface
	dest := genValueOutput(n, typ)
	c0, c1 := n.child[0], n.child[1]

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch {
		case isInterface:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i / j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vInt(c0.rval)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetInt(i / j)
				return next
			}
		case c1.rval.IsValid():
			v0 := genValueInt(c0)
			j := vInt(c1.rval)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetInt(i / j)
				return next
			}
		default:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetInt(i / j)
				return next
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch {
		case isInterface:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i / j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vUint(c0.rval)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetUint(i / j)
				return next
			}
		case c1.rval.IsValid():
			j := vUint(c1.rval)
			v0 := genValueUint(c0)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetUint(i / j)
				return next
			}
		default:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetUint(i / j)
				return next
			}
		}
	case reflect.Float32, reflect.Float64:
		switch {
		case isInterface:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i / j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vFloat(c0.rval)
			v1 := genValueFloat(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetFloat(i / j)
				return next
			}
		case c1.rval.IsValid():
			j := vFloat(c1.rval)
			v0 := genValueFloat(c0)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetFloat(i / j)
				return next
			}
		default:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetFloat(i / j)
				return next
			}
		}
	case reflect.Complex64, reflect.Complex128:
		switch {
		case isInterface:
			v0 := genComplex(c0)
			v1 := genComplex(c1)
			n.exec = func(f *frame) bltn {
				dest(f).Set(reflect.ValueOf(v0(f) / v1(f)).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			r0 := vComplex(c0.rval)
			v1 := genComplex(c1)
			n.exec = func(f *frame) bltn {
				dest(f).SetComplex(r0 / v1(f))
				return next
			}
		case c1.rval.IsValid():
			r1 := vComplex(c1.rval)
			v0 := genComplex(c0)
			n.exec = func(f *frame) bltn {
				dest(f).SetComplex(v0(f) / r1)
				return next
			}
		default:
			v0 := genComplex(c0)
			v1 := genComplex(c1)
			n.exec = func(f *frame) bltn {
				dest(f).SetComplex(v0(f) / v1(f))
				return next
			}
		}
	}
}

func quoConst(n *node) {
	v0, v1 := n.child[0].rval, n.child[1].rval
	isConst := (v0.IsValid() && isConstantValue(v0.Type())) && (v1.IsValid() && isConstantValue(v1.Type()))
	t := n.typ.rtype
	if isConst {
		t = constVal
	}
	n.rval = reflect.New(t).Elem()
	switch {
	case isConst:
		var operator token.Token
		// When the result of the operation is expected to be an int (because both
		// operands are ints), we want to force the type of the whole expression to be an
		// int (and not a float), which is achieved by using the QUO_ASSIGN operator.
		if n.typ.untyped && isInt(n.typ.rtype) {
			operator = token.QUO_ASSIGN
		} else {
			operator = token.QUO
		}
		v := constant.BinaryOp(vConstantValue(v0), operator, vConstantValue(v1))
		n.rval.Set(reflect.ValueOf(v))
	case isComplex(t):
		n.rval.SetComplex(vComplex(v0) / vComplex(v1))
	case isFloat(t):
		n.rval.SetFloat(vFloat(v0) / vFloat(v1))
	case isUint(t):
		n.rval.SetUint(vUint(v0) / vUint(v1))
	case isInt(t):
		n.rval.SetInt(vInt(v0) / vInt(v1))
	}
}

func rem(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.concrete().TypeOf()
	isInterface := n.typ.TypeOf().Kind() == reflect.Interface
	dest := genValueOutput(n, typ)
	c0, c1 := n.child[0], n.child[1]

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch {
		case isInterface:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i % j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vInt(c0.rval)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetInt(i % j)
				return next
			}
		case c1.rval.IsValid():
			v0 := genValueInt(c0)
			j := vInt(c1.rval)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetInt(i % j)
				return next
			}
		default:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetInt(i % j)
				return next
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch {
		case isInterface:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i % j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vUint(c0.rval)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetUint(i % j)
				return next
			}
		case c1.rval.IsValid():
			j := vUint(c1.rval)
			v0 := genValueUint(c0)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetUint(i % j)
				return next
			}
		default:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetUint(i % j)
				return next
			}
		}
	}
}

func remConst(n *node) {
	v0, v1 := n.child[0].rval, n.child[1].rval
	isConst := (v0.IsValid() && isConstantValue(v0.Type())) && (v1.IsValid() && isConstantValue(v1.Type()))
	t := n.typ.rtype
	if isConst {
		t = constVal
	}
	n.rval = reflect.New(t).Elem()
	switch {
	case isConst:
		v := constant.BinaryOp(constant.ToInt(vConstantValue(v0)), token.REM, constant.ToInt(vConstantValue(v1)))
		n.rval.Set(reflect.ValueOf(v))
	case isUint(t):
		n.rval.SetUint(vUint(v0) % vUint(v1))
	case isInt(t):
		n.rval.SetInt(vInt(v0) % vInt(v1))
	}
}

func shl(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.concrete().TypeOf()
	isInterface := n.typ.TypeOf().Kind() == reflect.Interface
	dest := genValueOutput(n, typ)
	c0, c1 := n.child[0], n.child[1]

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch {
		case isInterface:
			v0 := genValueInt(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i << j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vInt(c0.rval)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetInt(i << j)
				return next
			}
		case c1.rval.IsValid():
			v0 := genValueInt(c0)
			j := vUint(c1.rval)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetInt(i << j)
				return next
			}
		default:
			v0 := genValueInt(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetInt(i << j)
				return next
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch {
		case isInterface:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i << j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vUint(c0.rval)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetUint(i << j)
				return next
			}
		case c1.rval.IsValid():
			j := vUint(c1.rval)
			v0 := genValueUint(c0)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetUint(i << j)
				return next
			}
		default:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetUint(i << j)
				return next
			}
		}
	}
}

func shlConst(n *node) {
	v0, v1 := n.child[0].rval, n.child[1].rval
	isConst := (v0.IsValid() && isConstantValue(v0.Type()))
	t := n.typ.rtype
	if isConst {
		t = constVal
	}
	n.rval = reflect.New(t).Elem()
	switch {
	case isConst:
		v := constant.Shift(vConstantValue(v0), token.SHL, uint(vUint(v1)))
		n.rval.Set(reflect.ValueOf(v))
	case isUint(t):
		n.rval.SetUint(vUint(v0) << vUint(v1))
	case isInt(t):
		n.rval.SetInt(vInt(v0) << vUint(v1))
	}
}

func shr(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.concrete().TypeOf()
	isInterface := n.typ.TypeOf().Kind() == reflect.Interface
	dest := genValueOutput(n, typ)
	c0, c1 := n.child[0], n.child[1]

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch {
		case isInterface:
			v0 := genValueInt(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i >> j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vInt(c0.rval)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetInt(i >> j)
				return next
			}
		case c1.rval.IsValid():
			v0 := genValueInt(c0)
			j := vUint(c1.rval)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetInt(i >> j)
				return next
			}
		default:
			v0 := genValueInt(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetInt(i >> j)
				return next
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch {
		case isInterface:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i >> j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vUint(c0.rval)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetUint(i >> j)
				return next
			}
		case c1.rval.IsValid():
			j := vUint(c1.rval)
			v0 := genValueUint(c0)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetUint(i >> j)
				return next
			}
		default:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetUint(i >> j)
				return next
			}
		}
	}
}

func shrConst(n *node) {
	v0, v1 := n.child[0].rval, n.child[1].rval
	isConst := (v0.IsValid() && isConstantValue(v0.Type()))
	t := n.typ.rtype
	if isConst {
		t = constVal
	}
	n.rval = reflect.New(t).Elem()
	switch {
	case isConst:
		v := constant.Shift(vConstantValue(v0), token.SHR, uint(vUint(v1)))
		n.rval.Set(reflect.ValueOf(v))
	case isUint(t):
		n.rval.SetUint(vUint(v0) >> vUint(v1))
	case isInt(t):
		n.rval.SetInt(vInt(v0) >> vUint(v1))
	}
}

func sub(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.concrete().TypeOf()
	isInterface := n.typ.TypeOf().Kind() == reflect.Interface
	dest := genValueOutput(n, typ)
	c0, c1 := n.child[0], n.child[1]

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch {
		case isInterface:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i - j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vInt(c0.rval)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetInt(i - j)
				return next
			}
		case c1.rval.IsValid():
			v0 := genValueInt(c0)
			j := vInt(c1.rval)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetInt(i - j)
				return next
			}
		default:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetInt(i - j)
				return next
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch {
		case isInterface:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i - j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vUint(c0.rval)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetUint(i - j)
				return next
			}
		case c1.rval.IsValid():
			j := vUint(c1.rval)
			v0 := genValueUint(c0)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetUint(i - j)
				return next
			}
		default:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetUint(i - j)
				return next
			}
		}
	case reflect.Float32, reflect.Float64:
		switch {
		case isInterface:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i - j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vFloat(c0.rval)
			v1 := genValueFloat(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetFloat(i - j)
				return next
			}
		case c1.rval.IsValid():
			j := vFloat(c1.rval)
			v0 := genValueFloat(c0)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetFloat(i - j)
				return next
			}
		default:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetFloat(i - j)
				return next
			}
		}
	case reflect.Complex64, reflect.Complex128:
		switch {
		case isInterface:
			v0 := genComplex(c0)
			v1 := genComplex(c1)
			n.exec = func(f *frame) bltn {
				dest(f).Set(reflect.ValueOf(v0(f) - v1(f)).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			r0 := vComplex(c0.rval)
			v1 := genComplex(c1)
			n.exec = func(f *frame) bltn {
				dest(f).SetComplex(r0 - v1(f))
				return next
			}
		case c1.rval.IsValid():
			r1 := vComplex(c1.rval)
			v0 := genComplex(c0)
			n.exec = func(f *frame) bltn {
				dest(f).SetComplex(v0(f) - r1)
				return next
			}
		default:
			v0 := genComplex(c0)
			v1 := genComplex(c1)
			n.exec = func(f *frame) bltn {
				dest(f).SetComplex(v0(f) - v1(f))
				return next
			}
		}
	}
}

func subConst(n *node) {
	v0, v1 := n.child[0].rval, n.child[1].rval
	isConst := (v0.IsValid() && isConstantValue(v0.Type())) && (v1.IsValid() && isConstantValue(v1.Type()))
	t := n.typ.rtype
	if isConst {
		t = constVal
	}
	n.rval = reflect.New(t).Elem()
	switch {
	case isConst:
		v := constant.BinaryOp(vConstantValue(v0), token.SUB, vConstantValue(v1))
		n.rval.Set(reflect.ValueOf(v))
	case isComplex(t):
		n.rval.SetComplex(vComplex(v0) - vComplex(v1))
	case isFloat(t):
		n.rval.SetFloat(vFloat(v0) - vFloat(v1))
	case isUint(t):
		n.rval.SetUint(vUint(v0) - vUint(v1))
	case isInt(t):
		n.rval.SetInt(vInt(v0) - vInt(v1))
	}
}

func xor(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.concrete().TypeOf()
	isInterface := n.typ.TypeOf().Kind() == reflect.Interface
	dest := genValueOutput(n, typ)
	c0, c1 := n.child[0], n.child[1]

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch {
		case isInterface:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i ^ j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vInt(c0.rval)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetInt(i ^ j)
				return next
			}
		case c1.rval.IsValid():
			v0 := genValueInt(c0)
			j := vInt(c1.rval)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetInt(i ^ j)
				return next
			}
		default:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetInt(i ^ j)
				return next
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch {
		case isInterface:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).Set(reflect.ValueOf(i ^ j).Convert(typ))
				return next
			}
		case c0.rval.IsValid():
			i := vUint(c0.rval)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, j := v1(f)
				dest(f).SetUint(i ^ j)
				return next
			}
		case c1.rval.IsValid():
			j := vUint(c1.rval)
			v0 := genValueUint(c0)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				dest(f).SetUint(i ^ j)
				return next
			}
		default:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, i := v0(f)
				_, j := v1(f)
				dest(f).SetUint(i ^ j)
				return next
			}
		}
	}
}

func xorConst(n *node) {
	v0, v1 := n.child[0].rval, n.child[1].rval
	isConst := (v0.IsValid() && isConstantValue(v0.Type())) && (v1.IsValid() && isConstantValue(v1.Type()))
	t := n.typ.rtype
	if isConst {
		t = constVal
	}
	n.rval = reflect.New(t).Elem()
	switch {
	case isConst:
		v := constant.BinaryOp(constant.ToInt(vConstantValue(v0)), token.XOR, constant.ToInt(vConstantValue(v1)))
		n.rval.Set(reflect.ValueOf(v))
	case isUint(t):
		n.rval.SetUint(vUint(v0) ^ vUint(v1))
	case isInt(t):
		n.rval.SetInt(vInt(v0) ^ vInt(v1))
	}
}

// Assign operators

func addAssign(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.TypeOf()
	c0, c1 := n.child[0], n.child[1]
	setMap := isMapEntry(c0)
	var mapValue, indexValue func(*frame) reflect.Value

	if setMap {
		mapValue = genValue(c0.child[0])
		indexValue = genValue(c0.child[1])
	}

	if c1.rval.IsValid() {
		switch typ.Kind() {
		case reflect.String:
			v0 := genValueString(c0)
			v1 := vString(c1.rval)
			n.exec = func(f *frame) bltn {
				v, s := v0(f)
				v.SetString(s + v1)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v0 := genValueInt(c0)
			j := vInt(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetInt(i + j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v0 := genValueUint(c0)
			j := vUint(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetUint(i + j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Float32, reflect.Float64:
			v0 := genValueFloat(c0)
			j := vFloat(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetFloat(i + j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Complex64, reflect.Complex128:
			v0 := genValue(c0)
			v1 := vComplex(c1.rval)
			n.exec = func(f *frame) bltn {
				v := v0(f)
				v.SetComplex(v.Complex() + v1)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		}
	} else {
		switch typ.Kind() {
		case reflect.String:
			v0 := genValueString(c0)
			v1 := genValue(c1)
			n.exec = func(f *frame) bltn {
				v, s := v0(f)
				v.SetString(s + v1(f).String())
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetInt(i + j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetUint(i + j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Float32, reflect.Float64:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetFloat(i + j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Complex64, reflect.Complex128:
			v0 := genValue(c0)
			v1 := genValue(c1)
			n.exec = func(f *frame) bltn {
				v := v0(f)
				v.SetComplex(v.Complex() + v1(f).Complex())
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		}
	}
}

func andAssign(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.TypeOf()
	c0, c1 := n.child[0], n.child[1]
	setMap := isMapEntry(c0)
	var mapValue, indexValue func(*frame) reflect.Value

	if setMap {
		mapValue = genValue(c0.child[0])
		indexValue = genValue(c0.child[1])
	}

	if c1.rval.IsValid() {
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v0 := genValueInt(c0)
			j := vInt(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetInt(i & j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v0 := genValueUint(c0)
			j := vUint(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetUint(i & j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		}
	} else {
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetInt(i & j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetUint(i & j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		}
	}
}

func andNotAssign(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.TypeOf()
	c0, c1 := n.child[0], n.child[1]
	setMap := isMapEntry(c0)
	var mapValue, indexValue func(*frame) reflect.Value

	if setMap {
		mapValue = genValue(c0.child[0])
		indexValue = genValue(c0.child[1])
	}

	if c1.rval.IsValid() {
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v0 := genValueInt(c0)
			j := vInt(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetInt(i &^ j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v0 := genValueUint(c0)
			j := vUint(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetUint(i &^ j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		}
	} else {
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetInt(i &^ j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetUint(i &^ j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		}
	}
}

func mulAssign(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.TypeOf()
	c0, c1 := n.child[0], n.child[1]
	setMap := isMapEntry(c0)
	var mapValue, indexValue func(*frame) reflect.Value

	if setMap {
		mapValue = genValue(c0.child[0])
		indexValue = genValue(c0.child[1])
	}

	if c1.rval.IsValid() {
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v0 := genValueInt(c0)
			j := vInt(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetInt(i * j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v0 := genValueUint(c0)
			j := vUint(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetUint(i * j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Float32, reflect.Float64:
			v0 := genValueFloat(c0)
			j := vFloat(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetFloat(i * j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Complex64, reflect.Complex128:
			v0 := genValue(c0)
			v1 := vComplex(c1.rval)
			n.exec = func(f *frame) bltn {
				v := v0(f)
				v.SetComplex(v.Complex() * v1)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		}
	} else {
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetInt(i * j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetUint(i * j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Float32, reflect.Float64:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetFloat(i * j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Complex64, reflect.Complex128:
			v0 := genValue(c0)
			v1 := genValue(c1)
			n.exec = func(f *frame) bltn {
				v := v0(f)
				v.SetComplex(v.Complex() * v1(f).Complex())
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		}
	}
}

func orAssign(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.TypeOf()
	c0, c1 := n.child[0], n.child[1]
	setMap := isMapEntry(c0)
	var mapValue, indexValue func(*frame) reflect.Value

	if setMap {
		mapValue = genValue(c0.child[0])
		indexValue = genValue(c0.child[1])
	}

	if c1.rval.IsValid() {
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v0 := genValueInt(c0)
			j := vInt(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetInt(i | j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v0 := genValueUint(c0)
			j := vUint(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetUint(i | j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		}
	} else {
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetInt(i | j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetUint(i | j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		}
	}
}

func quoAssign(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.TypeOf()
	c0, c1 := n.child[0], n.child[1]
	setMap := isMapEntry(c0)
	var mapValue, indexValue func(*frame) reflect.Value

	if setMap {
		mapValue = genValue(c0.child[0])
		indexValue = genValue(c0.child[1])
	}

	if c1.rval.IsValid() {
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v0 := genValueInt(c0)
			j := vInt(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetInt(i / j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v0 := genValueUint(c0)
			j := vUint(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetUint(i / j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Float32, reflect.Float64:
			v0 := genValueFloat(c0)
			j := vFloat(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetFloat(i / j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Complex64, reflect.Complex128:
			v0 := genValue(c0)
			v1 := vComplex(c1.rval)
			n.exec = func(f *frame) bltn {
				v := v0(f)
				v.SetComplex(v.Complex() / v1)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		}
	} else {
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetInt(i / j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetUint(i / j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Float32, reflect.Float64:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetFloat(i / j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Complex64, reflect.Complex128:
			v0 := genValue(c0)
			v1 := genValue(c1)
			n.exec = func(f *frame) bltn {
				v := v0(f)
				v.SetComplex(v.Complex() / v1(f).Complex())
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		}
	}
}

func remAssign(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.TypeOf()
	c0, c1 := n.child[0], n.child[1]
	setMap := isMapEntry(c0)
	var mapValue, indexValue func(*frame) reflect.Value

	if setMap {
		mapValue = genValue(c0.child[0])
		indexValue = genValue(c0.child[1])
	}

	if c1.rval.IsValid() {
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v0 := genValueInt(c0)
			j := vInt(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetInt(i % j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v0 := genValueUint(c0)
			j := vUint(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetUint(i % j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		}
	} else {
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetInt(i % j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetUint(i % j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		}
	}
}

func shlAssign(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.TypeOf()
	c0, c1 := n.child[0], n.child[1]
	setMap := isMapEntry(c0)
	var mapValue, indexValue func(*frame) reflect.Value

	if setMap {
		mapValue = genValue(c0.child[0])
		indexValue = genValue(c0.child[1])
	}

	if c1.rval.IsValid() {
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v0 := genValueInt(c0)
			j := vUint(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetInt(i << j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v0 := genValueUint(c0)
			j := vUint(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetUint(i << j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		}
	} else {
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v0 := genValueInt(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetInt(i << j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetUint(i << j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		}
	}
}

func shrAssign(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.TypeOf()
	c0, c1 := n.child[0], n.child[1]
	setMap := isMapEntry(c0)
	var mapValue, indexValue func(*frame) reflect.Value

	if setMap {
		mapValue = genValue(c0.child[0])
		indexValue = genValue(c0.child[1])
	}

	if c1.rval.IsValid() {
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v0 := genValueInt(c0)
			j := vUint(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetInt(i >> j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v0 := genValueUint(c0)
			j := vUint(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetUint(i >> j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		}
	} else {
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v0 := genValueInt(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetInt(i >> j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetUint(i >> j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		}
	}
}

func subAssign(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.TypeOf()
	c0, c1 := n.child[0], n.child[1]
	setMap := isMapEntry(c0)
	var mapValue, indexValue func(*frame) reflect.Value

	if setMap {
		mapValue = genValue(c0.child[0])
		indexValue = genValue(c0.child[1])
	}

	if c1.rval.IsValid() {
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v0 := genValueInt(c0)
			j := vInt(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetInt(i - j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v0 := genValueUint(c0)
			j := vUint(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetUint(i - j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Float32, reflect.Float64:
			v0 := genValueFloat(c0)
			j := vFloat(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetFloat(i - j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Complex64, reflect.Complex128:
			v0 := genValue(c0)
			v1 := vComplex(c1.rval)
			n.exec = func(f *frame) bltn {
				v := v0(f)
				v.SetComplex(v.Complex() - v1)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		}
	} else {
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetInt(i - j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetUint(i - j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Float32, reflect.Float64:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetFloat(i - j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Complex64, reflect.Complex128:
			v0 := genValue(c0)
			v1 := genValue(c1)
			n.exec = func(f *frame) bltn {
				v := v0(f)
				v.SetComplex(v.Complex() - v1(f).Complex())
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		}
	}
}

func xorAssign(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.TypeOf()
	c0, c1 := n.child[0], n.child[1]
	setMap := isMapEntry(c0)
	var mapValue, indexValue func(*frame) reflect.Value

	if setMap {
		mapValue = genValue(c0.child[0])
		indexValue = genValue(c0.child[1])
	}

	if c1.rval.IsValid() {
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v0 := genValueInt(c0)
			j := vInt(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetInt(i ^ j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v0 := genValueUint(c0)
			j := vUint(c1.rval)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				v.SetUint(i ^ j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		}
	} else {
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetInt(i ^ j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				v, i := v0(f)
				_, j := v1(f)
				v.SetUint(i ^ j)
				if setMap {
					mapValue(f).SetMapIndex(indexValue(f), v)
				}
				return next
			}
		}
	}
}

func dec(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.TypeOf()
	c0 := n.child[0]
	setMap := isMapEntry(c0)
	var mapValue, indexValue func(*frame) reflect.Value

	if setMap {
		mapValue = genValue(c0.child[0])
		indexValue = genValue(c0.child[1])
	}

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v0 := genValueInt(c0)
		n.exec = func(f *frame) bltn {
			v, i := v0(f)
			v.SetInt(i - 1)
			if setMap {
				mapValue(f).SetMapIndex(indexValue(f), v)
			}
			return next
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v0 := genValueUint(c0)
		n.exec = func(f *frame) bltn {
			v, i := v0(f)
			v.SetUint(i - 1)
			if setMap {
				mapValue(f).SetMapIndex(indexValue(f), v)
			}
			return next
		}
	case reflect.Float32, reflect.Float64:
		v0 := genValueFloat(c0)
		n.exec = func(f *frame) bltn {
			v, i := v0(f)
			v.SetFloat(i - 1)
			if setMap {
				mapValue(f).SetMapIndex(indexValue(f), v)
			}
			return next
		}
	case reflect.Complex64, reflect.Complex128:
		v0 := genValue(c0)
		n.exec = func(f *frame) bltn {
			v := v0(f)
			v.SetComplex(v.Complex() - 1)
			if setMap {
				mapValue(f).SetMapIndex(indexValue(f), v)
			}
			return next
		}
	}
}

func inc(n *node) {
	next := getExec(n.tnext)
	typ := n.typ.TypeOf()
	c0 := n.child[0]
	setMap := isMapEntry(c0)
	var mapValue, indexValue func(*frame) reflect.Value

	if setMap {
		mapValue = genValue(c0.child[0])
		indexValue = genValue(c0.child[1])
	}

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v0 := genValueInt(c0)
		n.exec = func(f *frame) bltn {
			v, i := v0(f)
			v.SetInt(i + 1)
			if setMap {
				mapValue(f).SetMapIndex(indexValue(f), v)
			}
			return next
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v0 := genValueUint(c0)
		n.exec = func(f *frame) bltn {
			v, i := v0(f)
			v.SetUint(i + 1)
			if setMap {
				mapValue(f).SetMapIndex(indexValue(f), v)
			}
			return next
		}
	case reflect.Float32, reflect.Float64:
		v0 := genValueFloat(c0)
		n.exec = func(f *frame) bltn {
			v, i := v0(f)
			v.SetFloat(i + 1)
			if setMap {
				mapValue(f).SetMapIndex(indexValue(f), v)
			}
			return next
		}
	case reflect.Complex64, reflect.Complex128:
		v0 := genValue(c0)
		n.exec = func(f *frame) bltn {
			v := v0(f)
			v.SetComplex(v.Complex() + 1)
			if setMap {
				mapValue(f).SetMapIndex(indexValue(f), v)
			}
			return next
		}
	}
}

func bitNotConst(n *node) {
	v0 := n.child[0].rval
	isConst := v0.IsValid() && isConstantValue(v0.Type())
	t := n.typ.rtype
	if isConst {
		t = constVal
	}
	n.rval = reflect.New(t).Elem()
	switch {
	case isConst:
		v := constant.UnaryOp(token.XOR, vConstantValue(v0), 0)
		n.rval.Set(reflect.ValueOf(v))
	case isUint(t):
		n.rval.SetUint(^v0.Uint())
	case isInt(t):
		n.rval.SetInt(^v0.Int())
	}
}

func negConst(n *node) {
	v0 := n.child[0].rval
	isConst := v0.IsValid() && isConstantValue(v0.Type())
	t := n.typ.rtype
	if isConst {
		t = constVal
	}
	n.rval = reflect.New(t).Elem()
	switch {
	case isConst:
		v := constant.UnaryOp(token.SUB, vConstantValue(v0), 0)
		n.rval.Set(reflect.ValueOf(v))
	case isUint(t):
		n.rval.SetUint(-v0.Uint())
	case isInt(t):
		n.rval.SetInt(-v0.Int())
	case isFloat(t):
		n.rval.SetFloat(-v0.Float())
	case isComplex(t):
		n.rval.SetComplex(-v0.Complex())
	}
}

func notConst(n *node) {
	v0 := n.child[0].rval
	isConst := v0.IsValid() && isConstantValue(v0.Type())
	t := n.typ.rtype
	if isConst {
		t = constVal
	}
	n.rval = reflect.New(t).Elem()
	if isConst {
		v := constant.UnaryOp(token.NOT, vConstantValue(v0), 0)
		n.rval.Set(reflect.ValueOf(v))
	} else {
		n.rval.SetBool(!v0.Bool())
	}
}

func posConst(n *node) {
	v0 := n.child[0].rval
	isConst := v0.IsValid() && isConstantValue(v0.Type())
	t := n.typ.rtype
	if isConst {
		t = constVal
	}
	n.rval = reflect.New(t).Elem()
	switch {
	case isConst:
		v := constant.UnaryOp(token.ADD, vConstantValue(v0), 0)
		n.rval.Set(reflect.ValueOf(v))
	case isUint(t):
		n.rval.SetUint(+v0.Uint())
	case isInt(t):
		n.rval.SetInt(+v0.Int())
	case isFloat(t):
		n.rval.SetFloat(+v0.Float())
	case isComplex(t):
		n.rval.SetComplex(+v0.Complex())
	}
}

func equal(n *node) {
	tnext := getExec(n.tnext)
	dest := genValueOutput(n, reflect.TypeOf(true))
	typ := n.typ.concrete().TypeOf()
	isInterface := n.typ.TypeOf().Kind() == reflect.Interface
	c0, c1 := n.child[0], n.child[1]
	t0, t1 := c0.typ.TypeOf(), c1.typ.TypeOf()

	if c0.typ.cat == linkedT || c1.typ.cat == linkedT {
		switch {
		case isInterface:
			v0 := genValue(c0)
			v1 := genValue(c1)
			dest := genValue(n)
			n.exec = func(f *frame) bltn {
				i0 := v0(f).Interface()
				i1 := v1(f).Interface()
				dest(f).Set(reflect.ValueOf(i0 == i1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			i0 := c0.rval.Interface()
			v1 := genValue(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					i1 := v1(f).Interface()
					if i0 == i1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					i1 := v1(f).Interface()
					dest(f).SetBool(i0 == i1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			i1 := c1.rval.Interface()
			v0 := genValue(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					i0 := v0(f).Interface()
					if i0 == i1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					i0 := v0(f).Interface()
					dest(f).SetBool(i0 == i1)
					return tnext
				}
			}
		default:
			v0 := genValue(c0)
			v1 := genValue(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					i0 := v0(f).Interface()
					i1 := v1(f).Interface()
					if i0 == i1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					i0 := v0(f).Interface()
					i1 := v1(f).Interface()
					dest(f).SetBool(i0 == i1)
					return tnext
				}
			}
		}
		return
	}

	// Do not attempt to optimize '==' or '!=' if an operand is an interface.
	// This will preserve proper dynamic type checking at runtime. For static types,
	// type checks are already performed, so bypass them if possible.
	if t0.Kind() == reflect.Interface || t1.Kind() == reflect.Interface {
		v0 := genValue(c0)
		v1 := genValue(c1)
		if n.fnext != nil {
			fnext := getExec(n.fnext)
			n.exec = func(f *frame) bltn {
				i0 := v0(f).Interface()
				i1 := v1(f).Interface()
				if i0 == i1 {
					dest(f).SetBool(true)
					return tnext
				}
				dest(f).SetBool(false)
				return fnext
			}
		} else {
			dest := genValue(n)
			n.exec = func(f *frame) bltn {
				i0 := v0(f).Interface()
				i1 := v1(f).Interface()
				dest(f).SetBool(i0 == i1)
				return tnext
			}
		}
		return
	}

	switch {
	case isString(t0) || isString(t1):
		switch {
		case isInterface:
			v0 := genValueString(c0)
			v1 := genValueString(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 == s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vString(c0.rval)
			v1 := genValueString(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 == s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 == s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vString(c1.rval)
			v0 := genValueString(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 == s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 == s1)
					return tnext
				}
			}
		default:
			v0 := genValueString(c0)
			v1 := genValueString(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 == s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 == s1)
					return tnext
				}
			}
		}
	case isFloat(t0) || isFloat(t1):
		switch {
		case isInterface:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 == s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vFloat(c0.rval)
			v1 := genValueFloat(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 == s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 == s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vFloat(c1.rval)
			v0 := genValueFloat(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 == s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 == s1)
					return tnext
				}
			}
		default:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 == s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 == s1)
					return tnext
				}
			}
		}
	case isUint(t0) || isUint(t1):
		switch {
		case isInterface:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 == s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vUint(c0.rval)
			v1 := genValueUint(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 == s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 == s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vUint(c1.rval)
			v0 := genValueUint(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 == s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 == s1)
					return tnext
				}
			}
		default:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 == s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 == s1)
					return tnext
				}
			}
		}
	case isInt(t0) || isInt(t1):
		switch {
		case isInterface:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 == s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vInt(c0.rval)
			v1 := genValueInt(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 == s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 == s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vInt(c1.rval)
			v0 := genValueInt(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 == s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 == s1)
					return tnext
				}
			}
		default:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 == s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 == s1)
					return tnext
				}
			}
		}
	case isComplex(t0) || isComplex(t1):
		switch {
		case isInterface:
			v0 := genComplex(c0)
			v1 := genComplex(c1)
			n.exec = func(f *frame) bltn {
				s0 := v0(f)
				s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 == s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vComplex(c0.rval)
			v1 := genComplex(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					s1 := v1(f)
					if s0 == s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					s1 := v1(f)
					dest(f).SetBool(s0 == s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vComplex(c1.rval)
			v0 := genComplex(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					s0 := v0(f)
					if s0 == s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					s0 := v0(f)
					dest(f).SetBool(s0 == s1)
					return tnext
				}
			}
		default:
			v0 := genComplex(c0)
			v1 := genComplex(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					s0 := v0(f)
					s1 := v1(f)
					if s0 == s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					s0 := v0(f)
					s1 := v1(f)
					dest(f).SetBool(s0 == s1)
					return tnext
				}
			}
		}
	default:
		switch {
		case isInterface:
			v0 := genValue(c0)
			v1 := genValue(c1)
			n.exec = func(f *frame) bltn {
				i0 := v0(f).Interface()
				i1 := v1(f).Interface()
				dest(f).Set(reflect.ValueOf(i0 == i1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			i0 := c0.rval.Interface()
			v1 := genValue(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					i1 := v1(f).Interface()
					if i0 == i1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					i1 := v1(f).Interface()
					dest(f).SetBool(i0 == i1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			i1 := c1.rval.Interface()
			v0 := genValue(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					i0 := v0(f).Interface()
					if i0 == i1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					i0 := v0(f).Interface()
					dest(f).SetBool(i0 == i1)
					return tnext
				}
			}
		default:
			v0 := genValue(c0)
			v1 := genValue(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					i0 := v0(f).Interface()
					i1 := v1(f).Interface()
					if i0 == i1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					i0 := v0(f).Interface()
					i1 := v1(f).Interface()
					dest(f).SetBool(i0 == i1)
					return tnext
				}
			}
		}
	}
}

func greater(n *node) {
	tnext := getExec(n.tnext)
	dest := genValueOutput(n, reflect.TypeOf(true))
	typ := n.typ.concrete().TypeOf()
	isInterface := n.typ.TypeOf().Kind() == reflect.Interface
	c0, c1 := n.child[0], n.child[1]
	t0, t1 := c0.typ.TypeOf(), c1.typ.TypeOf()

	switch {
	case isString(t0) || isString(t1):
		switch {
		case isInterface:
			v0 := genValueString(c0)
			v1 := genValueString(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 > s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vString(c0.rval)
			v1 := genValueString(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 > s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 > s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vString(c1.rval)
			v0 := genValueString(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 > s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 > s1)
					return tnext
				}
			}
		default:
			v0 := genValueString(c0)
			v1 := genValueString(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 > s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 > s1)
					return tnext
				}
			}
		}
	case isFloat(t0) || isFloat(t1):
		switch {
		case isInterface:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 > s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vFloat(c0.rval)
			v1 := genValueFloat(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 > s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 > s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vFloat(c1.rval)
			v0 := genValueFloat(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 > s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 > s1)
					return tnext
				}
			}
		default:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 > s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 > s1)
					return tnext
				}
			}
		}
	case isUint(t0) || isUint(t1):
		switch {
		case isInterface:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 > s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vUint(c0.rval)
			v1 := genValueUint(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 > s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 > s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vUint(c1.rval)
			v0 := genValueUint(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 > s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 > s1)
					return tnext
				}
			}
		default:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 > s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 > s1)
					return tnext
				}
			}
		}
	case isInt(t0) || isInt(t1):
		switch {
		case isInterface:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 > s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vInt(c0.rval)
			v1 := genValueInt(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 > s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 > s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vInt(c1.rval)
			v0 := genValueInt(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 > s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 > s1)
					return tnext
				}
			}
		default:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 > s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 > s1)
					return tnext
				}
			}
		}
	}
}

func greaterEqual(n *node) {
	tnext := getExec(n.tnext)
	dest := genValueOutput(n, reflect.TypeOf(true))
	typ := n.typ.concrete().TypeOf()
	isInterface := n.typ.TypeOf().Kind() == reflect.Interface
	c0, c1 := n.child[0], n.child[1]
	t0, t1 := c0.typ.TypeOf(), c1.typ.TypeOf()

	switch {
	case isString(t0) || isString(t1):
		switch {
		case isInterface:
			v0 := genValueString(c0)
			v1 := genValueString(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 >= s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vString(c0.rval)
			v1 := genValueString(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 >= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 >= s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vString(c1.rval)
			v0 := genValueString(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 >= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 >= s1)
					return tnext
				}
			}
		default:
			v0 := genValueString(c0)
			v1 := genValueString(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 >= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 >= s1)
					return tnext
				}
			}
		}
	case isFloat(t0) || isFloat(t1):
		switch {
		case isInterface:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 >= s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vFloat(c0.rval)
			v1 := genValueFloat(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 >= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 >= s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vFloat(c1.rval)
			v0 := genValueFloat(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 >= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 >= s1)
					return tnext
				}
			}
		default:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 >= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 >= s1)
					return tnext
				}
			}
		}
	case isUint(t0) || isUint(t1):
		switch {
		case isInterface:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 >= s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vUint(c0.rval)
			v1 := genValueUint(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 >= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 >= s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vUint(c1.rval)
			v0 := genValueUint(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 >= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 >= s1)
					return tnext
				}
			}
		default:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 >= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 >= s1)
					return tnext
				}
			}
		}
	case isInt(t0) || isInt(t1):
		switch {
		case isInterface:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 >= s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vInt(c0.rval)
			v1 := genValueInt(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 >= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 >= s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vInt(c1.rval)
			v0 := genValueInt(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 >= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 >= s1)
					return tnext
				}
			}
		default:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 >= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 >= s1)
					return tnext
				}
			}
		}
	}
}

func lower(n *node) {
	tnext := getExec(n.tnext)
	dest := genValueOutput(n, reflect.TypeOf(true))
	typ := n.typ.concrete().TypeOf()
	isInterface := n.typ.TypeOf().Kind() == reflect.Interface
	c0, c1 := n.child[0], n.child[1]
	t0, t1 := c0.typ.TypeOf(), c1.typ.TypeOf()

	switch {
	case isString(t0) || isString(t1):
		switch {
		case isInterface:
			v0 := genValueString(c0)
			v1 := genValueString(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 < s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vString(c0.rval)
			v1 := genValueString(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 < s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 < s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vString(c1.rval)
			v0 := genValueString(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 < s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 < s1)
					return tnext
				}
			}
		default:
			v0 := genValueString(c0)
			v1 := genValueString(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 < s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 < s1)
					return tnext
				}
			}
		}
	case isFloat(t0) || isFloat(t1):
		switch {
		case isInterface:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 < s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vFloat(c0.rval)
			v1 := genValueFloat(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 < s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 < s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vFloat(c1.rval)
			v0 := genValueFloat(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 < s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 < s1)
					return tnext
				}
			}
		default:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 < s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 < s1)
					return tnext
				}
			}
		}
	case isUint(t0) || isUint(t1):
		switch {
		case isInterface:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 < s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vUint(c0.rval)
			v1 := genValueUint(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 < s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 < s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vUint(c1.rval)
			v0 := genValueUint(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 < s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 < s1)
					return tnext
				}
			}
		default:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 < s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 < s1)
					return tnext
				}
			}
		}
	case isInt(t0) || isInt(t1):
		switch {
		case isInterface:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 < s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vInt(c0.rval)
			v1 := genValueInt(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 < s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 < s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vInt(c1.rval)
			v0 := genValueInt(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 < s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 < s1)
					return tnext
				}
			}
		default:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 < s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 < s1)
					return tnext
				}
			}
		}
	}
}

func lowerEqual(n *node) {
	tnext := getExec(n.tnext)
	dest := genValueOutput(n, reflect.TypeOf(true))
	typ := n.typ.concrete().TypeOf()
	isInterface := n.typ.TypeOf().Kind() == reflect.Interface
	c0, c1 := n.child[0], n.child[1]
	t0, t1 := c0.typ.TypeOf(), c1.typ.TypeOf()

	switch {
	case isString(t0) || isString(t1):
		switch {
		case isInterface:
			v0 := genValueString(c0)
			v1 := genValueString(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 <= s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vString(c0.rval)
			v1 := genValueString(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 <= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 <= s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vString(c1.rval)
			v0 := genValueString(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 <= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 <= s1)
					return tnext
				}
			}
		default:
			v0 := genValueString(c0)
			v1 := genValueString(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 <= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 <= s1)
					return tnext
				}
			}
		}
	case isFloat(t0) || isFloat(t1):
		switch {
		case isInterface:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 <= s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vFloat(c0.rval)
			v1 := genValueFloat(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 <= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 <= s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vFloat(c1.rval)
			v0 := genValueFloat(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 <= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 <= s1)
					return tnext
				}
			}
		default:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 <= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 <= s1)
					return tnext
				}
			}
		}
	case isUint(t0) || isUint(t1):
		switch {
		case isInterface:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 <= s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vUint(c0.rval)
			v1 := genValueUint(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 <= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 <= s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vUint(c1.rval)
			v0 := genValueUint(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 <= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 <= s1)
					return tnext
				}
			}
		default:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 <= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 <= s1)
					return tnext
				}
			}
		}
	case isInt(t0) || isInt(t1):
		switch {
		case isInterface:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 <= s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vInt(c0.rval)
			v1 := genValueInt(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 <= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 <= s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vInt(c1.rval)
			v0 := genValueInt(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 <= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 <= s1)
					return tnext
				}
			}
		default:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 <= s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 <= s1)
					return tnext
				}
			}
		}
	}
}

func notEqual(n *node) {
	tnext := getExec(n.tnext)
	dest := genValueOutput(n, reflect.TypeOf(true))
	typ := n.typ.concrete().TypeOf()
	isInterface := n.typ.TypeOf().Kind() == reflect.Interface
	c0, c1 := n.child[0], n.child[1]
	t0, t1 := c0.typ.TypeOf(), c1.typ.TypeOf()

	if c0.typ.cat == linkedT || c1.typ.cat == linkedT {
		switch {
		case isInterface:
			v0 := genValue(c0)
			v1 := genValue(c1)
			dest := genValue(n)
			n.exec = func(f *frame) bltn {
				i0 := v0(f).Interface()
				i1 := v1(f).Interface()
				dest(f).Set(reflect.ValueOf(i0 != i1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			i0 := c0.rval.Interface()
			v1 := genValue(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					i1 := v1(f).Interface()
					if i0 != i1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					i1 := v1(f).Interface()
					dest(f).SetBool(i0 != i1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			i1 := c1.rval.Interface()
			v0 := genValue(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					i0 := v0(f).Interface()
					if i0 != i1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					i0 := v0(f).Interface()
					dest(f).SetBool(i0 != i1)
					return tnext
				}
			}
		default:
			v0 := genValue(c0)
			v1 := genValue(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					i0 := v0(f).Interface()
					i1 := v1(f).Interface()
					if i0 != i1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					i0 := v0(f).Interface()
					i1 := v1(f).Interface()
					dest(f).SetBool(i0 != i1)
					return tnext
				}
			}
		}
		return
	}

	// Do not attempt to optimize '==' or '!=' if an operand is an interface.
	// This will preserve proper dynamic type checking at runtime. For static types,
	// type checks are already performed, so bypass them if possible.
	if t0.Kind() == reflect.Interface || t1.Kind() == reflect.Interface {
		v0 := genValue(c0)
		v1 := genValue(c1)
		if n.fnext != nil {
			fnext := getExec(n.fnext)
			n.exec = func(f *frame) bltn {
				i0 := v0(f).Interface()
				i1 := v1(f).Interface()
				if i0 != i1 {
					dest(f).SetBool(true)
					return tnext
				}
				dest(f).SetBool(false)
				return fnext
			}
		} else {
			dest := genValue(n)
			n.exec = func(f *frame) bltn {
				i0 := v0(f).Interface()
				i1 := v1(f).Interface()
				dest(f).SetBool(i0 != i1)
				return tnext
			}
		}
		return
	}

	switch {
	case isString(t0) || isString(t1):
		switch {
		case isInterface:
			v0 := genValueString(c0)
			v1 := genValueString(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 != s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vString(c0.rval)
			v1 := genValueString(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 != s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 != s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vString(c1.rval)
			v0 := genValueString(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 != s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 != s1)
					return tnext
				}
			}
		default:
			v0 := genValueString(c0)
			v1 := genValueString(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 != s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 != s1)
					return tnext
				}
			}
		}
	case isFloat(t0) || isFloat(t1):
		switch {
		case isInterface:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 != s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vFloat(c0.rval)
			v1 := genValueFloat(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 != s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 != s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vFloat(c1.rval)
			v0 := genValueFloat(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 != s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 != s1)
					return tnext
				}
			}
		default:
			v0 := genValueFloat(c0)
			v1 := genValueFloat(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 != s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 != s1)
					return tnext
				}
			}
		}
	case isUint(t0) || isUint(t1):
		switch {
		case isInterface:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 != s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vUint(c0.rval)
			v1 := genValueUint(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 != s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 != s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vUint(c1.rval)
			v0 := genValueUint(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 != s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 != s1)
					return tnext
				}
			}
		default:
			v0 := genValueUint(c0)
			v1 := genValueUint(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 != s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 != s1)
					return tnext
				}
			}
		}
	case isInt(t0) || isInt(t1):
		switch {
		case isInterface:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			n.exec = func(f *frame) bltn {
				_, s0 := v0(f)
				_, s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 != s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vInt(c0.rval)
			v1 := genValueInt(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					if s0 != s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s1 := v1(f)
					dest(f).SetBool(s0 != s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vInt(c1.rval)
			v0 := genValueInt(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					if s0 != s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					dest(f).SetBool(s0 != s1)
					return tnext
				}
			}
		default:
			v0 := genValueInt(c0)
			v1 := genValueInt(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					if s0 != s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					_, s0 := v0(f)
					_, s1 := v1(f)
					dest(f).SetBool(s0 != s1)
					return tnext
				}
			}
		}
	case isComplex(t0) || isComplex(t1):
		switch {
		case isInterface:
			v0 := genComplex(c0)
			v1 := genComplex(c1)
			n.exec = func(f *frame) bltn {
				s0 := v0(f)
				s1 := v1(f)
				dest(f).Set(reflect.ValueOf(s0 != s1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			s0 := vComplex(c0.rval)
			v1 := genComplex(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					s1 := v1(f)
					if s0 != s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					s1 := v1(f)
					dest(f).SetBool(s0 != s1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			s1 := vComplex(c1.rval)
			v0 := genComplex(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					s0 := v0(f)
					if s0 != s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					s0 := v0(f)
					dest(f).SetBool(s0 != s1)
					return tnext
				}
			}
		default:
			v0 := genComplex(c0)
			v1 := genComplex(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					s0 := v0(f)
					s1 := v1(f)
					if s0 != s1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				n.exec = func(f *frame) bltn {
					s0 := v0(f)
					s1 := v1(f)
					dest(f).SetBool(s0 != s1)
					return tnext
				}
			}
		}
	default:
		switch {
		case isInterface:
			v0 := genValue(c0)
			v1 := genValue(c1)
			n.exec = func(f *frame) bltn {
				i0 := v0(f).Interface()
				i1 := v1(f).Interface()
				dest(f).Set(reflect.ValueOf(i0 != i1).Convert(typ))
				return tnext
			}
		case c0.rval.IsValid():
			i0 := c0.rval.Interface()
			v1 := genValue(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					i1 := v1(f).Interface()
					if i0 != i1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					i1 := v1(f).Interface()
					dest(f).SetBool(i0 != i1)
					return tnext
				}
			}
		case c1.rval.IsValid():
			i1 := c1.rval.Interface()
			v0 := genValue(c0)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					i0 := v0(f).Interface()
					if i0 != i1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					i0 := v0(f).Interface()
					dest(f).SetBool(i0 != i1)
					return tnext
				}
			}
		default:
			v0 := genValue(c0)
			v1 := genValue(c1)
			if n.fnext != nil {
				fnext := getExec(n.fnext)
				n.exec = func(f *frame) bltn {
					i0 := v0(f).Interface()
					i1 := v1(f).Interface()
					if i0 != i1 {
						dest(f).SetBool(true)
						return tnext
					}
					dest(f).SetBool(false)
					return fnext
				}
			} else {
				dest := genValue(n)
				n.exec = func(f *frame) bltn {
					i0 := v0(f).Interface()
					i1 := v1(f).Interface()
					dest(f).SetBool(i0 != i1)
					return tnext
				}
			}
		}
	}
}
