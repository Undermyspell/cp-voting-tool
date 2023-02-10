package roles

type Role string

const (
	Admin        Role = "admin"
	SessionAdmin Role = "session_admin"
	Contributor  Role = "contributor"
)
