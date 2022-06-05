package custom_listener

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/rishabh-arya95/antlr_poc/core_domain"
	"github.com/rishabh-arya95/antlr_poc/parser"
)

type CustomListener struct {
	*parser.BaseJavaParserListener
}

var isTest = false
var currentClassName string
var testList []TestCase
var classStringQueue []string
var fileName string

type TestCase struct {
	Name     string
	Position core_domain.CodePosition
	ID       string
	Class    string
	FileName string
}

func NewCustomListener(testFileName string) *CustomListener {
	fileName = testFileName
	testList = nil
	classStringQueue = nil
	return &CustomListener{}
}

func (c *CustomListener) EnterAnnotation(ctx *parser.AnnotationContext) {
	if ctx.QualifiedName() == nil {
		return
	}
	if ctx.QualifiedName().GetText() == "Test" || ctx.QualifiedName().GetText() == "ParameterizedTest" {
		isTest = true
	}
}

func (c *CustomListener) EnterMethodDeclaration(ctx *parser.MethodDeclarationContext) {
	if isTest {
		startLine := ctx.GetStart().GetLine()
		startLinePosition := ctx.GetStart().GetColumn()
		stopLine := ctx.GetStop().GetLine()
		stopLinePosition := ctx.GetStop().GetColumn()
		name := ""
		if ctx.Identifier() != nil {
			name = ctx.Identifier().GetText()
		}
		position := core_domain.CodePosition{
			StartLine:         startLine,
			StartLinePosition: startLinePosition,
			StopLine:          stopLine,
			StopLinePosition:  stopLinePosition,
		}
		h := md5.New()
		io.WriteString(h, fmt.Sprintf("%s%s%s", fileName, currentClassName, name))

		testCase := TestCase{
			Name:     name,
			Position: position,
			ID:       hex.EncodeToString(h.Sum(nil)),
			FileName: fileName,
			Class:    currentClassName,
		}
		testList = append(testList, testCase)
	}
	isTest = false
}

func (c *CustomListener) EnterClassDeclaration(ctx *parser.ClassDeclarationContext) {
	// fmt.Println("Entering class declaration")
	if ctx.Identifier() != nil {
		currentClassName = ctx.Identifier().GetText()
	}
}

func (c *CustomListener) GetTestCases() []TestCase {
	return testList
}

func (c *CustomListener) EnterInnerCreator(ctx *parser.InnerCreatorContext) {
	if ctx.Identifier() != nil {
		currentClassName = ctx.Identifier().GetText()
		classStringQueue = append(classStringQueue, currentClassName)
	}
}

func (c *CustomListener) ExitInnerCreator(ctx *parser.InnerCreatorContext) {
	if classStringQueue == nil || len(classStringQueue) <= 1 {
		return
	}

	classStringQueue = classStringQueue[0 : len(classStringQueue)-1]
	currentClassName = classStringQueue[len(classStringQueue)-1]
}