package main

import (
	"errors"
	"fmt"
)

type Token interface {}

func parse(str string) (Expr, error) {
	rpn, err := toRpn(str)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(rpn))
	stack := []Token{}
	var now Expr
	for _, t := range rpn {
		switch t {
		case '\\':
			end := false
			for {
				xs, x, err := popTk(stack)
				if err != nil {
					return nil, errors.New("token exhausted")
				}
				stack = xs
				switch y := x.(type) {
				case Expr:
					if now == nil {
						now = y
					} else {
						now = App{y, now}
					}
				case rune:
					if y == '.' {
						end = true
						break
					}
					if now == nil {
						now = Var{Name(string([]rune{y}))}
					} else {
						now = App{Var{Name(string([]rune{y}))}, now}
					}
				}
				if end {
					break
				}
			}
			xs, x, err := popTk(stack)
			if err != nil {
				return nil, errors.New("parameters don't exist")
			}
			stack = xs
			switch y := x.(type) {
			case Expr:
				return nil, errors.New("expression in lambda")
			case rune:
				now = Lam{Name(string([]rune{y})), now}
			}

		default:
			stack = pushTk(t, stack)
		}
	}
	end := false
	for {
		xs, x, err := popTk(stack)
		if err != nil {
			break
		}
		stack = xs
		switch y := x.(type) {
		case Expr:
			if now == nil {
				now = y
			} else {
				now = App{y, now}
			}
		case rune:
			if y == '.' {
				end = true
				break
			}
			if now == nil {
				now = Var{Name(string([]rune{y}))}
			} else {
				now = App{now, Var{Name(string([]rune{y}))}}
			}
		}
		if end {
			break
		}
	}

	return now, nil
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
		if err != nil {
			break
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

func pushTk(x Token, xs []Token) []Token {
	return append(xs, x)
}

func popTk(xs []Token) ([]Token, Token, error) {
	if len(xs) == 0 {
		return nil, ' ', errors.New("empty stack")
	}
	last := xs[len(xs) - 1]
	init := xs[:len(xs) - 1]
	return init, last, nil
}
