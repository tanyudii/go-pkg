package go_tex

import (
	"context"
	"fmt"
	goerr "pkg.tanyudii.me/go-pkg/go-err"
	"strings"
)

var (
	ErrContextNotFound        = goerr.NewInternalServerErrorWithName("[ERROR]: Context not found", "GOTEX_NOT_FOUND")
	ErrUnauthorized           = goerr.NewUnauthorizedErrorWithName("[ERROR]: Unauthorized", "GOTEX_UNAUTHORIZED")
	ErrUnauthorizedUserType   = fmt.Errorf("%w: user type", ErrUnauthorized)
	ErrUnauthorizedPermission = fmt.Errorf("%w: permission", ErrUnauthorized)
	ErrUnauthorizedScope      = fmt.Errorf("%w: scope", ErrUnauthorized)
)

const (
	ContextKey                           = "gotex"
	RequestHeaderKeyUserID               = "UserID"
	RequestHeaderKeyUserName             = "UserName"
	RequestHeaderKeyUserEmail            = "UserEmail"
	RequestHeaderKeyUserType             = "UserType"
	RequestHeaderKeyCompanyID            = "CompanyID"
	RequestHeaderKeyCompanyName          = "CompanyName"
	RequestHeaderKeyPermissions          = "Permissions"
	RequestHeaderKeyScopes               = "Scopes"
	RequestHeaderKeyClientID             = "ClientID"
	RequestHeaderKeyClientName           = "ClientName"
	RequestHeaderKeyInternalCallPassword = "InternalCallPassword"
	RequestHeaderKeyAuthorization        = "Authorization"
	RequestHeaderKeyRequestID            = "RequestID"
	RequestHeaderKeyAcceptLanguage       = "Accept-Language"
	RequestHeaderKeyXForwardedFor        = "X-Forwarded-For"
	RequestHeaderUserAgent               = "User-Agent"

	ScopeSeparator      = " "
	PermissionSeparator = ";"
)

type Gotex struct {
	UserID               string
	UserName             string
	UserEmail            string
	UserType             string
	CompanyID            string
	CompanyName          string
	Permissions          string
	ClientID             string
	ClientName           string
	Scopes               string
	InternalCallPassword string
	Authorization        string
	RequestID            string
	AcceptLanguage       string
	XForwardedFor        string
	UserAgent            string
}

func NewGotex(md ContextMD) *Gotex {
	return &Gotex{
		UserID:               md.Get(strings.ToLower(RequestHeaderKeyUserID)),
		UserName:             md.Get(strings.ToLower(RequestHeaderKeyUserName)),
		UserEmail:            md.Get(strings.ToLower(RequestHeaderKeyUserEmail)),
		UserType:             md.Get(strings.ToLower(RequestHeaderKeyUserType)),
		CompanyID:            md.Get(strings.ToLower(RequestHeaderKeyCompanyID)),
		CompanyName:          md.Get(strings.ToLower(RequestHeaderKeyCompanyName)),
		Permissions:          md.Get(strings.ToLower(RequestHeaderKeyPermissions)),
		ClientID:             md.Get(strings.ToLower(RequestHeaderKeyClientID)),
		ClientName:           md.Get(strings.ToLower(RequestHeaderKeyClientName)),
		Scopes:               md.Get(strings.ToLower(RequestHeaderKeyScopes)),
		InternalCallPassword: md.Get(strings.ToLower(RequestHeaderKeyInternalCallPassword)),
		Authorization:        md.Get(strings.ToLower(RequestHeaderKeyAuthorization)),
		RequestID:            md.Get(strings.ToLower(RequestHeaderKeyRequestID)),
		AcceptLanguage:       md.Get(strings.ToLower(RequestHeaderKeyAcceptLanguage)),
		XForwardedFor:        md.Get(strings.ToLower(RequestHeaderKeyXForwardedFor)),
		UserAgent:            md.Get(strings.ToLower(RequestHeaderUserAgent)),
	}
}

func (c *Gotex) ToContextMD(ctx context.Context) context.Context {
	md := FromIncoming(ctx)
	md.Set(strings.ToLower(RequestHeaderKeyUserID), c.UserID)
	md.Set(strings.ToLower(RequestHeaderKeyUserName), c.UserName)
	md.Set(strings.ToLower(RequestHeaderKeyUserEmail), c.UserEmail)
	md.Set(strings.ToLower(RequestHeaderKeyUserType), c.UserType)
	md.Set(strings.ToLower(RequestHeaderKeyCompanyID), c.CompanyID)
	md.Set(strings.ToLower(RequestHeaderKeyCompanyName), c.CompanyName)
	md.Set(strings.ToLower(RequestHeaderKeyPermissions), c.Permissions)
	md.Set(strings.ToLower(RequestHeaderKeyClientID), c.ClientID)
	md.Set(strings.ToLower(RequestHeaderKeyClientName), c.ClientName)
	md.Set(strings.ToLower(RequestHeaderKeyScopes), c.Scopes)
	md.Set(strings.ToLower(RequestHeaderKeyInternalCallPassword), c.InternalCallPassword)
	md.Set(strings.ToLower(RequestHeaderKeyAuthorization), c.Authorization)
	md.Set(strings.ToLower(RequestHeaderKeyRequestID), c.RequestID)
	md.Set(strings.ToLower(RequestHeaderKeyAcceptLanguage), c.AcceptLanguage)
	md.Set(strings.ToLower(RequestHeaderKeyXForwardedFor), c.XForwardedFor)
	md.Set(strings.ToLower(RequestHeaderUserAgent), c.UserAgent)
	ctx = NewContext(ctx, c)
	return md.ToIncoming(ctx)
}

func (c *Gotex) HasPermission(codes []string) (bool, error) {
	if len(codes) == 0 {
		return false, nil
	}
	permissions := splitString(c.Permissions, PermissionSeparator)
	if len(permissions) != 0 {
		mapCode := sliceStringsToMap(codes)
		for _, p := range permissions {
			if mapCode[p] {
				return true, nil
			}
		}
	}
	return false, ErrUnauthorizedPermission
}

func (c *Gotex) HasScope(codes []string) (bool, error) {
	if len(codes) == 0 {
		return false, nil
	}
	scopes := splitString(c.Scopes, ScopeSeparator)
	if len(scopes) != 0 {
		mapCode := sliceStringsToMap(codes)
		for _, s := range scopes {
			if mapCode[s] {
				return true, nil
			}
		}
	}
	return false, ErrUnauthorizedScope
}

func (c *Gotex) HasUserType(codes []string) (bool, error) {
	if len(codes) == 0 {
		return false, nil
	}
	for _, code := range codes {
		if code == c.UserType {
			return true, nil
		}
	}
	return false, ErrUnauthorizedUserType
}

func (c *Gotex) HasUserTypeByMapCode(codes map[string]bool) (bool, error) {
	if len(codes) == 0 {
		return false, nil
	}
	if c.UserType != "" && codes[c.UserType] {
		return true, nil
	}
	return false, ErrUnauthorizedUserType
}

func (c *Gotex) ToRequestHeaders() map[string]string {
	return map[string]string{
		RequestHeaderKeyUserID:               c.UserID,
		RequestHeaderKeyUserName:             c.UserName,
		RequestHeaderKeyUserEmail:            c.UserEmail,
		RequestHeaderKeyUserType:             c.UserType,
		RequestHeaderKeyCompanyID:            c.CompanyID,
		RequestHeaderKeyCompanyName:          c.CompanyName,
		RequestHeaderKeyPermissions:          c.Permissions,
		RequestHeaderKeyClientID:             c.ClientID,
		RequestHeaderKeyClientName:           c.ClientName,
		RequestHeaderKeyScopes:               c.Scopes,
		RequestHeaderKeyInternalCallPassword: c.InternalCallPassword,
		RequestHeaderKeyAuthorization:        c.Authorization,
		RequestHeaderKeyRequestID:            c.RequestID,
		RequestHeaderKeyAcceptLanguage:       c.AcceptLanguage,
		RequestHeaderKeyXForwardedFor:        c.XForwardedFor,
		RequestHeaderUserAgent:               c.UserAgent,
	}
}

func NewContext(ctx context.Context, eCtx *Gotex) context.Context {
	if eCtx == nil {
		return ctx
	}
	return context.WithValue(ctx, ContextKey, eCtx)
}

func FromContext(ctx context.Context) (*Gotex, bool) {
	rc, ok := ctx.Value(ContextKey).(*Gotex)
	return rc, ok
}

func FromContextWithErr(ctx context.Context) (*Gotex, error) {
	val, ok := FromContext(ctx)
	if !ok {
		return nil, ErrContextNotFound
	}
	return val, nil
}

func ParseToGrpcCtx(ctx context.Context, pwd ...string) context.Context {
	if r, ok := FromContext(ctx); ok {
		newCtx := FromIncoming(ctx)
		internalPwd := append([]string{r.InternalCallPassword}, pwd...)
		newCtx.Add(strings.ToLower(RequestHeaderKeyUserID), r.UserID)
		newCtx.Add(strings.ToLower(RequestHeaderKeyUserName), r.UserName)
		newCtx.Add(strings.ToLower(RequestHeaderKeyUserEmail), r.UserEmail)
		newCtx.Add(strings.ToLower(RequestHeaderKeyUserType), r.UserType)
		newCtx.Add(strings.ToLower(RequestHeaderKeyCompanyID), r.CompanyID)
		newCtx.Add(strings.ToLower(RequestHeaderKeyCompanyName), r.CompanyName)
		newCtx.Add(strings.ToLower(RequestHeaderKeyPermissions), r.Permissions)
		newCtx.Add(strings.ToLower(RequestHeaderKeyClientID), r.ClientID)
		newCtx.Add(strings.ToLower(RequestHeaderKeyClientName), r.ClientName)
		newCtx.Add(strings.ToLower(RequestHeaderKeyScopes), r.Scopes)
		newCtx.Add(strings.ToLower(RequestHeaderKeyInternalCallPassword), firstOrDefault(internalPwd...))
		newCtx.Add(strings.ToLower(RequestHeaderKeyAuthorization), r.Authorization)
		newCtx.Add(strings.ToLower(RequestHeaderKeyRequestID), r.RequestID)
		newCtx.Add(strings.ToLower(RequestHeaderKeyAcceptLanguage), r.AcceptLanguage)
		newCtx.Add(strings.ToLower(RequestHeaderKeyXForwardedFor), r.XForwardedFor)
		newCtx.Add(strings.ToLower(RequestHeaderUserAgent), r.UserAgent)
		return newCtx.ToOutgoing(ctx)
	}
	return ctx
}
