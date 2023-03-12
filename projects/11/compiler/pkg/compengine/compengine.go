package compengine

import (
	"compiler/pkg/symtable"
	"compiler/pkg/tokenizer"
	"compiler/pkg/vmwriter"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

type CompilationEngine struct {
	tokenizer                *tokenizer.Tokenizer
	vmWriter                 *vmwriter.VMWriter
	classSymTable            *symtable.SymbolTable
	subroutineSymTable       *symtable.SymbolTable
	isIdentifierDeclaration  bool
	outputFile               *os.File
	className                string
	functionName             string
	ifLabelCounter           int
	whileLabelCounter        int
	isConstructorCompilation bool
	isMethodCompilation      bool
}

func New(tokenizer *tokenizer.Tokenizer, vmWriter *vmwriter.VMWriter, classSymTable *symtable.SymbolTable, subroutineSymTable *symtable.SymbolTable, outputFile *os.File) *CompilationEngine {
	return &CompilationEngine{
		tokenizer:                tokenizer,
		vmWriter:                 vmWriter,
		classSymTable:            classSymTable,
		subroutineSymTable:       subroutineSymTable,
		outputFile:               outputFile,
		className:                strings.Split(path.Base(outputFile.Name()), ".")[0],
		functionName:             "",
		ifLabelCounter:           -1,
		whileLabelCounter:        -1,
		isConstructorCompilation: false,
		isMethodCompilation:      false,
	}
}

func (c *CompilationEngine) nextUniqueIfLabelTuple() (string, string, string) {
	c.ifLabelCounter++
	counter := strconv.Itoa(c.ifLabelCounter)
	return "IF_TRUE" + counter, "IF_FALSE" + counter, "IF_END" + counter
}

func (c *CompilationEngine) nextUniqueWhileLabelTuple() (string, string) {
	c.whileLabelCounter++
	counter := strconv.Itoa(c.whileLabelCounter)
	return "WHILE_EXP" + counter, "WHILE_END" + counter
}

func (c *CompilationEngine) getSymbolTable(identifier string) *symtable.SymbolTable {
	if c.subroutineSymTable.KindOf(identifier) != symtable.None {
		return c.subroutineSymTable
	}
	return c.classSymTable
}

func (c *CompilationEngine) getCurrentToken() string {
	switch c.tokenizer.TokenType() {
	case tokenizer.Keyword:
		return c.tokenizer.KeyWord()
	case tokenizer.Symbol:
		return c.tokenizer.Symbol()
	case tokenizer.Identifier:
		return c.tokenizer.Identifier()
	case tokenizer.IntConst:
		return strconv.Itoa(c.tokenizer.IntVal())
	case tokenizer.StringConst:
		return c.tokenizer.StringVal()
	}
	return ""
}

func (c *CompilationEngine) process(str string) {
	if str != c.getCurrentToken() {
		log.Println("syntax error:", str)
	}
	c.tokenizer.Advance()
}

func (c *CompilationEngine) processCurrentToken() {
	c.process(c.getCurrentToken())
}

func (c *CompilationEngine) getIdentifierSymtableOutput(identifier string) string {
	var category string
	entryKind := c.getSymbolTable(identifier).KindOf(identifier)
	switch entryKind {
	case symtable.Static:
		category = "static"
	case symtable.Field:
		category = "field"
	case symtable.Arg:
		category = "arg"
	case symtable.Var:
		category = "var"
	case symtable.None:
		if identifier == c.className {
			category = "class"
		} else {
			category = "subroutine"
		}
	}

	output := []string{identifier, category}
	if entryKind != symtable.None {
		output = append(output, strconv.Itoa(c.getSymbolTable(identifier).IndexOf(identifier)))

		if c.isIdentifierDeclaration {
			output = append(output, "declared")
		} else {
			output = append(output, "used")
		}
	}
	return strings.Join(output, " - ")
}

func createXMLToken(tokenName, value string) string {
	sanitizedValue := value
	if value == "<" {
		sanitizedValue = "&lt;"
	} else if value == ">" {
		sanitizedValue = "&gt;"
	} else if value == "&" {
		sanitizedValue = "&amp;"
	} else if value == "\"" {
		sanitizedValue = "&quot;"
	}

	var sb strings.Builder
	sb.WriteString("<" + tokenName + ">")
	sb.WriteString(" " + sanitizedValue + " ")
	sb.WriteString("</" + tokenName + ">\n")
	return sb.String()
}

func (c *CompilationEngine) CompileClass() {
	c.process("class")
	c.processCurrentToken()
	c.process("{")
	for c.getCurrentToken() == "static" || c.getCurrentToken() == "field" {
		c.CompileClassVarDec()
	}
	for c.getCurrentToken() == "constructor" || c.getCurrentToken() == "function" || c.getCurrentToken() == "method" {
		c.CompileSubroutine()
	}
	c.process("}")
}

func (c *CompilationEngine) CompileClassVarDec() {
	var kind symtable.SymbolTableEntryKind
	if c.getCurrentToken() == "static" {
		c.process("static")
		kind = symtable.Static
	} else {
		c.process("field")
		kind = symtable.Field
	}
	entryType := c.getCurrentToken()
	c.processCurrentToken()
	c.isIdentifierDeclaration = true
	c.classSymTable.Define(c.getCurrentToken(), entryType, kind)
	c.processCurrentToken()
	for c.getCurrentToken() == "," {
		c.process(",")
		c.classSymTable.Define(c.getCurrentToken(), entryType, kind)
		c.processCurrentToken()
	}
	c.process(";")
	c.isIdentifierDeclaration = false
}

func (c *CompilationEngine) CompileSubroutine() {
	c.subroutineSymTable.Reset()
	c.ifLabelCounter = -1
	c.whileLabelCounter = -1
	if c.getCurrentToken() == "constructor" {
		c.isConstructorCompilation = true
		c.process("constructor")
	} else if c.getCurrentToken() == "function" {
		c.process("function")
	} else {
		c.isMethodCompilation = true
		c.process("method")
	}
	isVoidSubroutine := false
	if c.getCurrentToken() == "void" {
		isVoidSubroutine = true
		c.process("void")
	} else {
		c.processCurrentToken()
	}
	c.functionName = c.className + "." + c.getCurrentToken()
	c.processCurrentToken()
	c.process("(")
	c.CompileParameterList()
	c.process(")")
	c.CompileSubroutineBody()
	if isVoidSubroutine {
		c.vmWriter.WritePush(vmwriter.Constant, 0)
		c.vmWriter.WriteReturn()
	}
	c.isConstructorCompilation = false
}

func (c *CompilationEngine) CompileParameterList() {
	token := c.getCurrentToken()
	isBuiltInType := token == "int" || token == "char" || token == "boolean"
	c.isIdentifierDeclaration = true
	if isBuiltInType || c.tokenizer.TokenType() == tokenizer.Identifier {
		entryType := c.getCurrentToken()
		c.processCurrentToken()
		c.subroutineSymTable.Define(c.getCurrentToken(), entryType, symtable.Arg)
		c.processCurrentToken()
	}
	for c.getCurrentToken() == "," {
		c.process(",")
		entryType := c.getCurrentToken()
		c.processCurrentToken()
		c.subroutineSymTable.Define(c.getCurrentToken(), entryType, symtable.Arg)
		c.processCurrentToken()
	}
	c.isIdentifierDeclaration = false
}

func (c *CompilationEngine) CompileSubroutineBody() {
	c.process("{")
	for c.getCurrentToken() == "var" {
		c.CompileVarDec()
	}
	c.vmWriter.WriteFunction(c.functionName, c.subroutineSymTable.VarCount(symtable.Var))
	if c.isMethodCompilation {
		c.vmWriter.WritePush(vmwriter.Argument, 0)
		c.vmWriter.WritePop(vmwriter.Pointer, 0)
	} else if c.isConstructorCompilation {
		c.vmWriter.WritePush(vmwriter.Constant, c.classSymTable.VarCount(symtable.Field))
		c.vmWriter.WriteCall("Memory.alloc", 1)
		c.vmWriter.WritePop(vmwriter.Pointer, 0)
	}
	c.CompileStatements()
	c.process("}")
}

func (c *CompilationEngine) CompileVarDec() {
	c.process("var")
	entryType := c.getCurrentToken()
	c.processCurrentToken()
	c.isIdentifierDeclaration = true
	c.subroutineSymTable.Define(c.getCurrentToken(), entryType, symtable.Var)
	c.processCurrentToken()
	for c.getCurrentToken() == "," {
		c.process(",")
		c.subroutineSymTable.Define(c.getCurrentToken(), entryType, symtable.Var)
		c.processCurrentToken()
	}
	c.process(";")
	c.isIdentifierDeclaration = false
}

func (c *CompilationEngine) CompileStatements() {
	stop := false
	for !stop {
		switch c.getCurrentToken() {
		case "let":
			c.CompileLet()
		case "if":
			c.CompileIf()
		case "while":
			c.CompileWhile()
		case "do":
			c.compileDo()
		case "return":
			c.CompileReturn()
		default:
			stop = true
		}
	}
}

func (c *CompilationEngine) CompileLet() {
	c.process("let")
	varName := c.getCurrentToken()
	c.processCurrentToken()
	isArrayAssignment := false
	if c.getCurrentToken() == "[" {
		isArrayAssignment = true
		c.process("[")
		c.CompileExpression()
		c.writePushForIdentifier(varName)
		c.vmWriter.WriteArithmetic(vmwriter.Add)
		c.process("]")
	}
	c.process("=")
	c.CompileExpression()
	if isArrayAssignment {
		c.vmWriter.WritePop(vmwriter.Temp, 0)
		c.vmWriter.WritePop(vmwriter.Pointer, 1)
		c.vmWriter.WritePush(vmwriter.Temp, 0)
		c.vmWriter.WritePop(vmwriter.That, 0)
	} else {
		c.writePopForIdentifier(varName)
	}
	c.process(";")
}

func (c *CompilationEngine) writePopForIdentifier(identifier string) {
	if segment, found := c.getMemorySegment(identifier); found {
		c.vmWriter.WritePop(segment, c.getSymbolTable(identifier).IndexOf(identifier))
	}
}

func (c *CompilationEngine) writePushForIdentifier(identifier string) {
	if segment, found := c.getMemorySegment(identifier); found {
		c.vmWriter.WritePush(segment, c.getSymbolTable(identifier).IndexOf(identifier))
	}
}

func (c *CompilationEngine) getMemorySegment(identifier string) (segment vmwriter.MemorySegment, found bool) {
	symTable := c.getSymbolTable(identifier)
	switch symTable.KindOf(identifier) {
	case symtable.Field:
		return vmwriter.This, true
	case symtable.Arg:
		return vmwriter.Argument, true
	case symtable.Var:
		return vmwriter.Local, true
	}
	return -1, false
}

func (c *CompilationEngine) CompileIf() {
	lt, lf, le := c.nextUniqueIfLabelTuple()
	c.process("if")
	c.process("(")
	c.CompileExpression()
	c.vmWriter.WriteIf(lt)
	c.vmWriter.WriteGoto(lf)
	c.vmWriter.WriteLabel(lt)
	c.process(")")
	c.process("{")
	c.CompileStatements()
	c.process("}")
	if c.getCurrentToken() == "else" {
		c.vmWriter.WriteGoto(le)
		c.vmWriter.WriteLabel(lf)
		c.process("else")
		c.process("{")
		c.CompileStatements()
		c.process("}")
		c.vmWriter.WriteLabel(le)
	} else {
		c.vmWriter.WriteLabel(lf)
	}
}

func (c *CompilationEngine) CompileWhile() {
	c.process("while")
	c.process("(")
	l1, l2 := c.nextUniqueWhileLabelTuple()
	c.vmWriter.WriteLabel(l1)
	c.CompileExpression()
	c.vmWriter.WriteArithmetic(vmwriter.Not)
	c.vmWriter.WriteIf(l2)
	c.process(")")
	c.process("{")
	c.CompileStatements()
	c.process("}")
	c.vmWriter.WriteGoto(l1)
	c.vmWriter.WriteLabel(l2)
}

func (c *CompilationEngine) compileDo() {
	c.process("do")
	functionName := c.getCurrentToken()
	var nArgs int
	c.processCurrentToken()
	if c.getCurrentToken() == "(" {
		c.process("(")
		functionName = c.className + "." + functionName
		c.vmWriter.WritePush(vmwriter.Pointer, 0)
		nArgs = c.CompileExpressionList() + 1
		c.process(")")
	} else {
		c.process(".")
		symTable := c.getSymbolTable(functionName)
		switch symTable.KindOf(functionName) {
		case symtable.Var:
			c.vmWriter.WritePush(vmwriter.Local, symTable.IndexOf(functionName))
			nArgs = 1
			functionName = symTable.TypeOf(functionName) + "." + c.getCurrentToken()
		case symtable.Field:
			c.vmWriter.WritePush(vmwriter.This, symTable.IndexOf(functionName))
			nArgs = 1
			functionName = symTable.TypeOf(functionName) + "." + c.getCurrentToken()
		default:
			functionName = functionName + "." + c.getCurrentToken()
		}
		c.processCurrentToken()
		c.process("(")
		nArgs += c.CompileExpressionList()
		c.process(")")
	}
	c.process(";")
	c.vmWriter.WriteCall(functionName, nArgs)
	c.vmWriter.WritePop(vmwriter.Temp, 0)
}

func (c *CompilationEngine) CompileReturn() {
	c.process("return")
	if c.getCurrentToken() != ";" {
		c.CompileExpression()
		c.vmWriter.WriteReturn()
	}
	c.process(";")
}

func (c *CompilationEngine) CompileExpression() {
	c.CompileTerm()
	for isOp(c.getCurrentToken()) {
		operator := c.getCurrentToken()
		c.processCurrentToken()
		c.CompileTerm()
		switch operator {
		case "+":
			c.vmWriter.WriteArithmetic(vmwriter.Add)
		case "-":
			c.vmWriter.WriteArithmetic(vmwriter.Sub)
		case "*":
			c.vmWriter.WriteCall("Math.multiply", 2)
		case "/":
			c.vmWriter.WriteCall("Math.divide", 2)
		case "<":
			c.vmWriter.WriteArithmetic(vmwriter.Lt)
		case ">":
			c.vmWriter.WriteArithmetic(vmwriter.Gt)
		case "=":
			c.vmWriter.WriteArithmetic(vmwriter.Eq)
		case "&":
			c.vmWriter.WriteArithmetic(vmwriter.And)
		}
	}
}

func isOp(token string) bool {
	operators := []string{"+", "-", "*", "/", "&", "|", "<", ">", "="}
	for _, op := range operators {
		if token == op {
			return true
		}
	}
	return false
}

func (c *CompilationEngine) CompileTerm() {
	if c.getCurrentToken() == "(" {
		c.process("(")
		c.CompileExpression()
		c.process(")")
	} else if c.getCurrentToken() == "-" {
		c.processCurrentToken()
		c.CompileTerm()
		c.vmWriter.WriteArithmetic(vmwriter.Neg)
	} else if c.getCurrentToken() == "~" {
		c.processCurrentToken()
		c.CompileTerm()
		c.vmWriter.WriteArithmetic(vmwriter.Not)
	} else {
		tokenType := c.tokenizer.TokenType()
		switch tokenType {
		case tokenizer.IntConst:
			c.vmWriter.WritePush(vmwriter.Constant, c.tokenizer.IntVal())
		case tokenizer.StringConst:
			c.vmWriter.WritePush(vmwriter.Constant, len(c.tokenizer.StringVal()))
			c.vmWriter.WriteCall("String.new", 1)
			for _, char := range c.tokenizer.StringVal() {
				c.vmWriter.WritePush(vmwriter.Constant, int(char))
				c.vmWriter.WriteCall("String.appendChar", 2)
			}
		}
		if c.getCurrentToken() == "true" {
			c.vmWriter.WritePush(vmwriter.Constant, 0)
			c.vmWriter.WriteArithmetic(vmwriter.Not)
		} else if c.getCurrentToken() == "null" || c.getCurrentToken() == "false" {
			c.vmWriter.WritePush(vmwriter.Constant, 0)
		} else if c.getCurrentToken() == "this" {
			c.vmWriter.WritePush(vmwriter.Pointer, 0)
		}

		identifier := c.getCurrentToken()
		c.processCurrentToken()
		if c.getCurrentToken() == "[" {
			c.process("[")
			c.CompileExpression()
			c.writePushForIdentifier(identifier)
			c.vmWriter.WriteArithmetic(vmwriter.Add)
			c.vmWriter.WritePop(vmwriter.Pointer, 1)
			c.vmWriter.WritePush(vmwriter.That, 0)
			c.process("]")
		} else if c.getCurrentToken() == "(" {
			c.process("(")
			c.CompileExpressionList()
			c.process(")")
		} else if c.getCurrentToken() == "." {
			c.process(".")
			classFunctionName := identifier + "." + c.getCurrentToken()
			c.processCurrentToken()
			c.process("(")
			nArgs := c.CompileExpressionList()
			c.vmWriter.WriteCall(classFunctionName, nArgs)
			c.process(")")
		} else if tokenType == tokenizer.Identifier {
			c.writePushForIdentifier(identifier)
		}
	}
}

func (c *CompilationEngine) CompileExpressionList() int {
	if c.getCurrentToken() == ")" {
		return 0
	}
	c.CompileExpression()
	i := 1
	for c.getCurrentToken() == "," {
		c.process(",")
		c.CompileExpression()
		i++
	}
	return i
}
