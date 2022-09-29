package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/google/uuid"
	parser "go-parser/goparser"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Function struct {
	Name    string
	Func    func(...float64) float64
	Mapping map[string]int
	Uid     string
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

func (F *Function) GenerateFunctions(s, name string) error {
	s1 := parser.ExpressionGenerator(s)
	f, m, err := parser.Generator(s1)
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

func main() {
	Fu := &Function{
		Func: func(f ...float64) float64 { return math.NaN() },
	}
	Fu.GenerateFunctions("a*c/((v+e)-(g+t)*r/(r+t))+t", "2test")
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
		ff, err := SplitStringToFloat(input)
		if err == nil {
			fmt.Println(Fu.Func(ff...))
		} else {
			fmt.Println("Input not slice of float64")
		}
	}
}
