package gtool

import (
	"context"
	"go.elastic.co/apm"
)

// example:
//  span, printCtx := gtool.StartSpan(ctx, "print article", "print.query")
//  defer span.End()
//  ...
// Start a span by transaction.
func StartSpan(ctx context.Context, name, spanType string) (*apm.Span, context.Context) {
	return apm.StartSpan(ctx, name, spanType)
}
