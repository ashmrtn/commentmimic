package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/ashmrtn/commenter/pkg/analyzer"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
