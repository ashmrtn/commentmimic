package analyzer

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	//nolint:lll
	commentMismatchTmpl = "first word of comment for element '%s' should be '%[1]s' not '%s'"
)

func checkComment(
	pass *analysis.Pass,
	commentExported bool,
	commentAllExported bool,
	elementName string,
	elementPos token.Pos,
	comment *ast.CommentGroup,
	elementExported bool,
	recvExported bool,
) {
	checkCommentMismatch(
		pass,
		elementName,
		comment,
	)
}

func checkCommentMismatch(
	pass *analysis.Pass,
	elementName string,
	comment *ast.CommentGroup,
) {
	if comment == nil {
		return
	}

	firstWord := strings.SplitN(comment.Text(), " ", 2)[0]
	if firstWord == elementName {
		return
	}

	pass.Reportf(
		comment.Pos(),
		commentMismatchTmpl,
		elementName,
		firstWord,
	)
}

func (m mimic) checkFuncDecl(pass *analysis.Pass, fun *ast.FuncDecl) {
	// Default to true so free functions will be marked as needing a comment if
	// commentExported is set.
	exportedRecv := true

	if fun.Recv != nil {
		r := fun.Recv.List[0]

		switch r.Type.(type) {
		case *ast.Ident:
			ident := r.Type.(*ast.Ident)
			exportedRecv = ident.IsExported()

		case *ast.StarExpr:
			star := r.Type.(*ast.StarExpr)

			ident, ok := star.X.(*ast.Ident)
			if !ok {
				break
			}

			exportedRecv = ident.IsExported()
		}
	}

	checkComment(
		pass,
		m.commentExportedFuncs,
		m.commentAllExportedFuncs,
		fun.Name.Name,
		fun.Pos(),
		fun.Doc,
		fun.Name.IsExported(),
		exportedRecv,
	)
}

func (m mimic) checkGenDecl(pass *analysis.Pass, decl *ast.GenDecl) {
	for _, s := range decl.Specs {
		ts, ok := s.(*ast.TypeSpec)
		if !ok {
			continue
		}

		exportedRecv := ts.Name.IsExported()

		iface, ok := ts.Type.(*ast.InterfaceType)
		if !ok {
			continue
		}

		// If the type-declaration a single declaration (i.e. not grouped by
		// parentheses), then the doc comment is attached to the GenDecl node. If
		// the interface is part of a grouped declaration, it's attached to the
		// TypeSpec node.
		doc := decl.Doc
		pos := decl.Pos()

		if decl.Lparen != token.NoPos {
			doc = ts.Doc
			pos = ts.Pos()
		}

		// Check if interface is commented properly.
		checkComment(
			pass,
			// Set to false so the flag completely controls output behavior.
			false,
			m.commentInterfaces,
			ts.Name.Name,
			pos,
			doc,
			exportedRecv,
			true,
		)

		for _, field := range iface.Methods.List {
			_, ok := field.Type.(*ast.FuncType)
			if !ok {
				continue
			}

			checkComment(
				pass,
				m.commentExportedFuncs,
				m.commentAllExportedFuncs,
				field.Names[0].Name,
				field.Pos(),
				field.Doc,
				field.Names[0].IsExported(),
				exportedRecv,
			)
		}
	}
}

func (m mimic) run(pass *analysis.Pass) (any, error) {
	inspec := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.GenDecl)(nil),
	}

	inspec.Nodes(nodeFilter, func(node ast.Node, push bool) bool {
		switch switched := node.(type) {
		case *ast.FuncDecl:
			m.checkFuncDecl(pass, switched)

		case *ast.GenDecl:
			if switched.Tok != token.TYPE {
				break
			}

			m.checkGenDecl(pass, switched)
		}

		return false
	})

	return nil, nil
}

type mimic struct {
}

func NewCommentMimic() *analysis.Analyzer {
	m := mimic{}

	return &analysis.Analyzer{
		Name: "CommentMimic",
		//nolint:lll
		Doc:      "Checks function/interface first words match the element name and exported element are commented",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run: func(pass *analysis.Pass) (any, error) {
			return m.run(pass)
		},
	}
}
