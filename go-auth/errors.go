package go_auth

import goerr "pkg.tanyudii.me/go-pkg/go-err"

var (
	ErrUnauthenticated = goerr.NewUnauthenticatedErrorWithName("unauthenticated", "UNAUTHENTICATED")
)
