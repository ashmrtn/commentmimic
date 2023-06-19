package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/ashmrtn/commentmimic/pkg/commentmimic"
)

func main() {
	a := analyzer.NewCommentMimic()
	singlechecker.Main(a)
}
