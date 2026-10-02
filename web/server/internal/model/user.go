package model

// User is the request-scoped identity loaded by middleware/auth.go.
// Role comes from profiles.role in the database, never from token metadata.
type User struct {
	ID     string
	Role   string
	Status string
}

// IsAdmin reports the admin role.
func (u *User) IsAdmin() bool { return u != nil && u.Role == "admin" }

// IsSuspended reports the suspended status (403 on every authenticated route).
func (u *User) IsSuspended() bool { return u != nil && u.Status == "suspended" }
