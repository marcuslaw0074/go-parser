package main

import (
	// "context"
	// "encoding/json"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"math"
	"sort"
	"strconv"
	"time"
	// "io/ioutil"
	// "math"
	// "os"
	"regexp"
	"strings"
)

type Node struct {
	Uid        string      `json:"uid"`
	Numerical  float64     `json:"numerical"`
	Expression string      `json:"expression"`
	Origin     *Expression `json:"origin"`
}

type Expression struct {
	RightNode       Node                           `json:"rightNode"`
	Operation       string                         `json:"operation"`
	LeftNode        Node                           `json:"leftNode"`
	WholeExpression string                         `json:"wholeExpression"`
	Function        func(...float64) float64       `json:"-"`
	LocalFunction   func(float64, float64) float64 `json:"-"`
	Mapping         map[string]int                 `json:"mapping"`
}

func generateFunction(s1, s2 Node, operation string) (func(...float64) float64, error) {
	switch operation {
	case "+":
		if s1.Expression != "" && s2.Expression != "" {
			return func(f ...float64) float64 {
				return f[0] + f[1]
			}, nil
		} else if s1.Expression != "" {
			return func(f ...float64) float64 {
				return f[0] + s2.Numerical
			}, nil
		} else if s2.Expression != "" {
			return func(f ...float64) float64 {
				return s1.Numerical + f[0]
			}, nil
		}
	case "-":
		if s1.Expression != "" && s2.Expression != "" {
			return func(f ...float64) float64 {
				return f[0] - f[1]
			}, nil
		} else if s1.Expression != "" {
			return func(f ...float64) float64 {
				return f[0] - s2.Numerical
			}, nil
		} else if s2.Expression != "" {
			return func(f ...float64) float64 {
				return s1.Numerical - f[0]
			}, nil
		}
	case "*":
		if s1.Expression != "" && s2.Expression != "" {
			return func(f ...float64) float64 {
				return f[0] * f[1]
			}, nil
		} else if s1.Expression != "" {
			return func(f ...float64) float64 {
				return f[0] * s2.Numerical
			}, nil
		} else if s2.Expression != "" {
			return func(f ...float64) float64 {
				return s1.Numerical * f[0]
			}, nil
		}
	case "/":
		if s1.Expression != "" && s2.Expression != "" {
			return func(f ...float64) float64 {
				return f[0] / f[1]
			}, nil
		} else if s1.Expression != "" {
			return func(f ...float64) float64 {
				return f[0] / s2.Numerical
			}, nil
		} else if s2.Expression != "" {
			return func(f ...float64) float64 {
				return s1.Numerical / f[0]
			}, nil
		}
	default:
		return func(f ...float64) float64 { return math.NaN() }, errors.New("no such operators")
	}
	return func(f ...float64) float64 { return math.NaN() }, errors.New("wrong expressions")
}

func matchNodesOperation(str string) (Expression, error) {
	str = strings.Replace(str, " ", "", -1)
	reg1 := `^[+-]?([0-9]+([.][0-9]*)?|[.][0-9]+) *[\*\+\-\/] *[_a-zA-Z]\w*$`
	reg2 := `^[_a-zA-Z]\w* *[\*\+\-\/] *[+-]?([0-9]+([.][0-9]*)?|[.][0-9]+)$`
	reg3 := `^[_a-zA-Z]\w* *[\*\+\-\/] *[_a-zA-Z]\w*`
	for ind, ele := range []string{reg1, reg2, reg3} {
		match, err := regexp.MatchString(ele, str)
		if err != nil {
			return Expression{}, err
		}
		if match {
			if ind == 0 || ind == 1 {
				res, err := regexp.Compile(`[+-]?([0-9]+([.][0-9]*)?|[.][0-9]+)`)
				if err != nil {
					return Expression{}, err
				}
				match := res.FindAllString(str, 1)
				val, err := strconv.ParseFloat(match[0], 64)
				if err != nil {
					return Expression{}, err
				}
				res2, err2 := regexp.Compile(`[\*\+\-\/]`)
				if err2 != nil {
					return Expression{}, err
				}
				match2 := res2.FindAllString(str, 1)
				res3, err3 := regexp.Compile(`[_a-zA-Z]\w*`)
				if err3 != nil {
					return Expression{}, err
				}
				match3 := res3.FindAllString(str, 1)
				f1, f2 := findFunction(match2[0])
				uid1 := uuid.New().String()
				uid2 := uuid.New().String()
				if ind == 0 {
					return Expression{
						RightNode: Node{
							Uid:       uid1,
							Numerical: val,
						},
						Operation: match2[0],
						LeftNode: Node{
							Uid:        uid2,
							Expression: match3[0],
						},
						Mapping:       map[string]int{match3[0]: 0},
						Function:      f1,
						LocalFunction: f2,
					}, nil
				} else {
					return Expression{
						LeftNode: Node{
							Uid:       uuid.New().String(),
							Numerical: val,
						},
						Operation: match2[0],
						RightNode: Node{
							Uid:        uuid.New().String(),
							Expression: match3[0],
						},
						Mapping:       map[string]int{match3[0]: 0},
						Function:      f1,
						LocalFunction: f2,
					}, nil
				}
			} else {
				res2, err2 := regexp.Compile(`[\*\+\-\/]`)
				if err2 != nil {
					return Expression{}, err
				}
				match2 := res2.FindAllString(str, 1)
				f1, f2 := findFunction(match2[0])
				return Expression{
					RightNode: Node{
						Uid:        uuid.New().String(),
						Expression: strings.Split(str, match2[0])[0],
					},
					Operation: match2[0],
					LeftNode: Node{
						Uid:        uuid.New().String(),
						Expression: strings.Split(str, match2[0])[1],
					},
					Function:      f1,
					LocalFunction: f2,
				}, nil
			}
		}
	}
	return Expression{}, errors.New("not match")
}

func GenerateExpression(str string) (Expression, error) {
	str = strings.Replace(str, " ", "", -1)
	reg1 := `^[+-]?([0-9]+([.][0-9]*)?|[.][0-9]+) *[\*\+\-\/] *[_a-zA-Z]\w*$`
	reg2 := `^[_a-zA-Z]\w* *[\*\+\-\/] *[+-]?([0-9]+([.][0-9]*)?|[.][0-9]+)$`
	reg3 := `^[_a-zA-Z]\w* *[\*\+\-\/] *[_a-zA-Z]\w*`
	for ind, ele := range []string{reg1, reg2, reg3} {
		match, err := regexp.MatchString(ele, str)
		if err != nil {
			return Expression{}, err
		}
		if match {
			if ind == 0 || ind == 1 {
				res, err := regexp.Compile(`[+-]?([0-9]+([.][0-9]*)?|[.][0-9]+)`)
				if err != nil {
					return Expression{}, err
				}
				match := res.FindAllString(str, 1)
				val, err := strconv.ParseFloat(match[0], 64)
				if err != nil {
					return Expression{}, err
				}
				res2, err2 := regexp.Compile(`[\*\+\-\/]`)
				if err2 != nil {
					return Expression{}, err
				}
				match2 := res2.FindAllString(str, 1)
				res3, err3 := regexp.Compile(`[_a-zA-Z]\w*`)
				if err3 != nil {
					return Expression{}, err
				}
				match3 := res3.FindAllString(str, 1)
				f1, f2 := findFunction(match2[0])
				uid1 := uuid.New().String()
				uid2 := uuid.New().String()
				if ind == 0 {
					return Expression{
						RightNode: Node{
							Uid:       uid1,
							Numerical: val,
						},
						Operation: match2[0],
						LeftNode: Node{
							Uid:        uid2,
							Expression: match3[0],
						},
						Mapping:       map[string]int{match3[0]: 0},
						Function:      f1,
						LocalFunction: f2,
					}, nil
				} else {
					return Expression{
						LeftNode: Node{
							Uid:       uuid.New().String(),
							Numerical: val,
						},
						Operation: match2[0],
						RightNode: Node{
							Uid:        uuid.New().String(),
							Expression: match3[0],
						},
						Mapping:       map[string]int{match3[0]: 0},
						Function:      f1,
						LocalFunction: f2,
					}, nil
				}
			} else {
				res2, err2 := regexp.Compile(`[\*\+\-\/]`)
				if err2 != nil {
					return Expression{}, err
				}
				match2 := res2.FindAllString(str, 1)
				f1, f2 := findFunction(match2[0])
				return Expression{
					RightNode: Node{
						Uid:        uuid.New().String(),
						Expression: strings.Split(str, match2[0])[0],
					},
					Operation: match2[0],
					LeftNode: Node{
						Uid:        uuid.New().String(),
						Expression: strings.Split(str, match2[0])[1],
					},
					Function:      f1,
					LocalFunction: f2,
				}, nil
			}
		}
	}
	return Expression{}, errors.New("not match")
}

func matchNodesOperationi(str string) (bool, error) {
	str = strings.Replace(str, " ", "", -1)
	// reg1 := `^[+-]?([0-9]+([.][0-9]*)?|[.][0-9]+) *[\*\+\-\/] *[_a-zA-Z]\w*`
	reg := `^[_a-zA-Z]\w* *[\*\+\-\/] *[_a-zA-Z]\w*`
	match, err := regexp.MatchString(reg, str)
	return match, err
}

func strContains(s []string, e string) int {
	for ind, a := range s {
		if a == e {
			return ind
		}
	}
	return -1
}

func trans(st string) []string {
	st = strings.Replace(st, " ", "", -1)
	stLs := strings.Split(st, "=")
	return stLs
}

func transNew(st string) ([]string, error) {
	st = strings.Replace(st, " ", "", -1)
	stLs := strings.Split(st, "=")
	if len(stLs) != 2 {
		return []string{}, errors.New("equal sign")
	}
	operateLs := []string{"+", "-", "*", "/"}
	ct := 0
	lss := []string{}
	for _, ele := range operateLs {
		if strings.Contains(stLs[1], ele) {
			ct++
			lss = strings.Split(stLs[1], ele)
		}
	}
	if ct != 1 {
		return []string{}, errors.New("can hv exactly one operation")
	}
	if len(lss) != 2 {
		return []string{}, errors.New("can hv exactly one operation")
	}
	return []string{}, errors.New("can hv exactly one operation")
}

// func calculate(st string, allPoints []string) {
// 	stLs := trans
// 	if strContains(allPoints, )
// }

const (
	DIFFERENCE = iota
	CUMULATIVE_SUM
	MOVING_AVERAGE
	SIN
	COS
	TAN
	SQRT
)

func matchFunction(str string) (bool, error) {
	str = strings.Replace(str, " ", "", -1)
	reg := `^[a-zA-Z]\w*\([_a-zA-Z][\w]*\)`
	match, err := regexp.MatchString(reg, str)
	return match, err
}

func getFunction(str string) (string, error) {
	match, err := matchFunction(str)
	if err != nil {
		return "", err
	}
	if match {
		str = strings.Replace(str, " ", "", -1)
		funcName := strings.ToUpper(strings.Split(str, "(")[0])
		switch funcName {
		case "DIFFERENCE":
			return "DIFFERENCE", nil
		case "CUMULATIVE_SUM":
			return "CUMULATIVE_SUM", nil
		case "SIN":
			return "SIN", nil
		case "COS":
			return "COS", nil
		case "SQRT":
			return "SQRT", nil
		default:
			return "", errors.New("no such functions")
		}
	} else {
		return "", errors.New("wrong format")
	}
}

func matchOperation(str string) (bool, error) {
	str = strings.Replace(str, " ", "", -1)
	reg := `^[_a-zA-Z]\w* *[\*\+\-\/] *[_a-zA-Z]\w*`
	match, err := regexp.MatchString(reg, str)
	return match, err
}

func getOperation(str string) (func(float64, float64) float64, map[string]int, error) {
	def := func(f1, f2 float64) float64 { return 0 }
	match, err := matchOperation(str)
	if err != nil {
		return def, map[string]int{}, err
	}
	if match {
		str = strings.Replace(str, " ", "", -1)
		var operands []string
		switch {
		case strings.Contains(str, "+"):
			operands = strings.Split(str, "+")
			return func(f1, f2 float64) float64 { return f1 + f2 }, map[string]int{operands[0]: 0, operands[1]: 1}, nil
		case strings.Contains(str, "-"):
			operands = strings.Split(str, "-")
			return func(f1, f2 float64) float64 { return f1 - f2 }, map[string]int{operands[0]: 0, operands[1]: 1}, nil
		case strings.Contains(str, "*"):
			operands = strings.Split(str, "*")
			return func(f1, f2 float64) float64 { return f1 * f2 }, map[string]int{operands[0]: 0, operands[1]: 1}, nil
		case strings.Contains(str, "/"):
			operands = strings.Split(str, "/")
			return func(f1, f2 float64) float64 { return f1 / f2 }, map[string]int{operands[0]: 0, operands[1]: 1}, nil
		default:
			return def, map[string]int{}, errors.New("wrong operation")
		}
	} else {
		return def, map[string]int{}, errors.New("wrong format")
	}
}

func findOperation(str string) string {
	for _, ele := range []string{"+", "-", "*", "/"} {
		if strings.Index(str, ele) > -1 {
			return ele
		}
	}
	return ""
}

func findOperand(str string) []string {
	str = strings.Replace(str, " ", "", -1)
	regexp, _ := regexp.Compile(`[_a-zA-Z]\w*`)
	match := regexp.FindAllString(str, 2)
	return match
}

func findFunction(operator string) (func(...float64) float64, func(float64, float64) float64) {
	switch operator {
	case "+":
		return func(f ...float64) float64 { return f[0] + f[1] }, func(f1, f2 float64) float64 { return f1 + f2 }
	case "-":
		return func(f ...float64) float64 { return f[0] - f[1] }, func(f1, f2 float64) float64 { return f1 - f2 }
	case "*":
		return func(f ...float64) float64 { return f[0] * f[1] }, func(f1, f2 float64) float64 { return f1 * f2 }
	case "/":
		return func(f ...float64) float64 { return f[0] / f[1] }, func(f1, f2 float64) float64 { return f1 / f2 }
	default:
		return func(f ...float64) float64 { return math.NaN() }, func(f1, f2 float64) float64 { return math.NaN() }
	}
}

func findMapKeysSorted(s map[string]int) []string {
	la := []string{}
	for key := range s {
		la = append(la, key)
	}

	sort.Slice(la, func(i, j int) bool { return s[la[i]] < s[la[j]] })
	// sort.Strings(la)
	fmt.Println(la, "findMapKeysSorted", s)
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

func ReplaceExpressionBoth(s, replaceRight, replaceLeft string) string {
	operator := ""
	switch {
	case strings.Contains(s, "+"):
		operator = "+"
	case strings.Contains(s, "-"):
		operator = "-"
	case strings.Contains(s, "*"):
		operator = "*"
	case strings.Contains(s, "/"):
		operator = "/"
	default:
		return ""
	}
	return fmt.Sprintf("(%s)%s(%s)", replaceLeft, operator, replaceRight)
}

func findMapKeys(s map[string]int) []string {
	la := []string{}
	for key := range s {
		la = append(la, key)
	}
	return la
}

func findMapValues(s map[string]int) []int {
	la := []int{}
	for _, val := range s {
		la = append(la, val)
	}
	return la
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

type Graph struct {
	nodes []*GraphNode
}

type GraphNode struct {
	id    int
	edges map[int]int
}

func New() *Graph {
	return &Graph{
		nodes: []*GraphNode{},
	}
}

func (g *Graph) AddNode() (id int) {
	id = len(g.nodes)
	g.nodes = append(g.nodes, &GraphNode{
		id:    id,
		edges: make(map[int]int),
	})
	return
}

func (g *Graph) AddEdge(n1, n2 int, w int) {
	g.nodes[n1].edges[n2] = w
}

func (e *Expression) findEndNode() *Expression {
	if e.RightNode.Origin == nil && e.LeftNode.Origin == nil {
		return e
	} else if e.RightNode.Origin == nil {
		return e.LeftNode.Origin.findEndNode()
	} else {
		return e.RightNode.Origin.findEndNode()
	}
}

func (e *Expression) existSecondEndNode() bool {
	if e.RightNode.Origin == nil && e.LeftNode.Origin == nil {
		return false
	} else {
		return true
	}
}

func (e *Expression) findSecondEndNode() *Expression {
	if e.RightNode.Origin == nil && e.LeftNode.Origin == nil {
		return &Expression{}
	} else if e.RightNode.Origin == nil {
		if e.LeftNode.Origin.RightNode.Origin == nil && e.LeftNode.Origin.LeftNode.Origin == nil {
			fmt.Println(e.RightNode.Uid, e.LeftNode.Uid)
			return e
		}
		return e.LeftNode.Origin.findSecondEndNode()
	} else if e.LeftNode.Origin == nil {
		if e.RightNode.Origin.RightNode.Origin == nil && e.RightNode.Origin.LeftNode.Origin == nil {
			fmt.Println(e.RightNode.Uid, e.LeftNode.Uid)
			return e
		}
		return e.RightNode.Origin.findSecondEndNode()
	} else {
		if e.RightNode.Origin.RightNode.Origin == nil && e.RightNode.Origin.LeftNode.Origin == nil && e.LeftNode.Origin.RightNode.Origin == nil && e.LeftNode.Origin.LeftNode.Origin == nil {
			fmt.Println(e.RightNode.Uid, e.LeftNode.Uid)
			return e
		} else if e.LeftNode.Origin.RightNode.Origin == nil && e.LeftNode.Origin.LeftNode.Origin == nil {
			return e.RightNode.Origin.findSecondEndNode()
		} else {
			return e.LeftNode.Origin.findSecondEndNode()
		}
	}
}

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

func findValuesByKeys(keys []string, mapping map[string]int) []int {
	la := make(([]int), 0)
	for _, key := range keys {
		for keyy, val := range mapping {
			if key == keyy {
				la = append(la, val)
			}
		}
	}
	// fmt.Println("findValuesByKeys", la, keys, mapping)
	return la
}

func (e *Expression) MergeNode(expList ...Expression) *Expression {
	localMap := map[string]int{}
	if e.LeftNode.Origin == nil {
		localMap = e.RightNode.Origin.Mapping
		// fmt.Println(localMap)
		exp := *e.RightNode.Origin
		newMap := e.Mapping
		// fmt.Println("twomaps", newMap, localMap)
		// delete(newMap, e.RightNode.Expression)
		// fmt.Println(e.RightNode.Uid)
		delete(newMap, e.RightNode.Uid)
		originMap := []string{}
		localMapKeys := findMapKeysSorted(localMap)
		sortedKeys := findMapKeysSorted(newMap)
		for _, key := range sortedKeys {
			if strContains(localMapKeys, key) == -1 {
				_, max := MinMax(findMapValues(localMap))
				localMap[key] = max + 1
			} else {
				originMap = append(originMap, key)
				fmt.Println("originMap", originMap)
			}
		}
		newlocalValues := findValuesByKeys(localMapKeys, localMap)
		newValues := findValuesByKeys(sortedKeys, localMap)
		// fmt.Println(localMapKeys, newlocalValues, "gaer")
		function := func(f ...float64) float64 {
			if len(originMap) == 0 {
				return e.LocalFunction(f[newValues[0]], exp.Function(SubSliceFloat(newlocalValues, f)...))
			}
			return e.LocalFunction(f[newValues[0]], exp.Function(SubSliceFloat(newlocalValues, f)...))
		}
		e.WholeExpression = ReplaceExpression(e.WholeExpression, e.RightNode.Origin.WholeExpression, false)
		e.Function = function
		e.RightNode.Origin = nil
		e.Mapping = localMap
		return e
	} else if e.RightNode.Origin == nil {
		localMap = e.LeftNode.Origin.Mapping
		exp := *e.LeftNode.Origin
		newMap := e.Mapping
		// delete(newMap, e.LeftNode.Expression)
		delete(newMap, e.LeftNode.Uid)
		originMap := []string{}
		localMapKeys := findMapKeysSorted(localMap)
		sortedKeys := findMapKeysSorted(newMap)
		for _, key := range sortedKeys {
			if strContains(localMapKeys, key) == -1 {
				_, max := MinMax(findMapValues(localMap))
				localMap[key] = max + 1
			} else {
				originMap = append(originMap, key)
			}
		}
		newlocalValues := findValuesByKeys(localMapKeys, localMap)
		newValues := findValuesByKeys(sortedKeys, localMap)
		function := func(f ...float64) float64 {
			if len(originMap) == 0 {
				return e.LocalFunction(exp.Function(SubSliceFloat(newlocalValues, f)...), f[newValues[0]])
				// return e.LocalFunction(f[newValues[0]], exp.Function(SubSliceFloat(newlocalValues, f)...))
			}
			return e.LocalFunction(exp.Function(SubSliceFloat(newlocalValues, f)...), f[newValues[0]])
			// return e.LocalFunction(f[newValues[0]], exp.Function(SubSliceFloat(newlocalValues, f)...))
		}
		e.WholeExpression = ReplaceExpression(e.WholeExpression, e.LeftNode.Origin.WholeExpression, true)
		e.Function = function
		e.LeftNode.Origin = nil
		e.Mapping = localMap
	} else {
		localMapLeft := e.LeftNode.Origin.Mapping
		expLeft := *e.LeftNode.Origin
		localMapRight := e.RightNode.Origin.Mapping
		expRight := *e.RightNode.Origin
		fmt.Println("localMapLeft", localMapLeft, "localMapRight", localMapRight)
		originMap := []string{}
		localMapLeftKeys := findMapKeysSorted(localMapLeft)
		localMapRightKeys := findMapKeysSorted(localMapRight)
		sortedKeys := findMapKeysSorted(localMapRight)
		for _, key := range sortedKeys {
			if strContains(localMapLeftKeys, key) == -1 {
				_, max := MinMax(findMapValues(localMapLeft))
				localMapLeft[key] = max + 1
			} else {
				originMap = append(originMap, key)
				fmt.Println("originMap", originMap)
			}
		}
		newRightValues := findValuesByKeys(localMapRightKeys, localMapLeft)
		newLeftValues := findValuesByKeys(localMapLeftKeys, localMapLeft)
		function := func(f ...float64) float64 {
			if len(originMap) == 0 {
				return e.LocalFunction(expLeft.Function(SubSliceFloat(newLeftValues, f)...), expRight.Function(SubSliceFloat(newRightValues, f)...))
			}
			return e.LocalFunction(expLeft.Function(SubSliceFloat(newLeftValues, f)...), expRight.Function(SubSliceFloat(newRightValues, f)...))
		}
		fmt.Println("function", function([]float64{1, 2, 3, 4, 5}...))
		e.WholeExpression = ReplaceExpressionBoth(e.WholeExpression, e.RightNode.Origin.WholeExpression, e.LeftNode.Origin.WholeExpression)
		e.Function = function
		e.RightNode.Origin = nil
		e.LeftNode.Origin = nil
		e.Mapping = localMapLeft
	}

	return e
}

func (e *Expression) generateFunctionMap() (func(...float64) float64, map[string]int) {
	for e.existSecondEndNode() {
		// fmt.Println(e.LeftNode.Uid, e.RightNode.Uid)
		e.findSecondEndNode().MergeNode()
	}
	return e.Function, e.Mapping
}

type Equation struct {
	RHS       string `json:"rhs"`
	LHS       string `json:"lhs"`
	Relation  string `json:"relation"`
	RightVar  string `json:"rightVar"`
	LeftVar   string `json:"leftVar"`
	Operation string `json:"operation"`
}

type EquationList struct {
	Equations     []string    `json:"equations"`
	Graph         *Expression `json:"graph"`
	EquationsList []Equation  `json:"equationlist"`
}

func reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func (q *EquationList) generateEquationsList() {
	reverse(q.Equations)
	relation := "="
	equationsList := make([]Equation, 0)
	for _, ele := range q.Equations {
		rhs := strings.Split(ele, relation)[1]
		operator := findOperation(rhs)
		equationsList = append(equationsList, Equation{
			RHS:       rhs,
			LHS:       strings.Split(ele, relation)[0],
			Relation:  relation,
			RightVar:  strings.Split(ele, operator)[1],
			LeftVar:   strings.Split(ele, operator)[0],
			Operation: operator,
		})
	}
	q.EquationsList = equationsList

}

func (* Expression) ExtendExpression() {

}

func (q *EquationList) generateChildNode() error {
	if q.Graph == nil {
		fmt.Println("First")
		res, err := GenerateExpression(q.EquationsList[0].RHS)
		if err != nil {
			return err
		}
		q.Graph = &res
		q.EquationsList = q.EquationsList[1:]
	} else {
		fmt.Println("Later")
		equ := q.EquationsList[0]
		fmt.Println(equ)
		res, err := GenerateExpression(equ.RHS)
		if err != nil {
			return err
		}
		if len(q.Graph.LeftNode.Expression) > 0 && q.Graph.LeftNode.Expression == equ.LHS {
			q.Graph.LeftNode.Origin = &res
		} else if len(q.Graph.RightNode.Expression) > 0 && q.Graph.RightNode.Expression == equ.LHS {
			q.Graph.RightNode.Origin = &res
		} else if len(q.Graph.LeftNode.Expression) == 0 {
			f, err := strconv.ParseFloat(equ.LHS, 64)
			if err != nil {
				return err
			}
			if q.Graph.LeftNode.Numerical == f {
				q.Graph.LeftNode.Origin = &res
			}
		} else if len(q.Graph.RightNode.Expression) == 0 {
			f, err := strconv.ParseFloat(equ.LHS, 64)
			if err != nil {
				return err
			}
			if q.Graph.RightNode.Numerical == f {
				q.Graph.RightNode.Origin = &res
			}
		} else {
			return errors.New("errorrrrrrrrrrrrrr")
		}
		q.EquationsList = q.EquationsList[1:]
	}
	return nil
}

func (q *EquationList) generateGraph() {
	// ls := q.EquationsList
	// e := &Expression{

	// }
	// for ind, ele := range ls {
	// 	for _, el := range ls[ind:] {
	// 		operator := findOperation(ele.RHS)
	// 		f1, f2 := findFunction(operator)
	// 		e.Operation = operator
	// 		e.LocalFunction = f2
	// 		e.Function = f1
	// 		e.WholeExpression = ele.RHS
	// 	}
	// }
}

func main() {

	resss, _ := matchNodesOperation("543.543+bndsfk")

	// sfs, _ := json.Marshal(resss)
	// fmt.Println(string(sfs))

	fmt.Println(resss)

	resd, _ := strconv.ParseFloat("5321.532", 64)
	fmt.Println(resd)
	d := EquationList{
		Equations: []string{"c=a+b", "d=y/c", "e=a+d"},
	}
	d.generateEquationsList()
	d.generateChildNode()
	d.generateChildNode()
	d.generateChildNode()
	fmt.Println(d.Graph)

	sfs, _ := json.Marshal(d.Graph)
	fmt.Println(string(sfs))

	// str := "example*jvdka"
	// match, err := regexp.MatchString(`^[_a-zA-Z]\w* *[\*\+\-\/] *[_a-zA-Z]\w*`, str)
	// fmt.Println("Match: ", match, " Error: ", err)

	// exxx := &Expression{
	// 	RightNode: Node{
	// 		Uid: "4",
	// 		Expression: "f",
	// 		Origin: &Expression{
	// 			RightNode: Node{
	// 				Uid:        "2",
	// 				Expression: "e",
	// 				Origin: &Expression{
	// 					RightNode: Node{
	// 						Uid:        "0",
	// 						Expression: "a",
	// 					},
	// 					LeftNode: Node{
	// 						Uid:        "1",
	// 						Expression: "b",
	// 					},
	// 					Operation:     "-",
	// 					WholeExpression: "b-a",
	// 					LocalFunction: func(f1, f2 float64) float64 { fmt.Println(f1, f2); return f1 - f2 },
	// 					Function:      func(f ...float64) float64 { fmt.Println(f); return f[0] - f[1] },
	// 					Mapping:       map[string]int{"0": 1, "1": 0},
	// 				},
	// 			},
	// 			LeftNode: Node{
	// 				Uid:        "3",
	// 				Expression: "f",
	// 				Origin: &Expression{
	// 					RightNode: Node{
	// 						Uid:        "1",
	// 						Expression: "b",
	// 					},
	// 					LeftNode: Node{
	// 						Uid:        "0",
	// 						Expression: "a",
	// 					},
	// 					Operation:     "+",
	// 					WholeExpression: "a+b",
	// 					LocalFunction: func(f1, f2 float64) float64 { fmt.Println(f1, f2); return f1 + f2 },
	// 					Function:      func(f ...float64) float64 { fmt.Println(f); return f[0] + f[1] },
	// 					Mapping:       map[string]int{"0": 0, "1": 1},
	// 				},
	// 			},
	// 			Operation:     "/",
	// 			WholeExpression: "f/e",
	// 			LocalFunction: func(f1, f2 float64) float64 { fmt.Println(f1, f2); return f1 / f2 },
	// 			Function:      func(f ...float64) float64 { fmt.Println(f); return f[0] / f[1] },
	// 			Mapping:       map[string]int{"2": 1, "3": 0},
	// 		},
	// 	},
	// 	LeftNode: Node{
	// 		Uid: "5",
	// 		Expression: "g",
	// 	},
	// 	Operation:       "/",
	// 	WholeExpression: "g/f",
	// 	LocalFunction:   func(f1, f2 float64) float64 { fmt.Println(f1, f2); return f1 / f2 },
	// 	Function:        func(f ...float64) float64 { fmt.Println(f); return f[0] / f[1] },
	// 	Mapping:         map[string]int{"4": 1, "5": 0},
	// }

	exxx := &Expression{
		RightNode: Node{
			Uid:        "4",
			Expression: "e",
			Origin: &Expression{
				RightNode: Node{
					Uid:        "2",
					Expression: "c",
					Origin: &Expression{
						RightNode: Node{
							Uid:        "0",
							Expression: "a",
						},
						LeftNode: Node{
							Uid:        "1",
							Expression: "b",
						},
						Operation:       "/",
						WholeExpression: "b/a",
						LocalFunction:   func(f1, f2 float64) float64 { fmt.Println(f1, f2); return f1 / f2 },
						Function:        func(f ...float64) float64 { fmt.Println(f); return f[0] / f[1] },
						Mapping:         map[string]int{"0": 1, "1": 0},
					},
				},
				LeftNode: Node{
					Uid:        "3",
					Expression: "d",
					Origin: &Expression{
						RightNode: Node{
							Uid:        "8",
							Expression: "j",
							Origin: &Expression{
								RightNode: Node{
									Uid:        "1",
									Expression: "b",
								},
								LeftNode: Node{
									Uid:        "0",
									Expression: "a",
								},
								Operation:       "/",
								WholeExpression: "a/b",
								LocalFunction:   func(f1, f2 float64) float64 { fmt.Println(f1, f2); return f1 / f2 },
								Function:        func(f ...float64) float64 { fmt.Println(f, f[0]/f[1]); return f[0] / f[1] },
								Mapping:         map[string]int{"1": 1, "0": 0},
							},
						},
						LeftNode: Node{
							Uid:        "9",
							Expression: "k",
							Origin: &Expression{
								RightNode: Node{
									Uid:        "1",
									Expression: "b",
								},
								LeftNode: Node{
									Uid:        "0",
									Expression: "a",
								},
								Operation:       "+",
								WholeExpression: "a+b",
								LocalFunction:   func(f1, f2 float64) float64 { fmt.Println(f1, f2); return f1 + f2 },
								Function:        func(f ...float64) float64 { fmt.Println(f, f[0]+f[1]); return f[0] + f[1] },
								Mapping:         map[string]int{"1": 1, "0": 0},
							},
						},
						Operation:       "*",
						WholeExpression: "k*j",
						LocalFunction:   func(f1, f2 float64) float64 { fmt.Println(f1, f2); return f1 * f2 },
						Function:        func(f ...float64) float64 { fmt.Println(f); return f[0] * f[1] },
						Mapping:         map[string]int{"8": 1, "9": 0},
					},
				},
				Operation:       "-",
				WholeExpression: "d-c",
				LocalFunction:   func(f1, f2 float64) float64 { fmt.Println(f1, f2); return f1 - f2 },
				Function:        func(f ...float64) float64 { fmt.Println(f); return f[0] - f[1] },
				Mapping:         map[string]int{"2": 1, "3": 0},
			},
		},
		LeftNode: Node{
			Uid:        "5",
			Expression: "f",
			Origin: &Expression{
				RightNode: Node{
					Uid:        "2",
					Expression: "c",
					Origin: &Expression{
						RightNode: Node{
							Uid:        "0",
							Expression: "a",
						},
						LeftNode: Node{
							Uid:        "1",
							Expression: "b",
						},
						Operation:       "*",
						WholeExpression: "b*a",
						LocalFunction:   func(f1, f2 float64) float64 { fmt.Println(f1, f2); return f1 * f2 },
						Function:        func(f ...float64) float64 { fmt.Println(f); return f[0] * f[1] },
						Mapping:         map[string]int{"0": 1, "1": 0},
					},
				},
				LeftNode: Node{
					Uid:        "7",
					Expression: "h",
				},
				Operation:       "+",
				WholeExpression: "h+c",
				LocalFunction:   func(f1, f2 float64) float64 { fmt.Println(f1, f2); return f1 + f2 },
				Function:        func(f ...float64) float64 { fmt.Println(f); return f[0] + f[1] },
				Mapping:         map[string]int{"2": 1, "7": 0},
			},
		},
		Operation:       "/",
		WholeExpression: "f/e",
		LocalFunction:   func(f1, f2 float64) float64 { fmt.Println(f1, f2); return f1 / f2 },
		Function:        func(f ...float64) float64 { fmt.Println(f); return f[0] / f[1] },
		Mapping:         map[string]int{"4": 1, "5": 0},
	}

	fun, mapp := exxx.generateFunctionMap()

	fmt.Println(fun([]float64{1, 2, 3, 4, 5, 6, 7}...), mapp, exxx.WholeExpression)

	fmt.Println("DONE!!!")

	sf, _ := json.Marshal(exxx)
	fmt.Println(string(sf))
	exxx.findSecondEndNode().MergeNode()
	exxx.findSecondEndNode().MergeNode()
	sf, _ = json.Marshal(exxx)
	fmt.Println(string(sf))
	fmt.Println(exxx.Mapping)
	fmt.Println(exxx.Function([]float64{1, 2, 3, 4}...))

	fmt.Println(exxx.findSecondEndNode())

	exx := &Expression{
		RightNode: Node{
			Uid:        "4",
			Expression: "j",
			Origin: &Expression{
				RightNode: Node{
					Uid:        "2",
					Expression: "dd",
				},
				LeftNode: Node{
					Uid:        "3",
					Expression: "Hh",
				},
				Operation:     "-",
				LocalFunction: func(f1, f2 float64) float64 { fmt.Println(f1, f2); return f1 - f2 },
				Function:      func(f ...float64) float64 { fmt.Println(f); return f[0] - f[1] },
				Mapping:       map[string]int{"2": 0, "3": 1},
			},
		},
		LeftNode: Node{
			Uid:        "2",
			Expression: "dd",
			Origin: &Expression{
				RightNode: Node{
					Uid:        "2",
					Expression: "dd",
				},
				LeftNode: Node{
					Uid:        "1",
					Expression: "H",
				},
				Operation:     "-",
				LocalFunction: func(f1, f2 float64) float64 { fmt.Println(f1, f2); return f1 - f2 },
				Function:      func(f ...float64) float64 { fmt.Println(f); return f[0] - f[1] },
				Mapping:       map[string]int{"2": 0, "1": 1},
			},
		},
		Operation:     "/",
		LocalFunction: func(f1, f2 float64) float64 { fmt.Println(f1, f2); return f1 / f2 },
		Function:      func(f ...float64) float64 { fmt.Println(f); return f[0] / f[1] },
		Mapping:       map[string]int{"4": 0, "2": 1},
	}

	fff3 := *exx.MergeNode()
	exx = &fff3
	fmt.Println(fff3.Function([]float64{2, 5, 8}...))
	fmt.Println(fff3.Mapping)
	// fmt.Println(ex.RightNode.Origin.MergeNode().Mapping)
	sss1, _ := json.Marshal(exx)
	fmt.Println(string(sss1))
	// fmt.Println(ex)
	fmt.Println("")

	time.Sleep(time.Hour)

	ex := &Expression{
		RightNode: Node{
			Uid:        "4",
			Expression: "j",
			Origin: &Expression{
				RightNode: Node{
					Uid:        "2",
					Expression: "d",
					Origin: &Expression{
						RightNode: Node{
							Uid:        "0",
							Expression: "fs",
						},
						LeftNode: Node{
							Uid:        "1",
							Expression: "Hgrfs",
						},
						Operation:     "*",
						LocalFunction: func(f1, f2 float64) float64 { fmt.Println(f1, f2); return f1 * f2 },
						Function:      func(f ...float64) float64 { fmt.Println(f); return f[0] * f[1] },
						Mapping:       map[string]int{"0": 0, "1": 1},
					},
				},
				LeftNode: Node{
					Uid:        "3",
					Expression: "H",
				},
				Operation:     "-",
				LocalFunction: func(f1, f2 float64) float64 { fmt.Println(f1, f2); return f1 - f2 },
				Function:      func(f ...float64) float64 { fmt.Println(f); return f[0] - f[1] },
				Mapping:       map[string]int{"2": 0, "3": 1},
			},
		},
		LeftNode: Node{
			Uid:        "5",
			Expression: "HAHA",
		},
		Operation:     "/",
		LocalFunction: func(f1, f2 float64) float64 { fmt.Println(f1, f2); return f1 / f2 },
		Function:      func(f ...float64) float64 { fmt.Println(f); return f[0] / f[1] },
		Mapping:       map[string]int{"4": 0, "5": 1},
	}

	// sss, _ := json.Marshal(*&ex.RightNode.Origin.MergeNode().Mapping)
	fff := *ex.RightNode.Origin.MergeNode()
	fmt.Println(fff.Function([]float64{2, 5, 9}...))
	ex.RightNode.Origin = &fff
	fff2 := *ex.MergeNode()
	ex = &fff2
	fmt.Println(fff2.Function([]float64{2, 5, 8, 12}...))
	fmt.Println(fff)
	// fmt.Println(ex.RightNode.Origin.MergeNode().Mapping)
	sss, _ := json.Marshal(ex)
	fmt.Println(string(sss))
	// fmt.Println(ex)
	fmt.Println("")

	time.Sleep(time.Hour)

	// sss, _ := json.Marshal(*ex.findSecondEndNode())
	// fmt.Println(string(sss))
	// fmt.Println("")

	fmt.Println(*ex.findEndNode())
	ss, _ := json.Marshal(ex)
	fmt.Println(string(ss))
	v1 := Node{
		Expression: "j",
		Origin: &Expression{
			RightNode: Node{
				Expression: "d",
			},
			LeftNode: Node{
				Expression: "H",
			},
			Operation: "-",
			Function:  func(f ...float64) float64 { return f[0] + f[1] },
			Mapping:   map[string]int{"d": 0, "H": 1},
		},
	}
	v2 := Node{
		Expression: "HAHA",
	}
	res, _ := json.Marshal(v1)
	fmt.Println(string(res))
	fmt.Println("")
	fmt.Println(v2.Origin == nil)
	f, _ := generateFunction(v1, v2, "+")

	fmt.Println(f([]float64{1, 2}...))

	// fmt.Println(findOperand("difference + nvdasj"))
	fmt.Println(matchNodesOperation("5.42345 + nvdasj"))
	str := "            example * jvdka    "
	str = strings.Replace(str, " ", "", -1)

	// res :=`^[a-zA-Z]\w*\([_a-zA-Z][\w]*\)`

	regexp, err := regexp.Compile(`[_a-zA-Z]\w*`)

	match := regexp.FindAllString(str, 2)

	fmt.Println("Match: ", match, " Error: ", err)

	ls := []string{"a=b-a", "c=a*c", "e=c/d"}
	gg := []func(float64, float64) float64{}
	mm := []map[string]int{}
	allM := map[string]int{}
	for _, ele := range ls {
		f, m, err := getOperation(strings.Split(ele, "=")[1])
		if err != nil {
			continue
		}
		gg = append(gg, f)
		mm = append(mm, m)
		for _, newKey := range findMapKeysSorted(m) {
			if strContains(findMapKeysSorted(allM), newKey) == -1 {
				fmt.Println(newKey)
				v := findMapValues(allM)
				if len(v) == 0 {
					allM[newKey] = 0
				} else {
					_, max := MinMax(v)
					allM[newKey] = max + 1
				}
			}
		}
		// fmt.Println(allM)
	}
	fmt.Println(allM, "ghieakjd", mm)

	lss := make([][]string, 0)
	// ls1 := make([]string, 0)
	// ls2 := make([]string, 0)
	for _, ele := range ls {
		lss = append(lss, trans(ele))
		// ls1 = append(ls1, trans(ele)[0])
		// ls2 = append(ls2, trans(ele)[1])
	}
	ff := func(s ...float64) float64 {
		fmt.Println(492)
		if len(s) >= 2 {
			return s[1] - s[0]
		}
		return 0
	}
	ctt := 0
	for ind, ele := range lss {
		for ind2, ele2 := range lss[ind:] {
			if ele[0] != ele2[0] && ele[1] != ele2[1] {
				if strings.Contains(ele2[1], ele[0]) {
					lss[ind+ind2] = []string{ele2[0], strings.Replace(ele2[1], ele[0], fmt.Sprintf("(%s)", ele[1]), -1)}
					f, m, err := getOperation(ele[1])
					f2, m2, err2 := getOperation(ele2[1])
					fmt.Println("ctt", ctt)
					if strings.Index(ele2[1], ele[0]) > strings.Index(ele2[1], findOperation(ele2[1])) {
						ff := func(s ...float64) float64 {
							return f2(s[len(s)-1], ff(s[:len(s)-1]...))
						}
						fmt.Println(ff([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}...), "HKH")
					} else {
						ff := func(s ...float64) float64 {
							return f2(ff(s[:len(s)-1]...), s[len(s)-1])
						}
						fmt.Println(ff([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}...), "HKH")
					}

					fmt.Println(f(1, 2), m, err)
					fmt.Println(f2(1, 2), m2, err2)
					fmt.Println(ele, ele2)
					fmt.Println(ind, ind2)
				}
			}
		}
	}

	fmt.Printf("%s=%s", lss[len(lss)-1][0], lss[len(lss)-1][1])
	fmt.Println()

}
