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
		inputPath      string
		outputPath     string
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
		"Sort file given by -in param")
	flag.StringVar(&inputPath, "in", "",
		"Input file path")
	flag.StringVar(&outputPath, "out", "",
		"Destination file path for operation")
	flag.Int64Var(&lines, "lines", 10000,
		"Set lines count for the file generating (default 10000)")
	flag.IntVar(&maxLength, "rowlen", 256,
		"Set max line length for the file generating (default 256")

	flag.Parse()

	if outputPath == "" {
		fmt.Println("Please, provide -out param with the address for result file")
		return
	}

	if isGenerateMode {
		fmt.Printf("Starting generating file by path %q\n", outputPath)
		if err := generator.GenerateFile(outputPath, lines, maxLength); err != nil {
			panic("Exit: Error happend")
		}

		fmt.Println("File has been generated")
		return
	}

	if isSortMode {
		if inputPath == "" {
			fmt.Println("Please, provide -in param")
			return
		}
		fmt.Printf("Starting sort of file %q\n", inputPath)

		msort := ext_merge_sort.New(inputPath, outputPath)
		if err := msort.Sort(); err != nil {
			panic("Exit: Error happend")
		}

		fmt.Printf("File has been sorted")
		return
	}
}
