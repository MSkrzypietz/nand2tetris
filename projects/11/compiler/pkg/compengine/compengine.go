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
	tokenizer               *tokenizer.Tokenizer
	vmWriter                *vmwriter.VMWriter
	classSymTable           *symtable.SymbolTable
	subroutineSymTable      *symtable.SymbolTable
	isIdentifierDeclaration bool
	outputFile              *os.File
	className               string
}

func New(tokenizer *tokenizer.Tokenizer, vmWriter *vmwriter.VMWriter, classSymTable *symtable.SymbolTable, subroutineSymTable *symtable.SymbolTable, outputFile *os.File) *CompilationEngine {
	return &CompilationEngine{
		tokenizer:          tokenizer,
		vmWriter:           vmWriter,
		classSymTable:      classSymTable,
		subroutineSymTable: subroutineSymTable,
		outputFile:         outputFile,
		className:          strings.Split(path.Base(outputFile.Name()), ".")[0],
	}
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
	if str == c.getCurrentToken() {
		c.writeXMLTokenOutput(str)
	} else {
		log.Println("syntax error:", str)
	}
	c.tokenizer.Advance()
}

func (c *CompilationEngine) processCurrentToken() {
	c.process(c.getCurrentToken())
}

func (c *CompilationEngine) writeXMLTokenOutput(token string) {
	switch c.tokenizer.TokenType() {
	case tokenizer.StringConst:
		c.outputFile.WriteString(createXMLToken("stringConstant", token))
	}
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
	c.outputFile.WriteString("<classVarDec>\n")
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
	c.outputFile.WriteString("</classVarDec>\n")
}

func (c *CompilationEngine) CompileSubroutine() {
	c.subroutineSymTable.Reset()
	if c.getCurrentToken() == "constructor" {
		c.process("constructor")
	} else if c.getCurrentToken() == "function" {
		c.process("function")
	} else {
		c.process("method")
	}
	isVoidSubroutine := false
	if c.getCurrentToken() == "void" {
		isVoidSubroutine = true
		c.process("void")
	} else {
		c.processCurrentToken()
	}
	c.vmWriter.WriteFunction(c.className+"."+c.getCurrentToken(), 0)
	c.processCurrentToken()
	c.process("(")
	c.CompileParameterList()
	c.process(")")
	c.CompileSubroutineBody()
	if isVoidSubroutine {
		c.vmWriter.WritePush(vmwriter.Constant, 0)
		c.vmWriter.WriteReturn()
	}
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
	c.CompileStatements()
	c.process("}")
}

func (c *CompilationEngine) CompileVarDec() {
	c.outputFile.WriteString("<varDec>\n")
	c.process("var")
	entryType := c.getCurrentToken()
	c.processCurrentToken()
	c.isIdentifierDeclaration = true
	c.subroutineSymTable.Define(c.getCurrentToken(), entryType, symtable.Var)
	c.processCurrentToken()
	if c.getCurrentToken() == "," {
		c.process(",")
		c.subroutineSymTable.Define(c.getCurrentToken(), entryType, symtable.Var)
		c.processCurrentToken()
	}
	c.process(";")
	c.isIdentifierDeclaration = false
	c.outputFile.WriteString("</varDec>\n")
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
	c.outputFile.WriteString("<letStatement>\n")
	c.process("let")
	c.processCurrentToken()
	if c.getCurrentToken() == "[" {
		c.process("[")
		c.CompileExpression()
		c.process("]")
	}
	c.process("=")
	c.CompileExpression()
	c.process(";")
	c.outputFile.WriteString("</letStatement>\n")
}

func (c *CompilationEngine) CompileIf() {
	c.outputFile.WriteString("<ifStatement>\n")
	c.process("if")
	c.process("(")
	c.CompileExpression()
	c.process(")")
	c.process("{")
	c.CompileStatements()
	c.process("}")
	if c.getCurrentToken() == "else" {
		c.process("else")
		c.process("{")
		c.CompileStatements()
		c.process("}")
	}
	c.outputFile.WriteString("</ifStatement>\n")
}

func (c *CompilationEngine) CompileWhile() {
	c.outputFile.WriteString("<whileStatement>\n")
	c.process("while")
	c.process("(")
	c.CompileExpression()
	c.process(")")
	c.process("{")
	c.CompileStatements()
	c.process("}")
	c.outputFile.WriteString("</whileStatement>\n")
}

func (c *CompilationEngine) compileDo() {
	c.process("do")
	functionName := c.getCurrentToken()
	var nArgs int
	c.processCurrentToken()
	if c.getCurrentToken() == "(" {
		c.process("(")
		nArgs = c.CompileExpressionList()
		c.process(")")
	} else {
		c.process(".")
		functionName = functionName + "." + c.getCurrentToken()
		c.processCurrentToken()
		c.process("(")
		nArgs = c.CompileExpressionList()
		c.process(")")
	}
	c.process(";")
	c.vmWriter.WriteCall(functionName, nArgs) // TODO: ref push && nArgs +1
	c.vmWriter.WritePop(vmwriter.Temp, 0)
}

func (c *CompilationEngine) CompileReturn() {
	c.process("return")
	if c.getCurrentToken() != ";" {
		c.CompileExpression()
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
		case "*":
			c.vmWriter.WriteCall("Math.multiply", 2)
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
	} else if c.getCurrentToken() == "-" || c.getCurrentToken() == "~" {
		c.processCurrentToken()
		c.CompileTerm()
	} else {
		switch c.tokenizer.TokenType() {
		case tokenizer.IntConst:
			c.vmWriter.WritePush(vmwriter.Constant, c.tokenizer.IntVal())
		}
		c.processCurrentToken()
		if c.getCurrentToken() == "[" {
			c.process("[")
			c.CompileExpression()
			c.process("]")
		} else if c.getCurrentToken() == "(" {
			c.process("(")
			c.CompileExpressionList()
			c.process(")")
		} else if c.getCurrentToken() == "." {
			c.process(".")
			c.processCurrentToken()
			c.process("(")
			c.CompileExpressionList()
			c.process(")")
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
