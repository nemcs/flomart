package identity

type contextKey string

const (
	CtxKey         contextKey = "user"
	CtxUserIDKey   contextKey = "userID"
	CtxUserRoleKey contextKey = "userRole"
)
