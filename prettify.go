package main

import (

)

func prettify(expr Expr) string {
	switch e := expr.(type) {
	case Var:
		return string(e.name)
	case Lam:
		return "\\" + string(e.name) + "." + prettify(e.expr)
	case App:
		return prettify(e.f) + prettify(e.arg)
	default:
		return ""
	}
}
