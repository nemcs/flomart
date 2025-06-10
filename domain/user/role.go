package user

type Role string

const (
	RoleClient  Role = "client"
	RoleCourier Role = "courier"
	RoleAdmin   Role = "admin"
)

func IsValidRole(r Role) bool {
	switch r {
	case RoleClient, RoleCourier, RoleAdmin:
		return true
	default:
		return false
	}
}
