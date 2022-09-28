package functiongenerator

import (
	"encoding/json"
	"errors"
	"fmt"
	parser "go-parser/expressionparser"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type SimpleNode struct {
	Uid        string  `json:"uid"`
	Numerical  float64 `json:"numerical"`
	Expression string  `json:"expression"`
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
	Graph         *parser.Expression          `json:"graph"`
	EquationsList []Equation                  `json:"equationlist"`
	AdjList       map[SimpleNode][]SimpleNode `json:"adjList"`
	AllNode       []SimpleNode                `json:"allNode"`
	StartNode     SimpleNode                  `json:"startNode"`
}

func findOperation(str string) string {
	for _, ele := range []string{"+", "-", "*", "/"} {
		if strings.Index(str, ele) > -1 {
			return ele
		}
	}
	return ""
}

func findAdjListKeys(m map[SimpleNode][]SimpleNode) []SimpleNode {
	ls := make([]SimpleNode, 0)
	for key := range m {
		ls = append(ls, key)
	}
	return ls
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

func matchAdjListKeys(keys []SimpleNode, exp string) (SimpleNode, error) {
	for _, ele := range keys {
		if ele.Expression == exp {
			return ele, nil
		}
	}
	return SimpleNode{}, errors.New("cannot find corresponding node")
}

func containsNode(values []SimpleNode, node SimpleNode) int {
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
			indKey := containsNode(allNode, keynode)
			indValL := containsNode(allNode, valnodeLeft)
			indValR := containsNode(allNode, valnodeRight)
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

func findEquationEquationList(eqs []Equation, node SimpleNode) (Equation, error) {
	for _, ele := range eqs {
		if ele.LHS == node.Expression {
			return ele, nil
		}
	}
	return Equation{}, errors.New("cannot find such equation")
}

func (q *EquationList) AddChildNode(ex *parser.Expression, nodePath []SimpleNode) *parser.Expression {
	ssss := nodePathtoList(nodePath)
	fmt.Println("nodePathtoList", ssss)
	r := ex
	for ind, ele := range nodePath {
		fmt.Println(ex, "INITIAL")
		if len(r.LeftNode.Uid) == 0 && len(r.RightNode.Uid) == 0 {
			start := q.StartNode
			startEqu, err := findEquationEquationList(q.EquationsList, start)
			if err != nil {
				return ex
			}
			f1, f2 := findFunction(startEqu.Operation)
			*ex = parser.Expression{
				RightNode: parser.Node{
					Uid:        q.AdjList[start][1].Uid,
					Expression: q.AdjList[start][1].Expression,
					Numerical:  q.AdjList[start][1].Numerical,
				},
				LeftNode: parser.Node{
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
						r.LeftNode.Origin = &parser.Expression{
							RightNode: parser.Node{
								Uid:        childrenNodes[1].Uid,
								Expression: childrenNodes[1].Expression,
								Numerical:  childrenNodes[1].Numerical,
							},
							LeftNode: parser.Node{
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
						r.RightNode.Origin = &parser.Expression{
							RightNode: parser.Node{
								Uid:        childrenNodes[1].Uid,
								Expression: childrenNodes[1].Expression,
								Numerical:  childrenNodes[1].Numerical,
							},
							LeftNode: parser.Node{
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

func (q *EquationList) GenerateExpression(ex *parser.Expression) *parser.Expression {
	visitedList := [][]SimpleNode{}
	AdjDFS(q.AdjList, q.StartNode, []SimpleNode{}, &visitedList)
	SortVisitedList(&visitedList)
	for ind, ele := range visitedList {
		q.AddChildNode(ex, ele)
		sfs, _ := json.Marshal(ex)
		fmt.Println(string(sfs), ind, "result")
	}
	return ex
}

func checkDuplicateVar(s []Equation) bool {
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

func Generator(equLs []string) (func(...float64) float64, map[string]int, error) {
	d := EquationList{
		Equations: equLs,
	}
	d.generateEquationsList()
	if checkDuplicateVar(d.EquationsList) {
		return func(f ...float64) float64 { return math.NaN() }, make(map[string]int), errors.New("equations contains deplicated variables")
	}
	d.GenerateAdjList()
	d.GenerateAdjList()
	exxe := &parser.Expression{}
	d.GenerateExpression(exxe)
	functions, mapping := exxe.GenerateFunctionMap()
	return functions, mapping, nil
}
