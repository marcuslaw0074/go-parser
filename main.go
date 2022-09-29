package main

import (
	"bufio"
	"fmt"
	parser "go-parser/goparser"
	"math"
	"os"
	"regexp"
	"strings"
)

var (
	Transformation map[string]string = map[string]string{
		"diff":   `(?i)diff *\([_a-z0-9\+\-\*\/ \(\)]+\)`,
		"cumsum": `(?i)cumsum *\([_a-z0-9\+\-\*\/ \(\)]+\)`,
		"ma":     `(?i)ma *\([_a-z0-9\+\-\*\/, \(\)]+\)`,
	}
)

type Trans struct {
	Uid   string `json:"uid"`
	RegexExp string `json:"regexExp"`
	RegexLocal string `json:"regexLocal"`
	Name  string `json:"name"`
}

type Transs struct {
	Transformation []Trans `json:"transformation"`
}

func (T *Transs) SplitTrans(s string) []string {
	for _, ele := range T.Transformation {
		match, err := regexp.MatchString(ele.RegexExp, s)
		if err == nil {
			if match {
				reg := regexp.MustCompile(ele.RegexLocal)
				res := reg.ReplaceAllString(s, "${1}")
				fmt.Println(res)
			}
		}
	}
	return make([]string, 0)
}

func main() {

	s := `Diff  (a+b*c)`
	match, err := regexp.MatchString(Transformation["diff"], s)
	if err == nil {
		if match {
			reg := regexp.MustCompile(`(?i)diff *`)
			res := reg.ReplaceAllString(s, "${1}")
			fmt.Println(res)
		}
	}

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
			fmt.Println(Fu.UseFunction(map[string]float64{"a": 1, "c": 2, "g": 3, "v": 4, "t": 5, "r": 6, "e": 7}))
		} else {
			fmt.Println("Input not slice of float64")
		}
	}
}
