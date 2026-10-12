package service

import (
	"context"
	"gocgrpcarchitecture/internal/app/resources/user/models"
)

type UserService interface {
	Create(ctx context.Context, user *models.User) error
	AddUserToRole(ctx context.Context, UserId uint64, RoleId uint64) (*models.UserRole, error)
	FindById(ctx context.Context, UUId string) (*models.User, error)
	FindByEmail(ctx context.Context, Email string) (*models.User, error)
	FindEmailVerification(ctx context.Context, Email string) (*models.User, error)

	FindEmailTokenVerification(ctx context.Context, UUId string, strToken string) (*models.User, error)
	UpdateEmailTokenVerification(ctx context.Context, Email string) (*models.User, error)
	IncrementFailedLoginAttempts(ctx context.Context, user *models.User) error
	LastLoginTimeStamps(ctx context.Context, user *models.User) error
	//CheckEmailTokenVerification(ctx context.Context, UUId string, strToken string) (bool, error)

	IsEnabledUser(ctx context.Context, Email string) (*models.User, error)
	IsDisabledUser(ctx context.Context, Email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, UUId string) error
	GetAllUsers(ctx context.Context) ([]*models.User, error)
}

/*
type UserService interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByEmail(string) (*models.User, error)
	VerifyUserByEmail(string) (*models.User, error)
	IsEnabledUser(string) (*models.User, error)
	IsDisabledUser(string) (*models.User, error)

	GetUserById(string) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	UpdateUserById(user *models.User) error
	DeleteUserById(string) error
	AddToUserRole(*models.UserRole) (*models.UserRole, error)
	GetUserTokenEmailVerification(string, string) (*models.User, error)
	UpdateUserEmailVerification(user *models.User) error
	ClearRefreshToken(string) error
	StoreRefreshToken(token *models.RefreshToken) (*models.RefreshToken, error)
	GetRefreshToken(token *models.RefreshToken) (*models.RefreshToken, error)
	RevokeRefreshToken(string) error
	DeleteRefreshToken(UUID string, email string) error
	ResetUserPasswordToken(user *models.User) error
	VerifyUserPasswordToken(string) (*models.User, error)
	UpdateUserPasswordToken(user *models.User) error
	IncrementFailedLoginAttempts(user *models.User) error
	ResetFailedLoginAttempts(user *models.User) error

}
*/
