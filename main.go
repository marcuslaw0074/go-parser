package main

import (
	"fmt"
	parser "go-parser/goparser"
)

// import (
// 	"fmt"
// 	parser "go-parser/goparser"
// )

// func main() {
// 	f, m, err := parser.Generator([]string{"c=b+a", "d=y/c", "e=a+d", "f=e+d"})
// 	if err == nil {
// 		fmt.Println(f([]float64{1, 2, 3}...), m)
// 	}
// }

func main() {

	s0 := "a*c/((v+e)-(g+t)*r)+t" //"a*c/((v+e)-(g+t)*r)+t" //"(dd+ f+  b)"//"(a+(c*f))+(1+b)"
	s1 := parser.ExpressionGenerator(s0)
	f, m, _ := parser.Generator(s1)
	fmt.Println(m, f([]float64{2, 3, 4, 5, 6, 7, 8}...))
	// fmt.Println(SplitRecuExpression(s0, 1))
	// fmt.Println(SplitRecuExpression(s0))
	// fmt.Println(SplitMultiExpression(s0))
	// fmt.Println(RemoveRecuParenthesis(s0))
}
