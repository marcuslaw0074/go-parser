package main

import (
	"fmt"
	"regexp"
	"strings"
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
	fmt.Println(newFirstParenthesis, "fde")
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
	fmt.Println(s, s[firstParenthesis:firstParenthesis+nd+1], "dfdfd")
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
			newVar := fmt.Sprintf("a%v", ind)
			s = strings.Replace(s, s1, newVar, 1)
			fmt.Println(s1, s)
			ls = append(ls, fmt.Sprintf("%s=%s", newVar, s1))
			ind++
		} else {
			newVar := fmt.Sprintf("a%v", ind)
			ls = append(ls, fmt.Sprintf("%s=%s", newVar, s))
			break
		}
	}
	return ls
}

func RemoveParenthesis(s string) []string {
	ls := make([]string, 0)
	ind := 0
	for {
		if ContainParenthesis(s) {
			s1 := ExtractExpParenthesis(s)
			newVar := fmt.Sprintf("a%v", ind)
			s = strings.Replace(s, s1, newVar, 1)
			ls = append(ls, fmt.Sprintf("%s=%s", newVar, s1))
			ind++
		} else {
			newVar := fmt.Sprintf("a%v", ind)
			ls = append(ls, fmt.Sprintf("%s=%s", newVar, s))
			break
		}
	}
	return ls
}

func SplitMultiExpression(s string) (string, error) {
	multiDivioper := `[_a-zA-Z]\w* *[\*\/] *[_a-zA-Z]\w*`
	matchStr := make([]string, 0)
	match, err := regexp.MatchString(multiDivioper, s)
	if err != nil {
		return "", err
	}
	if match {
		res, err := regexp.Compile(multiDivioper)
		if err != nil {
			return "", err
		}
		matchStr = res.FindAllString(s, 1)
	}
	return matchStr[0], nil
}

func SplitAddExpression(s string) (string, error) {
	multiDivioper := `[_a-zA-Z]\w* *[\+\-] *[_a-zA-Z]\w*`
	matchStr := make([]string, 0)
	match, err := regexp.MatchString(multiDivioper, s)
	if err != nil {
		return "", err
	}
	if match {
		res, err := regexp.Compile(multiDivioper)
		if err != nil {
			return "", err
		}
		matchStr = res.FindAllString(s, 1)
	}
	return matchStr[0], nil
}

func SplitRecuMultiExpression(s string) ([]string, error) {
	multiDivioper := `[_a-zA-Z]\w* *[\*\/] *[_a-zA-Z]\w*`
	matchStr := make([]string, 0)
	ind := 0
	for {
		match, err := regexp.MatchString(multiDivioper, s)
		if err != nil {
			return matchStr, err
		}
		if match {
			s1, err := SplitMultiExpression(s)
			if err != nil {
				return matchStr, err
			}
			newVar := fmt.Sprintf("MULTI_DIVID%v", ind)
			s = strings.Replace(s, s1, newVar, 1)
			matchStr = append(matchStr, fmt.Sprintf("%s=%s", newVar, s1))
		} else {
			break
		}
		ind++
	}
	newVar := fmt.Sprintf("MULTI_DIVID%v", ind)
	matchStr = append(matchStr, fmt.Sprintf("%s=%s", newVar, s))
	// fmt.Println(s)
	return matchStr, nil
}

func SplitRecuAddExpression(s string, ini int) ([]string, error) {
	fmt.Println(s)
	multiDivioper := `[_a-zA-Z]\w* *[\+\-] *[_a-zA-Z]\w*`
	matchStr := make([]string, 0)
	ind := 0
	for {
		match, err := regexp.MatchString(multiDivioper, s)
		if err != nil {
			return matchStr, err
		}
		if match {
			s1, err := SplitAddExpression(s)
			if err != nil {
				return matchStr, err
			}
			newVar := fmt.Sprintf("ADD_SUB%v_%v", ini, ind)
			fmt.Println(newVar)
			s = strings.Replace(s, s1, newVar, 1)
			matchStr = append(matchStr, fmt.Sprintf("%s=%s", newVar, s1))
		} else {
			break
		}
		ind++
	}
	newVar := fmt.Sprintf("ADD_SUB%v_%v", ini, ind)
	matchStr = append(matchStr, fmt.Sprintf("%s=%s", newVar, s))
	// fmt.Println(s)
	return matchStr, nil
}

func SplitRecuExpression(s string) ([]string, error) {
	matchStr, err := SplitRecuMultiExpression(s)
	newMatchStr := []string{}
	if err != nil {
		return make([]string, 0), err
	} else {
		for ind, ele := range matchStr {
			fmt.Println(ele, ind)
			newAdd, err := SplitRecuAddExpression(strings.Split(ele, "=")[1], ind)
			if err != nil {
				return make([]string, 0), err
			}
			newMatchStr = append(newMatchStr, newAdd...)
		}
	}
	return newMatchStr, nil
}

func main() {

	s0 := "a*c*v" //"(dd+ f+  b)"//"(a+(c*f))+(1+b)"
	fmt.Println(SplitRecuExpression(s0))
	// fmt.Println(SplitMultiExpression(s0))
	fmt.Println(RemoveRecuParenthesis(s0))
}
