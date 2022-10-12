package goparser

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

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

type Node struct {
	Uid        string      `json:"uid"`
	Numerical  float64     `json:"numerical"`
	Expression string      `json:"expression"`
	Origin     *Expression `json:"origin"`
}

type SimpleNode struct {
	Uid        string  `json:"uid"`
	Numerical  float64 `json:"numerical"`
	Expression string  `json:"expression"`
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

type Equation struct {
	RHS       string `json:"rhs"`
	LHS       string `json:"lhs"`
	Relation  string `json:"relation"`
	RightVar  string `json:"rightVar"`
	LeftVar   string `json:"leftVar"`
	Operation string `json:"operation"`
}

type EquationList struct {
	Equations     []string                    `json:"equations"`
	Graph         *Expression                 `json:"graph"`
	EquationsList []Equation                  `json:"equationlist"`
	AdjList       map[SimpleNode][]SimpleNode `json:"adjList"`
	AllNode       []SimpleNode                `json:"allNode"`
	StartNode     SimpleNode                  `json:"startNode"`
}

const (
	DIFFERENCE = iota
	CUMULATIVE_SUM
	MOVING_AVERAGE
	SIN
	COS
	TAN
	SQRT
)

func reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func Contain[K comparable](l []K, k K) int {
	for ind, ele := range l {
		if ele == k {
			return ind
		}
	}
	return -1
}

func findOperation(str string) string {
	for _, ele := range []string{"+", "-", "*", "/"} {
		if strings.Index(str, ele) > -1 {
			return ele
		}
	}
	return ""
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

func findFunctionTwin(operator string) (func(...float64) float64, func(float64, float64) float64) {
	switch operator {
	case "+":
		return func(f ...float64) float64 { return 2 * f[0] }, func(f1, f2 float64) float64 { return 2 * f1 }
	case "-":
		return func(f ...float64) float64 { return 0 }, func(f1, f2 float64) float64 { return 0 }
	case "*":
		return func(f ...float64) float64 { return f[0] * f[0] }, func(f1, f2 float64) float64 { return f1 * f1 }
	case "/":
		return func(f ...float64) float64 { return 1 }, func(f1, f2 float64) float64 { return 1 }
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

func findMapValues[K comparable, V any](s map[K]V) []V {
	la := []V{}
	for _, val := range s {
		la = append(la, val)
	}
	return la
}

func SubSliceFloat[K comparable](s []int, ls []K) []K {
	max := len(ls)
	lss := []K{}
	for _, ele := range s {
		if ele < max {
			lss = append(lss, ls[ele])
		}
	}
	return lss
}

func findValuesByKeys[K comparable, V any](keys []K, mapping map[K]V) []V {
	la := make(([]V), 0)
	for _, key := range keys {
		for keyy, val := range mapping {
			if key == keyy {
				la = append(la, val)
			}
		}
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

func findAdjListKeys(m map[SimpleNode][]SimpleNode) []SimpleNode {
	ls := make([]SimpleNode, 0)
	for key := range m {
		ls = append(ls, key)
	}
	return ls
}

func matchAdjListKeys(keys []SimpleNode, exp string) (SimpleNode, error) {
	for _, ele := range keys {
		if ele.Expression == exp {
			return ele, nil
		}
	}
	return SimpleNode{}, errors.New("cannot find corresponding node")
}

func ContainsNode(values []SimpleNode, node SimpleNode) int {
	for ind, ele := range values {
		if len(node.Expression) > 0 {
			if ele.Expression == node.Expression {
				return ind
			}
		} else if ele.Numerical == node.Numerical {
			return ind
		}
	}
	return -1
}

func findEquationEquationList(eqs []Equation, node SimpleNode) (Equation, error) {
	for _, ele := range eqs {
		if ele.LHS == node.Expression {
			return ele, nil
		}
	}
	return Equation{}, errors.New("cannot find such equation")
}

func ContainsSimpleNode(ss []SimpleNode, s SimpleNode) int {
	for ind, ele := range ss {
		if ele == s {
			return ind
		}
	}
	return -1
}

func AdjDFS(adjList map[SimpleNode][]SimpleNode, startIndex SimpleNode, visited []SimpleNode, visitedList *[][]SimpleNode) {
	visited = append(visited, startIndex)
	for _, ele := range adjList[startIndex] {
		if ContainsSimpleNode(visited, ele) == -1 {
			viss := make([]SimpleNode, len(visited))
			copy(viss, visited)
			AdjDFS(adjList, ele, viss, visitedList)
		}
	}
	*visitedList = append(*visitedList, visited)
}

func SortVisitedList(visitedList *[][]SimpleNode) {
	vis := *visitedList
	sort.Slice(*visitedList, func(i, j int) bool {
		return len(vis[i]) < len(vis[j])
	})
}

func nodePathtoList(nodePath []SimpleNode) []string {
	ls := []string{}
	for _, ele := range nodePath {
		ls = append(ls, ele.Expression)
	}
	return ls
}

func CheckDuplicateVar(s []Equation) bool {
	allKeys := make(map[string]bool)
	for _, item := range s {
		if _, value := allKeys[item.LHS]; !value {
			allKeys[item.LHS] = true
		} else {
			return true
		}
	}
	return false
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
			return e
		}
		return e.LeftNode.Origin.findSecondEndNode()
	} else if e.LeftNode.Origin == nil {
		if e.RightNode.Origin.RightNode.Origin == nil && e.RightNode.Origin.LeftNode.Origin == nil {
			return e
		}
		return e.RightNode.Origin.findSecondEndNode()
	} else {
		if e.RightNode.Origin.RightNode.Origin == nil && e.RightNode.Origin.LeftNode.Origin == nil && e.LeftNode.Origin.RightNode.Origin == nil && e.LeftNode.Origin.LeftNode.Origin == nil {
			return e
		} else if e.LeftNode.Origin.RightNode.Origin == nil && e.LeftNode.Origin.LeftNode.Origin == nil {
			return e.RightNode.Origin.findSecondEndNode()
		} else {
			return e.LeftNode.Origin.findSecondEndNode()
		}
	}
}

func (e *Expression) MergeNode() {
	localMap := map[string]int{}
	if e.LeftNode.Origin == nil {
		localMap = e.RightNode.Origin.Mapping
		exp := *e.RightNode.Origin
		newMap := e.Mapping
		delete(newMap, e.RightNode.Uid)
		originMap := []string{}
		localMapKeys := findMapKeysSorted(localMap)
		sortedKeys := findMapKeysSorted(newMap)
		for _, key := range sortedKeys {
			if Contain(localMapKeys, key) == -1 {
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
				return e.LocalFunction(f[newValues[0]], exp.Function(SubSliceFloat(newlocalValues, f)...))
			}
			return e.LocalFunction(f[newValues[0]], exp.Function(SubSliceFloat(newlocalValues, f)...))
		}
		e.WholeExpression = ReplaceExpression(e.WholeExpression, e.RightNode.Origin.WholeExpression, false)
		e.Function = function
		e.RightNode.Origin = nil
		e.Mapping = localMap
	} else if e.RightNode.Origin == nil {
		localMap = e.LeftNode.Origin.Mapping
		exp := *e.LeftNode.Origin
		newMap := e.Mapping
		delete(newMap, e.LeftNode.Uid)
		originMap := []string{}
		localMapKeys := findMapKeysSorted(localMap)
		sortedKeys := findMapKeysSorted(newMap)
		for _, key := range sortedKeys {
			if Contain(localMapKeys, key) == -1 {
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
			}
			return e.LocalFunction(exp.Function(SubSliceFloat(newlocalValues, f)...), f[newValues[0]])
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
		originMap := []string{}
		localMapLeftKeys := findMapKeysSorted(localMapLeft)
		localMapRightKeys := findMapKeysSorted(localMapRight)
		sortedKeys := findMapKeysSorted(localMapRight)
		for _, key := range sortedKeys {
			if Contain(localMapLeftKeys, key) == -1 {
				_, max := MinMax(findMapValues(localMapLeft))
				localMapLeft[key] = max + 1
			} else {
				originMap = append(originMap, key)
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
		e.WholeExpression = ReplaceExpressionBoth(e.WholeExpression, e.RightNode.Origin.WholeExpression, e.LeftNode.Origin.WholeExpression)
		e.Function = function
		e.RightNode.Origin = nil
		e.LeftNode.Origin = nil
		e.Mapping = localMapLeft
	}
}

func (e *Expression) GenerateFunctionMap() (func(...float64) float64, map[string]int) {
	for e.existSecondEndNode() {
		e.findSecondEndNode().MergeNode()
	}
	return e.Function, e.Mapping
}

func (q *EquationList) generateEquationsList() {
	relation := "="
	equationsList := make([]Equation, 0)
	for _, ele := range q.Equations {
		rhs := strings.Split(ele, relation)[1]
		operator := findOperation(rhs)
		equationsList = append(equationsList, Equation{
			RHS:       rhs,
			LHS:       strings.Split(ele, relation)[0],
			Relation:  relation,
			RightVar:  strings.Split(rhs, operator)[1],
			LeftVar:   strings.Split(rhs, operator)[0],
			Operation: operator,
		})
	}
	if q.EquationsList == nil {
		q.EquationsList = equationsList
	} else {
		fmt.Println("Already generated EquationList!")
	}
}

func (q *EquationList) GenerateAdjList() {
	q.generateEquationsList()
	adjList := map[SimpleNode][]SimpleNode{}
	allNode := []SimpleNode{}
	for _, ele := range q.EquationsList {
		keys := findAdjListKeys(adjList)
		key, err := matchAdjListKeys(keys, ele.LHS)
		if err != nil {
			num1, err1 := strconv.ParseFloat(ele.RightVar, 64)
			num2, err2 := strconv.ParseFloat(ele.LeftVar, 64)
			keynode := SimpleNode{
				Uid:        uuid.New().String(),
				Expression: ele.LHS,
			}
			valnodeLeft := SimpleNode{}
			valnodeRight := SimpleNode{}
			if err1 == nil {
				valnodeRight = SimpleNode{
					Uid:       uuid.New().String(),
					Numerical: num1,
				}
				valnodeLeft = SimpleNode{
					Uid:        uuid.New().String(),
					Expression: ele.LeftVar,
				}
			} else if err2 == nil {
				valnodeRight = SimpleNode{
					Uid:        uuid.New().String(),
					Expression: ele.RightVar,
				}
				valnodeLeft = SimpleNode{
					Uid:       uuid.New().String(),
					Numerical: num2,
				}
			} else {
				if ele.RightVar == ele.LeftVar {
					uid := uuid.New().String()
					valnodeRight = SimpleNode{
						Uid:        uid,
						Expression: ele.RightVar,
					}
					valnodeLeft = SimpleNode{
						Uid:        uid,
						Expression: ele.LeftVar,
					}
				} else {
					valnodeRight = SimpleNode{
						Uid:        uuid.New().String(),
						Expression: ele.RightVar,
					}
					valnodeLeft = SimpleNode{
						Uid:        uuid.New().String(),
						Expression: ele.LeftVar,
					}
				}

			}
			indValL := ContainsNode(allNode, valnodeLeft)
			if indValL > -1 {
				valnodeLeft = allNode[indValL]
			} else {
				allNode = append(allNode, valnodeLeft)
			}
			indValR := ContainsNode(allNode, valnodeRight)
			if indValR > -1 {
				valnodeRight = allNode[indValR]
			} else {
				allNode = append(allNode, valnodeRight)
			}
			indKey := ContainsNode(allNode, keynode)
			if indKey > -1 {
				keynode = allNode[indKey]
			} else {
				allNode = append(allNode, keynode)
			}
			adjList[keynode] = []SimpleNode{valnodeLeft, valnodeRight}
		} else {
			fmt.Printf("Exists LHS Exxpression: %s", key.Expression)
		}
	}
	q.AdjList = adjList
	q.AllNode = allNode
	q.StartNode = allNode[len(allNode)-1]
}

func ExpressionChild(children []SimpleNode, equa Equation) Expression {
	if children[0].Uid == children[1].Uid {
		f1, f2 := findFunctionTwin(equa.Operation)
		return Expression{
			RightNode: Node{
				Uid:        children[1].Uid,
				Expression: children[1].Expression,
				Numerical:  children[1].Numerical,
			},
			LeftNode: Node{
				Uid:        children[0].Uid,
				Expression: children[0].Expression,
				Numerical:  children[0].Numerical,
			},
			Operation:       equa.Operation,
			WholeExpression: equa.RHS,
			Mapping:         map[string]int{children[0].Uid: 0},
			Function:        f1,
			LocalFunction:   f2,
		}
	} else {
		f1, f2 := findFunction(equa.Operation)
		return Expression{
			RightNode: Node{
				Uid:        children[1].Uid,
				Expression: children[1].Expression,
				Numerical:  children[1].Numerical,
			},
			LeftNode: Node{
				Uid:        children[0].Uid,
				Expression: children[0].Expression,
				Numerical:  children[0].Numerical,
			},
			Operation:       equa.Operation,
			WholeExpression: equa.RHS,
			Mapping:         map[string]int{children[0].Uid: 0, children[1].Uid: 1},
			Function:        f1,
			LocalFunction:   f2,
		}
	}
}

func (q *EquationList) AddChildNodeNew(ex *Expression, nodePath []SimpleNode) {
	r := ex
	for ind, ele := range nodePath {
		if len(r.LeftNode.Uid) == 0 && len(r.RightNode.Uid) == 0 {
			start := q.StartNode
			startEqu, err := findEquationEquationList(q.EquationsList, start)
			if err != nil {
				return
			}
			*ex = ExpressionChild(q.AdjList[start], startEqu)
		} else {
			if ind == 0 {
				continue
			}
			parentNode, err := findEquationEquationList(q.EquationsList, nodePath[ind-1])
			currentNode, _ := findEquationEquationList(q.EquationsList, ele)
			if err != nil {
				return
			} else {
				if parentNode.LeftVar == ele.Expression {
					if r.LeftNode.Origin == nil {
						childrenNodes := q.AdjList[ele]
						if len(childrenNodes) == 0 {
							return
						}
						res := ExpressionChild(childrenNodes, currentNode)
						r.LeftNode.Origin = &res
					} else {
						r = r.LeftNode.Origin
					}
				} else {
					if r.RightNode.Origin == nil {
						childrenNodes := q.AdjList[ele]
						if len(childrenNodes) == 0 {
							return
						}
						res := ExpressionChild(childrenNodes, currentNode)
						r.RightNode.Origin = &res
					} else {
						r = r.RightNode.Origin
					}
				}
			}
		}
	}
}

func (q *EquationList) GenerateNew(ex *Expression) {
	visitedList := [][]SimpleNode{}
	AdjDFS(q.AdjList, q.StartNode, []SimpleNode{}, &visitedList)
	SortVisitedList(&visitedList)
	for _, ele := range visitedList {
		q.AddChildNodeNew(ex, ele)
	}
}

func (q *EquationList) uidToExpression(m map[string]int) map[string]int {
	mm := make(map[string]int)
	for _, ele := range q.AllNode {
		for key, val := range m {
			if key == ele.Uid {
				mm[ele.Expression] = val
			}
		}
	}
	return mm
}

func FunctionGenerator(equaLs []string) (func(...float64) float64, map[string]int, error) {
	d := EquationList{
		Equations: equaLs,
	}
	d.generateEquationsList()
	if CheckDuplicateVar(d.EquationsList) {
		return func(f ...float64) float64 { return math.NaN() }, make(map[string]int), errors.New("equations contains deplicated variables")
	}
	d.GenerateAdjList()
	exxe := &Expression{}
	d.GenerateNew(exxe)
	functions, mapping := exxe.GenerateFunctionMap()
	return functions, d.uidToExpression(mapping), nil
}

type Function struct {
	Name       string                   `json:"name"`
	Expression string                   `json:"expression"`
	Func       func(...float64) float64 `json:"-"`
	Mapping    map[string]int           `json:"-"`
	Uid        string                   `json:"uid"`
}

type FunctionStore struct {
	UserName string      `json:"username"`
	UserId   string      `json:"userid"`
	Methods  []*Function `json:"methods"`
}

func (S *FunctionStore) AddFunctionByName(s, name string) error {
	if S.GetFunctionByName(name) > -1 {
		return errors.New("function name already exists, please enter another function name.")
	}
	F := &Function{}
	err := F.GenerateFunctions(s, name)
	if err != nil {
		return err
	} else {
		S.Methods = append(S.Methods, F)
		return nil
	}
}

func (S *FunctionStore) AddFunction(F *Function) error {
	name := F.Name
	s := F.Expression
	return S.AddFunctionByName(s, name)
}

func (S *FunctionStore) RemoveFunctionByName(name string) error {
	s := S.GetFunctionByName(name)
	if s == -1 {
		return errors.New(fmt.Sprintf("cannot delete function by name: %s", name))
	}
	S.Methods = append(S.Methods[:s], S.Methods[s+1:]...)
	return nil
}

func (S *FunctionStore) GetFunctionByName(name string) int {
	for ind, ele := range S.Methods {
		if ele.Name == name {
			return ind
		}
	}
	return -1
}

func (S *FunctionStore) UpdateFunctionByName(s, name string) error {
	index := S.GetFunctionByName(name)
	F := S.Methods[index]
	f, m, err := FunctionGenerator(ExpressionGenerator(s))
	if err != nil {
		return err
	}
	F.Func = f
	F.Mapping = m
	return nil
}

func (F *Function) ValidName() error {
	pattern := `^[_a-zA-Z]\w*`
	res, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	if len(res.FindAllString(F.Name, 1)) == 0 {
		err := errors.New("wrong name format, Can't save function!!")
		fmt.Println(err)
		return err
	} else {
		return nil
	}
}

func (F *Function) SaveFunctions(f func(...float64) float64, m map[string]int, uid, name, exp string) error {
	F.Func = f
	F.Mapping = m
	F.Uid = uid
	F.Name = name
	F.Expression = exp
	err := F.ValidName()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func SplitStringToFloat(s string) ([]float64, error) {
	ls := make([]float64, 0)
	ss := strings.Split(s, ",")
	for _, ele := range ss {
		pattern := `[+-]?([0-9]*[.])?[0-9]+`
		res, err := regexp.Compile(pattern)
		if err != nil {
			return []float64{}, err
		}
		ele = res.FindAllString(ele, 1)[0]
		f, err := strconv.ParseFloat(ele, 64)
		if err != nil {
			return []float64{}, err
		}
		ls = append(ls, f)
	}
	return ls, nil
}

func (F *Function) GenerateFunctions(s, name string) error {
	f, m, err := FunctionGenerator(ExpressionGenerator(s))
	if err == nil {
		uid := uuid.New().String()
		err := F.SaveFunctions(f, m, uid, name, s)
		if err != nil {
			return err
		}
		return nil
	} else {
		return err
	}
}

func GenerateFloatSlice(m map[string]int, f map[string]float64) []float64 {
	s := findMapKeysSorted(m)
	k := []float64{}
	for _, ele := range s {
		k = append(k, f[ele])
	}
	return k
}

func (F *Function) CallFunctionByMap(f map[string]float64) float64 {
	return F.Func(GenerateFloatSlice(F.Mapping, f)...)
}

func (F *Function) UpdateFunctions(s, name string) error {
	return nil
}

func (F *Function) UseFunction(input map[string]float64) (float64, error) {
	keys := findMapKeysSorted(F.Mapping)
	sort.Slice(keys, func(i, j int) bool { return F.Mapping[keys[i]] < F.Mapping[keys[j]] })
	ls := []float64{}
	for _, ele := range keys {
		val, exist := input[ele]
		if exist {
			ls = append(ls, val)
		} else {
			return math.NaN(), errors.New(fmt.Sprintf("input map doesnt contain key: %v", val))
		}
	}
	res := F.Func(ls...)
	if math.IsNaN(res) {
		return math.NaN(), errors.New("function is not initialized yet")
	}
	return F.Func(ls...), nil
}

const (
	LE  = "<"
	GE  = ">"
	LEQ = "<="
	GEQ = ">="
)

type InEquaExpression struct {
	Inequality string                `json:"inequality"`
	Func       func(...float64) bool `json:"-"`
	Mapping    map[string]int        `json:"-"`
	Operator   string                `json:"operator"`
	Expression string                `json:"expression"`
	Result     bool                  `json:"result"`
}

type InEquaExpressions struct {
	Inequalities []InEquaExpression    `json:"inequalities"`
	Func         func(...float64) bool `json:"-"`
	Mapping      map[string]int        `json:"-"`
	Expression   string                `json:"expression"`
	Result       bool                  `json:"result"`
}

func GenerateFunctions(sts []string) {
	m := map[string]*InEquaExpression{}
	for _, ele := range sts {
		con := false
		for _, el := range []string{LEQ, GEQ, LE, GE} {
			if strings.Contains(ele, el) {
				m[strings.Split(ele, "=")[0]] = &InEquaExpression{
					Inequality: strings.Split(ele, "=")[1],
				}
				m[strings.Split(ele, "=")[0]].GenerateFunction()
				con = true
			}
		}
		if !con {
			m[strings.Split(ele, "=")[0]] = &InEquaExpression{
				Expression: strings.Split(ele, "=")[1],
			}
		}
	}
}

func (i *InEquaExpression) GenerateFunction() error {
	s := []string{}
	for _, ele := range []string{LEQ, GEQ, LE, GE} {
		s = strings.Split(i.Inequality, ele)
		if len(s) > 1 {
			i.Operator = ele
			break
		}
	}
	if len(s) != 2 {
		return errors.New("more than one inequality operator")
	}
	fu, m, err := FunctionGenerator(ExpressionGenerator(s[1]))
	if err != nil {
		return err
	}
	_, max := MinMax(findMapValues(m))
	m[s[0]] = max + 1
	i.Mapping = m
	i.Func = func(f ...float64) bool {
		switch i.Operator {
		case LEQ:
			// fmt.Println(f[len(f)-1] ,LEQ, fu(f[:len(f)-1]...))
			return f[len(f)-1] <= fu(f[:len(f)-1]...)
		case GEQ:
			// fmt.Println(f[len(f)-1] ,GEQ, fu(f[:len(f)-1]...))
			return f[len(f)-1] >= fu(f[:len(f)-1]...)
		case LE:
			// fmt.Println(f[len(f)-1] ,LE, fu(f[:len(f)-1]...))
			return f[len(f)-1] < fu(f[:len(f)-1]...)
		case GE:
			// fmt.Println(f[len(f)-1] ,GE, fu(f[:len(f)-1]...))
			return f[len(f)-1] > fu(f[:len(f)-1]...)
		default:
			panic("invalid operator")
		}
	}
	return nil
}

func (i *InEquaExpression) CallFunctionByMap(f map[string]float64) bool {
	return i.Func(GenerateFloatSlice(i.Mapping, f)...)
}

func ContainAndOr(s string) bool {
	pattern := `(?i)[ \(\)]+(and|or)[ \(\)]+`
	match, err := regexp.MatchString(pattern, s)
	if err != nil {
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

func ExtractSubSlice[K comparable](i []K, j []int) []K {
	k := []K{}
	for _, ele := range j {
		k = append(k, i[ele])
	}
	return k
}

func CopyMap[K comparable, V any](old map[K]V) map[K]V {
	l := map[K]V{}
	for key, val := range old {
		l[key] = val
	}
	return l
}

func GenerateIneq(m map[string]string) *InEquaExpressions {
	ls := []string{}
	for key := range m {
		ls = append(ls, key)
	}
	ii := &InEquaExpressions{
		Inequalities: []InEquaExpression{},
	}
	sort.Strings(ls)
	newMap := map[string]int{}
	for _, ele := range ls {
		i := &InEquaExpression{
			Inequality: m[ele],
		}
		j := &InEquaExpression{
			Inequality: m[ele],
			Expression: ele,
		}
		i.GenerateFunction()
		keys := findMapKeysSorted(i.Mapping)
		for _, key := range keys {
			values := findMapValues(newMap)
			_, exists := newMap[key]
			if !exists {
				oldVval := i.Mapping[key]
				if Contain(values, oldVval) > -1 {
					_, max := MinMax(values)
					newMap[key] = max + 1
				} else {
					newMap[key] = oldVval
				}
			}
		}
		j.Mapping = CopyMap(i.Mapping)
		for key := range i.Mapping {
			i.Mapping[key] = newMap[key]
		}
		j.Func = func(f ...float64) bool {
			g := map[string]float64{}
			for o, s := range i.Mapping {
				g[o] = f[s]
			}
			return i.Func(GenerateFloatSlice(j.Mapping, g)...)
		}
		ii.Inequalities = append(ii.Inequalities, *j)
	}
	ii.Mapping = newMap
	return ii
}

func (e *Express) CallFunctionByMap(f map[string]float64) bool {
	return e.Function(GenerateFloatSlice(e.Mapping, f)...)
}

func InputExpression(s string) (*Express, error) {
	e := &Express{}
	f, m, err := EnterExpression(s)
	if err != nil {
		return e, err
	} else {
		e.Function = f
		e.Mapping = m
		e.FunctionMap = CallFunctionByMap(f, m)
		return e, nil
	}
}

func EnterExpression(s string) (func(...float64) bool, map[string]int, error) {
	res, m, err := ReplaceInequality(s)
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
				ss0 := InEquaExpression{}
				ss1 := InEquaExpression{}
				for _, el := range H.Inequalities {
					if el.Expression == s0 {
						ss0 = el
					}
					if el.Expression == s1 {
						ss1 = el
					}
				}
				H.Inequalities = append(H.Inequalities, InEquaExpression{
					Func: func(f ...float64) bool {
						return ss0.Func(f...) && ss1.Func(f...)
					},
					Expression: exp,
					Mapping:    H.Mapping,
				})
			} else if strings.Index(s, " or ") > -1 {
				s0 := strings.Split(s, " or ")[0]
				s1 := strings.Split(s, " or ")[1]
				ss0 := InEquaExpression{}
				ss1 := InEquaExpression{}
				for _, el := range H.Inequalities {
					if el.Expression == s0 {
						ss0 = el
					}
					if el.Expression == s1 {
						ss1 = el
					}
				}
				H.Inequalities = append(H.Inequalities, InEquaExpression{
					Func: func(f ...float64) bool {
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
		return f(GenerateFloatSlice(m, mm)...)
	}
}

type Funcs interface {
	func(...float64) bool | func(map[string]float64) bool
}

type Express struct {
	Function    func(...float64) bool
	Mapping     map[string]int
	FunctionMap func(map[string]float64) bool
}

type ConditionFunctions struct {
	ExpressionFunction func(map[string]float64) float64
	ExpressionMapping  map[string]int
	InequalityFunction func(map[string]float64) bool
	InequalityMapping  map[string]int
}

type ConditionExpression struct {
	Inequality string
	Expression string
}

type IfElseCondition struct {
	Uid        int
	Name       string
	Conditions []ConditionExpression
}

func (i *IfElseCondition) ConditionFunction() (func(map[string]float64) float64, []string, error) {
	allFunction := make([]*ConditionFunctions, 0)
	allPoints := []string{}
	for ind, mapping := range i.Conditions {
		ex := &Function{}
		ex.GenerateFunctions(mapping.Expression, fmt.Sprintf("%v", ind))
		inq, _ := InputExpression(mapping.Inequality)
		for key := range ex.Mapping {
			if Contain(allPoints, key) == -1 {
				allPoints = append(allPoints, key)
			}
		}
		for key := range inq.Mapping {
			if Contain(allPoints, key) == -1 {
				allPoints = append(allPoints, key)
			}
		}
		allFunction = append(allFunction, &ConditionFunctions{
			ExpressionFunction: ex.CallFunctionByMap,
			ExpressionMapping:  ex.Mapping,
			InequalityFunction: inq.FunctionMap,
			InequalityMapping:  inq.Mapping,
		})
	}
	return func(input map[string]float64) float64 {
		for _, ele := range allFunction {
			if ele.InequalityFunction(input) {
				return ele.ExpressionFunction(input)
			}
		}
		return math.NaN()
	}, allPoints, nil
}

func (s *ConditionExpression) FillStruct(m map[string]interface{}) error {
	for k, v := range m {
		err := SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}
