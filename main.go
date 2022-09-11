package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/ashmrtn/commentmimic/pkg/analyzer"
)

func main() {
	a := analyzer.NewCommentMimic()
	singlechecker.Main(a)
}
