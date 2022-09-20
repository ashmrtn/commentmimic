# CommentMimic

CommentMimic is a golang linter that enforces starting comments for interfaces
and functions with the name of the interface or function. The case of the first
word must match the case of the element being commented. There are also a set of
flags for requiring comments on exported interfaces or functions.

## Why do I need this?
Some go best-practices, like
[effective-go](https://github.com/golovers/effective-go#comment-sentences), say
comments should be written in full sentences and should start with the thing
being commented. Following the best-practices can make reading documentation
easier, but it makes upkeep a lot harder, especially when code is refactored or
a well-meaning colleague asks you to rename an interface or function to better
match what it does. When that happens, it's easy to forget to update the comment
to match the new interface or function name. In the end, the code ends up with
some elements having comments that match their names and others are clearly the
left-overs of previous itertions of the code.

CommentMimic was created to help developers by scanning the code and comments
and outputting notices when the first word of a comment for an interface or
function does not match the item being commented. It can be run from the command
line like other go tools, making it easy to integrate into your workflow.

CommentMimic aims to be unobtrusive by only warning about things that already
have comments attached to them. Some other linters expect comments on all
exported functions and only output warnings about unmatched comment first words
and item names on exported functions. While that's a step in the right
direction, those linters can be overwhelming, because not many codebases have
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
CommentMimic optionally enforces comments on all exported interfaces and
functions depending flags passed to it.

`--comment-exported` requires comments on all functions that have exported
interfaces or exported receivers (i.e. the struct the function belongs to is
exported as well).

`--comment-all-exported` requires comments on all functions regardless of
whether their receiver is exported.

`--comment-interfaces` requires comments on all exported interfaces.

## Limitations
CommentMimic has the following limitations and oddities:

* ignores leading whitespace in comments
* only checks interfaces and functions

Leading whitespace is ignored because it doesn't play well with comments
starting with `/* text starts after a space...`. Go expects comments using
`/* */` to have the `/*` and `*/` on lines of their own. To be a little more
flexible, CommentMimic can handle when they share a line with comment text, but
the linter can no longer check if there's leading whitepace.

Currently CommentMimic only checks interface and function comments. This may be
expanded later on, but the [official pages](https://tip.golang.org/doc/comment)
on golang doc comments don't explicitly state standards for comment formats. The
struct comments written on that page consistently start with `A` or `An` but
comments on `var` and `const` decalarations show much more variability in their
format. Eventually CommentMimic may be expanded to cover struct comments as
well, but support for others seems unlikely.
