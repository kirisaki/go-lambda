package main

import (
	"errors"
)

type Token interface {}

func parse(str string) (Expr, error) {
	rpn, err := toRpn(str)
	if err != nil {
		return nil, err
	}
	stack := []Expr{}
	for _, t := range rpn {
		switch t {
		case '\\':
			lam := true
			for lam {
				xs, x, _ := popExpr(stack)
				ys, y, err := popExpr(xs)
				if err != nil {
					return nil, errors.New("token exhausted")
				}
				stack = ys
				if w, ok := x.(Sym); ok {
					x = Var{Name(string([]rune{w.symbol}))}
				}
				if z, ok := y.(Sym); ok {
					if z.symbol == '.' {
						vs, v, err := popExpr(stack)
						if err != nil {
							return nil, errors.New("argument notfound")
						}
						stack = vs
						if u, ok := v.(Sym); ok {
							stack = pushExpr(Lam{Name(string([]rune{u.symbol})), x}, stack)
							lam = false
						} else {
							return nil, errors.New("argument must be symbol")
						}
					} else {
						y = Var{Name(string([]rune{z.symbol}))}
						stack = pushExpr(App{y, x}, stack)
					}
				}
			}

		default:
			stack = pushExpr(Sym{t}, stack)
		}
	}
	for {
		xs, x, err1 := popExpr(stack)
		if err1 != nil {
			return nil, errors.New("no result")
		}
		if x0, ok := x.(Sym); ok {
			x = Var{Name(string([]rune{x0.symbol}))}
		} 
		ys, y, err2 := popExpr(xs)
		if err2 != nil {
			return x, nil
		} else {
			if y0, ok := y.(Sym); ok {
				y = Var{Name(string([]rune{y0.symbol}))}
			} 
			stack = pushExpr(App{y, x}, ys)
		}
	}
}

func toRpn(str string) ([]rune, error) {
	stack := []rune{}
	rpn := []rune{}
	isParam := false
	lamCnt := 0
	for _, c := range str {
		switch c {
		case '\\', 'λ':
			if isParam {
				return nil, errors.New("lmbda in parameters")
			}
			stack = push('\\', stack)
			isParam = true
		case '.':
			for {
				xs, x, err := pop(stack)
				if err != nil {
					return nil, errors.New("mismatched lambda")
				}
				stack = xs
				if x == '\\' {
					break
				}
				rpn = append(rpn, x)
			}
			isParam = false
		case '(':
			if isParam {
				return nil, errors.New("parens in parameters")
			}
			stack = push('(', stack)
		case ')':
			if isParam {
				return nil, errors.New("parens in parameters")
			}
			for {
				xs, x, err := pop(stack)
				if err != nil {
					return nil, errors.New("mismatched parens")
				}
				stack = xs
				if x == '(' {
					break
				}
				rpn = append(rpn, x)
			}
			for i := 0; i < lamCnt; i++ {
				rpn = append(rpn, '\\')
			}
			lamCnt = 0

		case ' ', '　':
		default:
			rpn = append(rpn, c)
			if isParam {
				lamCnt++
				rpn = append(rpn, '.')
			}
		}
	}
	if isParam {
		return nil, errors.New("lacking expression")
	}
	for {
		xs, x, err := pop(stack)
		if err != nil{
			break
		}
		if x == '(' || x == ')' || x == '\\' || x == '.' {
			return nil, errors.New("invalid tokens remain")
		}
		stack = xs
		rpn = append(rpn, x)
	}
	for i := 0; i < lamCnt; i++ {
		rpn = append(rpn, '\\')
	}
	return rpn, nil
}

func push(x rune, xs []rune) []rune {
	return append(xs, x)
}

func pop(xs []rune) ([]rune, rune, error) {
	if len(xs) == 0 {
		return nil, ' ', errors.New("empty stack")
	}
	last := xs[len(xs) - 1]
	init := xs[:len(xs) - 1]
	return init, last, nil
}

func pushExpr(x Expr, xs []Expr) []Expr {
	return append(xs, x)
}

func popExpr(xs []Expr) ([]Expr, Expr, error) {
	if len(xs) == 0 {
		return nil, nil, errors.New("empty stack")
	}
	last := xs[len(xs) - 1]
	init := xs[:len(xs) - 1]
	return init, last, nil
}
