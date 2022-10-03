package main

import (
	"bufio"
	"fmt"
	"go-parser/goparser"
	"math"
	"os"
	"regexp"
	"strings"
	// "time"
	"unsafe"

	// "github.com/google/uuid"
	// "strconv"
)

var (
	Transformation map[string]string = map[string]string{
		"diff":   `(?i)diff *\([_a-z0-9\+\-\*\/ \(\)]+\)`,
		"cumsum": `(?i)cumsum *\([_a-z0-9\+\-\*\/ \(\)]+\)`,
		"ma":     `(?i)ma *\([_a-z0-9\+\-\*\/, \(\)]+\)`,
	}
)

type Trans struct {
	Uid        string `json:"uid"`
	RegexExp   string `json:"regexExp"`
	RegexLocal string `json:"regexLocal"`
	Name       string `json:"name"`
}

type Transs struct {
	Transformation []Trans  `json:"transformation"`
	EquationList   []string `json:"equationList"`
}

func (T *Transs) GenerateEquationList(s string) {
	la := make([]string, 0)
	for _, ele := range T.Transformation {
		match, err := regexp.MatchString(ele.RegexExp, s)
		if err == nil {
			if match {
				// reg := regexp.MustCompile(ele.RegexLocal)
				// res := reg.ReplaceAllString(s, "${1}")
				reg := regexp.MustCompile(ele.RegexExp)
				matchStr := reg.FindAllString(s, -1)
				fmt.Println(matchStr, 343)
				for _, el := range matchStr {
					la = append(la, el)
				}
				res := reg.ReplaceAllString(s, "a")
				fmt.Println(res)
			}
		}
	}
	T.EquationList = la
}

func ptrToFunction(ptr uintptr) *func(...float64) float64 {
	p := unsafe.Pointer(ptr)
	fmt.Println(p, "pointer")
	return (*func(...float64) float64)(p)
}

///////////////////////////////////////////////////



func ExtractExpParenthesis(s string) string {
	ndx := 0
	firstParenthesis := strings.Index(s, "(")
	balance := 1
	lastParenthesis := 0
	for ndx < len(s) {
		if ndx == firstParenthesis {
			ndx++
			continue
		}
		if string(s[ndx]) == "(" {
			balance++
		} else if string(s[ndx]) == ")" {
			balance--
		}
		if balance == 0 {
			lastParenthesis = ndx
			break
		}
		ndx++
	}
	return s[firstParenthesis : lastParenthesis+1]
}

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

func RemoveRecuParenthesis(s string) []string {
	ls := make([]string, 0)
	ind := 0
	for {
		if ContainParenthesis(s) {
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

func SplitMultiExpression(s string) (string, error) {
	pattern := `[_a-zA-Z]\w* *[\*\/] *[_a-zA-Z]\w*`
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

func SplitAddExpression(s string) (string, error) {
	pattern := `[_a-zA-Z]\w* *[\+\-] *[_a-zA-Z]\w*`
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

func SplitRecuMultiExpression(s string, ini int) ([]string, error) {
	if strings.Index(s, "*") == -1 && strings.Index(s, "/") == -1 {
		return []string{s}, nil
	}
	pattern := `[_a-zA-Z]\w* *[\*\/] *[_a-zA-Z]\w*`
	sInit := strings.Split(s, "=")[0]
	s = strings.Split(s, "=")[1]
	matchStr := make([]string, 0)
	ind := 0
	for {
		match, err := regexp.MatchString(pattern, s)
		if err != nil {
			return matchStr, err
		}
		if match {
			s1, err := SplitMultiExpression(s)
			if err != nil {
				return matchStr, err
			}
			newVar := fmt.Sprintf("MULTI_DIVID%v_%v", ini, ind)
			s = strings.Replace(s, s1, newVar, 1)
			matchStr = append(matchStr, fmt.Sprintf("%s=%s", newVar, s1))
		} else {
			break
		}
		ind++
	}
	patt := `[_a-zA-Z]\w* *[\+\-] *[_a-zA-Z]\w*`
	match, _ := regexp.MatchString(patt, s)
	if match {
		newVar := fmt.Sprintf("MULTI_DIVID%v_%v", ini, ind)
		matchStr = append(matchStr, fmt.Sprintf("%s=%s", newVar, s))
	}
	lhs := strings.Split(matchStr[len(matchStr)-1], "=")[0]
	matchStr[len(matchStr)-1] = strings.Replace(matchStr[len(matchStr)-1], lhs, sInit, 1)
	return matchStr, nil
}

func SplitRecuAddExpression(s string, ini, init2 int) ([]string, error) {
	if strings.Index(s, "+") == -1 && strings.Index(s, "-") == -1 {
		return []string{s}, nil
	}
	sInit := strings.Split(s, "=")[0]
	s = strings.Split(s, "=")[1]
	pattern := `[_a-zA-Z]\w* *[\+\-] *[_a-zA-Z]\w*`
	matchStr := make([]string, 0)
	ind := 0
	for {
		match, err := regexp.MatchString(pattern, s)
		if err != nil {
			return matchStr, err
		}
		if match {
			s1, err := SplitAddExpression(s)
			if err != nil {
				return matchStr, err
			}
			newVar := fmt.Sprintf("ADD_SUB%v_%v_%v", ini, init2, ind)
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

func SplitRecuExpression(s string, ini int) ([]string, error) {
	matchStr, err := SplitRecuMultiExpression(s, ini)
	newMatchStr := []string{}
	if err != nil {
		return make([]string, 0), err
	} else {
		for ind, ele := range matchStr {
			newAdd, err := SplitRecuAddExpression(ele, ind, ini)
			if err != nil {
				return make([]string, 0), err
			}
			newMatchStr = append(newMatchStr, newAdd...)
		}
	}
	return newMatchStr, nil
}

func ExpressionGenerator(s string) []string {
	stLs := RemoveRecuParenthesis(s)
	ls := make([]string, 0)
	for ind, ele := range stLs {
		res, err := SplitRecuExpression(ele, ind)
		if err == nil {
			ls = append(ls, res...)
		}
	}
	return ls
}

///////////////////////////////////////////////////////////////////////////////////////

func main() {

	fmt.Println(ExpressionGenerator("((f>a+b+v) and (g>b+v)) or (u<v)"))

	Fu := &goparser.Function{
		Func: func(f ...float64) float64 { return math.NaN() },
	}
	Fu.GenerateFunctions("b+a+b/(a*a)+c", "test")

	fmt.Println(Fu.CallFunctionByMap(map[string]float64{"a":1, "b":2, "c":3, "d": 4}))

	i := &goparser.InEquaExpression{
		Inequality:"d<b+a+b/(a*a)+c",
	}
	i.GenerateFunction()
	fmt.Println(i)
	fmt.Println(i.CallFunctionByMap(map[string]float64{"a":1, "b":2, "c":3, "d": 4}))

	for {
		mm := Fu.Mapping
		fmt.Print(fmt.Sprintf("%v number: ", len(mm)))
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input.", err)
			return
		}
		input = strings.TrimSuffix(input, "\n")
		ff, err := goparser.SplitStringToFloat(input)
		if err == nil {
			fmt.Println(Fu.Func(ff...), Fu.Mapping)
		} else {
			fmt.Println("Input not slice of float64")
		}
	}
}
