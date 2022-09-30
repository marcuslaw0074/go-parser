package main

import (
	"bufio"
	"fmt"
	parser "go-parser/goparser"
	"math"
	"os"
	"regexp"
	"strings"
	"time"
	"unsafe"

	"github.com/google/uuid"
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

func ptrToString(ptr uintptr) string {
	p := unsafe.Pointer(ptr)
	return *(*string)(p)
}

func ptrToFunction(ptr uintptr) *func(...float64) float64 {
	p := unsafe.Pointer(ptr)
	fmt.Println(p, "pointer")
	return (*func(...float64) float64)(p)
}

func main() {

	// fn := func (s float64)  {
	// 	fmt.Println(s)
	// }
	// adrr, _ := strconv.ParseUint(fmt.Sprint(fn), 0, 64)
	// faked := *(*func(float64))(unsafe.Pointer(uintptr(adrr)))
	// faked(1.0)

	t := &Transs{
		Transformation: []Trans{{
			Name:       "diff",
			RegexExp:   `(?i)diff *\([_a-z0-9\+\-\*\/ \(\)]+\)`,
			RegexLocal: `(?i)diff *`,
			Uid:        uuid.NewString(),
		},
		},
	}
	fmt.Println(t)
	t.GenerateEquationList("diff (s+4-diff(t))")
	fmt.Println(t)

	time.Sleep(time.Hour)

	Fu := &parser.Function{
		Func: func(f ...float64) float64 { return math.NaN() },
	}
	Fu.GenerateFunctions("a*c/((v+e)-(g+t)*r/(r+t))+t/(a-c)", "test")

	// hi := "HI"

	// getting address as string dynamically to preserve compatibility
	// address := fmt.Sprint(Fu.Func)

	// fmt.Printf("Address of var hi: %s\n", address)

	// convert to uintptr
	// var adr uint64
	// adr, err := strconv.ParseUint(address, 0, 64)
	// if err != nil {
	// 	panic(err)
	// }
	// var ptr uintptr = uintptr(adr)
	// f := *ptrToFunction(ptr)
	// fmt.Printf("String at address: %s\n", address)
	// fmt.Printf("Value: %s\n", (f)([]float64{1,2,3,4,5,6,7}...))

	a := 10

LOOP:
	for a < 10 {
		if a == 4 {
			a = a + 1
			goto LOOP
		}
		fmt.Println(a)
		a++
	}

	s := `Diff  (a+b*c)`
	match, err := regexp.MatchString(Transformation["diff"], s)
	if err == nil {
		if match {
			reg := regexp.MustCompile(`(?i)diff *`)
			res := reg.ReplaceAllString(s, "${1}")
			fmt.Println(res)
		}
	}

	// Fu := &parser.Function{
	// 	Func: func(f ...float64) float64 { return math.NaN() },
	// }
	// Fu.GenerateFunctions("a*c/((v+e)-(g+t)*r/(r+t))+t", "test")

	for {
		mm := Fu.Mapping
		fmt.Print(fmt.Sprintf("%v number: ", len(mm)))
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			return
		}
		input = strings.TrimSuffix(input, "\n")
		ff, err := parser.SplitStringToFloat(input)
		if err == nil {
			fmt.Println(Fu.Func(ff...), Fu.Mapping)
			fmt.Println(Fu.UseFunction(map[string]float64{"a": 1, "c": 2, "g": 3, "v": 4, "t": 5, "r": 6, "e": 7}))
		} else {
			fmt.Println("Input not slice of float64")
		}
	}
}
