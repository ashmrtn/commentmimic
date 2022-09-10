package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/ashmrtn/commenter/pkg/analyzer"
)

func main() {
	a := analyzer.NewCommentMimic()
	singlechecker.Main(a)
}
