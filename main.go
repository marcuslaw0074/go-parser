package main

import (
	// "bufio"
	"fmt"
	"go-parser/goparser"
	// "sync"
	// "math"
	// "os"
	"regexp"
	"sort"
	"strings"
	// "time"
	// "unsafe"
	// "github.com/google/uuid"
	// "strconv"
)

///////////////////////////////////////////////////

func MaxDepthParenthesis(s string) (int, int) {
	balance := 1
	maxDepth := 0
	newFirstParenthesis := 0
	for ind, ele := range s {
		if string(ele) == "(" {
			balance++
		} else if string(ele) == ")" {
			balance--
		}
		if balance > maxDepth {
			maxDepth = balance
			newFirstParenthesis = ind
		}
	}
	return maxDepth, newFirstParenthesis
}

func ExtractMaxDepthParenthesis(s string) string {
	nd := 0
	balance := 0
	depth, firstParenthesis := MaxDepthParenthesis(s)
	for _, ele := range s[firstParenthesis:] {
		if string(ele) == ")" {
			break
		}
		if balance == -depth {
			break
		}
		nd++
	}
	return s[firstParenthesis : firstParenthesis+nd+1]
}

func ContainParenthesis(s string) bool {
	if strings.Index(s, "(") > -1 && strings.Index(s, ")") > -1 {
		return true
	} else {
		return false
	}
}

func ContainAndOr(s string) bool {
	pattern := `(?i)[ \(\)]+(and|or)[ \(\)]+`
	match, err := regexp.MatchString(pattern, s)
	if err != nil {
		fmt.Println(err)
		return false
	} else if match {
		return true
	} else {
		return false
	}
}

func RemoveRecuAndOrParenthesis(s string) []string {
	ls := make([]string, 0)
	ind := 0
	for {
		if ContainParenthesis(s) && ContainAndOr(s) {
			s1 := ExtractMaxDepthParenthesis(s)
			newVar := fmt.Sprintf("EXPP%v", ind)
			s = strings.Replace(s, s1, newVar, 1)
			ls = append(ls, fmt.Sprintf("%s=%s", newVar, s1[1:len(s1)-1]))
			ind++
		} else {
			newVar := fmt.Sprintf("EXPP%v", ind)
			ls = append(ls, fmt.Sprintf("%s=%s", newVar, s))
			break
		}
	}
	return ls
}

func ExtractBalanceParenthesis(s string) string {
	balance := 0
	for ind, ele := range s {
		if string(ele) == "(" {
			balance++
		} else if string(ele) == ")" {
			balance--
		}
		if balance == 0 {
			s = s[:ind+1]
			break
		}
	}
	return s
}

func ReplaceInequality(s string) (string, map[string]string, error) {
	pattern := `\( *[_a-zA-Z]\w*[\>\<][\= ]*[_a-zA-Z \(\)\+\-\*\/]*\)`
	res, err := regexp.Compile(pattern)
	m := map[string]string{}
	if err != nil {
		return "", m, err
	} else {
		matchStr := res.FindAllString(s, -1)
		for ind, ele := range matchStr {
			ele = ExtractBalanceParenthesis(ele)
			key := fmt.Sprintf("Inequa_%v", ind)
			m[key] = ele[1 : len(ele)-1]
			s = strings.ReplaceAll(s, ele, key)
		}
		return s, m, nil
	}
}

func SplitRecuAndOrExpression(s string, ini, init2 int) ([]string, error) {
	if strings.Index(strings.ToLower(s), "and") == -1 && strings.Index(strings.ToLower(s), "or") == -1 {
		return []string{s}, nil
	}
	sInit := strings.Split(s, "=")[0]
	s = strings.Split(s, "=")[1]
	pattern := `(?i)[_a-z]\w* *(and|or) *[_a-z]\w*`
	matchStr := make([]string, 0)
	ind := 0
	for {
		match, err := regexp.MatchString(pattern, s)
		if err != nil {
			return matchStr, err
		}
		if match {
			s1, err := SplitAndOrExpression(s)
			if err != nil {
				return matchStr, err
			}
			newVar := fmt.Sprintf("ADD_SUB%v", ind)
			s = strings.Replace(s, s1, newVar, 1)
			matchStr = append(matchStr, fmt.Sprintf("%s=%s", newVar, s1))
		} else {
			break
		}
		ind++
	}
	lhs := strings.Split(matchStr[len(matchStr)-1], "=")[0]
	matchStr[len(matchStr)-1] = strings.Replace(matchStr[len(matchStr)-1], lhs, sInit, 1)
	return matchStr, nil
}

func SplitAndOrExpression(s string) (string, error) {
	pattern := `(?i)[_a-z]\w* *(and|or) *[_a-z]\w*`
	matchStr := make([]string, 0)
	match, err := regexp.MatchString(pattern, s)
	if err != nil {
		return "", err
	}
	if match {
		res, err := regexp.Compile(pattern)
		if err != nil {
			return "", err
		}
		matchStr = res.FindAllString(s, 1)
	}
	return matchStr[0], nil
}

func InequalityExpressionGenerator(s string) []string {
	stLs := RemoveRecuAndOrParenthesis(s)
	ls := make([]string, 0)
	for ind, ele := range stLs {
		res, err := SplitRecuAndOrExpression(ele, ind, 0)
		if err == nil {
			ls = append(ls, res...)
		}
	}
	return ls
}

func findMapKeysSorted(s map[string]int) []string {
	la := []string{}
	for key := range s {
		la = append(la, key)
	}
	sort.Slice(la, func(i, j int) bool { return s[la[i]] < s[la[j]] })
	return la
}

func findMapValues(s map[string]int) []int {
	la := []int{}
	for _, val := range s {
		la = append(la, val)
	}
	return la
}

func ContainInt(l []int, k int) int {
	for ind, ele := range l {
		if ele == k {
			return ind
		}
	}
	return -1
}

func ExtractSubSlice(i []float64, j []int) []float64 {
	k := []float64{}
	for _, ele := range j {
		k = append(k, i[ele])
	}
	return k
}

func CopyMap(old map[string]int) map[string]int {
	l := map[string]int{}
	for key, val := range old {
		l[key] = val
	}
	return l
}

func GenerateIneq(m map[string]string) *goparser.InEquaExpressions {
	ls := []string{}
	for key := range m {
		ls = append(ls, key)
	}
	ii := &goparser.InEquaExpressions{
		Inequalities: []goparser.InEquaExpression{},
	}
	sort.Strings(ls)
	newMap := map[string]int{}
	for _, ele := range ls {
		fmt.Println(m[ele], "m[ele]")
		i := &goparser.InEquaExpression{
			Inequality: m[ele],
		}
		j := &goparser.InEquaExpression{
			Inequality: m[ele],
			Expression: ele,
		}
		i.GenerateFunction()
		fmt.Println(i.Mapping, "DDDDDDDDD")
		keys := findMapKeysSorted(i.Mapping)
		for _, key := range keys {
			values := findMapValues(newMap)
			_, exists := newMap[key]
			if !exists {
				oldVval := i.Mapping[key]
				if ContainInt(values, oldVval) > -1 {
					_, max := goparser.MinMax(values)
					newMap[key] = max + 1
				} else {
					newMap[key] = oldVval
				}
			}
		}
		fmt.Println(m[ele], i.Mapping, "HHHHHHHH")
		j.Mapping = CopyMap(i.Mapping)
		for key := range i.Mapping {
			i.Mapping[key] = newMap[key]
		}
		j.Func = func(f ...float64) bool {
			fmt.Println(f, "FFFFFFFFFFFFFFFFF")
			g := map[string]float64{}
			for o, s := range i.Mapping {
				g[o] = f[s]
			}
			// i.Mapping = j.Mapping
			fmt.Println(i.Mapping, j.Mapping, g, "i.Mapping")
			return i.Func(goparser.GenerateFloatSlice(j.Mapping, g)...)
		}
		ii.Inequalities = append(ii.Inequalities, *j)
	}
	ii.Mapping = newMap
	return ii
}

func EnterExpression(s string) (func(...float64) bool, map[string]int, error) {
	res, m, err := ReplaceInequality("((f>a+(b+v)) and (g>b+v)) or (u>b+v*o)")
	if err != nil {
		return func(f ...float64) bool { return false }, map[string]int{}, err
	} else {
		res2 := InequalityExpressionGenerator(res)
		H := GenerateIneq(m)
		for _, ele := range res2 {
			exp := strings.Split(ele, "=")[0]
			s := strings.Split(ele, "=")[1]
			if strings.Index(s, " and ") > -1 {
				s0 := strings.Split(s, " and ")[0]
				s1 := strings.Split(s, " and ")[1]
				ss0 := goparser.InEquaExpression{}
				ss1 := goparser.InEquaExpression{}
				for _, el := range H.Inequalities {
					if el.Expression == s0 {
						ss0 = el
					}
					if el.Expression == s1 {
						ss1 = el
					}
				}
				H.Inequalities = append(H.Inequalities, goparser.InEquaExpression{
					Func: func(f ...float64) bool {
						fmt.Println(f, ss0.Func(f...), "and", ss1.Func(f...), f)
						return ss0.Func(f...) && ss1.Func(f...)
					},
					Expression: exp,
					Mapping:    H.Mapping,
				})
			} else if strings.Index(s, " or ") > -1 {
				s0 := strings.Split(s, " or ")[0]
				s1 := strings.Split(s, " or ")[1]
				ss0 := goparser.InEquaExpression{}
				ss1 := goparser.InEquaExpression{}
				for _, el := range H.Inequalities {
					if el.Expression == s0 {
						ss0 = el
					}
					if el.Expression == s1 {
						ss1 = el
					}
				}
				H.Inequalities = append(H.Inequalities, goparser.InEquaExpression{
					Func: func(f ...float64) bool {
						fmt.Println(f, ss0.Func(f...), "or", ss1.Func(f...), f)
						return ss0.Func(f...) || ss1.Func(f...)
					},
					Expression: exp,
				})
			}
		}
		return H.Inequalities[len(H.Inequalities)-1].Func, H.Mapping, nil
	}
}

func CallFunctionByMap(f func(...float64) bool, m map[string]int) func(map[string]float64) bool {
	return func(mm map[string]float64) bool {
		fmt.Println(mm, goparser.GenerateFloatSlice(m, mm))
		return f(goparser.GenerateFloatSlice(m, mm)...)
	}
}

///////////////////////////////////////////////////////////////////////////////////////

func main() {

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

	f, m, _ := EnterExpression("((f>a+(b+v)) and (g>b+v)) or (u>b+v*o)")
	ff := CallFunctionByMap(f, m)
	fmt.Println(ff(map[string]float64{"a": 3, "b": 1, "f": 34, "g": 8, "o": 6, "u": 12, "v": 2}))
	// fmt.Println(f([]float64{2, 10, 3, 4, 5, 12, 7}...), m)

	// Fu := &goparser.Function{
	// 	Func: func(f ...float64) float64 { return math.NaN() },
	// }
	// Fu.GenerateFunctions("b+v*o", "test")

	// fmt.Println(Fu.CallFunctionByMap(map[string]float64{"b": 1, "v": 2, "o": 6}))

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
}
