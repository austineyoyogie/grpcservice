package implementation

import (
	"context"
	"gocgrpcarchitecture/internal/app/resources/user/models"
	"gocgrpcarchitecture/internal/app/resources/user/service"
	"gocgrpcarchitecture/internal/pkg/utils"
	"time"

	"gorm.io/gorm"
)

type UserServiceImpl struct {
	G *gorm.DB
}

func NewUserServiceImpl(G *gorm.DB) service.UserService {
	return &UserServiceImpl{G: G}
}
func (U UserServiceImpl) AddUserToRole(ctx context.Context, UserId uint64, RoleId uint64) (*models.UserRole, error) {
	var user models.User
	tx := U.G.Begin()
	U.G.WithContext(ctx).Model(&models.User{}).Where("Id = ?", UserId).First(&user)
	columns := map[string]interface{}{
		"user_id":   user.Id,
		"role_id":   true,
		"CreatedAt": time.Now(),
		"UpdatedAt": time.Now(),
	}
	err := U.G.WithContext(ctx).Model(&models.UserRole{}).Create(&columns).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return nil, tx.Commit().Error
}
func (U UserServiceImpl) Create(ctx context.Context, user *models.User) error {
	//strTime := time.Now()
	tx := U.G.Begin()
	user.FirstName = utils.Escape(user.FirstName)
	user.LastName = utils.Escape(user.LastName)
	user.Email = utils.IsToLower(user.Email)
	hash, _ := utils.HashPassword(user.Password)
	user.Password = string(hash)
	strToken := utils.RandomString(255)
	user.EmailToken = strToken
	//user.PasswordLastChangedAt = strTime
	if err := U.G.WithContext(ctx).Model(&models.User{}).Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}
	//subject := "Account Email Verification"
	//sendEmail := message.Deliver([]string{user.Email}, subject)
	//verification := fmt.Sprintf("http://localhost:4200/v2/verify/account?uuid=%s&token=%s", user.UUId, user.EmailToken)
	//sendEmail.EmailTemplate("aPI/messages/verification.html", verification)
	return tx.Commit().Error
}
func (U UserServiceImpl) FindById(ctx context.Context, UUId string) (*models.User, error) {
	tx := U.G.Begin()
	user := &models.User{}
	if err := U.G.WithContext(ctx).Model(&models.User{}).Where("UUId = ?", UUId).Preload("Roles").Find(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return user, tx.Commit().Error
}
func (U UserServiceImpl) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	tx := U.G.Begin()
	user := &models.User{}
	if err := U.G.WithContext(ctx).Model(&models.User{}).Where("Email = ?", email).Preload("Roles").Find(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return user, tx.Commit().Error
}
func (U UserServiceImpl) FindEmailVerification(ctx context.Context, Email string) (*models.User, error) {
	tx := U.G.Begin()
	user := &models.User{}
	if err := U.G.WithContext(ctx).Model(&models.User{}).Where("Email = ?", Email).Where("EmailVerified = ?", true).Find(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return user, tx.Commit().Error
}
func (U UserServiceImpl) FindEmailTokenVerification(ctx context.Context, UUId string, strToken string) (*models.User, error) {
	tx := U.G.Begin()
	user := &models.User{}
	if err := U.G.WithContext(ctx).Model(&models.User{}).Where("UUId = ?", UUId).Where("EmailToken = ?", strToken).Find(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return user, tx.Commit().Error
}
func (U UserServiceImpl) UpdateEmailTokenVerification(ctx context.Context, Email string) (*models.User, error) {
	tx := U.G.Begin()
	columns := map[string]interface{}{
		"EmailToken":        "",
		"EmailVerified":     true,
		"AccountEnabled":    true,
		"AccountDisabled":   false,
		"AccountDisabledAt": nil,
	}
	if err := U.G.WithContext(ctx).Model(&models.User{}).Where("Email = ?", Email).UpdateColumns(&columns).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return nil, tx.Commit().Error
}
func (U UserServiceImpl) IncrementFailedLoginAttempts(ctx context.Context, user *models.User) error {
	tx := U.G.Begin()
	columns := map[string]interface{}{
		"FailedAttempts":        user.FailedAttempts + 1,
		"LastFailedAttemptTime": time.Now(),
	}
	if err := U.G.WithContext(ctx).Model(&models.User{}).Where("Email = ?", user.Email).UpdateColumns(&columns).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
func (U UserServiceImpl) LastLoginTimeStamps(ctx context.Context, user *models.User) error {
	tx := U.G.Begin()
	columns := map[string]interface{}{
		"LastLoginTimestamp": time.Now(),
	}
	if err := U.G.WithContext(ctx).Model(&models.User{}).Where("Email = ?", user.Email).UpdateColumns(&columns).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (U UserServiceImpl) IsEnabledUser(ctx context.Context, Email string) (*models.User, error) {
	tx := U.G.Begin()
	user := &models.User{}
	if err := U.G.WithContext(ctx).Model(&models.User{}).Where("Email = ?", Email).Where("AccountEnabled = ?", true).Find(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return user, tx.Commit().Error
}
func (U UserServiceImpl) IsDisabledUser(ctx context.Context, Email string) (*models.User, error) {
	tx := U.G.Begin()
	user := &models.User{}
	if err := U.G.WithContext(ctx).Model(&models.User{}).Where("Email = ?", Email).Where("AccountDisabled = ?", true).Find(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return user, tx.Commit().Error
}
func (U UserServiceImpl) Update(ctx context.Context, user *models.User) error {
	tx := U.G.Begin()
	columns := map[string]interface{}{
		"FirstName": user.FirstName,
		"LastName":  user.LastName,
		"Phone":     user.Phone,
		"UpdatedAt": time.Now(),
	}
	if err := U.G.WithContext(ctx).Model(&models.User{}).Where("Email = ?", user.Email).Where("UUId = ?", user.UUId).UpdateColumns(columns).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error

	// NEED TO SEND UPDATE EMAIL TO VERIFY CREDENTIAL AND NOTIFICATION
}
func (U UserServiceImpl) Delete(ctx context.Context, UUId string) error {
	tx := U.G.Begin()
	columns := map[string]interface{}{
		"EmailVerified":     false,
		"AccountEnabled":    false,
		"AccountDisabled":   true,
		"AccountDisabledAt": time.Now(),
	}
	if err := U.G.WithContext(ctx).Model(&models.User{}).Where("UUId = ?", UUId).UpdateColumns(columns).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error

	// NEED TO SEND UPDATE EMAIL TO VERIFY CREDENTIAL AND NOTIFICATION
}
func (U UserServiceImpl) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	//TODO implement me
	panic("implement me")
}

//func (U UserServiceImpl) CheckEmailTokenVerification(ctx context.Context, UUId string, strToken string) (bool, error) {
//
//	return true, nil
//}
