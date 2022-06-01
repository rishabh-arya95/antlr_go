package core_domain

import (
	"strings"
)

var ASSERTION_LIST = []string{
	"assert",
	"should",
	"check",    // ArchUnit,
	"maynotbe", // ArchUnit,
	"is",       // RestAssured,
	"spec",     // RestAssured,
	"verify",   // Mockito,
}

type CodeCall struct {
	Package      string
	Type         string
	NodeName     string
	FunctionName string
	Parameters   []CodeProperty
	Position     CodePosition
}

func NewCodeMethodCall() CodeCall {
	return CodeCall{}
}

func (c *CodeCall) BuildFullMethodName() string {
	isConstructor := c.FunctionName == ""
	if isConstructor {
		return c.Package + "." + c.NodeName
	}
	return c.Package + "." + c.NodeName + "." + c.FunctionName
}

func (c *CodeCall) BuildClassFullName() string {
	return c.Package + "." + c.NodeName
}

func (c *CodeCall) IsSystemOutput() bool {
	return c.NodeName == "System.out" && (c.FunctionName == "println" || c.FunctionName == "printf" || c.FunctionName == "print")
}

func (c *CodeCall) IsThreadSleep() bool {
	return c.FunctionName == "sleep" && c.NodeName == "Thread"
}

func (c *CodeCall) HasAssertion() bool {
	methodName := strings.ToLower(c.FunctionName)
	for _, assertion := range ASSERTION_LIST {
		if strings.HasPrefix(methodName, assertion) {
			return true
		}
	}

	return false
}
