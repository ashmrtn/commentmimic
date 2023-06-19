package testdata

const FileGenTmpls = `{{define "ReceiverFunction"}}// {{ .Receiver }}
type {{ .Receiver }} struct{
}

{{template "maybeCommentWithError" .}}
func (r {{ .Receiver -}}) {{template "freeFuncDef" .}} { {{template "maybeElementError" .}}
  return false
}{{end}}

{{define "ReceiverPtrFunction"}}// {{ .Receiver }}
type {{ .Receiver }} struct{
}

{{template "maybeCommentWithError" .}}
func (r *{{- .Receiver -}}) {{template "freeFuncDef" .}} { {{template "maybeElementError" .}}
  return false
}{{end}}

{{define "FreeFunction"}}{{template "maybeCommentWithError" .}}
func {{template "freeFuncDef" .}} { {{template "maybeElementError" .}}
  return false
}{{end}}

{{define "Struct"}}{{template "maybeCommentWithError" .}}
type {{ .Element }} struct {} {{template "maybeElementError" .}}{{end}}

{{define "BlockStruct"}}{{template "maybeCommentWithError" }}
type (
{{if .Confusing}}
  testConfusingStruct1 struct {}
{{end}}

  {{template "maybeCommentWithError" .}}
  {{ .Element }} struct {} {{template "maybeElementError" .}}

{{if .Confusing}}
  testConfusingStruct2 struct {}
{{end}}
){{end}}

{{define "Interface"}}{{template "maybeCommentWithError" .}}
type {{template "interfaceInner" .}}{{end}}

{{define "BlockInterface"}}{{template "maybeComment" .InterfaceBlockFirstWord}}
type (
{{if .Confusing}}
  testConfusingInterface1 interface{}
{{end}}

  {{template "maybeCommentWithError" .}}
  {{template "interfaceInner" .}}

{{if .Confusing}}
  testConfusingInterface2 interface{}
{{end}}
){{end}}

{{define "freeFuncDef"}}{{ .Element -}}() bool{{end}}

{{define "interfaceInner"}}{{ .Element }} interface { {{template "maybeElementError" .}}
{{if .InterfaceFunc}}{{if .InterfaceFunc.Confusing}}  testConfusingInterfaceFunc1() bool{{end}}
  {{template "maybeCommentWithError" .InterfaceFunc}}
  {{template "freeFuncDef" .InterfaceFunc}} {{template "maybeElementError" .InterfaceFunc}}
{{if .InterfaceFunc.Confusing}}  testConfusingInterfaceFunc2() bool{{end}}{{end}}
}{{end}}

{{define "regularComment"}}// {{if .Text0}}{{ .Text0 }} {{end}}{{ .Text }} has a comment.{{if .Multiline}}
// This is another line of comment text.{{end}}{{end}}

{{define "blockCommentInline"}}/* {{if .Text0}}{{ .Text0 }} {{end}}{{ .Text }} has a comment.{{if .Multiline}}
This is another line of comment text.{{end}} */{{end}}

{{define "blockCommentMultiline"}}/*
{{if .Text0}}{{ .Text0 }} {{end}}{{ .Text }} has a comment.{{if .Multiline}}
This is another line of comment text.{{end}}
*/{{end}}

{{define "maybeCommentWithError"}}{{if .FirstWord.Text}}{{if eq .FirstWord.Type.String "InlineComment"}}{{template "regularCommentWithError" .}}{{else if eq .FirstWord.Type.String "BlockInlineComment"}}{{template "blockCommentInlineWithError" .}}{{else}}{{template "blockCommentMultilineWithError" .}}{{end}}{{end}}{{end}}

{{define "maybeComment"}}{{if .Text}}{{if eq .Type.String "InlineComment"}}{{template "regularComment" .}}{{else if eq .Type.String "BlockInlineComment"}}{{template "blockCommentInline" .}}{{else}}{{template "blockCommentMultiline" .}}{{end}}{{end}}{{end}}

{{define "regularCommentWithError"}}// {{if .FirstWord.Text0}}{{ .FirstWord.Text0 }} {{end}}{{ .FirstWord.Text }} has a comment.{{if .CommentError}} {{template "commentMismatch" .}}{{end}}{{if .FirstWord.Multiline}}
// This is another line of comment text.{{end}}{{end}}

{{define "blockCommentInlineWithError"}}/*{{if .FirstWord.Text0}}{{ .FirstWord.Text0 }} {{end}}{{ .FirstWord.Text}} has a comment.{{if .FirstWord.Multiline}}
This is another line of comment text.{{end}}{{if .CommentError}} {{template "commentMismatch" .}}{{end}} */{{end}}

{{define "blockCommentMultilineWithError"}}/*
{{if .FirstWord.Text0}}{{ .FirstWord.Text0 }} {{end}}{{ .FirstWord.Text }} has a comment.{{if .FirstWord.Multiline}}
This is another line of comment text.{{end}}{{if .CommentError}} {{template "commentMismatch" .}}{{end}}
*/{{end}}

{{define "commentMismatch"}}// want "first word of comment is '{{if .FirstWord.Text0 }}{{- .FirstWord.Text0 -}}{{else}}{{- .FirstWord.Text -}}{{end}}' instead of '{{- .Element -}}'" {{end}}

{{define "commentMissing"}}// want "exported element '{{- .Element -}}' should be commented"{{end}}

{{define "maybeElementError"}}{{if .ElementError}} {{template "commentMissing" .}}{{end}}{{end}}`
