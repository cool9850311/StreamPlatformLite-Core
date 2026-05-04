package repository

import "github.com/cool9850311/StreamPlatformLite-Core/internal/domain/entity/account"

type AccountRepository interface {
	Create(account account.Account) error
	GetAll() ([]account.Account, error)
	GetByUsername(username string) (*account.Account, error)
	Update(account account.Account) error
	Delete(username string) error
}
