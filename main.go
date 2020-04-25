package main

import (
	"flag"
	"fmt"
	"os"

	"joom-sort/ext_merge_sort"
	"joom-sort/generator"

	"github.com/sirupsen/logrus"
)

func main() {
	var (
		isGenerateMode bool
		isSortMode     bool
		filePath       string
		lines          int64
		maxLength      int
	)

	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)
	//logrus.SetFormatter(&(logrus.JSONFormatter{}))
	logrus.SetLevel(logrus.DebugLevel)

	flag.BoolVar(&isGenerateMode, "generate", false,
		"Generate new file for sorting")
	flag.BoolVar(&isSortMode, "sort", true,
		"Sort file given by -file param")
	flag.StringVar(&filePath, "file", "",
		"File path for operation")
	flag.Int64Var(&lines, "lines", 10000,
		"Set lines count for the file generating (default 10000)")
	flag.IntVar(&maxLength, "rowlen", 256,
		"Set max line length for the file generating (default 256")

	flag.Parse()

	if filePath == "" {
		fmt.Println("Please, provide -file param")
		return
	}
	if isGenerateMode {
		fmt.Printf("Starting generating file by path %q\n", filePath)
		if err := generator.GenerateFile(filePath, lines, maxLength); err != nil {
			panic("Exit: Error happend")
		}

		fmt.Printf("File has been generated %q\n", filePath)
		return
	}

	if isSortMode {
		fmt.Printf("Starting sort of file %q\n", filePath)

		msort := ext_merge_sort.New(filePath)
		if err := msort.Sort(); err != nil {
			panic("Exit: Error happend")
		}

		fmt.Printf("File has been sorted")
		return
	}
}
