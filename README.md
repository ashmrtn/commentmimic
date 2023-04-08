# CommentMimic

CommentMimic is a golang linter that enforces starting comments for interfaces,
functions, and structs with the name of the interface, function, or struct. The
case of the first word must match the case of the element being commented. There
are also a set of flags for requiring comments on exported interfaces,
functions, or structs. Struct comments may start with an "A" or "An" followed by
the struct name (in the proper case).

## Why should I use CommentMimic?
Some go best-practices, like
[effective-go](https://github.com/golovers/effective-go#comment-sentences), say
comments should be written in full sentences and should start with the thing
being commented. Following the best-practices can make reading documentation
easier, but it makes upkeep a lot harder, especially when code is refactored or
a well-meaning colleague asks you to rename an interface, function, or struct to
better match what it does. When that happens, it's easy to forget to update the
comment to match the new name. In the end, the code ends up with some elements
having comments that match their names and others are clearly the left-overs of
previous iterations of the code.

CommentMimic was created to help developers by scanning the code and comments
and outputting notices when the first word of a comment for an interface,
function, or struct doesn't match the item being commented. It can be run from
the command line like other go tools, making it easy to integrate into your
workflow.

CommentMimic aims to be unobtrusive by only warning about things that already
have comments attached to them. Some other linters expect comments on all
exported functions and only output warnings about unmatched comment first words
and item names on exported functions. While that's a step in the right
direction, those linters can be overwhelming because not many codebases have
all exported items commented, and only find issues on exported items.

## Installing
CommentMimic is provided as a go module and can be installed by running
`go install github.com/ashmrtn/commentmimic@latest`

## Running
CommentMimic uses a similar front-end to `go vet`. You can run CommentMimic by
executing `commentmimic <flags> <packages>` after installation. You can pass
multiple packages to CommentMimic by using the `...` expression, just like you
would with other tools. For example, to check your whole project with
CommentMimic just run `commentmimic ./...`.

### Flags
CommentMimic optionally enforces comments on all exported interfaces,
functions, and structs depending flags passed to it.

`--comment-exported` requires comments on all functions that have exported
interfaces or exported receivers (i.e. the struct the function belongs to is
exported as well).

`--comment-all-exported` requires comments on all functions regardless of
whether their receiver is exported.

`--comment-interfaces` requires comments on all exported interfaces.

`--comment-structs` requires comments on all exported structs.

## Limitations
CommentMimic has the following limitations and oddities:

* ignores leading whitespace in comments
* doesn't lint comments on consts or vars
* comments on a type-block with a single type definition won't be applied to the
  type defined in the block

Leading whitespace is ignored because it doesn't play well with comments
starting with `/* text starts after a space...`. Go expects comments using
`/* */` to have the `/*` and `*/` on lines of their own. To be a little more
flexible, CommentMimic can handle when they share a line with comment text, but
the linter can no longer check if there's leading whitespace.

Currently CommentMimic only checks interface, function, and struct comments.
This may be expanded later on, but the
[official pages](https://tip.golang.org/doc/comment) on golang doc comments
don't explicitly state standards for comment formats. The `var` and `const`
declarations on that page show lots of variability in their format, making it
hard to enforce a standard.

Golang allows declaring one or more types in a type-block like shown below.
Comments can also be associated with the type-block, in addition to or as a
replacement for comments on the individual type declarations. CommentMimic won't
apply comments on the type-block to the type definition(s) in the block, even if
there's only one of them.

```go
// A someType is a struct used for passing data.
//
// This comment layout will cause a lint error because the comment is on the
// type-block instead of the type declaration.
type (
    someType struct{}
)

type (
    // A someOtherType is a struct used for passing data.
    //
    // This comment layout will not cause a lint error because the comment is
    // associated with the type declaration.
    someOtherType struct{}
)
```
