package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "commenter",
	Doc:      "Checks for the golang comment best-practices in comments",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspec := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspec.Nodes(nodeFilter, func(node ast.Node, push bool) bool {
		fd, ok := node.(*ast.FuncDecl)
		if !ok {
			return false
		}

		if fd.Doc == nil {
			return false
		}

		text := fd.Doc.Text()
		if len(text) == 0 {
			return false
		}

		firstWord := strings.SplitN(text, " ", 2)[0]
		if firstWord == fd.Name.Name {
			return false
		}

		pass.Reportf(
			fd.Doc.Pos(),
			"first word of doc comment for function '%s' should be '%[1]s'",
			fd.Name.Name,
		)

		return false
	})

	return nil, nil
}
