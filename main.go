package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/rishabh-arya95/antlr_poc/parser"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type CodePosition struct {
	StartLine         int
	StartLinePosition int
	StopLine          int
	StopLinePosition  int
}

type CodeFunction struct {
	Name       string
	ReturnType string
	Modifiers  []string
	Position   CodePosition
}

type TreeShapeListener struct {
	*parser.BaseJavaParserListener
}

var isTest, isIgnore = false, false
var imports []string
var testFunctions []CodeFunction

func NewTreeShapeListener() *TreeShapeListener {
	return new(TreeShapeListener)
}

func GetNodeIndex(node antlr.ParseTree) int {
	if node == nil || node.GetParent() == nil {
		return -1
	}
	parent := node.GetParent()

	for i := 0; i < parent.GetChildCount(); i++ {
		if parent.GetChild(i) == node {
			return i
		}
	}
	return 0
}

func (t *TreeShapeListener) EnterImportDeclaration(ctx *parser.ImportDeclarationContext) {
	importText := ctx.QualifiedName().GetText()
	imports = append(imports, importText)
}

func (t *TreeShapeListener) EnterAnnotation(ctx *parser.AnnotationContext) {
	if ctx.QualifiedName() == nil {
		return
	}
	if ctx.QualifiedName().GetText() == "Test" {
		isTest = true
	}

	if ctx.QualifiedName().GetText() == "Ignore" {
		isIgnore = true
	}
}

func (t *TreeShapeListener) EnterMethodDeclaration(ctx *parser.MethodDeclarationContext) {

	if isTest && !isIgnore {
		startLine := ctx.GetStart().GetLine()
		startLinePosition := ctx.GetStart().GetColumn()
		stopLine := ctx.GetStop().GetLine()
		stopLinePosition := ctx.GetStop().GetColumn()
		name := ""

		if ctx.Identifier() != nil {
			name = ctx.Identifier().GetText()
		}
		typeType := ctx.TypeTypeOrVoid().GetText()
		position := CodePosition{
			StartLine:         startLine,
			StartLinePosition: startLinePosition,
			StopLine:          stopLine,
			StopLinePosition:  stopLinePosition,
		}
		currentMethod := CodeFunction{
			Name:       name,
			ReturnType: typeType,
			Position:   position,
		}
		if reflect.TypeOf(ctx.GetParent().GetParent()).String() == "*parser.ClassBodyDeclarationContext" {
			bodyCtx := ctx.GetParent().GetParent().(*parser.ClassBodyDeclarationContext)
			for _, modifier := range bodyCtx.AllModifier() {
				if !strings.Contains(modifier.GetText(), "@") {
					currentMethod.Modifiers = append(currentMethod.Modifiers, modifier.GetText())
				}
			}
		}
		testFunctions = append(testFunctions, currentMethod)

	}
	isTest = false
	isIgnore = false
}

func main() {
	input, err := antlr.NewFileStream("./testdata/CalculatorTest.java")
	if err != nil {
		panic(err)
	}

	lexer := parser.NewJavaLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewJavaParser(stream)

	antlr.NewParseTreeWalker().Walk(NewTreeShapeListener(), p.CompilationUnit())

	fmt.Printf("imports: %+v\n", imports)
	fmt.Printf("Test Cases: %+v\n", testFunctions)
	fmt.Printf("Total Executable Test Cases: %+v\n", len(testFunctions))

}
