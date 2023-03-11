package compengine

import (
	"compiler/pkg/symtable"
	"compiler/pkg/tokenizer"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

type CompilationEngine struct {
	tokenizer               *tokenizer.Tokenizer
	classSymTable           *symtable.SymbolTable
	subroutineSymTable      *symtable.SymbolTable
	isIdentifierDeclaration bool
	outputFile              *os.File
	className               string
}

func New(tokenizer *tokenizer.Tokenizer, classSymTable *symtable.SymbolTable, subroutineSymTable *symtable.SymbolTable, outputFile *os.File) *CompilationEngine {
	return &CompilationEngine{
		tokenizer:          tokenizer,
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
	case tokenizer.Keyword:
		c.outputFile.WriteString(createXMLToken("keyword", token))
	case tokenizer.Symbol:
		c.outputFile.WriteString(createXMLToken("symbol", token))
	case tokenizer.Identifier:
		c.outputFile.WriteString(createXMLToken("identifier", c.getIdentifierSymtableOutput(token)))
	case tokenizer.IntConst:
		c.outputFile.WriteString(createXMLToken("integerConstant", token))
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
	c.outputFile.WriteString("<class>\n")
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
	c.outputFile.WriteString("</class>")
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
	c.outputFile.WriteString("<subroutineDec>\n")
	if c.getCurrentToken() == "constructor" {
		c.process("constructor")
	} else if c.getCurrentToken() == "function" {
		c.process("function")
	} else {
		c.process("method")
	}
	if c.getCurrentToken() == "void" {
		c.process("void")
	} else {
		c.processCurrentToken()
	}
	c.processCurrentToken()
	c.process("(")
	c.CompileParameterList()
	c.process(")")
	c.CompileSubroutineBody()
	c.outputFile.WriteString("</subroutineDec>\n")
}

func (c *CompilationEngine) CompileParameterList() {
	c.outputFile.WriteString("<parameterList>\n")
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
	c.outputFile.WriteString("</parameterList>\n")
}

func (c *CompilationEngine) CompileSubroutineBody() {
	c.outputFile.WriteString("<subroutineBody>\n")
	c.process("{")
	for c.getCurrentToken() == "var" {
		c.CompileVarDec()
	}
	c.CompileStatements()
	c.process("}")
	c.outputFile.WriteString("</subroutineBody>\n")
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
	c.outputFile.WriteString("<statements>\n")
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
	c.outputFile.WriteString("</statements>\n")
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
	c.outputFile.WriteString("<doStatement>\n")
	c.process("do")
	c.processCurrentToken()
	if c.getCurrentToken() == "(" {
		c.process("(")
		c.CompileExpressionList()
		c.process(")")
	} else {
		c.process(".")
		c.processCurrentToken()
		c.process("(")
		c.CompileExpressionList()
		c.process(")")
	}
	c.process(";")
	c.outputFile.WriteString("</doStatement>\n")
}

func (c *CompilationEngine) CompileReturn() {
	c.outputFile.WriteString("<returnStatement>\n")
	c.process("return")
	if c.getCurrentToken() != ";" {
		c.CompileExpression()
	}
	c.process(";")
	c.outputFile.WriteString("</returnStatement>\n")
}

func (c *CompilationEngine) CompileExpression() {
	c.outputFile.WriteString("<expression>\n")
	c.CompileTerm()
	for isOp(c.getCurrentToken()) {
		c.processCurrentToken()
		c.CompileTerm()
	}
	c.outputFile.WriteString("</expression>\n")
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
	c.outputFile.WriteString("<term>\n")
	if c.getCurrentToken() == "(" {
		c.process("(")
		c.CompileExpression()
		c.process(")")
	} else if c.getCurrentToken() == "-" || c.getCurrentToken() == "~" {
		c.processCurrentToken()
		c.CompileTerm()
	} else {
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
	c.outputFile.WriteString("</term>\n")
}

func (c *CompilationEngine) CompileExpressionList() int {
	c.outputFile.WriteString("<expressionList>\n")
	if c.getCurrentToken() == ")" {
		c.outputFile.WriteString("</expressionList>\n")
		return 0
	}
	c.CompileExpression()
	i := 1
	for c.getCurrentToken() == "," {
		c.process(",")
		c.CompileExpression()
		i++
	}
	c.outputFile.WriteString("</expressionList>\n")
	return i
}
