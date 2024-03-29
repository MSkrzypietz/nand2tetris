package main

import (
	"compiler/pkg/compengine"
	"compiler/pkg/symtable"
	"compiler/pkg/tokenizer"
	"compiler/pkg/vmwriter"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	inputPath := path.Clean(os.Args[1])

	var filePaths []string
	inputPathStats, _ := os.Stat(inputPath)
	if inputPathStats.IsDir() {
		filepath.WalkDir(inputPath, func(filePath string, d fs.DirEntry, err error) error {
			if path.Ext(filePath) == ".jack" {
				filePaths = append(filePaths, filePath)
			}
			return nil
		})
	} else {
		filePaths = append(filePaths, inputPath)
	}

	os.RemoveAll("out")
	os.Mkdir("out", os.ModePerm)
	for _, filePath := range filePaths {
		filename := path.Base(filePath)
		xmlOutputFilename := filename[:strings.LastIndex(filename, ".")] + ".vm"
		outputFile, err := os.Create(path.Join("out", xmlOutputFilename))
		if err != nil {
			log.Fatal(err)
		}

		println("compiling", filePath)
		t := tokenizer.New(filePath)
		w := vmwriter.New(outputFile)
		classSymTable := symtable.New()
		subroutineSymTable := symtable.New()
		c := compengine.New(t, w, classSymTable, subroutineSymTable, outputFile)
		c.CompileClass()

		outputFile.Sync()
		outputFile.Close()
	}
}
