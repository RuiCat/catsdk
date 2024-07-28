package eval

import (
	"fmt"
	"path"
	"reflect"
)

// Symbols returns a map of interpreter exported symbol values for the given
// import path. If the argument is the empty string, all known symbols are
// returned.
func (interp *Interpreter) Symbols(importPath string) Exports {
	m := map[string]map[string]reflect.Value{}
	interp.mutex.RLock()
	defer interp.mutex.RUnlock()

	for k, v := range interp.srcPkg {
		if importPath != "" && k != importPath {
			continue
		}
		syms := map[string]reflect.Value{}
		for n, s := range v {
			if !canExport(n) {
				// Skip private non-exported symbols.
				continue
			}
			switch s.kind {
			case constSym:
				syms[n] = s.rval
			case funcSym:
				syms[n] = genFunctionWrapper(s.node)(interp.frame)
			case varSym:
				syms[n] = interp.frame.data[s.index]
			case typeSym:
				syms[n] = reflect.New(s.typ.TypeOf())
			}
		}

		if len(syms) > 0 {
			m[k] = syms
		}

		if importPath != "" {
			return m
		}
	}

	if importPath != "" && len(m) > 0 {
		return m
	}

	for k, v := range interp.binPkg {
		if importPath != "" && k != importPath {
			continue
		}
		m[k] = v
		if importPath != "" {
			return m
		}
	}

	return m
}

// getWrapper returns the wrapper type of the corresponding interface, trying
// first the composed ones, or nil if not found.
func getWrapper(n *node, t reflect.Type) reflect.Type {
	p, ok := n.interp.binPkg[t.PkgPath()]
	if !ok {
		return nil
	}
	w := p["_"+t.Name()]
	lm := n.typ.methods()

	// mapTypes may contain composed interfaces wrappers to test against, from
	// most complex to simplest (guaranteed by construction of mapTypes). Find the
	// first for which the interpreter type has all the methods.
	for _, rt := range n.interp.mapTypes[w] {
		match := true
		for i := 1; i < rt.NumField(); i++ {
			// The interpreter type must have all required wrapper methods.
			if _, ok := lm[rt.Field(i).Name[1:]]; !ok {
				match = false
				break
			}
		}
		if match {
			return rt
		}
	}

	// Otherwise return the direct "non-composed" interface.
	return w.Type().Elem()
}

// Use loads binary runtime symbols in the interpreter context so
// they can be used in interpreted code.
func (interp *Interpreter) Use(values Exports) error {
	for k, v := range values {
		importPath := path.Dir(k)
		packageName := path.Base(k)

		if k == "." && v["MapTypes"].IsValid() {
			// Use mapping for special interface wrappers.
			for kk, vv := range v["MapTypes"].Interface().(map[reflect.Value][]reflect.Type) {
				interp.mapTypes[kk] = vv
			}
			continue
		}

		if importPath == "." {
			return fmt.Errorf("export path %[1]q is missing a package name; did you mean '%[1]s/%[1]s'?", k)
		}

		if importPath == selfPrefix {
			interp.hooks.Parse(v)
			continue
		}

		if interp.binPkg[importPath] == nil {
			interp.binPkg[importPath] = make(map[string]reflect.Value)
			interp.pkgNames[importPath] = packageName
		}

		for s, sym := range v {
			interp.binPkg[importPath][s] = sym
		}
		if k == selfPath {
			interp.binPkg[importPath]["Self"] = reflect.ValueOf(interp)
		}
	}

	// Checks if input values correspond to stdlib packages by looking for one
	// well known stdlib package path.
	return nil
}
