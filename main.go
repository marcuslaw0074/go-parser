package main

import (
	"bufio"
	"fmt"
	parser "go-parser/goparser"
	"math"
	"os"
	"strings"
)

func main() {
	Fu := &parser.Function{
		Func: func(f ...float64) float64 { return math.NaN() },
	}
	fmt.Println(Fu.Func([]float64{1,2,3,4,5,6}...))
	Fu.GenerateFunctions("a*c/((v+e)-(g+t)*r/(r+t))+t", "test")
	for {
		mm := Fu.Mapping
		fmt.Print(fmt.Sprintf("%v number: ", len(mm)))
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			return
		}
		input = strings.TrimSuffix(input, "\n")
		ff, err := parser.SplitStringToFloat(input)
		if err == nil {
			fmt.Println(Fu.Func(ff...))
			fmt.Println(Fu.UseFunction(map[string]float64{"a":1, "c":2, "g":3, "v":4, "t":5, "r":6, "e":7}))
		} else {
			fmt.Println("Input not slice of float64")
		}
	}
}
