package go_tex

import (
	"context"
	"google.golang.org/grpc/metadata"
)

type ContextMD metadata.MD

func FromIncoming(ctx context.Context) ContextMD {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		return ContextMD(md.Copy())
	}
	return ContextMD(metadata.Pairs())
}

func FromOutgoing(ctx context.Context) ContextMD {
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		return ContextMD(md)
	}
	return ContextMD(metadata.Pairs())
}

func (c ContextMD) ToOutgoing(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.MD(c))
}

func (c ContextMD) ToIncoming(ctx context.Context) context.Context {
	return metadata.NewIncomingContext(ctx, metadata.MD(c))
}

func (c ContextMD) Get(key string) string {
	if v, ok := c[key]; ok {
		return v[0]
	}
	return ""
}

func (c ContextMD) Delete(key string) ContextMD {
	delete(c, key)
	return c
}

func (c ContextMD) Set(key string, val string) ContextMD {
	c[key] = []string{val}
	return c
}

func (c ContextMD) Add(key string, val string) ContextMD {
	c[key] = append(c[key], val)
	return c
}
