package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/rishabh-arya95/antlr_poc/ast_java"
	"github.com/rishabh-arya95/antlr_poc/parser"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

func main() {
	input, err := antlr.NewFileStream("./testdata/CalculatorTest.java")
	if err != nil {
		panic(err)
	}

	lexer := parser.NewJavaLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewJavaParser(stream)

	jListener := ast_java.NewJavaFullListener(filepath.Base("./testdata/CalculatorTest.java"))
	antlr.NewParseTreeWalker().Walk(jListener, p.CompilationUnit())

	fmt.Printf("%+v\n", jListener.GetNodeInfo())

	identModel, _ := json.MarshalIndent(jListener.GetNodeInfo(), "", "\t")
	ioutil.WriteFile("ast.json", []byte(identModel), 0644)

}
