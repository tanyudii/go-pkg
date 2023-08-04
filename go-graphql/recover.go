package go_graphql

import (
	"context"
	"errors"
	"fmt"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"runtime/debug"
)

func Recover(_ context.Context, err interface{}) error {
	if _, ok := err.(error); !ok {
		err = errors.New("unexpected error happened")
	}
	fmt.Printf("panic recovered: %v; stacktrace: %s\n", err, string(debug.Stack()))
	return gqlerror.Errorf("internal system error")
}
