package testdata

const FileGenTmpls = `{{define "ReceiverFunction"}}type {{ .Receiver }} struct{
}

{{template "maybeCommentWithError" .}}
func (r {{ .Receiver -}}) {{template "freeFuncDef" .}} { {{template "maybeElementError" .}}
  return false
}{{end}}

{{define "ReceiverPtrFunction"}}type {{ .Receiver }} struct{
}

{{template "maybeCommentWithError" .}}
func (r *{{- .Receiver -}}) {{template "freeFuncDef" .}} { {{template "maybeElementError" .}}
  return false
}{{end}}

{{define "FreeFunction"}}{{template "maybeCommentWithError" .}}
func {{template "freeFuncDef" .}} { {{template "maybeElementError" .}}
  return false
}{{end}}

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

{{define "comment"}}// {{ . }} has a comment.{{end}}

{{define "maybeCommentWithError"}}{{if len .FirstWord}}{{template "comment" .FirstWord}}{{if .CommentError}} {{template "commentMismatch" .}}{{end}}{{end}}{{end}}

{{define "maybeComment"}}{{if len .}}{{template "comment" .}}{{end}}{{end}}

{{define "commentMismatch"}}// want "first word of comment for element '{{- .Element -}}' should be '{{- .Element -}}' not '{{- .FirstWord -}}'"{{end}}

{{define "commentMissing"}}// want "exported element '{{- .Element -}}' should be commented"{{end}}

{{define "maybeElementError"}}{{if .ElementError}} {{template "commentMissing" .}}{{end}}{{end}}`
