package expressionparser

import (
	"go-parser/tool"
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
	Mapping         map[string]int                 `json:"-"`
}

func (e *Expression) FindEndNode() *Expression {
	if e.RightNode.Origin == nil && e.LeftNode.Origin == nil {
		return e
	} else if e.RightNode.Origin == nil {
		return e.LeftNode.Origin.FindEndNode()
	} else {
		return e.RightNode.Origin.FindEndNode()
	}
}

func (e *Expression) FindSecondEndNode() *Expression {
	if e.RightNode.Origin == nil && e.LeftNode.Origin == nil {
		return &Expression{}
	} else if e.RightNode.Origin == nil {
		if e.LeftNode.Origin.RightNode.Origin == nil && e.LeftNode.Origin.LeftNode.Origin == nil {
			return e
		}
		return e.LeftNode.Origin.FindSecondEndNode()
	} else if e.LeftNode.Origin == nil {
		if e.RightNode.Origin.RightNode.Origin == nil && e.RightNode.Origin.LeftNode.Origin == nil {
			return e
		}
		return e.RightNode.Origin.FindSecondEndNode()
	} else {
		if e.RightNode.Origin.RightNode.Origin == nil && e.RightNode.Origin.LeftNode.Origin == nil && e.LeftNode.Origin.RightNode.Origin == nil && e.LeftNode.Origin.LeftNode.Origin == nil {
			return e
		} else if e.LeftNode.Origin.RightNode.Origin == nil && e.LeftNode.Origin.LeftNode.Origin == nil {
			return e.RightNode.Origin.FindSecondEndNode()
		} else {
			return e.LeftNode.Origin.FindSecondEndNode()
		}
	}
}

func (e *Expression) ExistSecondEndNode() bool {
	if e.RightNode.Origin == nil && e.LeftNode.Origin == nil {
		return false
	} else {
		return true
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
		localMapKeys := tool.FindMapKeysSorted(localMap)
		sortedKeys := tool.FindMapKeysSorted(newMap)
		for _, key := range sortedKeys {
			if tool.StrContains(localMapKeys, key) == -1 {
				_, max := tool.MinMax(tool.FindMapValues(localMap))
				localMap[key] = max + 1
			} else {
				originMap = append(originMap, key)
			}
		}
		newlocalValues := tool.FindValuesByKeys(localMapKeys, localMap)
		newValues := tool.FindValuesByKeys(sortedKeys, localMap)
		function := func(f ...float64) float64 {
			if len(originMap) == 0 {
				return e.LocalFunction(f[newValues[0]], exp.Function(tool.SubSliceFloat(newlocalValues, f)...))
			}
			return e.LocalFunction(f[newValues[0]], exp.Function(tool.SubSliceFloat(newlocalValues, f)...))
		}
		e.WholeExpression = tool.ReplaceExpression(e.WholeExpression, e.RightNode.Origin.WholeExpression, false)
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
		localMapKeys := tool.FindMapKeysSorted(localMap)
		sortedKeys := tool.FindMapKeysSorted(newMap)
		for _, key := range sortedKeys {
			if tool.StrContains(localMapKeys, key) == -1 {
				_, max := tool.MinMax(tool.FindMapValues(localMap))
				localMap[key] = max + 1
			} else {
				originMap = append(originMap, key)
			}
		}
		newlocalValues := tool.FindValuesByKeys(localMapKeys, localMap)
		newValues := tool.FindValuesByKeys(sortedKeys, localMap)
		function := func(f ...float64) float64 {
			if len(originMap) == 0 {
				return e.LocalFunction(exp.Function(tool.SubSliceFloat(newlocalValues, f)...), f[newValues[0]])
			}
			return e.LocalFunction(exp.Function(tool.SubSliceFloat(newlocalValues, f)...), f[newValues[0]])
		}
		e.WholeExpression = tool.ReplaceExpression(e.WholeExpression, e.LeftNode.Origin.WholeExpression, true)
		e.Function = function
		e.LeftNode.Origin = nil
		e.Mapping = localMap
	} else {
		localMapLeft := e.LeftNode.Origin.Mapping
		expLeft := *e.LeftNode.Origin
		localMapRight := e.RightNode.Origin.Mapping
		expRight := *e.RightNode.Origin
		originMap := []string{}
		localMapLeftKeys := tool.FindMapKeysSorted(localMapLeft)
		localMapRightKeys := tool.FindMapKeysSorted(localMapRight)
		sortedKeys := tool.FindMapKeysSorted(localMapRight)
		for _, key := range sortedKeys {
			if tool.StrContains(localMapLeftKeys, key) == -1 {
				_, max := tool.MinMax(tool.FindMapValues(localMapLeft))
				localMapLeft[key] = max + 1
			} else {
				originMap = append(originMap, key)
			}
		}
		newRightValues := tool.FindValuesByKeys(localMapRightKeys, localMapLeft)
		newLeftValues := tool.FindValuesByKeys(localMapLeftKeys, localMapLeft)
		function := func(f ...float64) float64 {
			if len(originMap) == 0 {
				return e.LocalFunction(expLeft.Function(tool.SubSliceFloat(newLeftValues, f)...), expRight.Function(tool.SubSliceFloat(newRightValues, f)...))
			}
			return e.LocalFunction(expLeft.Function(tool.SubSliceFloat(newLeftValues, f)...), expRight.Function(tool.SubSliceFloat(newRightValues, f)...))
		}
		e.WholeExpression = tool.ReplaceExpressionBoth(e.WholeExpression, e.RightNode.Origin.WholeExpression, e.LeftNode.Origin.WholeExpression)
		e.Function = function
		e.RightNode.Origin = nil
		e.LeftNode.Origin = nil
		e.Mapping = localMapLeft
	}

	return e
}

func (e *Expression) GenerateFunctionMap() (func(...float64) float64, map[string]int) {
	for e.ExistSecondEndNode() {
		e.FindSecondEndNode().MergeNode()
	}
	return e.Function, e.Mapping
}
