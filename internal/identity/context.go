package identity

import (
	"context"
	"errors"
)

func GetUserFromCtx(ctx context.Context) (*UserClaim, error) {
	claims, ok := ctx.Value(CtxKey).(*UserClaim)
	if !ok {
		return nil, errors.New("user not found in context")
	}
	return claims, nil
}
