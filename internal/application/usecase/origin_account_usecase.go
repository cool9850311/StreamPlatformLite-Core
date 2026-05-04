package usecase

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/cool9850311/StreamPlatformLite-Core/internal/application/dto"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/application/dto/config"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/application/interface/jwt"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/application/interface/repository"
	domainAccount "github.com/cool9850311/StreamPlatformLite-Core/internal/domain/entity/account"
	innerErrors "github.com/cool9850311/StreamPlatformLite-Core/internal/domain/entity/errors"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/domain/interface/libarary/bcrypt"
	"github.com/cool9850311/StreamPlatformLite-Core/internal/domain/interface/logger"
	"github.com/cool9850311/StreamPlatformLite-Core/pkg/role"
)

type OriginAccountUseCase struct {
	accountRepo  repository.AccountRepository
	log          logger.Logger
	bcrypt       bcrypt.BcryptGenerator
	config       config.Config
	jwtGenerator jwt.JWTGenerator
}

func NewOriginAccountUseCase(accountRepo repository.AccountRepository, log logger.Logger, bcrypt bcrypt.BcryptGenerator, config config.Config, jwtGenerator jwt.JWTGenerator) *OriginAccountUseCase {
	return &OriginAccountUseCase{
		accountRepo:  accountRepo,
		log:          log,
		bcrypt:       bcrypt,
		config:       config,
		jwtGenerator: jwtGenerator,
	}
}

func (uc *OriginAccountUseCase) Login(ctx context.Context, username, password string) (string, error) {
	acc, err := uc.accountRepo.GetByUsername(username)
	if err != nil {
		uc.log.Error(ctx, "Login failed: "+err.Error())
		return "", innerErrors.ErrPassword
	}

	// Origin accounts cannot have Admin role
	if acc.Role == role.Admin {
		uc.log.Error(ctx, "Login failed: Admin role is not allowed for origin accounts")
		return "", innerErrors.ErrUnauthorized
	}

	if !uc.bcrypt.CheckPasswordHash(password, acc.Password) {
		uc.log.Error(ctx, "Login failed CheckPasswordHash: ")
		return "", innerErrors.ErrPassword
	}
	token, err := uc.generateToken(ctx, strconv.FormatUint(uint64(acc.ID), 10), username, acc.Role)
	if err != nil {
		return token, err
	}
	return token, nil
}

func (uc *OriginAccountUseCase) CreateAccount(ctx context.Context, operatorRole role.Role, username string, userRole role.Role) (domainAccount.Account, error) {

	if err := uc.checkAdminRole(operatorRole); err != nil {
		uc.log.Error(ctx, "Unauthorized access to CreateAccount")
		return domainAccount.Account{}, err
	}
	if userRole == role.Admin {
		uc.log.Error(ctx, "Admin role cannot be assigned to a user")
		return domainAccount.Account{}, errors.New("admin role cannot be assigned to a user")
	}
	// Generate a random password of 12 characters
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	password := make([]byte, 12)
	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}
	hashedPassword, err := uc.bcrypt.HashPassword(string(password))
	if err != nil {
		uc.log.Error(ctx, "CreateAccount failed: "+err.Error())
		return domainAccount.Account{}, err
	}

	acc := domainAccount.Account{
		Username: username,
		Password: hashedPassword,
		Role:     userRole,
	}

	err = uc.accountRepo.Create(acc)
	if err != nil {
		uc.log.Error(ctx, "CreateAccount failed: "+err.Error())
		return domainAccount.Account{}, err
	}
	acc.Password = string(password)

	return acc, nil
}

func (uc *OriginAccountUseCase) GetAccountList(ctx context.Context, userRole role.Role) ([]dto.AccountListDTO, error) {
	if err := uc.checkAdminRole(userRole); err != nil {
		uc.log.Error(ctx, "Unauthorized access to GetAccountList")
		return nil, err
	}
	accounts, err := uc.accountRepo.GetAll()
	if err != nil {
		uc.log.Error(ctx, "GetAccountList failed: "+err.Error())
		return nil, err
	}

	// Convert []account.Account to []dto.AccountListDTO
	accountListDTOs := make([]dto.AccountListDTO, len(accounts))
	for i, acc := range accounts {
		accountListDTOs[i] = dto.AccountListDTO{
			Username: acc.Username,
			Role:     acc.Role,
		}
	}

	return accountListDTOs, nil
}

func (uc *OriginAccountUseCase) DeleteAccount(ctx context.Context, userRole role.Role, username string) error {
	if err := uc.checkAdminRole(userRole); err != nil {
		uc.log.Error(ctx, "Unauthorized access to DeleteAccount")
		return err
	}
	err := uc.accountRepo.Delete(username)
	if err != nil {
		uc.log.Error(ctx, "DeleteAccount failed: "+err.Error())
		return err
	}

	return nil
}

func (uc *OriginAccountUseCase) ChangePassword(ctx context.Context, userRole role.Role, username string, oldPassword, newPassword string) error {
	if err := uc.checkUserRole(userRole); err != nil {
		uc.log.Error(ctx, "Unauthorized access to ChangePassword")
		return err
	}
	// Retrieve the account using the username
	acc, err := uc.accountRepo.GetByUsername(username)
	if err != nil {
		uc.log.Error(ctx, "ChangePassword failed: "+err.Error())
		return innerErrors.ErrNotFound
	}

	// Verify the old password
	if !uc.bcrypt.CheckPasswordHash(oldPassword, acc.Password) {
		return innerErrors.ErrPassword
	}

	// Update the account with the new password
	hashedPassword, err := uc.bcrypt.HashPassword(newPassword)
	if err != nil {
		uc.log.Error(ctx, "ChangePassword failed: "+err.Error())
		return err
	}
	acc.Password = hashedPassword
	err = uc.accountRepo.Update(*acc)
	if err != nil {
		uc.log.Error(ctx, "ChangePassword failed: "+err.Error())
		return errors.New("failed to update password")
	}

	uc.log.Info(ctx, "Password changed successfully")
	return nil
}

func (u *OriginAccountUseCase) checkAdminRole(userRole role.Role) error {
	if userRole != role.Admin {
		return innerErrors.ErrUnauthorized
	}
	return nil
}

func (u *OriginAccountUseCase) checkUserRole(userRole role.Role) error {
	if userRole > role.User {
		return innerErrors.ErrUnauthorized
	}
	return nil
}

func (u *OriginAccountUseCase) generateToken(ctx context.Context, userID string, username string, userRole role.Role) (string, error) {
	jwt, err := u.jwtGenerator.GenerateOriginToken(ctx, userID, username, userRole, u.config.JWT.SecretKey)
	if err != nil {
		u.log.Error(ctx, "Error generating JWT: "+err.Error())
		return "", innerErrors.ErrInternal
	}
	return jwt, nil
}

func (u *OriginAccountUseCase) generateRedirectURL(token string) (string, error) {
	var redirectURL string
	if u.config.Server.HTTPS {
		redirectURL = fmt.Sprintf("https://%s?token=%s", u.config.Frontend.Domain, token)
	} else {
		redirectURL = fmt.Sprintf("http://%s:%s?token=%s", u.config.Frontend.Domain, strconv.Itoa(u.config.Frontend.Port), token)
	}
	return redirectURL, nil
}
