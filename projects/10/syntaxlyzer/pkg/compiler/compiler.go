package compiler

import (
	"log"
	"os"
	"strconv"
	"strings"
	"syntaxlyzer/pkg/tokenizer"
)

type Compiler struct {
	tokenizer  *tokenizer.Tokenizer
	outputFile *os.File
}

func New(tokenizer *tokenizer.Tokenizer, outputFile *os.File) *Compiler {
	return &Compiler{
		tokenizer:  tokenizer,
		outputFile: outputFile,
	}
}

func (c *Compiler) getCurrentToken() string {
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

func (c *Compiler) process(str string) {
	if str == c.getCurrentToken() {
		c.writeXMLTokenOutput(str)
	} else {
		log.Println("syntax error:", str)
	}
	c.tokenizer.Advance()
}

func (c *Compiler) processCurrentToken() {
	c.process(c.getCurrentToken())
}

func (c *Compiler) writeXMLTokenOutput(token string) {
	switch c.tokenizer.TokenType() {
	case tokenizer.Keyword:
		c.outputFile.WriteString(createXMLToken("keyword", token))
	case tokenizer.Symbol:
		c.outputFile.WriteString(createXMLToken("symbol", token))
	case tokenizer.Identifier:
		c.outputFile.WriteString(createXMLToken("identifier", token))
	case tokenizer.IntConst:
		c.outputFile.WriteString(createXMLToken("integerConstant", token))
	case tokenizer.StringConst:
		c.outputFile.WriteString(createXMLToken("stringConstant", token))
	}
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

func (c *Compiler) CompileClass() {
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

func (c *Compiler) CompileClassVarDec() {
	c.outputFile.WriteString("<classVarDec>\n")
	if c.getCurrentToken() == "static" {
		c.process("static")
	} else {
		c.process("field")
	}
	c.processCurrentToken()
	c.processCurrentToken()
	if c.getCurrentToken() == "," {
		c.process(",")
		c.processCurrentToken()
	}
	c.process(";")
	c.outputFile.WriteString("</classVarDec>\n")
}

func (c *Compiler) CompileSubroutine() {
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

func (c *Compiler) CompileParameterList() {
	c.outputFile.WriteString("<parameterList>\n")
	token := c.getCurrentToken()
	isBuiltInType := token == "int" || token == "char" || token == "boolean"
	if isBuiltInType || c.tokenizer.TokenType() == tokenizer.Identifier {
		c.processCurrentToken()
		c.processCurrentToken()
	}
	for c.getCurrentToken() == "," {
		c.process(",")
		c.processCurrentToken()
		c.processCurrentToken()
	}
	c.outputFile.WriteString("</parameterList>\n")
}

func (c *Compiler) CompileSubroutineBody() {
	c.outputFile.WriteString("<subroutineBody>\n")
	c.process("{")
	for c.getCurrentToken() == "var" {
		c.CompileVarDec()
	}
	c.CompileStatements()
	c.process("}")
	c.outputFile.WriteString("</subroutineBody>\n")
}

func (c *Compiler) CompileVarDec() {
	c.outputFile.WriteString("<varDec>\n")
	c.process("var")
	c.processCurrentToken()
	c.processCurrentToken()
	if c.getCurrentToken() == "," {
		c.process(",")
		c.processCurrentToken()
	}
	c.process(";")
	c.outputFile.WriteString("</varDec>\n")
}

func (c *Compiler) CompileStatements() {
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

func (c *Compiler) CompileLet() {
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

func (c *Compiler) CompileIf() {
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

func (c *Compiler) CompileWhile() {
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

func (c *Compiler) compileDo() {
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

func (c *Compiler) CompileReturn() {
	c.outputFile.WriteString("<returnStatement>\n")
	c.process("return")
	if c.getCurrentToken() != ";" {
		c.CompileExpression()
	}
	c.process(";")
	c.outputFile.WriteString("</returnStatement>\n")
}

func (c *Compiler) CompileExpression() {
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

func (c *Compiler) CompileTerm() {
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

func (c *Compiler) CompileExpressionList() int {
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
