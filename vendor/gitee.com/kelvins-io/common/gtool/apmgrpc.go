package gtool

import (
	"context"
	"go.elastic.co/apm"
)

func WithGrpcContext(ctx context.Context) context.Context {
	return WithAPMContext(ctx)
}

func WithAPMContext(ctx context.Context) context.Context {
	return apm.DetachedContext(ctx)
}
