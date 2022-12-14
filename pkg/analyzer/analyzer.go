package analyzer

import (
	"flag"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	//nolint:lll
	commentMismatchTmpl = "first word of comment is '%s' instead of '%s'"
	commentEmptyTmpl    = "empty comment on '%s'"
	commentMissingTmpl  = "exported element '%s' should be commented"

	testFileNameSuffix = "_test.go"
)

var testCommentPrefix = []string{
	"Benchmark",
	"Example",
	"Fuzz",
	"Test",
}

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
		elementPos,
	)
	checkExported(
		pass,
		commentExported,
		commentAllExported,
		elementName,
		comment,
		elementPos,
		elementExported,
		recvExported,
	)
}

func extractSingleCommentText(input string) string {
	// Comment of the form /**/.
	if input[1] == '*' {
		return input[2 : len(input)-2]
	}

	// comment of the form //.
	return input[2:]
}

// containsOnlyMachineReadableComment returns true if the CommentGroup
// contains only a machine-readable comment.
func containsOnlyMachineReadableComment(comment *ast.CommentGroup) bool {
	onlyMachine := true

	for _, lc := range comment.List {
		lineText := strings.TrimSpace(extractSingleCommentText(lc.Text))
		// Machine-readable comment.
		if len(lineText) > 0 {
			continue
		}

		onlyMachine = false
	}

	return onlyMachine
}

func checkCommentMismatch(
	pass *analysis.Pass,
	elementName string,
	comment *ast.CommentGroup,
	elementPos token.Pos,
) {
	if comment == nil {
		return
	}

	text := comment.Text()

	// This comment could be a machine-readable comment of the form
	// //something:else, an empty comment, or a comment containing only
	// whitespace.
	if len(text) == 0 {
		if !containsOnlyMachineReadableComment(comment) {
			// Empty comment.
			pass.Reportf(
				elementPos,
				commentEmptyTmpl,
				elementName,
			)
		}

		return
	}

	words := strings.Fields(strings.TrimSpace(text))

	// Set to empty if there's no non-whitespace characters in the comment.
	firstWord := ""
	if len(words) > 0 {
		firstWord = words[0]
	}

	if firstWord == elementName {
		return
	}

	pass.Reportf(
		comment.Pos(),
		commentMismatchTmpl,
		firstWord,
		elementName,
	)
}

func checkExported(
	pass *analysis.Pass,
	commentExported bool,
	commentAllExported bool,
	elementName string,
	comment *ast.CommentGroup,
	elementPos token.Pos,
	elementExported bool,
	recvExported bool,
) {
	commented := false
	// We only want to report missing comments if we haven't already reported the
	// comment is empty.
	if comment != nil {
		commented = len(comment.Text()) > 0 ||
			!containsOnlyMachineReadableComment(comment)
	}

	if commented || !elementExported {
		return
	}

	// Either we're commenting everything or the receiver is exported and we're
	// only commenting things with exported receivers and elements.
	if commentAllExported || (recvExported && commentExported) {
		pass.Reportf(
			elementPos,
			commentMissingTmpl,
			elementName,
		)
	}
}

func isTestFunc(fset *token.FileSet, pos token.Pos, elementName string) bool {
	fName := fset.Position(pos).Filename
	if !strings.HasSuffix(fName, testFileNameSuffix) {
		return false
	}

	for _, name := range testCommentPrefix {
		if strings.HasPrefix(elementName, name) {
			return true
		}
	}

	return false
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

	commentExported := m.commentExportedFuncs
	commentAllExported := m.commentAllExportedFuncs

	if isTestFunc(pass.Fset, fun.Pos(), fun.Name.Name) {
		commentExported = false
		commentAllExported = false
	}

	checkComment(
		pass,
		commentExported,
		commentAllExported,
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
	commentExportedFuncs    bool
	commentAllExportedFuncs bool
	commentInterfaces       bool
	noTestComments          bool
}

func NewCommentMimic() *analysis.Analyzer {
	m := mimic{}

	fs := flag.NewFlagSet("CommentMimicFlags", flag.ExitOnError)
	fs.BoolVar(
		&m.commentExportedFuncs,
		"comment-exported",
		false,
		"require comments on exported functions if their receiver is also exported",
	)

	fs.BoolVar(
		&m.commentAllExportedFuncs,
		"comment-all-exported",
		false,
		"require comments on all exported functions",
	)

	fs.BoolVar(
		&m.commentInterfaces,
		"comment-interfaces",
		false,
		"require comments on all exported interfaces",
	)

	fs.BoolVar(
		&m.noTestComments,
		"no-test-comments",
		true,
		"don't require comments on tests, benchmarks, examples, and fuzz tests",
	)

	return &analysis.Analyzer{
		Name: "commentmimic",
		//nolint:lll
		Doc:      "Checks function/interface first words match the element name and exported element are commented",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Flags:    *fs,
		Run: func(pass *analysis.Pass) (any, error) {
			return m.run(pass)
		},
	}
}
