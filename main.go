package main

import (
	"fmt"
	parser "go-parser/goparser"
)

func main() {
	s0 := "a*c/((v+e)-(g+t)*r/(r+t))+t"
	s1 := parser.ExpressionGenerator(s0)
	f, m, _ := parser.Generator(s1)
	fmt.Println(m, f([]float64{2, 3, 4, 5, 6, 7, 8}...))
}
