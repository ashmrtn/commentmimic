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
struct declarations -- should be ignored in all cases.
*/

// structUnexportedCorrectComment has a correctly formatted comment.
type structUnexportedCorrectComment struct{
}

// This structUnexportedWrongComment has an incorrectly formatted comment.
type structUnexportedWrongComment struct {
}

// StructExportedCorrectComment has a correctly formatted comment.
type StructExportedCorrectComment struct{
}

// This StructExportedWrongComment has an incorrectly formatted comment.
type StructExportedWrongComment struct {
}

type (
  // structUnexportedCorrectComment has a correctly formatted comment.
  blockStructUnexportedCorrectComment struct{
  }

  // This structUnexportedWrongComment has an incorrectly formatted comment.
  blockStructUnexportedWrongComment struct {
  }

  // StructExportedCorrectComment has a correctly formatted comment.
  BlockStructExportedCorrectComment struct{
  }

  // This StructExportedWrongComment has an incorrectly formatted comment.
  BlockStructExportedWrongComment struct {
  }
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

	EmptyComments = `package a

  //  // want "first word of comment for element 'a' should be 'a' not ''"
  func a() bool {
    return false
  }

  type b struct {}

  //  // want "first word of comment for element 'c' should be 'c' not ''"
  func (ab b) c() bool {
    return false
  }

  //  // want "first word of comment for element 'd' should be 'd' not ''"
  func (ab *b) d() bool {
    return false
  }

  type e interface {
    //  // want "first word of comment for element 'f' should be 'f' not ''"
    f() bool
  }
`
)
