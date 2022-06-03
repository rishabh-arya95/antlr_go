package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/rishabh-arya95/antlr_poc/ast_java"
	"github.com/rishabh-arya95/antlr_poc/parser"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

func main() {

	// byt, err := ioutil.ReadFile("./testdata/example.diff")
	// if err != nil {
	// 	panic(err)
	// }

	// diff, err := diff.ParseFileDiff(byt)
	// if err != nil {
	// 	panic(err)
	// }

	var linesChanged []int
	// for _, h := range diff.Hunks {
	// 	linesChanged = append(linesChanged, int(h.NewStartLine))
	// }

	file := "/Users/rishabharya/Desktop/Projects/java/junit-java-example/src/main/java/com/javacodegeeks/examples/junitmavenexample/Calculator.java"
	input, err := antlr.NewFileStream(file)
	if err != nil {
		panic(err)
	}

	lexer := parser.NewJavaLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewJavaParser(stream)

	jListener := ast_java.NewJavaFullListener(filepath.Base(file), linesChanged)
	antlr.NewParseTreeWalker().Walk(jListener, p.CompilationUnit())

	identModel, _ := json.MarshalIndent(jListener.GetNodeInfo(), "", "\t")
	ioutil.WriteFile("ast.json", []byte(identModel), 0644)

	changedMethods, _ := json.MarshalIndent(jListener.GetChangedMethods(), "", "\t")
	ioutil.WriteFile("changed_methods.json", []byte(changedMethods), 0644)

}

// interface
// imports
// package
// import .*
// import com.text.Formatter;
// private Formatter textFormatter;
// private com.json.Formatter jsonFormatter;
// Support for dependency injection
