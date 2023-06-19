package testdata

//go:generate stringer -type=CommentType
type CommentType int

const (
	InlineComment CommentType = iota
	BlockInlineComment
	BlockMultilineComment
)
