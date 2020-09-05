package gtool

import (
	"context"
	"google.golang.org/grpc/metadata"
)

func GetMetaData(ctx context.Context, key string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", nil
	}

	vals := md.Get(key)
	if len(vals) == 0 {
		return "", nil
	}

	return vals[0], nil
}

func AddMetaData(ctx context.Context, key string, val string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, key, val)
}
