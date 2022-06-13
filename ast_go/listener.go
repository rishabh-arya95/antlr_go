package ast_go

import (
	"regexp"

	"github.com/rishabh-arya95/antlr_poc/ast_go/parser"
	"github.com/rishabh-arya95/antlr_poc/core_domain"
)

var regex = regexp.MustCompile(`\b([^()]+)\((.*)\)$`)

type GoListener struct {
	parser.BaseGoParserListener
}

// ignoreFuncs := map[string]struct{}{
// 	"append": struct{}{},
// 	"copy":   struct{}{},
// 	"delete": struct{}{},
// 	"len":    struct{}{},
// 	"make":   struct{}{},
// 	"new":    struct{}{},
// 	"panic":  struct{}{},

// }
// var methodMap = make(map[string]core_domain.CodeFunction)
var currentNode *core_domain.CodeDataStruct
var currentMethod core_domain.CodeFunction

func NewGoListener(file string) *GoListener {
	currentNode = core_domain.NewDataStruct()
	currentNode.FilePath = file
	currentMethod = core_domain.NewJMethod()
	return &GoListener{}
}

func (l *GoListener) EnterPackageClause(c *parser.PackageClauseContext) {
	currentNode.Package = c.IDENTIFIER().GetText()
}

func (l *GoListener) EnterImportDecl(c *parser.ImportDeclContext) {
	for _, imp := range c.AllImportSpec() {
		var aliasName string
		if imp.GetAlias() != nil {
			aliasName = imp.GetAlias().GetText()
		}
		currentNode.Imports = append(currentNode.Imports,
			core_domain.NewGoImport(imp.GetText(), aliasName))
	}
}
func (l *GoListener) EnterFunctionDecl(c *parser.FunctionDeclContext) {
	currentMethod.Name = c.IDENTIFIER().GetText()
	if c.Signature() != nil {
		ExtractSignature(c.Signature().(*parser.SignatureContext))
	}
	currentMethod.Position = core_domain.CodePosition{
		StartLine:         c.GetStart().GetLine(),
		StartLinePosition: c.GetStart().GetColumn(), // different
		StopLine:          c.GetStop().GetLine(),
		StopLinePosition:  c.GetStop().GetColumn(),
	}
	// fmt.Println(c.Block().GetText())

}

func (l *GoListener) EnterSimpleStmt(c *parser.SimpleStmtContext) {
	if c.ShortVarDecl() != nil {
		shortVar := c.ShortVarDecl().(*parser.ShortVarDeclContext)
		exprList := shortVar.ExpressionList().(*parser.ExpressionListContext)
		for _, expr := range exprList.AllExpression() {
			codeCall := ExtractExpr(expr.(*parser.ExpressionContext))
			if codeCall.Type == "function" {
				currentMethod.FunctionCalls = append(currentMethod.FunctionCalls, codeCall)
			}
		}
	}

	if c.Assignment() != nil {
		assign := c.Assignment().(*parser.AssignmentContext)
		exprList := assign.ExpressionList(1).(*parser.ExpressionListContext)
		for _, expr := range exprList.AllExpression() {
			codeCall := ExtractExpr(expr.(*parser.ExpressionContext))
			if codeCall.Type == "function" {
				currentMethod.FunctionCalls = append(currentMethod.FunctionCalls, codeCall)
			}
		}

	}
	if c.ExpressionStmt() != nil {
		expr := c.ExpressionStmt().(*parser.ExpressionStmtContext)
		codeCall := ExtractExpr(expr.Expression().(*parser.ExpressionContext))
		if codeCall.Type == "function" {
			currentMethod.FunctionCalls = append(currentMethod.FunctionCalls, codeCall)
		}
	}
}

func (l *GoListener) EnterReturnStmt(c *parser.ReturnStmtContext) {
	exprList := c.ExpressionList().(*parser.ExpressionListContext)
	for _, expr := range exprList.AllExpression() {
		codeCall := ExtractExpr(expr.(*parser.ExpressionContext))
		if codeCall.Type == "function" {
			currentMethod.FunctionCalls = append(currentMethod.FunctionCalls, codeCall)
		}
	}
}

func (l *GoListener) ExitFunctionDecl(c *parser.FunctionDeclContext) {
	currentNode.Functions = append(currentNode.Functions, currentMethod)
	currentMethod = core_domain.NewJMethod()
}

func (l *GoListener) EnterMethodDecl(c *parser.MethodDeclContext) {
	currentMethod.Name = c.IDENTIFIER().GetText()
	currentMethod.Position = core_domain.CodePosition{
		StartLine:         c.GetStart().GetLine(),
		StartLinePosition: c.GetStart().GetColumn(), // different
		StopLine:          c.GetStop().GetLine(),
		StopLinePosition:  c.GetStop().GetColumn(),
	}
	if c.Signature() != nil {
		ExtractSignature(c.Signature().(*parser.SignatureContext))
	}
	if c.Receiver() != nil {
		receiverCtx := c.Receiver().(*parser.ReceiverContext)
		currentMethod.Receiver = ExtractParameters(receiverCtx.Parameters().(*parser.ParametersContext))[0]
	}
}

func (l *GoListener) ExitMethodDecl(c *parser.MethodDeclContext) {
	currentNode.Functions = append(currentNode.Functions, currentMethod)
	currentMethod = core_domain.NewJMethod()
}

func ExtractSignature(c *parser.SignatureContext) {
	if c.Parameters() != nil {
		currentMethod.Parameters = ExtractParameters(c.Parameters().(*parser.ParametersContext))
	}
	if c.Result() != nil {
		ExtractResult(c.Result().(*parser.ResultContext))
	}
}

func ExtractResult(c *parser.ResultContext) {
	if c.Type_() != nil {
		currentMethod.MultipleReturns = []core_domain.CodeProperty{core_domain.NewCodeParameter(c.Type_().GetText(), "")}
	}
	if c.Parameters() != nil {
		currentMethod.MultipleReturns = ExtractParameters(c.Parameters().(*parser.ParametersContext))
	}
}
func ExtractParameters(c *parser.ParametersContext) []core_domain.CodeProperty {
	var parameters = make([]core_domain.CodeProperty, 0, len(c.AllParameterDecl()))
	for _, param := range c.AllParameterDecl() {
		p := param.(*parser.ParameterDeclContext)

		var typeType, typeValue string
		if p.Type_() != nil {
			typeType = p.Type_().GetText()
		}
		if p.IdentifierList() != nil {
			typeValue = p.IdentifierList().GetText()
		}
		parameters = append(parameters,
			core_domain.NewCodeParameter(typeType, typeValue))
	}
	return parameters
}

func GetFunctionCall(str string) string {
	matches := regex.FindStringSubmatch(str)
	if len(matches) == 3 {
		return matches[1]
	}
	return ""
}

func ExtractExpr(expr *parser.ExpressionContext) core_domain.CodeCall {
	funcName := GetFunctionCall(expr.GetText())
	if funcName != "" {
		position := core_domain.CodePosition{
			StartLine:         expr.GetStart().GetLine(),
			StartLinePosition: expr.GetStart().GetColumn(),
			StopLine:          expr.GetStop().GetLine(),
			StopLinePosition:  expr.GetStop().GetColumn(),
		}
		return core_domain.CodeCall{
			NodeName:     expr.GetText(),
			FunctionName: funcName,
			Position:     position,
			Type:         "function",
		}
	}
	return core_domain.CodeCall{}
}

// func ExtractExpr(c *parser.ExpressionStmtContext) (funcName string) {

// }

// func (l *GoListener) EnterSourceFile(ctx *parser.SourceFileContext) {
// 	// TODO: top level declaration and consts
// 	for _, e := range ctx.AllDeclaration() {
// 		fmt.Println(e.GetText())
// 	}
// 	// fmt.Println(currentNode)
// }

func (l *GoListener) GetNodeInfo() core_domain.CodeDataStruct {
	return *currentNode
}



// func (this *OgVisitor) VisitParameterDecl(ctx *parser.ParameterDeclContext, delegate antlr.ParseTreeVisitor) interface{} {
// 	node := &Parameter{
// 		Node:       common.NewNode(ctx, this.File, &Parameter{}),
// 		Type:       this.VisitType_(ctx.Type_().(*parser.Type_Context), delegate).(*Type),
// 		IsVariadic: ctx.RestOp() != nil,
// 	}
// 	if ctx.IdentifierList() != nil {
// 		node.IdentifierList = this.VisitIdentifierList(ctx.IdentifierList().(*parser.IdentifierListContext), delegate).(*IdentifierList)
// 	}
// 	return node
// }

// func (this *OgVisitor) VisitIdentifierList(ctx *parser.IdentifierListContext, delegate antlr.ParseTreeVisitor) interface{} {
// 	return &IdentifierList{
// 		Node: common.NewNode(ctx, this.File, &IdentifierList{}),
// 		List: strings.Split(ctx.GetText(), ","),
// 	}
// }

// func (this *OgVisitor) VisitType_(ctx *parser.Type_Context, delegate antlr.ParseTreeVisitor) interface{} {
// 	node := &Type{Node: common.NewNode(ctx, this.File, &Type{})}
// 	if ctx.TypeName() != nil {
// 		node.TypeName = this.VisitTypeName(ctx.TypeName().(*parser.TypeNameContext), delegate).(string)
// 	}
// 	if ctx.TypeLit() != nil {
// 		node.TypeLit = this.VisitTypeLit(ctx.TypeLit().(*parser.TypeLitContext), delegate).(*TypeLit)
// 	}
// 	if ctx.Type_() != nil {
// 		node.Type = this.VisitType_(ctx.Type_().(*parser.Type_Context), delegate).(*Type)
// 	}
// 	return node
// }
// func (this *OgVisitor) VisitTypeName(ctx *parser.TypeNameContext, delegate antlr.ParseTreeVisitor) interface{} {
// 	return ctx.GetText()
// }

// func (this *OgVisitor) VisitLiteralType(ctx *parser.LiteralTypeContext, delegate antlr.ParseTreeVisitor) interface{} {
// 	node := &LiteralType{Node: common.NewNode(ctx, this.File, &LiteralType{})}
// 	if ctx.StructType() != nil {
// 		node.Struct = this.VisitStructType(ctx.StructType().(*parser.StructTypeContext), delegate).(*StructType)
// 	}
// 	if ctx.ArrayType() != nil {
// 		node.Array = this.VisitArrayType(ctx.ArrayType().(*parser.ArrayTypeContext), delegate).(*ArrayType)
// 	}
// 	if ctx.ElementType() != nil {
// 		node.Element = this.VisitElementType(ctx.ElementType().(*parser.ElementTypeContext), delegate).(*Type)
// 	}
// 	if ctx.SliceType() != nil {
// 		node.Slice = this.VisitSliceType(ctx.SliceType().(*parser.SliceTypeContext), delegate).(*SliceType)
// 	}
// 	if ctx.MapType() != nil {
// 		node.Map = this.VisitMapType(ctx.MapType().(*parser.MapTypeContext), delegate).(*MapType)
// 	}
// 	if ctx.TypeName() != nil {
// 		node.Type = this.VisitTypeName(ctx.TypeName().(*parser.TypeNameContext), delegate).(string)
// 	}
// 	return node
// }
