package exports

import (
	"go/ast"
	"strings"
)

type Comment struct {
	Doc  []string // Documentation comments associated with the entity.
	Line string   // Associated line comment.
}

// ParseComment parses the documentation and line comments from AST comment groups.
func ParseComment(doc *ast.CommentGroup, lineComment *ast.CommentGroup) Comment {
	c := Comment{}

	if lineComment != nil && len(lineComment.List) > 0 {
		c.Line = strings.TrimSpace(lineComment.Text())
	}

	if doc != nil && len(doc.List) > 0 {
		c.Doc = strings.Split(strings.TrimSpace(doc.Text()), "\n")
	}

	return c
}
