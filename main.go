package main

import (
	// "context"
	// "encoding/json"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"
	// "io/ioutil"
	// "math"
	// "os"
	"regexp"
	"strings"
)

type Node struct {
	Uid        string         `json:"uid"`
	Numerical  float64     `json:"numerical"`
	Expression string      `json:"expression"`
	Origin     *Expression `json:"origin"`
}

type Expression struct {
	RightNode     Node                           `json:"rightNode"`
	Operation     string                         `json:"operation"`
	LeftNode      Node                           `json:"leftNode"`
	Function      func(...float64) float64       `json:"-"`
	LocalFunction func(float64, float64) float64 `json:"-"`
	Mapping       map[string]int                 `json:"-"`
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
				if ind == 0 {
					return Expression{
						RightNode: Node{
							Numerical: val,
						},
						Operation: match2[0],
						LeftNode: Node{
							Expression: match3[0],
						},
						Mapping: map[string]int{match3[0]: 0},
						Function: func(f ...float64) float64 {
							if len(f) != 1 {
								return math.NaN()
							}
							return f[0]
						},
					}, nil
				} else {
					return Expression{
						LeftNode: Node{
							Numerical: val,
						},
						Operation: match2[0],
						RightNode: Node{
							Expression: match3[0],
						},
						Mapping: map[string]int{match3[0]: 0},
					}, nil
				}
			} else {
				res2, err2 := regexp.Compile(`[\*\+\-\/]`)
				if err2 != nil {
					return Expression{}, err
				}
				match2 := res2.FindAllString(str, 1)
				return Expression{
					RightNode: Node{
						Expression: strings.Split(str, match2[0])[0],
					},
					Operation: match2[0],
					LeftNode: Node{
						Expression: strings.Split(str, match2[0])[1],
					},
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

func (e *Expression) MergeNode() (*Expression) {
	localMap := map[string]int{}
	if e.LeftNode.Origin == nil {
		localMap = e.RightNode.Origin.Mapping
		exp := *e.RightNode.Origin
		newMap := e.Mapping
		delete(newMap, e.RightNode.Expression)
		originMap := []string{}
		localMapKeys := findMapKeys(localMap)
		for key := range newMap {
			if strContains(localMapKeys, key) == -1 {
				_, max := MinMax(findMapValues(localMap))
				localMap[key] = max+1
			} else {
				originMap = append(originMap, key)
			}
		}
		function := func(f ...float64) float64 {
			fmt.Println(f, localMap)
			if len(originMap) == 0 {
				return e.LocalFunction(f[len(f)-1], exp.Function(f[:len(f)-1]...))
			}
			return e.LocalFunction(f[len(f)-1], exp.Function(f[:len(f)-1]...))
		}
		e.Function = function
		e.RightNode.Origin = nil
		e.Mapping = localMap
		return e
	}

	return e
}

func main() {
	// str := "example*jvdka"
	// match, err := regexp.MatchString(`^[_a-zA-Z]\w* *[\*\+\-\/] *[_a-zA-Z]\w*`, str)
	// fmt.Println("Match: ", match, " Error: ", err)
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
							Uid:       "0",
							// Numerical: 4.2,
							Expression: "fs",
						},
						LeftNode: Node{
							Uid:        "1",
							Expression: "Hgrfs",
						},
						Operation:     "*",
						LocalFunction: func(f1, f2 float64) float64 { return f1 * f2 },
						Function:      func(f ...float64) float64 {fmt.Println(f); return f[0] * f[1] },
						Mapping:       map[string]int{"fs": 0, "Hgrfs": 1},
					},
				},
				LeftNode: Node{
					Uid:        "3",
					Expression: "H",
				},
				Operation:     "-",
				LocalFunction: func(f1, f2 float64) float64 {fmt.Println(f1, f2); return f1 - f2 },
				Function:      func(f ...float64) float64 { return f[0] - f[1] },
				Mapping:       map[string]int{"d": 0, "H": 1},
			},
		},
		LeftNode: Node{
			Uid:        "5",
			Expression: "HAHA",
		},
		Operation:     "/",
		LocalFunction: func(f1, f2 float64) float64 { return f1 / f2 },
		Function:      func(f ...float64) float64 { return f[0] / f[1] },
		Mapping:       map[string]int{"j": 0, "HAHA": 1},
	}

	// sss, _ := json.Marshal(*&ex.RightNode.Origin.MergeNode().Mapping)
	fff := *ex.RightNode.Origin.MergeNode()
	fmt.Println(fff.Function([]float64{2, 5, 8}...))
	ex.RightNode.Origin = &fff
	fff2 := *ex.MergeNode()
	ex = &fff2
	fmt.Println(fff2.Function([]float64{2, 5, 8, 8}...))
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
		for _, newKey := range findMapKeys(m) {
			if strContains(findMapKeys(allM), newKey) == -1 {
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
