package auth

import (
	"github.com/consensys/quorum-key-manager/src/auth/authenticator/oicd"
)

type Config struct {
	OICD oicd.Config
}
