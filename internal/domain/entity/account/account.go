package account

import "github.com/cool9850311/StreamPlatformLite-Core/pkg/role"

type Account struct {
	ID       uint      `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Role     role.Role `json:"role"`
}
