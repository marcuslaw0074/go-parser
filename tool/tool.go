package tool

import (
	"fmt"
	"sort"
	"strings"
)

func SubSliceFloat(s []int, ls []float64) []float64 {
	max := len(ls)
	lss := []float64{}
	for _, ele := range s {
		if ele < max {
			lss = append(lss, ls[ele])
		}
	}
	return lss
}

func FindMapKeysSorted(s map[string]int) []string {
	la := []string{}
	for key := range s {
		la = append(la, key)
	}
	sort.Strings(la)
	return la
}

func StrContains(s []string, e string) int {
	for ind, a := range s {
		if a == e {
			return ind
		}
	}
	return -1
}

func MinMax(array []int) (int, int) {
	var max int = array[0]
	var min int = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

func FindMapValues(s map[string]int) []int {
	la := []int{}
	for _, val := range s {
		la = append(la, val)
	}
	return la
}

func FindValuesByKeys(keys []string, mapping map[string]int) []int {
	la := []int{}
	for _, key := range keys {
		for keyy, val := range mapping {
			if key == keyy {
				la = append(la, val)
			}
		}
	}
	return la
}

func ReplaceExpression(s, replace string, replaceLeft bool) string {
	operands := []string{}
	operator := ""
	switch {
	case strings.Contains(s, "+"):
		operands = strings.Split(s, "+")
		operator = "+"
	case strings.Contains(s, "-"):
		operands = strings.Split(s, "-")
		operator = "-"
	case strings.Contains(s, "*"):
		operands = strings.Split(s, "*")
		operator = "*"
	case strings.Contains(s, "/"):
		operands = strings.Split(s, "/")
		operator = "/"
	default:
		return ""
	}
	if replaceLeft {
		return fmt.Sprintf("(%s)%s%s", replace, operator, operands[1])
	} else {
		return fmt.Sprintf("%s%s(%s)", operands[0], operator, replace)
	}
}

