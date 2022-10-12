package main

import (
	// "bufio"
	"fmt"
	"go-parser/goparser"
	"math"
	"strings"

	// "sync"
	// "math"
	// "os"
	"regexp"
	// "sort"
	// "strings"
	// "time"
	// "unsafe"
	// "github.com/google/uuid"
	// "strconv"
)

///////////////////////////////////////////////////

// func keys[K comparable, V any](m map[K]V) []K {
// 	keys := make([]K, 0, len(m))
// 	for k := range m {
// 		keys = append(keys, k)
// 	}
// 	return keys
// }

// type s[T any] struct {
// 	t T
// }

// type MyConstraint interface {
// 	int | int8 | int16 | int32 | int64
// }

// func MyFunc[T MyConstraint](input T) {
// 	// ...
// }

func GenerateIfElseFunction(s, Name string, Uid int) func(map[string]float64) float64 {
	i := &goparser.IfElseCondition{
		Uid: Uid,
		Name: Name,
		Conditions: GenerateConditionExpression(s),
	}
	f, _, _ := i.ConditionFunction()
	return f
}

func GenerateConditionExpression(s string) []goparser.ConditionExpression {
	res := []goparser.ConditionExpression{}
	pattern := `\{([^\{\}]|\n|\r)*\}`
	pattern_ := `\([^\{\}]*\)`
	matches := regexp.MustCompile(pattern).FindAllString(s, -1)
	split := regexp.MustCompile(pattern).Split(s, -1)
	for ind, ele := range split {
		fmt.Println(ind, ele)
		val := regexp.MustCompile(pattern_).FindAllString(ele, 1)
		if len(val) == 1 {
			fmt.Println(matches, split, val)
			sp := strings.ReplaceAll(regexp.MustCompile(`(?i)return`).Split(matches[ind], -1)[1], "}", "")
			res = append(res, goparser.ConditionExpression{
				Inequality: val[0][1 : len(val[0])-1],
				Expression: strings.ReplaceAll(sp, " ", ""),
			})
		}

	}
	return res
}

func main() {

	// MyFunc[int8](1)

	// vegetableSet := map[string]bool{
	// 	"potato":  true,
	// 	"cabbage": true,
	// 	"carrot":  true,
	// }
	// myInstance := s[int]{t: 1}
	// mm := keys(vegetableSet)
	// fmt.Println(mm, myInstance)

	// ls := []int{}
	// wg := sync.WaitGroup{}
	// wg.Add(4)
	// for ind := range []int{1,2,3,4} {
	// 	go func (i int)  {
	// 		ls = append(ls, i)
	// 		wg.Done()
	// 	}(ind)
	// }
	// wg.Wait()

	// fmt.Println(ls)

	// fmt.Println(SplitRecuAndOrExpression("f=g and o and h", 1, 2))

	// fmt.Println(InequalityExpressionGenerator("(Inequa_0 and Inequa_1) or (Inequa_2 or (Inequa_3 and Inequa_4 and Inequa_5))"))

	// H := GenerateIneq(map[string]string{"Inequa_0": "f>a+(b+v)", "Inequa_1": "g>b+v", "Inequa_2": "u>b+v*o"})
	// fmt.Println(H, 12345)
	// fmt.Println(H.Inequalities[2].Func([]float64{1, 2, 3, 4, 5, 6, 7}...))

	// fmt.Println(H.Inequalities[1].CallFunctionByMap(map[string]float64{"g":1, "b":1, "v":2}))
	// fmt.Println(InequalityExpressionGenerator("(Inequa_0 and Inequa_1 and Inequa_4) or (Inequa_2 or Inequa_3)"))
	// res, m, _ := ReplaceInequality("((f>a+(b+v)) and (g>b+v)) or (u>b+v*o)")
	// res2 := InequalityExpressionGenerator(res)
	// H := GenerateIneq(m)
	// for _, ele := range res2 {
	// 	exp := strings.Split(ele, "=")[0]
	// 	s := strings.Split(ele, "=")[1]
	// 	if strings.Index(s, " and ") > -1 {
	// 		s0 := strings.Split(s, " and ")[0]
	// 		s1 := strings.Split(s, " and ")[1]
	// 		ss0 := goparser.InEquaExpression{}
	// 		ss1 := goparser.InEquaExpression{}
	// 		for _, el := range H.Inequalities {
	// 			if el.Expression == s0 {
	// 				ss0 = el
	// 			}
	// 			if el.Expression == s1 {
	// 				ss1 = el
	// 			}
	// 		}
	// 		H.Inequalities = append(H.Inequalities, goparser.InEquaExpression{
	// 			Func: func(f ...float64) bool {
	// 				return ss0.Func(f...) && ss1.Func(f...)
	// 			},
	// 			Expression: exp,
	// 			Mapping:    H.Mapping,
	// 		})
	// 	} else if strings.Index(s, " or ") > -1 {
	// 		s0 := strings.Split(s, " or ")[0]
	// 		s1 := strings.Split(s, " or ")[1]
	// 		ss0 := goparser.InEquaExpression{}
	// 		ss1 := goparser.InEquaExpression{}
	// 		for _, el := range H.Inequalities {
	// 			if el.Expression == s0 {
	// 				ss0 = el
	// 			}
	// 			if el.Expression == s1 {
	// 				ss1 = el
	// 			}
	// 		}
	// 		H.Inequalities = append(H.Inequalities, goparser.InEquaExpression{
	// 			Func: func(f ...float64) bool {
	// 				return ss0.Func(f...) || ss1.Func(f...)
	// 			},
	// 			Expression: exp,
	// 		})
	// 	}
	// }
	// fmt.Println(res, m, res2, H)
	// fmt.Println(H.Mapping)
	// fmt.Println(H.Inequalities[len(H.Inequalities)-1].Func([]float64{1, 2, 3, 4, 5, 6, 117}...))
	// fmt.Println(goparser.ReplaceInequality("((f>a+(b+v)*a+a/b/a/a) and (g>b+v)) or (u>b+v*o)"))

	// f, m, _ := goparser.EnterExpression("((f>a+(b+v)*a+a/b/a/a) and (g>b+v)) or (u>b+v*o)")
	// ff := goparser.CallFunctionByMap(f, m)
	// fmt.Println(ff(map[string]float64{"a": 3, "b": 1, "f": 13, "g": 8, "o": 6, "u": 12, "v": 2}))

	fmt.Println(GenerateConditionExpression("if ((e<b-a) and (e<d-f)) { return (b-a)*c*d } else if ((e>b-a) or (e<d-f)) { return e*e }"))

	i := &goparser.IfElseCondition{
		Uid:  0,
		Name: "test",
		Conditions: []goparser.ConditionExpression{
			{
				Expression: "b*a*d",
				Inequality: "(c>b/a)",
			},
			{
				Expression: "a/b*d",
				Inequality: "(b>c/a)",
			},
		},
	}
	fff, _, _ := i.ConditionFunction()
	fmt.Println(fff(map[string]float64{
		"b": 10,
		"a": 3,
		"c": 3,
		"d": 5,
	}))

	fmt.Println(GenerateIfElseFunction(" if ((c>b/a)) { return b*a*d } else if ((b>c/a)) { return a/b*d }", "t", 0)(map[string]float64{
		"b": 10,
		"a": 3,
		"c": 3,
		"d": 5,
	}))

	ee, _ := goparser.InputExpression("(f>a+(b+v)*a+a/b/a/a) and (g>b+v) or (u>b+v*o)")
	fmt.Println(ee.CallFunctionByMap(map[string]float64{"a": 3, "b": 1, "f": 13, "g": 4, "o": 6, "u": 12, "v": 2}))

	Fu := &goparser.Function{}
	Fu.GenerateFunctions("a+(b+v)*a+a/b/a/a", "test")

	fmt.Println(Fu.CallFunctionByMap(map[string]float64{"a": 3, "b": 1, "v": 2}))

	// i := &goparser.InEquaExpression{
	// 	Inequality: "u>b+v*o",
	// }
	// i.GenerateFunction()
	// fmt.Println(i)
	// fmt.Println(i.CallFunctionByMap(map[string]float64{"b": 1, "v": 2, "o": 6, "u": 11}))

	// for {
	// 	mm := Fu.Mapping
	// 	fmt.Print(fmt.Sprintf("%v number: ", len(mm)))
	// 	reader := bufio.NewReader(os.Stdin)
	// 	input, err := reader.ReadString('\n')
	// 	if err != nil {
	// 		fmt.Println("An error occured while reading input.", err)
	// 		return
	// 	}
	// 	input = strings.TrimSuffix(input, "\n")
	// 	ff, err := goparser.SplitStringToFloat(input)
	// 	if err == nil {
	// 		fmt.Println(Fu.Func(ff...), Fu.Mapping)
	// 	} else {
	// 		fmt.Println("Input not slice of float64")
	// 	}
	// }
	fmt.Println(math.NaN())
}
