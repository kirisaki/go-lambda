package main

import (
	"errors"
	"fmt"
)

func parse(str string) (Expr, error){
	//var prev, now Expr
	stack := []rune{}
	rpn := []rune{}
	isParam := false
	prevParam := false
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
			prevParam = true
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
			if isParam {
				lamCnt++
			}
			if prevParam {
				rpn = append(rpn, '.')
				prevParam = false
			}
			rpn = append(rpn, c)
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
	fmt.Println(string(rpn))
	return nil, nil
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
