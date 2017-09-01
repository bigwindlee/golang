/*  ==== Knowledge Points ====
1. tags in Go terminology
2. Anonymous field (embedding)
3. interface embedding (Will be assigned any value that satisfies the interface)
4. Partially initialize struct by using the syntax fieldName: fieldValue
5. idiom of for...range loop handling interface slice
6. type switch
*/

package main

import (
	"fmt"
)

type Optioner interface {
	Name() string
	IsValid() bool
}

type OptionCommon struct {
	ShortName string "short option name"
	LongName  string "long option name"
}

type IntOption struct {
	OptionCommon        // Anonymous field (embedding)
	Value, Min, Max int // Named fields (aggregation)
}

func (option IntOption) Name() string {
	return name(option.ShortName, option.LongName)
}

func (option IntOption) IsValid() bool {
	return option.Min <= option.Value && option.Value <= option.Max
}

type StringOption struct {
	OptionCommon        // Anonymous field (embedding)
	Value        string // Named fields (aggregation)
}

func (option StringOption) Name() string {
	return option.Value
}

func (option StringOption) IsValid() bool {
	return option.Value != ""
}

func name(shortName, longName string) string {
	if longName == "" {
		return shortName
	}
	return longName
}

type FloatOption struct {
	Optioner // Anonymous field (interface embedding: needs concrete type)
	// When we create FloatOption values we must assign to the embedded field
	// any value that satisfies the Optioner interface;

	Value float64 // Named field (aggregation)
}

type GenericOption struct {
	OptionCommon // Anonymous field (embedding)
}

func (option GenericOption) Name() string {
	return name(option.ShortName, option.LongName)
}

func (option GenericOption) IsValid() bool {
	return true
}

func main() {
	fileOption := StringOption{OptionCommon{"f", "file"}, "index.html"}

	// Partially initialize
	topOption := IntOption{
		OptionCommon: OptionCommon{"t", "top"},
		Max:          100,
	}
	sizeOption := FloatOption{GenericOption{OptionCommon{"s", "size"}}, 19.5}

	for _, option := range []Optioner{topOption, fileOption, sizeOption} {
		fmt.Print("name=", option.Name(), " • valid=", option.IsValid())
		fmt.Print(" • value=")
		switch option := option.(type) {
		case IntOption:
			fmt.Print(option.Value, " • min=", option.Min, " • max=", option.Max, "\n")
		case StringOption:
			fmt.Println(option.Value) // cannot fallthrough in type switch
		case FloatOption:
			fmt.Println(option.Value)
		}
	}
}

/*
[output]
name=top • valid=true • value=0 • min=0 • max=100
name=index.html • valid=true • value=index.html
name=size • valid=true • value=19.5
*/
