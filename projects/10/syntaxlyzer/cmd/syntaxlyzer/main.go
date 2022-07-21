package main

import (
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"syntaxlyzer/pkg/tokenizer"
)

func main() {
	filePath := os.Args[1]

	filename := path.Base(filePath)
	tokenOutputFilename := filename[:strings.LastIndex(filename, ".")] + "T.xml"
	os.Mkdir("out", os.ModePerm)
	outputFile, err := os.Create(path.Join("out", tokenOutputFilename))
	if err != nil {
		log.Fatal(err)
	}

	writeTokenOutput(filePath, outputFile)

	outputFile.Sync()
	outputFile.Close()
}

func writeTokenOutput(inputFilePath string, outputFile *os.File) {
	t := tokenizer.New(inputFilePath)
	outputFile.WriteString("<tokens>\n")
	for t.HasMoreTokens() {
		t.Advance()

		switch t.TokenType() {
		case tokenizer.Keyword:
			outputFile.WriteString(createXMLToken("keyword", t.KeyWord()))
		case tokenizer.Symbol:
			outputFile.WriteString(createXMLToken("symbol", t.Symbol()))
		case tokenizer.Identifier:
			outputFile.WriteString(createXMLToken("identifier", t.Identifier()))
		case tokenizer.IntConst:
			outputFile.WriteString(createXMLToken("integerConstant", strconv.Itoa(t.IntVal())))
		case tokenizer.StringConst:
			outputFile.WriteString(createXMLToken("stringConstant", t.StringVal()))
		}
	}
	outputFile.WriteString("</tokens>\n")
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
