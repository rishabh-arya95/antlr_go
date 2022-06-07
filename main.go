package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/rishabh-arya95/antlr_poc/ast_java"
	"github.com/rishabh-arya95/antlr_poc/core_domain"

	"github.com/rishabh-arya95/antlr_poc/parser"
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

	// var linesChanged []int
	// defer profile.Start().Stop()
	// for _, h := range diff.Hunks {
	// 	linesChanged = append(linesChanged, int(h.NewStartLine))
	// }

	t1 := time.Now()
	repoDir := "/Users/rishabharya/Desktop/Projects/java/junit-java-example/"
	// repoDir := "/Users/rishabharya/Desktop/Projects/logging-log4j2/"
	// os.DirFS(global.RepoDir), strings.TrimSuffix(path, ext)+".{yml,yaml}"
	files, err := doublestar.Glob(os.DirFS(repoDir), "**/src/test/**/*.java")
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Since(t1))
	fmt.Println(len(files))
	var nodeInfos = make([]core_domain.CodeDataStruct, 0, len(files))
	// var testDiscovery = make([]custom_listener.TestCase, 0, len(files))
	for _, file := range files {
		if strings.HasPrefix(file, "log4j-api-java9") {
			continue
		}
		input, err := antlr.NewFileStream(repoDir + file)
		if err != nil {
			panic(err)
		}
		//
		lexer := parser.NewJavaLexer(input)
		stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
		p := parser.NewJavaParser(stream)
		p.BuildParseTrees = true
		tree := p.CompilationUnit()
		jListener := ast_java.NewJavaFullListener(file, []int{})
		// jListener := custom_listener.NewCustomListener(file)
		// w := NewIterativeParseTreeWalker()
		// w.Walk(jListener, tree)
		antlr.NewParseTreeWalker().Walk(jListener, tree)
		nodeInfos = append(nodeInfos, jListener.GetNodeInfo()...)
		jListener.GetCallGraph().String()
	}

	fmt.Println(time.Since(t1))
	fmt.Println(len(nodeInfos))
	fmt.Println(len(files))

	identModel, _ := json.MarshalIndent(nodeInfos, "", "\t")
	ioutil.WriteFile("testDiscovery.json", []byte(identModel), 0644)
	fmt.Println(time.Since(t1))
	// changedMethods, _ := json.MarshalIndent(jListener.GetChangedMethods(), "", "\t")
	// ioutil.WriteFile("changed_methods.json", []byte(changedMethods), 0644)

}

// interface
// imports
// package
// import .*
// import com.text.Formatter;
// private Formatter textFormatter;
// private com.json.Formatter jsonFormatter;
// Support for dependency injection
// Inner Classes
// class Outer{

// 	final int z=10;
// 	class Inner extends HasStatic {
// 	  static final int x = 3;
// 	  static int y = 4;
// 	}
// }

type IWalker struct {
	*antlr.ParseTreeWalker
}

func NewIterativeParseTreeWalker() *IWalker {
	return new(IWalker)
}
func (i *IWalker) Walk(listener antlr.ParseTreeListener, t antlr.Tree) {
	var stack []antlr.Tree
	var indexStack []int
	currentNode := t
	currentIndex := 0

	for currentNode != nil {
		// pre-order visit
		switch tt := currentNode.(type) {
		case antlr.ErrorNode:
			listener.VisitErrorNode(tt)
		case antlr.TerminalNode:
			listener.VisitTerminal(tt)
		default:
			i.EnterRule(listener, currentNode.(antlr.RuleNode))
		}
		// Move down to first child, if exists
		if currentNode.GetChildCount() > 0 {
			stack = append(stack, currentNode)
			indexStack = append(indexStack, currentIndex)
			currentIndex = 0
			currentNode = currentNode.GetChild(0)
			continue
		}

		for {
			// post-order visit
			if ruleNode, ok := currentNode.(antlr.RuleNode); ok {
				i.ExitRule(listener, ruleNode)
			}
			// No parent, so no siblings
			if len(stack) == 0 {
				currentNode = nil
				currentIndex = 0
				break
			}
			// Move to next sibling if possible
			currentIndex++
			if stack[len(stack)-1].GetChildCount() > currentIndex {
				currentNode = stack[len(stack)-1].GetChild(currentIndex)
				break
			}
			// No next, sibling, so move up
			currentNode, stack = stack[len(stack)-1], stack[:len(stack)-1]
			currentIndex, indexStack = indexStack[len(indexStack)-1], indexStack[:len(indexStack)-1]
		}
	}

}
