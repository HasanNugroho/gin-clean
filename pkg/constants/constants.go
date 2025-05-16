package constants

type Role string

const (
	ROLE_USER  Role = "user"
	ROLE_ADMIN Role = "admin"
)

func (r *Role) IsValidRole() bool {
	switch *r {
	case ROLE_USER, ROLE_ADMIN:
		return true
	}
	return false
}
