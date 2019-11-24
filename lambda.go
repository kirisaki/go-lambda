package main

import (
	"sort"
	"strconv"
)

type Name string

type Names []Name

func (xs Names) Len() int{
	return len(xs)
}

func (xs Names) Less(i, j int) bool {
	return xs[i] < xs[j]
}

func (xs Names) Swap(i, j int) {
	xs[i], xs[j] = xs[j], xs[i]
}

type Var struct {
	name Name
}

type Lam struct {
	name Name
	expr Expr
}

type App struct {
	f Expr
	arg Expr
}

type Expr interface {
	reduce() Expr
}

func (x Var) reduce() Expr {
	return x
}

func (x Lam) reduce() Expr {
	return x
}

func (x App) reduce() Expr {
	switch g := x.f.(type) {
	case Var:
		return App{g, x.arg.reduce()}
	case Lam:
		return subst(g.name, x.arg.reduce(), g.expr.reduce())
	case App:
		return App{x.f.reduce(), x.arg.reduce()}.reduce()
	default:
		panic("")
	}
}

var cnt = 0

func subst(x Name, s Expr, y Expr) Expr {
	switch z := y.(type) {
	case Var:
		if x == z.name {
			return s
		} else {
			return y
		}
	case Lam:
		if x == z.name {
			return z
		} else {
			if !elem(z.name, free(s)) {
				return Lam{z.name, subst(x, s, z.expr)}
			} else {
				n := Name((string)(z.name) + strconv.Itoa(cnt))
				cnt++
				return  Lam{n, subst(x, s, subst(z.name, Var{n}, z.expr))}
			}
		}
	case App:
		return App{subst(x, s, z.f), subst(x, s, z.arg)}

	default:
		panic("")
	}
}

func free(x Expr) Names {
	switch y := x.(type) {
	case Var:
		return Names{y.name}
	case Lam:
		return remove(y.name, free(y.expr))
	case App:
		return union(free(y.f), free(y.arg))
	default:
		return Names{}
	}
}


func remove(x Name, xs Names) Names {
	res := Names{}
	for _, v := range xs {
		if v != x {
			res = append(res, v)
		}
	}
	return res
}

func union(xs Names, ys Names) Names {
	zs := append(xs, ys...)
	sort.Sort(zs)
	res := Names{}
	prev := Name("")
	for _, v := range zs{
		if prev != v{
			res = append(res, v)
		}
		prev = v
	}
	return res
}

func elem(x Name, xs Names) bool {
	for _, v := range xs {
		if x == v {
			return true
		}
	}
	return false
}
