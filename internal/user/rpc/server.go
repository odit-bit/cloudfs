package rpc

import (
	"context"

	"github.com/odit-bit/cloudfs/internal/user"
	"github.com/odit-bit/cloudfs/internal/user/userpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ userpb.AuthServiceServer = (*AuthService)(nil)

type AuthService struct {
	accounts *user.Users
	logger   *logrus.Logger
	userpb.UnimplementedAuthServiceServer
}

func NewAuthService(accounts *user.Users, logger *logrus.Logger) *AuthService {
	return &AuthService{accounts: accounts, logger: logger}
}

// BasicAuth implements userpb.AuthServiceServer.
func (a *AuthService) BasicAuth(ctx context.Context, req *userpb.BasicAuthRequest) (*userpb.BasicAuthResponse, error) {
	res, err := a.accounts.BasicAuth(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return &userpb.BasicAuthResponse{
		Token:  res.Token.Key,
		UserID: res.ID,
	}, nil
}

// Register implements userpb.AuthServiceServer.
func (a *AuthService) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	acc, err := a.accounts.Register(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return &userpb.RegisterResponse{
		UserID: acc.ID.String(),
	}, nil
}

// TokenAuth implements userpb.AuthServiceServer.
func (a *AuthService) TokenAuth(ctx context.Context, req *userpb.TokenAuthRequest) (*userpb.TokenAuthResponse, error) {
	res, err := a.accounts.TokenAuth(ctx, req.Token)
	if err != nil {
		return nil, err
	}

	return &userpb.TokenAuthResponse{
		UserID:     res.UserID,
		ValidUntil: timestamppb.New(res.ValidUntil()),
	}, nil
}

// mustEmbedUnimplementedAuthServiceServer implements userpb.AuthServiceServer.
func (a *AuthService) mustEmbedUnimplementedAuthServiceServer() {
	panic("unimplemented")
}
