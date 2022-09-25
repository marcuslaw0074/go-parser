package expressionparser

import (
	"encoding/json"
	"testing"
)

type valueTest struct {
	input  []float64
	output float64
}

type formulaTest struct {
	values  []valueTest
	formula func(...float64) float64
}

type allFormulaTest struct {
	samples []formulaTest
}

var (
	test_formula_1 formulaTest = formulaTest{
		[]valueTest{
			{[]float64{2, 5, 9}, -3},
			{[]float64{2, 4, 8}, -4},
			{[]float64{1, 0, 9}, 9},
			{[]float64{10, 5, 5}, 1},
		}, func(s ...float64) float64 {
			return s[2] / (s[0] - s[1])
		},
	}
	test_formula_2 formulaTest = formulaTest{
		[]valueTest{
			{[]float64{2, 5, 9}, 19},
			{[]float64{2, 4, 8}, 16},
			{[]float64{1, 0, -9}, -9},
			{[]float64{-10, 5, 95}, 45},
		}, func(s ...float64) float64 {
			return (s[0] * s[1]) + s[2]
		},
	}
	test_formula_3 formulaTest = formulaTest{
		[]valueTest{
			{[]float64{10, 1, 1, 2}, 3},
			{[]float64{2, 4, 8, -6}, -1},
			{[]float64{18, 0, -9, 0}, -2},
			{[]float64{247, 7, 31, -7}, 10},
		}, func(s ...float64) float64 {
			return (s[0] - s[1]) / (s[2] + s[3])
		},
	}
	test_formula_4 formulaTest = formulaTest{
		[]valueTest{
			{[]float64{1, 2}, -1},
			{[]float64{2, 4}, -1},
			{[]float64{-9, 0}, 1},
			{[]float64{30, 25}, 6},
		}, func(s ...float64) float64 {
			return (s[0] - s[1]) / (s[2] + s[3])
		},
	}
)

var allSample allFormulaTest = allFormulaTest{
	samples: []formulaTest{
		test_formula_1, 
		test_formula_2, 
		test_formula_3,
	},
}

func TestValidSample(t *testing.T) {
	for _, sample := range allSample.samples {
		for _, ele := range sample.values {
			sol := sample.formula(ele.input...)
			if sol != ele.output {
				t.Errorf("Incorrect, got: %f, want: %f.", sol, ele.output)
			}
		}
	}
}

func TestMergeNodeR(t *testing.T) {
	ex := &Expression{
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
				LocalFunction: func(f1, f2 float64) float64 { return f1 - f2 },
				Function:      func(f ...float64) float64 { return f[0] - f[1] },
				Mapping:       map[string]int{"2": 0, "3": 1},
			},
		},
		LeftNode: Node{
			Uid:        "1",
			Expression: "d",
		},
		Operation:     "/",
		LocalFunction: func(f1, f2 float64) float64 { return f1 / f2 },
		Function:      func(f ...float64) float64 { return f[0] / f[1] },
		Mapping:       map[string]int{"1": 0, "4": 1},
	}
	res := *ex.MergeNode()
	for _, ele := range test_formula_1.values {
		sol := res.Function(ele.input...)
		if sol != ele.output {
			t.Errorf("Incorrect, got: %f, want: %f.", sol, ele.output)
		}
	}
}

func TestMergeNodeL(t *testing.T) {
	ex := &Expression{
		RightNode: Node{
			Uid:        "3",
			Expression: "j",
		},
		LeftNode: Node{
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
				LocalFunction: func(f1, f2 float64) float64 { return f1 * f2 },
				Function:      func(f ...float64) float64 { return f[0] * f[1] },
				Mapping:       map[string]int{"0": 0, "1": 1},
			},
		},
		Operation:     "+",
		LocalFunction: func(f1, f2 float64) float64 { return f1 + f2 },
		Function:      func(f ...float64) float64 { return f[0] + f[1] },
		Mapping:       map[string]int{"3": 0, "2": 1},
	}
	res := *ex.MergeNode()
	for _, ele := range test_formula_2.values {
		sol := res.Function(ele.input...)
		if sol != ele.output {
			t.Errorf("Incorrect, got: %f, want: %f.", sol, ele.output)
		}
	}
}

func TestMergeNodeRL(t *testing.T) {
	ex := &Expression{
		RightNode: Node{
			Uid:        "4",
			Expression: "j",
			Origin: &Expression{
				RightNode: Node{
					Uid:        "0",
					Expression: "dd",
				},
				LeftNode: Node{
					Uid:        "1",
					Expression: "Hh",
				},
				Operation:     "+",
				LocalFunction: func(f1, f2 float64) float64 { return f1 + f2 },
				Function:      func(f ...float64) float64 { return f[0] + f[1] },
				Mapping:       map[string]int{"0": 0, "1": 1},
			},
		},
		LeftNode: Node{
			Uid:        "5",
			Expression: "dh",
			Origin: &Expression{
				RightNode: Node{
					Uid:        "2",
					Expression: "df",
				},
				LeftNode: Node{
					Uid:        "3",
					Expression: "H",
				},
				Operation:     "-",
				LocalFunction: func(f1, f2 float64) float64 { return f1 - f2 },
				Function:      func(f ...float64) float64 { return f[0] - f[1] },
				Mapping:       map[string]int{"2": 0, "3": 1},
			},
		},
		Operation:     "/",
		LocalFunction: func(f1, f2 float64) float64 { return f1 / f2 },
		Function:      func(f ...float64) float64 { return f[0] / f[1] },
		Mapping:       map[string]int{"4": 0, "5": 1},
	}
	res := *ex.MergeNode()
	for _, ele := range test_formula_3.values {
		sol := res.Function(ele.input...)
		if sol != ele.output {
			t.Errorf("Incorrect, got: %f, want: %f.", sol, ele.output)
		}
	}
}

func TestFindSecondEndNode(t *testing.T) {
	ex := &Expression{
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
				LocalFunction: func(f1, f2 float64) float64 { return f1 - f2 },
				Function:      func(f ...float64) float64 { return f[0] - f[1] },
				Mapping:       map[string]int{"2": 0, "3": 1},
			},
		},
		LeftNode: Node{
			Uid:        "1",
			Expression: "d",
		},
		Operation:     "/",
		LocalFunction: func(f1, f2 float64) float64 { return f1 / f2 },
		Function:      func(f ...float64) float64 { return f[0] / f[1] },
		Mapping:       map[string]int{"1": 0, "4": 1},
	}
	res := ex.FindSecondEndNode()
	sf, err := json.Marshal(res)
	if err != nil {
		t.Errorf("Error: %v, fail to transform to json string", err)
	} else {
		res := string(sf)
		expect := `{"rightNode":{"uid":"4","numerical":0,"expression":"j","origin":{"rightNode":{"uid":"2","numerical":0,"expression":"dd","origin":null},"operation":"-","leftNode":{"uid":"3","numerical":0,"expression":"Hh","origin":null}}},"operation":"/","leftNode":{"uid":"1","numerical":0,"expression":"d","origin":null}}`
		if res != expect {
			t.Errorf("Incorrect, got: %s, want: %s.", res, expect)
		}
	}
}

func TestMergeNodeRRepeat(t *testing.T) {
	ex := &Expression{
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
				LocalFunction: func(f1, f2 float64) float64 { return f1 - f2 },
				Function:      func(f ...float64) float64 { return f[0] - f[1] },
				Mapping:       map[string]int{"2": 0, "3": 1},
			},
		},
		LeftNode: Node{
			Uid:        "2",
			Expression: "dd",
		},
		Operation:     "/",
		LocalFunction: func(f1, f2 float64) float64 { return f1 / f2 },
		Function:      func(f ...float64) float64 { return f[0] / f[1] },
		Mapping:       map[string]int{"2": 0, "4": 1},
	}
	res := *ex.MergeNode()
	for _, ele := range test_formula_4.values {
		sol := res.Function(ele.input...)
		if sol != ele.output {
			t.Errorf("Incorrect, got: %f, want: %f.", sol, ele.output)
		}
	}
}