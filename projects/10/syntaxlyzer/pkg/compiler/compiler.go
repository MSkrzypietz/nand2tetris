package compiler

import "syntaxlyzer/pkg/tokenizer"

type Compiler struct {
	tokenizer *tokenizer.Tokenizer
}

func New(filename string) *Compiler {
	return &Compiler{
		tokenizer: tokenizer.New(filename),
	}
}

func (c *Compiler) process(s string) {
	// c.tokenizer.
}

func (c *Compiler) compileClass() {}

func (c *Compiler) compileClassVarDec() {}

func (c *Compiler) compileSubroutine() {}

func (c *Compiler) compileParameterList() {}

func (c *Compiler) compileSubroutinBody() {}

func (c *Compiler) compileVarDec() {}

func (c *Compiler) compileStatements() {}

func (c *Compiler) compileLet() {}

func (c *Compiler) compileIf() {}

func (c *Compiler) compileWhile() {}

func (c *Compiler) compileDo() {}

func (c *Compiler) compileReturn() {}

func (c *Compiler) compileExpression() {}

func (c *Compiler) compileTerm() {}

func (c *Compiler) compileExpressionList() int {
	return 0
}
