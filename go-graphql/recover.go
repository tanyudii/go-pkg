package go_graphql

import (
	"context"
	"errors"
	"github.com/vektah/gqlparser/v2/gqlerror"
	gologger "pkg.tanyudii.me/go-pkg/go-logger"
	"runtime/debug"
)

func Recover(_ context.Context, err interface{}) error {
	if _, ok := err.(error); !ok {
		err = errors.New("unexpected error happened")
	}
	gologger.WithField("stacktrace", string(debug.Stack())).Errorf("panic recovered: %v", err)
	return gqlerror.Errorf("internal system error")
}
