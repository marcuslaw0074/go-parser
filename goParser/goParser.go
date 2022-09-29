package goparser

import (
	// "encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
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

func SplitRecuMultiExpression(s string, ini int) ([]string, error) {
	if strings.Index(s, "*") == -1 && strings.Index(s, "/") == -1 {
		return []string{s}, nil
	}
	multiDivioper := `[_a-zA-Z]\w* *[\*\/] *[_a-zA-Z]\w*`
	sInit := strings.Split(s, "=")[0]
	s = strings.Split(s, "=")[1]
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

func strContains(s []string, e string) int {
	for ind, a := range s {
		if a == e {
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

func findMapValues(s map[string]int) []int {
	la := []int{}
	for _, val := range s {
		la = append(la, val)
	}
	return la
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

func (e *Expression) MergeNode(expList ...Expression) *Expression {
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
			if strContains(localMapLeftKeys, key) == -1 {
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

	return e
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
				valnodeRight = SimpleNode{
					Uid:        uuid.New().String(),
					Expression: ele.RightVar,
				}
				valnodeLeft = SimpleNode{
					Uid:        uuid.New().String(),
					Expression: ele.LeftVar,
				}
			}
			indKey := ContainsNode(allNode, keynode)
			indValL := ContainsNode(allNode, valnodeLeft)
			indValR := ContainsNode(allNode, valnodeRight)
			if indValL > -1 {
				valnodeLeft = allNode[indValL]
			} else {
				allNode = append(allNode, valnodeLeft)
			}
			if indValR > -1 {
				valnodeRight = allNode[indValR]
			} else {
				allNode = append(allNode, valnodeRight)
			}
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

func (q *EquationList) AddChildNodeNew(ex *Expression, nodePath []SimpleNode) *Expression {
	r := ex
	for ind, ele := range nodePath {
		if len(r.LeftNode.Uid) == 0 && len(r.RightNode.Uid) == 0 {
			start := q.StartNode
			startEqu, err := findEquationEquationList(q.EquationsList, start)
			if err != nil {
				return ex
			}
			f1, f2 := findFunction(startEqu.Operation)
			*ex = Expression{
				RightNode: Node{
					Uid:        q.AdjList[start][1].Uid,
					Expression: q.AdjList[start][1].Expression,
					Numerical:  q.AdjList[start][1].Numerical,
				},
				LeftNode: Node{
					Uid:        q.AdjList[start][0].Uid,
					Expression: q.AdjList[start][0].Expression,
					Numerical:  q.AdjList[start][0].Numerical,
				},
				Operation:       startEqu.Operation,
				WholeExpression: startEqu.RHS,
				Mapping:         map[string]int{q.AdjList[start][0].Uid: 0, q.AdjList[start][1].Uid: 1},
				Function:        f1,
				LocalFunction:   f2,
			}
		} else {
			if ind == 0 {
				continue
			}
			parentNode, err := findEquationEquationList(q.EquationsList, nodePath[ind-1])
			currentNode, _ := findEquationEquationList(q.EquationsList, ele)
			if err != nil {
				return ex
			} else {
				if parentNode.LeftVar == ele.Expression {
					if r.LeftNode.Origin == nil {
						childrenNodes := q.AdjList[ele]
						if len(childrenNodes) == 0 {
							return ex
						}

						f1, f2 := findFunction(currentNode.Operation)
						r.LeftNode.Origin = &Expression{
							RightNode: Node{
								Uid:        childrenNodes[1].Uid,
								Expression: childrenNodes[1].Expression,
								Numerical:  childrenNodes[1].Numerical,
							},
							LeftNode: Node{
								Uid:        childrenNodes[0].Uid,
								Expression: childrenNodes[0].Expression,
								Numerical:  childrenNodes[0].Numerical,
							},
							Operation:       currentNode.Operation,
							WholeExpression: currentNode.RHS,
							Function:        f1,
							LocalFunction:   f2,
							Mapping:         map[string]int{childrenNodes[0].Uid: 0, childrenNodes[1].Uid: 1},
						}
					} else {
						r = r.LeftNode.Origin
					}
				} else {
					if r.RightNode.Origin == nil {
						childrenNodes := q.AdjList[ele]
						if len(childrenNodes) == 0 {
							return ex
						}
						f1, f2 := findFunction(currentNode.Operation)
						r.RightNode.Origin = &Expression{
							RightNode: Node{
								Uid:        childrenNodes[1].Uid,
								Expression: childrenNodes[1].Expression,
								Numerical:  childrenNodes[1].Numerical,
							},
							LeftNode: Node{
								Uid:        childrenNodes[0].Uid,
								Expression: childrenNodes[0].Expression,
								Numerical:  childrenNodes[0].Numerical,
							},
							Operation:       currentNode.Operation,
							WholeExpression: currentNode.RHS,
							Function:        f1,
							LocalFunction:   f2,
							Mapping:         map[string]int{childrenNodes[0].Uid: 0, childrenNodes[1].Uid: 1},
						}
					} else {
						r = r.RightNode.Origin
					}
				}
			}
		}
	}
	return ex
}

func (q *EquationList) GenerateNew(ex *Expression) *Expression {
	visitedList := [][]SimpleNode{}
	AdjDFS(q.AdjList, q.StartNode, []SimpleNode{}, &visitedList)
	SortVisitedList(&visitedList)
	for _, ele := range visitedList {
		q.AddChildNodeNew(ex, ele)
	}
	return ex
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

func Generator(equaLs []string) (func(...float64) float64, map[string]int, error) {
	d := EquationList{
		Equations: equaLs,
	}
	d.generateEquationsList()
	if CheckDuplicateVar(d.EquationsList) {
		return func(f ...float64) float64 { return math.NaN() }, make(map[string]int), errors.New("equations contains deplicated variables")
	}
	d.GenerateAdjList()
	d.GenerateAdjList()
	exxe := &Expression{}
	d.GenerateNew(exxe)
	functions, mapping := exxe.GenerateFunctionMap()
	return functions, d.uidToExpression(mapping), nil
}

type Function struct {
	Name    string                   `json:"name"`
	Func    func(...float64) float64 `json:"-"`
	Mapping map[string]int           `json:"-"`
	Uid     string                   `json:"uid"`
}

type FunctionStore struct {
	UserName string      `json:"username"`
	UserId   string      `json:"userid"`
	Methods  []*Function `json:"methods"`
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

func (F *Function) SaveFunctions(f func(...float64) float64, m map[string]int, uid, name string) error {
	F.Func = f
	F.Mapping = m
	F.Uid = uid
	F.Name = name
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
	f, m, err := Generator(ExpressionGenerator(s))
	if err == nil {
		uid := uuid.New().String()
		err := F.SaveFunctions(f, m, uid, name)
		if err != nil {
			return err
		}
		return nil
	} else {
		return err
	}
}
