package contextutils

import (
	"context"
	"errors"
	"flomart/internal/identity"
)

func GetUserFromCtx(ctx context.Context) (*identity.UserClaim, error) {
	claims, ok := ctx.Value(identity.CtxKey).(*identity.UserClaim)
	if !ok {
		return nil, errors.New("user not found in context")
	}
	return claims, nil
}
