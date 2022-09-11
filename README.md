# CommentMimic

CommentMimic is a golang linter that enforces starting comments on interfaces
and functions with the name of the interface or function. The case of the first
word must match the case of the element being commented. There are also a set of
flags requiring comments on exported interfaces or functions.

## Installing
CommentMimic is provided as a go module and can be installed by running
`go install github.com/ashmrtn/commentmimic@latest`

## Running
CommentMimic uses a similar front-end to `go vet`. You can run CommentMimic by
executing `commentmimic <flags> <packages>` after installation.
