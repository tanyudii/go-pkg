package go_tex

import (
	"context"
	"strings"
)

func GetUserID(ctx context.Context) (string, error) {
	eCtx, err := FromContextWithErr(ctx)
	if err != nil {
		return "", err
	}
	return eCtx.UserID, nil
}

func GetUserType(ctx context.Context) (string, error) {
	eCtx, err := FromContextWithErr(ctx)
	if err != nil {
		return "", err
	}
	return eCtx.UserType, nil
}

func GetCompanyID(ctx context.Context) (string, error) {
	eCtx, err := FromContextWithErr(ctx)
	if err != nil {
		return "", err
	}
	return eCtx.CompanyID, nil
}

func GetAcceptLanguage(ctx context.Context, defaultVal ...string) string {
	var acceptLang string
	if eCtx, ok := FromContext(ctx); ok {
		acceptLang = eCtx.AcceptLanguage
	}
	if acceptLang == "" && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return acceptLang
}

func DuplicateCtx(ctx context.Context) (context.Context, error) {
	eCtx, err := FromContextWithErr(ctx)
	if err != nil {
		return nil, err
	}
	return NewContext(context.Background(), eCtx), nil
}

func CreateInternalEContextDummy(pwd ...string) *Gotex {
	return &Gotex{
		UserID:               "DummyUserID",
		UserName:             "DummyUserName",
		UserEmail:            "DummyUserEmail",
		UserType:             "DummyUserType",
		CompanyID:            "DummyCompanyID",
		CompanyName:          "DummyCompanyName",
		Permissions:          "DummyPermissions",
		ClientID:             "DummyClientID",
		ClientName:           "DummyClientName",
		Scopes:               "*",
		InternalCallPassword: firstOrDefault(pwd...),
	}
}

func CreateGRPCContextDummy(ctx context.Context) context.Context {
	eCtxDummy := CreateInternalEContextDummy()
	return ParseToGrpcCtx(NewContext(ctx, eCtxDummy))
}

func firstOrDefault[T any](v ...T) T {
	if len(v) != 0 {
		return v[0]
	}
	var dfl T
	return dfl
}

func splitString(s string, sep string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, sep)
}

func sliceStringsToMap(v []string) map[string]bool {
	if len(v) == 0 {
		return nil
	}
	m := make(map[string]bool)
	for _, val := range v {
		m[val] = true
	}
	return m
}
