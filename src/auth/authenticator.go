package auth

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
)

//go:generate mockgen -source=authenticator.go -destination=mock/authenticator.go -package=mock

type Authenticator interface {
	AuthenticateJWT(ctx context.Context, token string) (*entities.UserInfo, error)
	AuthenticateAPIKey(ctx context.Context, apiKey []byte) (*entities.UserInfo, error)
}
