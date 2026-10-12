package resource

import (
	"context"
	"gocgrpcarchitecture/generate/proto/userProto"
	"gocgrpcarchitecture/internal/app/resources/user/models"
	"gocgrpcarchitecture/internal/app/resources/user/service"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	//"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"gorm.io/gorm"
)

type UserServerResource struct {
	userProto.UnimplementedUserServiceServer
	UserService service.UserService
	DB          *gorm.DB
}

func (S *UserServerResource) AddUserToRole(ctx context.Context, req *userProto.CreateUserRoleRequest) (*userProto.CreateUserRoleResponse, error) {
	if err, _ := S.UserService.AddUserToRole(ctx, req.UserRole.GetUserId(), req.UserRole.GetRoleId()); err != nil {
		return nil, status.Errorf(codes.Internal, "error adding user to role: %v", err)
	}
	return &userProto.CreateUserRoleResponse{UserRole: &userProto.UserRole{
		UserId: req.UserRole.GetUserId(),
		RoleId: req.UserRole.GetRoleId(),
	}}, nil
}
func (S *UserServerResource) CreateUser(ctx context.Context, req *userProto.CreateUserRequest) (*userProto.CreateUserResponse, error) {
	U := &models.User{
		FirstName:             req.User.GetFirstName(),
		LastName:              req.User.GetLastName(),
		Email:                 req.User.GetEmail(),
		Password:              req.User.GetPassword(),
		Phone:                 req.User.GetPhone(),
		PasswordLastChangedAt: req.User.GetPasswordLastChangedAt().AsTime(),
	}
	if err := S.UserService.Create(ctx, U); err != nil {
		return nil, status.Errorf(codes.Internal, "Error creating user: %v", err)
	}
	return &userProto.CreateUserResponse{User: &userProto.User{
		Id:                    U.Id,
		UUId:                  U.UUId,
		UniqueId:              U.UniqueId,
		FirstName:             U.FirstName,
		LastName:              U.LastName,
		Email:                 U.Email,
		Phone:                 U.Phone,
		PasswordLastChangedAt: timestamppb.New(U.PasswordLastChangedAt),
	}}, nil
}
func (S *UserServerResource) FindUserById(ctx context.Context, req *userProto.FindUserByIdRequest) (*userProto.FindUserByIdResponse, error) {
	UI, err := S.UserService.FindById(ctx, req.GetUUId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "error getting user: %v", err)
	}

	var URoles []*userProto.Role
	for _, role := range UI.Roles {
		URoles = append(URoles, &userProto.Role{
			Id:          role.Id,
			Name:        role.Name,
			Permissions: role.Permissions,
		})
	}
	return &userProto.FindUserByIdResponse{User: &userProto.User{
		Id:        UI.Id,
		UUId:      UI.UUId,
		UniqueId:  UI.UniqueId,
		FirstName: UI.FirstName,
		LastName:  UI.LastName,
		Email:     UI.Email,
		Phone:     UI.Phone,
		Roles:     URoles,
	}}, nil
}
func (S *UserServerResource) FindUserByEmail(ctx context.Context, req *userProto.FindUserByEmailRequest) (*userProto.FindUserByEmailResponse, error) {
	StrEmail := strings.TrimSpace(strings.ToLower(req.User.GetEmail()))
	EM, err := S.UserService.FindByEmail(ctx, StrEmail)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check email status: %v", err)
	}
	var URoles []*userProto.Role
	for _, role := range EM.Roles {
		URoles = append(URoles, &userProto.Role{
			Id:          role.Id,
			Name:        role.Name,
			Permissions: role.Permissions,
		})
	}
	return &userProto.FindUserByEmailResponse{User: &userProto.User{
		Id:                    EM.Id,
		UUId:                  EM.UUId,
		UniqueId:              EM.UniqueId,
		FirstName:             EM.FirstName,
		LastName:              EM.LastName,
		Password:              EM.Password,
		Email:                 EM.Email,
		Phone:                 EM.Phone,
		FailedAttempts:        uint64(EM.FailedAttempts),
		LastLoginTimestamp:    EM.LastLoginTimeStamp.Time.String(),
		LastFailedAttemptTime: timestamppb.New(EM.LastFailedAttemptTime.Time),
		PasswordLastChangedAt: timestamppb.New(EM.PasswordLastChangedAt.UTC()),
	}}, nil
}
func (S *UserServerResource) FindUserEmailVerification(ctx context.Context, req *userProto.FindUserEmailVerificationRequest) (*userProto.FindUserEmailVerificationResponse, error) {
	StrEmail := strings.TrimSpace(strings.ToLower(req.User.GetEmail()))
	EM, err := S.UserService.FindEmailVerification(ctx, StrEmail)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error verifying user: %v", err)
	}

	return &userProto.FindUserEmailVerificationResponse{User: &userProto.User{
		EmailVerified: wrapperspb.Bool(EM.EmailVerified.Bool),
	}}, nil
}
func (S *UserServerResource) FindEmailTokenVerification(ctx context.Context, req *userProto.FindEmailTokenVerificationRequest) (*userProto.FindEmailTokenVerificationResponse, error) {
	EM, err := S.UserService.FindEmailTokenVerification(ctx, req.User.GetUUId(), req.User.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error verifying user: %v", err)
	}
	return &userProto.FindEmailTokenVerificationResponse{User: &userProto.User{
		UUId:       EM.UUId,
		EmailToken: EM.EmailToken,
	}}, nil
}
func (S *UserServerResource) UpdateEmailTokenVerification(ctx context.Context, req *userProto.UpdateEmailTokenVerificationRequest) (*userProto.UpdateEmailTokenVerificationResponse, error) {
	StrEmail := strings.TrimSpace(strings.ToLower(req.User.GetEmail()))
	EM, err := S.UserService.UpdateEmailTokenVerification(ctx, StrEmail)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error updating email verification: %v", err)
	}

	return &userProto.UpdateEmailTokenVerificationResponse{User: &userProto.User{
		EmailToken:        "",
		EmailVerified:     wrapperspb.Bool(EM.EmailVerified.Bool),
		AccountEnabled:    wrapperspb.Bool(EM.AccountEnabled.Bool),
		AccountDisabled:   wrapperspb.Bool(EM.AccountDisabled.Bool),
		AccountDisabledAt: "",
	}}, nil
}
func (S *UserServerResource) IncrementFailedLoginAttempts(ctx context.Context, req *userProto.IncrementFailedLoginAttemptsRequest) (*userProto.IncrementFailedLoginAttemptsResponse, error) {
	StrEmail := strings.TrimSpace(strings.ToLower(req.User.GetEmail()))
	EM, err := S.UserService.FindByEmail(ctx, StrEmail)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error finding user by email: %v", err)
	}

	err = S.UserService.IncrementFailedLoginAttempts(ctx, EM)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error incrementing failed login attempts: %v", err)
	}

	return &userProto.IncrementFailedLoginAttemptsResponse{User: &userProto.User{
		Email:                 req.User.GetEmail(),
		FailedAttempts:        req.User.GetFailedAttempts(),
		LastFailedAttemptTime: timestamppb.New(time.Now()),
	}}, nil
}
func (S *UserServerResource) LastLoginTimeStamps(ctx context.Context, req *userProto.LastLoginTimeStampsRequest) (*userProto.LastLoginTimeStampsResponse, error) {
	StrEmail := strings.TrimSpace(strings.ToLower(req.User.GetEmail()))
	EM, err := S.UserService.FindByEmail(ctx, StrEmail)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error finding user by email: %v", err)
	}

	err = S.UserService.LastLoginTimeStamps(ctx, EM)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error updating last login timestamp: %v", err)
	}

	return &userProto.LastLoginTimeStampsResponse{User: &userProto.User{
		Email:              EM.Email,
		LastLoginTimestamp: time.Now().String(),
	}}, nil
}

func (S *UserServerResource) IsEnabledUser(ctx context.Context, req *userProto.IsEnabledUserRequest) (*userProto.IsEnabledUserResponse, error) {
	EM, err := S.UserService.IsEnabledUser(ctx, req.User.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error enabled user: %v", err)
	}
	return &userProto.IsEnabledUserResponse{User: &userProto.User{
		AccountEnabled: wrapperspb.Bool(EM.AccountEnabled.Bool),
	}}, nil
}
func (S *UserServerResource) IsDisabledUser(ctx context.Context, req *userProto.IsDisabledUserRequest) (*userProto.IsDisabledUserResponse, error) {
	EM, err := S.UserService.IsDisabledUser(ctx, req.User.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error disabled user: %v", err)
	}
	return &userProto.IsDisabledUserResponse{User: &userProto.User{
		AccountDisabled: wrapperspb.Bool(EM.AccountDisabled.Bool),
	}}, nil
}
func (S *UserServerResource) UpdateUserById(ctx context.Context, req *userProto.UpdateUserByIdRequest) (*userProto.UpdateUserByIdResponse, error) {
	US := &models.User{
		UUId:      req.User.GetUUId(),
		FirstName: req.User.GetFirstName(),
		LastName:  req.User.GetLastName(),
		Email:     req.User.GetEmail(),
		Phone:     req.User.GetPhone(),
	}
	if err := S.UserService.Update(ctx, US); err != nil {
		return nil, status.Errorf(codes.Internal, "error updating user: %v", err)
	}
	return &userProto.UpdateUserByIdResponse{User: &userProto.User{
		Id:        US.Id,
		UUId:      US.UUId,
		FirstName: US.FirstName,
		LastName:  US.LastName,
		Email:     US.Email,
		Phone:     US.Phone,
	}}, nil
}
func (S *UserServerResource) DeleteUserById(ctx context.Context, req *userProto.DeleteUserByIdRequest) (*userProto.DeleteUserByIdResponse, error) {
	if err := S.UserService.Delete(ctx, req.GetUUId()); err != nil {
		return nil, status.Errorf(codes.Internal, "error deleting user: %v", err)
	}
	return &userProto.DeleteUserByIdResponse{Success: true}, nil
}
