package go_auth

import (
	"context"
)

type TokenService interface {
	TokenInfo(ctx context.Context, jwtToken string) (*TokenInfoResponse, error)
}

type TokenInfoResponse struct {
	TokenInfo  *TokenInfo
	ClientInfo *ClientInfo
	Scope      string
}

type TokenInfo struct {
	UserID         string
	UserSerial     string
	UserName       string
	UserEmail      string
	UserType       string
	CompanyID      string
	CompanySerial  string
	CompanyName    string
	Permissions    []string
	IsInternalCall bool
}

type ClientInfo struct {
	ClientID   string
	ClientName string
}
