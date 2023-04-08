package testdata

const (
	OutOfScopePatterns = `package a

/*
Var declarations -- should be ignored in all cases.
*/

// varUnexportedCorrectComment has a correctly formatted comment.
var varUnexportedCorrectComment = 43
// This varUnexportedWrongComment has an incorrectly formatted comment.
var varUnexportedWrongComment = 43
// VarExportedCorrectComment has a correctly formatted comment.
var VarExportedCorrectComment = 43
// This VarExportedWrongComment has an incorrectly formatted comment.
var VarExportedWrongComment = 43

var (
  // blockVarUnexportedCorrectComment has a correctly formatted comment.
  blockVarUnexportedCorrectComment = 43
  // This blockVarUnexportedWrongComment has an incorrectly formatted comment.
  blockVarUnexportedWrongComment = 43
  // BlockVarExportedCorrectComment has a correctly formatted comment.
  BlockVarExportedCorrectComment = 43
  // This BlockVarExportedWrongComment has an incorrectly formatted comment.
  BlockVarExportedWrongComment = 43
)

/*
Const declarations -- should be ignored in all cases.
*/

// constUnexportedCorrectComment has a correctly formatted comment.
const constUnexportedCorrectComment = 43
// This constUnexportedWrongComment has an incorrectly formatted comment.
const constUnexportedWrongComment = 43
// ConstExportedCorrectComment has a correctly formatted comment.
const ConstExportedCorrectComment = 43
// This ConstExportedWrongComment has an incorrectly formatted comment.
const ConstExportedWrongComment = 43

const (
  // blockConstUnexportedCorrectComment has a correctly formatted comment.
  blockConstUnexportedCorrectComment = 43
  // This blockConstUnexportedWrongComment has an incorrectly formatted comment.
  blockConstUnexportedWrongComment = 43
  // BlockConstExportedCorrectComment has a correctly formatted comment.
  BlockConstExportedCorrectComment = 43
  // This BlockConstExportedWrongComment has an incorrectly formatted comment.
  BlockConstExportedWrongComment = 43
)

/*
Type equivalences -- should be ignored in all cases.
*/

type (
  // equivalenceUnexportedCorrectComment has a correctly formatted comment.
  blockEquivalenceUnexportedCorrectComment int

  // This equivalenceUnexportedWrongComment has an incorrectly formatted comment.
  blockEquivalenceUnexportedWrongComment int

  // EquivalenceExportedCorrectComment has a correctly formatted comment.
  BlockEquivalenceExportedCorrectComment int

  // This EquivalenceExportedWrongComment has an incorrectly formatted comment.
  BlockEquivalenceExportedWrongComment int
)
`

	ExtraWhitespace = `package a

// 	Element0 has a comment.
func Element0() bool {
  return false
}

/*
	Element1 has a comment.
*/
func Element1() bool {
  return false
}

/*
Element2
has a comment.
*/
func Element2() bool {
  return false
}

/*
Element3	has a comment.
*/
func Element3() bool {
  return false
}
`

	EmptyComments = `package a

//nolint:commentmimic
func ignoreMachineReadable() bool {
  return false
}

//nolint:commentmimic // want "first word of comment is 'This' instead of 'CommentMismatch'"
// This function has a comment.
func CommentMismatch() bool {
  return false
}

//nolint:commentmimic
//
func EmptyComment() bool { // want "empty comment on 'EmptyComment'"
  return false
}

//nolint:commentmimic
//	
func EmptyComment2() bool { // want "empty comment on 'EmptyComment2'"
  return false
}

//nolint:commentmimic
/*
*/
func EmptyComment3() bool { // want "empty comment on 'EmptyComment3'"
  return false
}

//nolint:commentmimic
/*
   
*/
func EmptyComment4() bool { // want "empty comment on 'EmptyComment4'"
  return false
}
`

	MachineReadableExported = `package a

//nolint:commentmimic
func FreeFunc() bool { // want "exported element 'FreeFunc' should be commented"
  return false
}

type testStruct struct {}

//nolint:commentmimic
func (t testStruct) PrivateReceiver() bool { // want "exported element 'PrivateReceiver' should be commented"
  return false
}

//nolint:commentmimic
func (t *testStruct) PrivatePtrReceiver() bool { // want "exported element 'PrivatePtrReceiver' should be commented"
  return false
}

//nolint:commentmimic
type TestInterface interface { // want "exported element 'TestInterface' should be commented"
  ExportedInterfaceFunc() bool // want "exported element 'ExportedInterfaceFunc' should be commented"
}

//nolint:commentmimic
type testIface interface {
  UnexportedInterfaceFunc() bool // want "exported element 'UnexportedInterfaceFunc' should be commented"
}
`

	SkipTestComments = `package a_test

  import (
    "testing"
  )

  type BenchmarkInterface interface {} // want "exported element 'BenchmarkInterface' should be commented"

  type ExampleInterface interface {} // want "exported element 'ExampleInterface' should be commented"

  type FuzzInterface interface {} // want "exported element 'FuzzInterface' should be commented"

  type TestInterface interface {} // want "exported element 'TestInterface' should be commented"

  func BenchmarkDoesntNeedComment(b *testing.B) {}

  func ExportedElement() bool { // want "exported element 'ExportedElement' should be commented"
    return false
  }

  func ExampleDoesntNeedComment() {}

  func FuzzDoesntNeedComment(f *testing.F) {}

  func TestDoesntNeedComment(t *testing.T) {}
`
)
