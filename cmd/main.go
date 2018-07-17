package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/midorigreen/gpipe"
)

type InputType struct {
	ID   int
	Name string
}

func (i InputType) String() string {
	return fmt.Sprintf("%v, %v", i.ID, i.Name)
}

type OutputType struct {
	ID        int
	Name      string
	UpperName string
}

func (o OutputType) String() string {
	return fmt.Sprintf("%v, %v, %v", o.ID, o.Name, o.UpperName)
}

type Upper struct{}

func (u Upper) Convert(out fmt.Stringer, params []interface{}) (fmt.Stringer, error) {
	o, ok := out.(InputType)
	if !ok {
		return nil, fmt.Errorf("cast error")
	}
	return OutputType{
		ID:        o.ID,
		Name:      o.Name,
		UpperName: strings.ToUpper(o.Name),
	}, nil
}

func main() {
	in := gpipe.Input(os.Stdin, 10)
	convs := []gpipe.Converter{
		Upper{},
	}
	gpipe.Output(os.Stdout, in, parse, convs)
}

func parse(b []byte) (fmt.Stringer, error) {
	r := csv.NewReader(bytes.NewBuffer(b))
	columns, err := r.Read()
	if err != nil {
		return InputType{}, err
	}
	id, err := strconv.Atoi(columns[0])
	if err != nil {
		return InputType{}, err
	}
	return InputType{
		ID:   id,
		Name: columns[1],
	}, nil
}
