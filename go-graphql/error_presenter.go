package go_graphql

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
	goerr "pkg.tanyudii.me/go-pkg/go-err"
)

func ErrorPresenter(ctx context.Context, e error) *gqlerror.Error {
	err := graphql.DefaultErrorPresenter(ctx, e)
	if err.Extensions == nil {
		err.Extensions = make(map[string]interface{})
	}
	realErr := err.Unwrap()
	customErr, isCustomErr := realErr.(goerr.CustomError)
	if isCustomErr {
		customExtension := goerr.CustomErrorToMapInterface(customErr)
		badRequestErr, isBadRequestErr := realErr.(*goerr.BadRequestError)
		if isBadRequestErr {
			if fields := badRequestErr.GetFields(); len(fields) != 0 {
				customExtension["fields"] = fields
			}
		}
		err.Extensions["errors"] = customExtension
	}

	return err
}
