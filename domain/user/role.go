package user

type Role string

// Безопасно ли хранить их в коде или в env выносить?
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
