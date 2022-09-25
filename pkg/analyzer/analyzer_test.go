package analyzer_test

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/ashmrtn/commentmimic/pkg/analyzer"
	"github.com/ashmrtn/commentmimic/pkg/analyzer/testdata"
)

const (
	lower = true
	upper = !lower

	exported   = true
	unexported = !exported
)

var (
	fileTmpl = template.Must(
		template.New("fileGenerator").Parse(testdata.FileGenTmpls),
	)
	flagProduct = []map[string]bool{
		{
			"comment-exported":     false,
			"comment-all-exported": false,
			"comment-interfaces":   false,
		},
		{
			"comment-exported":     false,
			"comment-all-exported": false,
			"comment-interfaces":   true,
		},
		{
			"comment-exported":     false,
			"comment-all-exported": true,
			"comment-interfaces":   false,
		},
		{
			"comment-exported":     false,
			"comment-all-exported": true,
			"comment-interfaces":   true,
		},
		{
			"comment-exported":     true,
			"comment-all-exported": false,
			"comment-interfaces":   false,
		},
		{
			"comment-exported":     true,
			"comment-all-exported": false,
			"comment-interfaces":   true,
		},
		{
			"comment-exported":     true,
			"comment-all-exported": true,
			"comment-interfaces":   false,
		},
		{
			"comment-exported":     true,
			"comment-all-exported": true,
			"comment-interfaces":   true,
		},
	}
)

type commentData struct {
	Type      testdata.CommentType
	Text      string
	Multiline bool
}

type templateData struct {
	// Name of the test case.
	name string
	// Template this should be executed with. For tests not the templating system.
	template string

	// If an error should be reported for the comment, and the first word of the
	// comment. If firstWord is empty then neither the comment nor the error check
	// will be output.
	CommentError bool
	FirstWord    commentData

	// Name of the element and whether an error should be reported for the
	// element, often if it is missing a comment.
	Element      string
	ElementError bool

	// Extra info that may be affected by flags or used to try to confuse the
	// linter.
	Receiver                string
	InterfaceBlockFirstWord commentData

	// Inner information for an interface; whether it should have other stuff
	// around it and if it has functions.
	Confusing     bool
	InterfaceFunc *templateData
}

func caseWord(w string, toLower bool) string {
	res := ""

	if len(w) == 0 {
		return res
	}

	if toLower {
		res += strings.ToLower(string(w[0]))
	} else {
		res += strings.ToUpper(string(w[0]))
	}

	if len(w) > 1 {
		res += w[1:]
	}

	return res
}

func lowerWord(w string) string {
	return caseWord(w, lower)
}

func capWord(w string) string {
	return caseWord(w, upper)
}

func kebabToCamel(f string) string {
	res := ""

	parts := strings.Split(f, "-")
	for _, p := range parts {
		res += capWord(p)
	}

	return res
}

func flagsToTestName(flags map[string]bool) string {
	items := make([]string, 0, len(flags))

	for k := range flags {
		items = append(items, k)
	}

	sort.Strings(items)

	res := ""

	for _, f := range items {
		if !flags[f] {
			continue
		}

		res += kebabToCamel(f)
	}

	if len(res) == 0 {
		return "AllDisabled"
	}

	return res
}

func generateCommentMimicCases(name string) []templateData {
	base := []templateData{
		{
			name: "UnexportedNoError",
			FirstWord: commentData{
				Text: lowerWord(name),
			},
			Element: lowerWord(name),
		},
		{
			name: "ExportedNoError",
			FirstWord: commentData{
				Text: capWord(name),
			},
			Element: capWord(name),
		},
		{
			name: "UnexportedWrongCase",
			FirstWord: commentData{
				Text: capWord(name),
			},
			Element:      lowerWord(name),
			CommentError: true,
		},
		{
			name: "ExportedWrongCase",
			FirstWord: commentData{
				Text: lowerWord(name),
			},
			Element:      capWord(name),
			CommentError: true,
		},
		{
			name: "UnexportedPostfixWord",
			FirstWord: commentData{
				Text: lowerWord(name) + "a",
			},
			Element:      lowerWord(name),
			CommentError: true,
		},
		{
			name: "ExportedPostfixWord",
			FirstWord: commentData{
				Text: capWord(name) + "a",
			},
			Element:      capWord(name),
			CommentError: true,
		},
	}

	commentTypes := []testdata.CommentType{
		testdata.InlineComment,
		testdata.BlockInlineComment,
		testdata.BlockMultilineComment,
	}
	res := make([]templateData, 0, 3*len(commentTypes)*len(base))

	for _, c := range base {
		for _, t := range commentTypes {
			for _, multiline := range []bool{false, true} {
				tmp := c
				tmp.FirstWord.Type = t
				tmp.FirstWord.Multiline = multiline
				tmp.name += t.String()

				if multiline {
					tmp.name += "Multiline"
				} else {
					tmp.name += "SingleLine"
				}

				res = append(res, tmp)
			}
		}
	}

	return res
}

// genFunctionCases takes a set of partially populated test cases and returns a
// set of fully populated test cases. For each generated case, the function and
// receiver has the same export status as the partially populated case (i.e.
// partially populated is unexported then the function will be unexported). The
// generated tests cover:
//   - free functions
//   - receiver functions
//   - pointer receiver functions
func genFunctionCases(
	tests []templateData,
	recvName string,
	funcExported bool,
	receiverExported bool,
) map[string][]templateData {
	res := map[string][]templateData{}
	templates := map[string]string{
		"_FreeFunc":         "FreeFunction",
		"-Receiver_Func":    "ReceiverFunction",
		"-ReceiverPtr_Func": "ReceiverPtrFunction",
	}

	for t, tmplName := range templates {
		for _, testCase := range tests {
			fullName := t
			if funcExported {
				fullName = strings.ReplaceAll(fullName, "_", "Exported")
			} else {
				fullName = strings.ReplaceAll(fullName, "_", "Unexported")
			}

			if receiverExported {
				fullName = strings.ReplaceAll(fullName, "-", "Exported")
			} else {
				fullName = strings.ReplaceAll(fullName, "-", "Unexported")
			}

			testCase.template = tmplName
			testCase.Receiver = recvName

			res[fullName] = append(res[fullName], testCase)
		}
	}

	return res
}

// genFunctionCasesWithExports takes a set of partially populated test cases and
// returns a set of fully populated test cases. For each generated case, the
// function has the same export status as the partially populated case (i.e.
// partially populated is unexported then the function will be unexported). The
// generated tests cover:
//   - free functions
//   - unexported receiver functions
//   - exported receiver functions
//   - unexported pointer receiver functions
//   - exported pointer receiver functions
func genFunctionCasesWithExports(
	tests []templateData,
	recvName string,
) map[string][]templateData {
	res := map[string][]templateData{
		"FreeFunction":                  make([]templateData, 0, len(tests)),
		"UnexportedReceiverFunction":    make([]templateData, 0, len(tests)),
		"ExportedReceiverFunction":      make([]templateData, 0, len(tests)),
		"UnexportedPtrReceiverFunction": make([]templateData, 0, len(tests)),
		"ExportedPtrReceiverFunction":   make([]templateData, 0, len(tests)),
	}

	templates := []struct {
		name     string
		template string
		recvFunc func(string) string
	}{
		{
			name:     "FreeFunction",
			template: "FreeFunction",
			recvFunc: func(s string) string {
				return s
			},
		},
		{
			name:     "UnexportedReceiverFunction",
			template: "ReceiverFunction",
			recvFunc: lowerWord,
		},
		{
			name:     "ExportedReceiverFunction",
			template: "ReceiverFunction",
			recvFunc: capWord,
		},
		{
			name:     "UnexportedPtrReceiverFunction",
			template: "ReceiverPtrFunction",
			recvFunc: lowerWord,
		},
		{
			name:     "ExportedPtrReceiverFunction",
			template: "ReceiverPtrFunction",
			recvFunc: capWord,
		},
	}

	for _, t := range templates {
		for _, testCase := range tests {
			testCase.template = t.template
			testCase.Receiver = t.recvFunc(recvName)

			res[t.name] = append(res[t.name], testCase)
		}
	}

	return res
}

func genInterfaceFuncCases(
	tests []templateData,
	interfaceName string,
	interfaceExported bool,
	funcExported bool,
) map[string][]templateData {
	type info struct {
		template      string
		confusing     bool
		confusingFunc bool
	}

	res := map[string][]templateData{}
	templates := map[string]info{
		"-InterfaceOne_Func": {
			template: "Interface",
		},
		"-InterfaceMany_Func": {
			template:      "Interface",
			confusingFunc: true,
		},
		"OneBlock-InterfaceOne_Func": {
			template: "BlockInterface",
		},
		"OneBlock-InterfaceMany_Func": {
			template:      "BlockInterface",
			confusing:     true,
			confusingFunc: true,
		},
		"ManyBlock-InterfaceOne_Func": {
			template: "BlockInterface",
		},
		"ManyBlock-InterfaceMany_Func": {
			template:      "BlockInterface",
			confusing:     true,
			confusingFunc: true,
		},
	}

	for t, tmpl := range templates {
		for _, testCase := range tests {
			tmpTestCase := testCase

			fullName := t
			if funcExported {
				fullName = strings.ReplaceAll(fullName, "_", "Exported")
			} else {
				fullName = strings.ReplaceAll(fullName, "_", "Unexported")
			}

			if interfaceExported {
				fullName = strings.ReplaceAll(fullName, "-", "Exported")
			} else {
				fullName = strings.ReplaceAll(fullName, "-", "Unexported")
			}

			iface := templateData{
				name:          tmpTestCase.name,
				template:      tmpl.template,
				Confusing:     tmpl.confusing,
				InterfaceFunc: &tmpTestCase,
				Element:       interfaceName,
			}

			iface.InterfaceFunc.Confusing = tmpl.confusingFunc

			res[fullName] = append(res[fullName], iface)
		}
	}

	return res
}

func genEmptyInterfaceCases(
	tests []templateData,
	interfaceExported bool,
) map[string][]templateData {
	res := map[string][]templateData{}
	templates := map[string]templateData{
		"-Interface": {
			template: "Interface",
		},
		"OneBlock-Interface": {
			template: "BlockInterface",
		},
		"ManyBlock-Interface": {
			template:  "BlockInterface",
			Confusing: true,
		},
	}

	for t, tmpl := range templates {
		for _, testCase := range tests {
			fullName := t
			if interfaceExported {
				fullName = strings.ReplaceAll(fullName, "-", "Exported")
			} else {
				fullName = strings.ReplaceAll(fullName, "-", "Unexported")
			}

			testCase.template = tmpl.template
			testCase.Confusing = tmpl.Confusing
			res[fullName] = append(res[fullName], testCase)
		}
	}

	return res
}

func genTestFiles(
	t *testing.T,
	tmpl *template.Template,
	tmplName string,
	test templateData,
) (string, func()) {
	t.Helper()

	buf := &bytes.Buffer{}
	require.NoError(t, tmpl.ExecuteTemplate(buf, tmplName, test))
	fullFile := fmt.Sprintf("package a\n\n%s\n", buf)

	fileMap := map[string]string{
		"a/a.go": fullFile,
	}

	dir, cleanup, err := analysistest.WriteFiles(fileMap)
	require.NoError(t, err)

	return dir, cleanup
}

func executeMimicWithFlagsOnFiles(
	t *testing.T,
	flags map[string]bool,
	dir string,
) {
	t.Helper()

	mimic := analyzer.NewCommentMimic()

	for flag, value := range flags {
		require.NoError(t, mimic.Flags.Set(flag, strconv.FormatBool(value)))
	}

	analysistest.Run(t, dir, mimic, "a")
}

func executeCommentMimicWithAllFlagCombos(
	t *testing.T,
	tmpl *template.Template,
	tmplName string,
	test templateData,
) {
	t.Helper()

	dir, cleanup := genTestFiles(t, tmpl, tmplName, test)
	defer cleanup()

	for _, flags := range flagProduct {
		flags := flags

		t.Run(flagsToTestName(flags), func(t1 *testing.T) {
			executeMimicWithFlagsOnFiles(t1, flags, dir)
		})
	}
}

func executeCommentMimic(
	t *testing.T,
	tmpl *template.Template,
	tmplName string,
	test templateData,
	flags map[string]bool,
) {
	t.Helper()

	dir, cleanup := genTestFiles(t, tmpl, tmplName, test)
	defer cleanup()

	executeMimicWithFlagsOnFiles(t, flags, dir)
}

type CommentMimicSuite struct {
	suite.Suite
}

func TestCommentMimic(t *testing.T) {
	suite.Run(t, new(CommentMimicSuite))
}

func (s *CommentMimicSuite) TestDoesNotErrorOnOutOfScope() {
	for _, flags := range flagProduct {
		flags := flags

		s.T().Run(flagsToTestName(flags), func(t *testing.T) {
			t.Parallel()

			fileMap := map[string]string{
				"a/a.go": testdata.OutOfScopePatterns,
			}

			dir, cleanup, err := analysistest.WriteFiles(fileMap)
			require.NoError(t, err)

			defer cleanup()

			mimic := analyzer.NewCommentMimic()

			for flag, value := range flags {
				require.NoError(s.T(), mimic.Flags.Set(flag, strconv.FormatBool(value)))
			}

			analysistest.Run(t, dir, mimic, "a")
		})
	}
}

func (s *CommentMimicSuite) TestHandlesExtraWhitespace() {
	t := s.T()

	fileMap := map[string]string{
		"a/a.go": testdata.ExtraWhitespace,
	}

	dir, cleanup, err := analysistest.WriteFiles(fileMap)
	require.NoError(t, err)

	defer cleanup()

	mimic := analyzer.NewCommentMimic()
	analysistest.Run(t, dir, mimic, "a")
}

func (s *CommentMimicSuite) TestMachineCommentsMismatch() {
	t := s.T()
	flags := map[string]bool{
		"comment-exported":     true,
		"comment-all-exported": true,
		"comment-interfaces":   true,
	}

	fileMap := map[string]string{
		"a/a.go": testdata.EmptyComments,
	}

	dir, cleanup, err := analysistest.WriteFiles(fileMap)
	require.NoError(t, err)

	defer cleanup()

	mimic := analyzer.NewCommentMimic()

	for flag, value := range flags {
		require.NoError(t, mimic.Flags.Set(flag, strconv.FormatBool(value)))
	}

	analysistest.Run(t, dir, mimic, "a")
}

func (s *CommentMimicSuite) TestMachineCommentsOnExported() {
	t := s.T()
	flags := map[string]bool{
		"comment-all-exported": true,
		"comment-interfaces":   true,
	}

	fileMap := map[string]string{
		"a/a.go": testdata.MachineReadableExported,
	}

	dir, cleanup, err := analysistest.WriteFiles(fileMap)
	require.NoError(t, err)

	defer cleanup()

	mimic := analyzer.NewCommentMimic()

	for flag, value := range flags {
		require.NoError(t, mimic.Flags.Set(flag, strconv.FormatBool(value)))
	}

	analysistest.Run(t, dir, mimic, "a")
}

func (s *CommentMimicSuite) TestFuncCommentErrors() {
	element := "element"
	base := generateCommentMimicCases(element)

	receiver := "receiver"
	all := genFunctionCasesWithExports(base, receiver)

	for name, tests := range all {
		name := name
		tests := tests

		s.T().Run(name, func(t1 *testing.T) {
			t1.Parallel()

			for _, test := range tests {
				test := test

				t1.Run(test.name, func(t *testing.T) {
					t.Parallel()
					executeCommentMimicWithAllFlagCombos(
						t,
						fileTmpl,
						test.template,
						test,
					)
				})
			}
		})
	}
}

func (s *CommentMimicSuite) TestEmptyInterfaceCommentErrors() {
	element := "element"
	table := generateCommentMimicCases(element)
	patterns := []struct {
		name      string
		template  string
		confusing bool
	}{
		{
			name:      "Interface",
			template:  "Interface",
			confusing: false,
		},
		{
			name:      "BlockOneInterface",
			template:  "BlockInterface",
			confusing: false,
		},
		{
			name:      "BlockSeveralInterfaces",
			template:  "BlockInterface",
			confusing: true,
		},
	}

	for _, pattern := range patterns {
		s.T().Run(pattern.name, func(t1 *testing.T) {
			pattern := pattern

			t1.Parallel()

			for _, test := range table {
				test := test

				t1.Run(test.name, func(t *testing.T) {
					t.Parallel()

					test.template = pattern.template
					test.Confusing = pattern.confusing

					executeCommentMimicWithAllFlagCombos(
						t,
						fileTmpl,
						test.template,
						test,
					)
				})
			}
		})
	}
}

func (s *CommentMimicSuite) TestInterfaceFuncCommentErrors() {
	element := "element"
	elementFunc := "elementFunc"
	funcs := generateCommentMimicCases(elementFunc)

	patterns := []struct {
		name      string
		template  string
		exported  bool
		confusing bool
	}{
		{
			name:      "UnexportedInterface",
			template:  "Interface",
			exported:  false,
			confusing: false,
		},
		{
			name:      "ExportedInterface",
			template:  "Interface",
			exported:  true,
			confusing: false,
		},
		{
			name:      "BlockOneUnexportedInterface",
			template:  "BlockInterface",
			exported:  false,
			confusing: false,
		},
		{
			name:      "BlockOneExportedInterface",
			template:  "BlockInterface",
			exported:  true,
			confusing: false,
		},
		{
			name:      "BlockSeveralUnexportedInterfaces",
			template:  "BlockInterface",
			exported:  false,
			confusing: true,
		},
		{
			name:      "BlockSeveralExportedInterfaces",
			template:  "BlockInterface",
			exported:  true,
			confusing: true,
		},
	}

	funcPatterns := []struct {
		name      string
		hasFunc   bool
		confusing bool
	}{
		{
			name:      "OneFunc",
			hasFunc:   true,
			confusing: false,
		},
		{
			name:      "SeveralFunc",
			hasFunc:   true,
			confusing: true,
		},
	}

	for _, pattern := range patterns {
		pattern := pattern

		s.T().Run(pattern.name, func(t1 *testing.T) {
			t1.Parallel()

			elementName := lowerWord(element)
			if pattern.exported {
				elementName = capWord(element)
			}

			test := templateData{
				template: pattern.template,
				Element:  elementName,
				FirstWord: commentData{
					Type: testdata.InlineComment,
					Text: elementName,
				},
				Confusing: pattern.confusing,
			}

			for _, funcPattern := range funcPatterns {
				funcPattern := funcPattern

				t1.Run(funcPattern.name, func(t2 *testing.T) {
					t2.Parallel()

					// Interface with functions.
					for _, funcCase := range funcs {
						funcCase := funcCase

						t2.Run(funcCase.name, func(t *testing.T) {
							t.Parallel()

							funcCase.Confusing = funcPattern.confusing
							test.InterfaceFunc = &funcCase

							executeCommentMimicWithAllFlagCombos(
								t,
								fileTmpl,
								test.template,
								test,
							)
						})
					}
				})
			}
		})
	}
}

func (s *CommentMimicSuite) TestCommentAccessibleExportedFuncs() {
	const (
		element  = "element"
		receiver = "receiver"
		iface    = "iface"
	)

	var (
		flags = map[string]bool{
			"comment-exported":     true,
			"comment-all-exported": false,
			// Turn off as this will be testing some exported interfaces that we don't
			// want to comment.
			"comment-interfaces": false,
		}

		funcCases = []templateData{
			{
				name: "NoError",
				FirstWord: commentData{
					Type: testdata.InlineComment,
					Text: capWord(element),
				},
				Element: capWord(element),
			},
			{
				name: "MimicError",
				FirstWord: commentData{
					Type: testdata.InlineComment,
					Text: "foo",
				},
				CommentError: true,
				Element:      capWord(element),
			},
			{
				name:         "Error",
				Element:      capWord(element),
				ElementError: true,
			},
		}
	)

	cases := genFunctionCases(funcCases, capWord(receiver), exported, exported)
	cases["UnexportedFreeFunc"] = append(
		cases["UnexportedFreeFunc"],
		templateData{
			name:     "NoError",
			template: "FreeFunction",
			Element:  lowerWord(element),
		},
	)
	cases["UnexportedReceiverExportedFunc"] = append(
		cases["UnexportedReceiverExportedFunc"],
		templateData{
			name:     "NoError",
			template: "ReceiverFunction",
			Element:  capWord(element),
			Receiver: lowerWord(receiver),
		},
	)
	cases["UnexportedReceiverPtrExportedFunc"] = append(
		cases["UnexportedReceiverPtrExportedFunc"],
		templateData{
			name:     "NoError",
			template: "ReceiverPtrFunction",
			Element:  capWord(element),
			Receiver: lowerWord(receiver),
		},
	)

	for k, v := range genInterfaceFuncCases(
		funcCases,
		capWord(iface),
		exported,
		exported,
	) {
		cases[k] = append(cases[k], v...)
	}

	for k, v := range genInterfaceFuncCases(
		[]templateData{
			{
				name:    "NoError",
				Element: capWord(element),
			},
		},
		lowerWord(iface),
		unexported,
		exported,
	) {
		cases[k] = append(cases[k], v...)
	}

	for name, caseList := range cases {
		name := name
		caseList := caseList

		s.T().Run(name, func(t1 *testing.T) {
			t1.Parallel()

			for _, test := range caseList {
				test := test

				t1.Run(test.name, func(t *testing.T) {
					t.Parallel()
					executeCommentMimic(
						t,
						fileTmpl,
						test.template,
						test,
						flags,
					)
				})
			}
		})
	}
}

func (s *CommentMimicSuite) TestCommentAllExportedFuncs() {
	const (
		element  = "element"
		receiver = "receiver"
		iface    = "iface"
	)

	var (
		flagSets = []map[string]bool{
			{
				"comment-exported":     false,
				"comment-all-exported": true,
				// Turn off as this will be testing some exported interfaces that we
				// don't want to comment.
				"comment-interfaces": false,
			},
			{
				"comment-exported":     true,
				"comment-all-exported": true,
				// Turn off as this will be testing some exported interfaces that we
				// don't want to comment.
				"comment-interfaces": false,
			},
		}

		funcCases = []templateData{
			{
				name: "NoError",
				FirstWord: commentData{
					Type: testdata.InlineComment,
					Text: capWord(element),
				},
				Element: capWord(element),
			},
			{
				name: "MimicError",
				FirstWord: commentData{
					Type: testdata.InlineComment,
					Text: "foo",
				},
				CommentError: true,
				Element:      capWord(element),
			},
			{
				name:         "Error",
				Element:      capWord(element),
				ElementError: true,
			},
		}
	)

	cases := genFunctionCases(funcCases, capWord(receiver), exported, exported)
	cases["UnexportedFreeFunc"] = append(
		cases["UnexportedFreeFunc"],
		templateData{
			name:     "NoError",
			template: "FreeFunction",
			Element:  lowerWord(element),
		},
	)
	cases["UnexportedReceiverExportedFunc"] = append(
		cases["UnexportedReceiverExportedFunc"],
		templateData{
			name:         "Error",
			template:     "ReceiverFunction",
			Element:      capWord(element),
			ElementError: true,
			Receiver:     lowerWord(receiver),
		},
	)
	cases["UnexportedReceiverPtrExportedFunc"] = append(
		cases["UnexportedReceiverPtrExportedFunc"],
		templateData{
			name:         "Error",
			template:     "ReceiverPtrFunction",
			Element:      capWord(element),
			ElementError: true,
			Receiver:     lowerWord(receiver),
		},
	)

	for k, v := range genInterfaceFuncCases(
		funcCases,
		capWord(iface),
		exported,
		exported,
	) {
		cases[k] = append(cases[k], v...)
	}

	for k, v := range genInterfaceFuncCases(
		[]templateData{
			{
				name:         "Error",
				Element:      capWord(element),
				ElementError: true,
			},
		},
		lowerWord(iface),
		unexported,
		exported,
	) {
		cases[k] = append(cases[k], v...)
	}

	for _, flags := range flagSets {
		flags := flags

		s.T().Run(flagsToTestName(flags), func(t1 *testing.T) {
			t1.Parallel()

			for name, caseList := range cases {
				name := name
				caseList := caseList

				t1.Run(name, func(t2 *testing.T) {
					t2.Parallel()

					for _, test := range caseList {
						test := test

						t2.Run(test.name, func(t *testing.T) {
							t.Parallel()
							executeCommentMimic(
								t,
								fileTmpl,
								test.template,
								test,
								flags,
							)
						})
					}
				})
			}
		})
	}
}

func (s *CommentMimicSuite) TestCommentExportedEmptyInterfaces() {
	const iface = "iface"

	var (
		flags = map[string]bool{
			"comment-exported":     false,
			"comment-all-exported": false,
			"comment-interfaces":   true,
		}

		ifaceCases = []templateData{
			{
				name: "NoError",
				FirstWord: commentData{
					Type: testdata.InlineComment,
					Text: capWord(iface),
				},
				Element: capWord(iface),
			},
			{
				name: "MimicError",
				FirstWord: commentData{
					Type: testdata.InlineComment,
					Text: "foo",
				},
				CommentError: true,
				Element:      capWord(iface),
			},
			{
				name:         "Error",
				Element:      capWord(iface),
				ElementError: true,
			},
		}
	)

	cases := genEmptyInterfaceCases(ifaceCases, exported)
	cases["UnexportedInterface"] = append(
		cases["UnexportedInterface"],
		templateData{
			name:     "NoError",
			template: "Interface",
			Element:  lowerWord(iface),
		},
	)
	cases["UnexportedInterface"] = append(
		cases["UnexportedInterface"],
		templateData{
			name:     "NoError",
			template: "BlockInterface",
			Element:  lowerWord(iface),
		},
	)
	cases["ManyBlockUnexportedInterface"] = append(
		cases["ManyBlockUnexportedInterface"],
		templateData{
			name:      "NoError",
			template:  "BlockInterface",
			Element:   lowerWord(iface),
			Confusing: true,
		},
	)
	cases["OneBlockExportedInterface"] = append(
		cases["OneBlockExportedInterface"],
		templateData{
			name:         "ErrorCommentedBlock",
			template:     "BlockInterface",
			Element:      capWord(iface),
			ElementError: true,
			InterfaceBlockFirstWord: commentData{
				Type: testdata.InlineComment,
				Text: capWord(iface),
			},
		},
	)

	for name, caseList := range cases {
		name := name
		caseList := caseList

		s.T().Run(name, func(t1 *testing.T) {
			t1.Parallel()

			for _, test := range caseList {
				test := test

				t1.Run(test.name, func(t *testing.T) {
					t.Parallel()
					executeCommentMimic(
						t,
						fileTmpl,
						test.template,
						test,
						flags,
					)
				})
			}
		})
	}
}
