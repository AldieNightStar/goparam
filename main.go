package main

import (
	"fmt"
	"strings"

	"github.com/AldieNightStar/goscriptable"
)

func main() {
	args := goscriptable.ParseArgs(goscriptable.GetOsArgs())
	if len(args) < 1 {
		printUsage()
		return
	}
	name, namePresent := args["name"]
	params, paramsPresent := args["params"]
	packg, packgPresent := args["package"]
	_, getters := args["get"]
	_, setters := args["set"]
	_, ctor := args["ctor"]
	_, builder := args["builder"]
	if !namePresent || !paramsPresent || !packgPresent {
		if !namePresent {
			println("ERR: Name is absent!")
		}
		if !paramsPresent {
			println("ERR: Parameters is absent!")
		}
		if !packgPresent {
			println("ERR: Package is absent!")
		}
		printUsage()
		return
	}
	sb := &strings.Builder{}
	GenerateStruct(sb, name, packg, params, getters, setters, ctor, builder)
	src := sb.String()
	if len(src) > 0 {
		f := goscriptable.CreateFile(name + "_goparam.go")
		f.WriteString(src)
		f.Close()
	}
	fmt.Println("Done!")
}

func GenerateStruct(sb *strings.Builder, name string, packg string, params string, getters, setters, ctor, builder bool) {
	fmt.Fprintf(sb, "package %s\n\n", packg)
	fmt.Fprintf(sb, "type %s struct {\n", name)

	paramsArr := ParseParams(params)
	for _, param := range paramsArr {
		fmt.Fprintf(sb, "\t%s %s\n", param.Name, param.Type)
	}

	fmt.Fprintf(sb, "}\n\n")

	if ctor {
		fmt.Fprintf(sb, "func New%s(", name)
		arr := make([]string, 0, 8)
		for _, param := range paramsArr {
			arr = append(arr, fmt.Sprintf("%s %s", param.Name, param.Type))
		}
		str := strings.Join(arr, ", ")
		fmt.Fprintf(sb, str)
		fmt.Fprintf(sb, ") *%s {\n\tself := &%s{}\n", name, name)
		for _, param := range paramsArr {
			fmt.Fprintf(sb, "\tself.%s = %s\n", param.Name, param.Name)
		}
		fmt.Fprintf(sb, "\treturn self\n")
		fmt.Fprintf(sb, "}\n\n")
	}
	if getters {
		for _, param := range paramsArr {
			fmt.Fprintf(sb, "func (self *%s) Get%s() %s {\n", name, param.Name, param.Type)
			fmt.Fprintf(sb, "\treturn self.%s\n}\n\n", param.Name)
		}
	}
	if setters {
		for _, param := range paramsArr {
			fmt.Fprintf(sb, "func (self *%s) Set%s(val %s) {\n", name, param.Name, param.Type)
			fmt.Fprintf(sb, "\tself.%s = val\n}\n\n", param.Name)
		}
	}
	if builder {
		GenerateBuilder(sb, name, paramsArr)
	}
}

func GenerateBuilder(sb *strings.Builder, typeName string, params []*Param) {
	fmt.Fprintf(sb, "type %sBuilder struct {\n\tBVal *%s\n}\n\nfunc New%sBuilder() *%sBuilder {\n\treturn &%sBuilder{&%s{}}\n}\n\n", typeName, typeName, typeName, typeName, typeName, typeName)
	for _, param := range params {
		fmt.Fprintf(sb, "func (self *%sBuilder) %s(%s %s) *%sBuilder {\n\tself.BVal.%s = %s\n\treturn self\n}\n\n", typeName, param.Name, param.Name, param.Type, typeName, param.Name, param.Name)
	}
	fmt.Fprintf(sb, "func (self *%sBuilder) Build() *%s {\n\treturn self.BVal\n}\n\n", typeName, typeName)
}

func ParseParams(line string) []*Param {
	arr := make([]*Param, 0, 8)
	params := strings.Split(line, ", ")
	for _, param := range params {
		if !strings.Contains(param, " ") {
			continue
		}
		pArr := strings.Split(param, " ")
		for i, p := range pArr {
			pArr[i] = strings.Trim(p, " \t")
		}
		arr = append(arr, &Param{Name: pArr[0], Type: pArr[1]})
	}
	return arr
}

type Param struct {
	Name string
	Type string
}

func printUsage() {
	fmt.Println("goparam -package main -name User -params \"Age int, Id int, Type UserType\" -get -set -ctor -builder\n")
	fmt.Println("Options:")
	fmt.Println("  -name    - Type/Struct name")
	fmt.Println("  -params  - Parameter list, separated by comma: Name Type, Name2 Type2, ...")
	fmt.Println("  -package - Package name")
	fmt.Println("  -get     - Generate getters")
	fmt.Println("  -set     - Generate setters")
	fmt.Println("  -ctor    - Generate constructor New...()")
	fmt.Println("  -builder - Adds builder feature to your struct")
}
