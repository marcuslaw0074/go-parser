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
		} else {
			fmt.Println("Input not slice of float64")
		}
	}
}
