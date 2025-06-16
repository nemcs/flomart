package identity

type contextKey string

const (
	CtxUserIDKey   contextKey = "userID"
	CtxUserRoleKey contextKey = "userRole"
)
