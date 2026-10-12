package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gocgrpcarchitecture/generate/proto/userProto"
	"gocgrpcarchitecture/internal/app/resources/user/models"
	"gocgrpcarchitecture/internal/app/resources/user/resource"
	"gocgrpcarchitecture/internal/app/resources/user/service"
	"gocgrpcarchitecture/internal/pkg/utils"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const MAXIMUM = 3
const AWAIT = 1

type UserResourceHandler struct {
	UserService    service.UserService
	ServerResource resource.UserServerResource
	conn           *grpc.ClientConn
}

func NewUserResourceHandler() *UserResourceHandler {
	return &UserResourceHandler{}
}

func (app *UserResourceHandler) UserResourceHandlerRoutes(mux *http.ServeMux, conn *grpc.ClientConn) {
	time.Sleep(2 * time.Second)
	mux.HandleFunc("POST /register", app.RegisterUserHandler(conn))
	mux.HandleFunc("POST /login", app.LoginUserHandler(conn))
	mux.HandleFunc("GET /verify/account", app.FindEmailTokenVerificationHandler(conn))

	mux.HandleFunc("GET /user/{uuid}", app.GetUserByUUIDHandler(conn))
	mux.HandleFunc("GET /email/{email}", app.GetUserByEmailHandler(conn))
	mux.HandleFunc("PUT /user/{uuid}", app.UpdateUserByIdHandler(conn))
	mux.HandleFunc("DELETE /user/{uuid}", app.DeleteUserByIdHandler(conn))
	return
}

func (app *UserResourceHandler) RegisterUserHandler(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var UR models.UserRole
		client := userProto.NewUserServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var US struct{ UniqueID, FirstName, LastName, Email, Password, Phone string }
		if err := json.NewDecoder(r.Body).Decode(&US); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"error": "Error - Invalid. All fields is required.",
			})
			return
		}
		_, err := utils.RegexValidate(US.FirstName)
		if err != nil {
			utils.WriteJSON(w, http.StatusNotAcceptable, map[string]any{
				"error": "Error - That field an invalid characters.",
			})
			return
		}
		_, err = utils.RegexValidate(US.LastName)
		if err != nil {
			utils.WriteJSON(w, http.StatusNotAcceptable, map[string]any{
				"error": "Error - That field an invalid characters.",
			})
			return
		}
		findEmail, _ := client.FindUserByEmail(ctx, &userProto.FindUserByEmailRequest{
			User: &userProto.User{
				Email: utils.IsToLower(US.Email),
			},
		})
		if findEmail.User.Email == utils.IsToLower(US.Email) {
			utils.WriteJSON(w, http.StatusNotAcceptable, map[string]any{
				"error": "Error - User with that email already exists.",
			})
			return
		}
		_, err = utils.IsMatch(US.Phone)
		if err != nil {
			utils.WriteJSON(w, http.StatusNotAcceptable, map[string]any{
				"error": "Error - That field an invalid phone numbers.",
			})
			return
		}
		response, err := client.CreateUser(ctx, &userProto.CreateUserRequest{
			User: &userProto.User{
				FirstName:             utils.Escape(US.FirstName),
				LastName:              utils.Escape(US.LastName),
				Email:                 utils.IsToLower(US.Email),
				Password:              US.Password,
				Phone:                 utils.Escape(US.Phone), // Need to add validate with regex
				PasswordLastChangedAt: timestamppb.New(time.Now().UTC()),
			},
		})
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"error": "Error - Create User: status internal server error."})
			return
		}
		UR.UserId = response.User.GetId()
		roles, err := client.AddUserToRole(ctx, &userProto.CreateUserRoleRequest{
			UserRole: &userProto.UserRole{
				UserId: UR.UserId,
				RoleId: 1,
			},
		})
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"error": "Error - Add To Role: status internal server error."})
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]any{
			"message": "Success - Your account have successful created",
			"data":    response,
			"roles":   roles,
		})
	}
}
func (app *UserResourceHandler) LoginUserHandler(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := userProto.NewUserServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var US struct{ Email, Password, Phone string }
		if err := json.NewDecoder(r.Body).Decode(&US); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"error": "Error - Invalid. All fields is required.",
			})
			return
		}
		findEmail, _ := client.FindUserByEmail(ctx, &userProto.FindUserByEmailRequest{
			User: &userProto.User{
				Email: utils.IsToLower(US.Email),
			},
		})
		if findEmail.User.GetEmail() != utils.IsToLower(US.Email) {
			utils.WriteJSON(w, http.StatusNotAcceptable, map[string]any{
				"error": "Error - invalid credentials email and password.",
			})
			return
		}
		findEmailVerify, _ := client.FindUserEmailVerification(ctx, &userProto.FindUserEmailVerificationRequest{
			User: &userProto.User{
				Email: US.Email,
			},
		})
		if findEmailVerify.User.GetEmailVerified().GetValue() == false {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]any{
				"error": "Error - Credentials error. You have not verified your account email.",
			})
			return
		}
		IsEnabled, _ := client.IsEnabledUser(ctx, &userProto.IsEnabledUserRequest{
			User: &userProto.User{
				Email: US.Email,
			},
		})

		if IsEnabled.User.GetAccountEnabled().GetValue() == false {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]any{
				"error": "Error - Your account is not enabled. Contact help and service before try again.",
			})
			return
		}
		IsDisabled, _ := client.IsDisabledUser(ctx, &userProto.IsDisabledUserRequest{
			User: &userProto.User{
				Email: US.Email,
			},
		})
		if IsDisabled.User.GetAccountDisabled().GetValue() == true {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]any{
				"error": "Error - Your account is disabled. Contact help and service before try again.",
			})
			return
		}

		password := US.Password
		hashedPassword := findEmail.User.GetPassword()
		if findEmail.User.GetFailedAttempts() >= MAXIMUM && time.Since(findEmail.User.GetLastFailedAttemptTime().AsTime()) < AWAIT*time.Minute {
			utils.WriteJSON(w, http.StatusForbidden, map[string]any{
				"error": "Error - Account is locked. Try again in 5 minutes.",
			})
			return
		}
		err := utils.ComparePassword(hashedPassword, password)
		if err != nil {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				_, _ = client.IncrementFailedLoginAttempts(ctx, &userProto.IncrementFailedLoginAttemptsRequest{
					User: &userProto.User{
						Email:          findEmail.User.GetEmail(),
						FailedAttempts: findEmail.User.GetFailedAttempts(),
					},
				})
				utils.WriteJSON(w, http.StatusUnauthorized, map[string]any{
					"error": "Error - An invalid unauthorized credentials.",
				})
				return
			}
		}

		fmt.Println("Hello world ", findEmail.User.PasswordLastChangedAt)

		if utils.HashPasswordIsExpired(findEmail.User.PasswordLastChangedAt.AsTime()) {
			utils.WriteJSON(w, http.StatusForbidden, map[string]any{
				"error": "Error - Hashed password is expired. Please reset your password.",
			})
			return
		}
		Time, _ := client.LastLoginTimeStamps(ctx, &userProto.LastLoginTimeStampsRequest{
			User: &userProto.User{
				Email: findEmail.User.GetEmail(),
			},
		})
		fmt.Println("LLastLoginTimeStamps ", Time.User.LastLoginTimestamp)

		utils.WriteJSON(w, http.StatusOK, map[string]any{
			"message": "Success - Your account have successful created",
			"data":    findEmail.GetUser(),
		})
	}
}
func (app *UserResourceHandler) GetUserByUUIDHandler(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := userProto.NewUserServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		UUID := r.PathValue("uuid")
		findId, _ := client.FindUserById(ctx, &userProto.FindUserByIdRequest{
			UUId: UUID,
		})
		if UUID != findId.User.GetUUId() {
			utils.WriteJSON(w, http.StatusNotFound, map[string]any{
				"error": "Error - User with that id doesn't exists.",
			})
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]any{
			"message": "Success",
			"data":    findId,
		})
	}
}
func (app *UserResourceHandler) GetUserByEmailHandler(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := userProto.NewUserServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		EM := r.PathValue("email")
		findEmail, _ := client.FindUserByEmail(ctx, &userProto.FindUserByEmailRequest{
			User: &userProto.User{
				Email: EM,
			},
		})
		if EM != findEmail.User.GetEmail() {
			utils.WriteJSON(w, http.StatusNotAcceptable, map[string]any{
				"error": "Error - User with that email doesn't exists.",
			})
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]any{
			"message": "Success",
			"data":    findEmail.GetUser(),
		})
	}
}
func (app *UserResourceHandler) UpdateUserByIdHandler(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := userProto.NewUserServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var US struct{ FirstName, LastName, Email, Phone string }
		if err := json.NewDecoder(r.Body).Decode(&US); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"error": "Error - Invalid. All fields is required.",
			})
			return
		}
		UId := r.PathValue("uuid")
		findId, _ := client.FindUserById(ctx, &userProto.FindUserByIdRequest{
			UUId: UId,
		})
		if UId != findId.User.GetUUId() {
			utils.WriteJSON(w, http.StatusNotFound, map[string]any{
				"error": "Error - User with that id doesn't exists.",
			})
			return
		}

		findEmail, _ := client.FindUserByEmail(ctx, &userProto.FindUserByEmailRequest{
			User: &userProto.User{
				Email: utils.IsToLower(US.Email),
			},
		})
		if findEmail.User.Email != utils.IsToLower(US.Email) {
			utils.WriteJSON(w, http.StatusNotAcceptable, map[string]any{
				"error": "Error - User with that email doesn't exists.",
			})
			return
		}

		userById, err := client.UpdateUserById(ctx, &userProto.UpdateUserByIdRequest{
			User: &userProto.User{
				UUId:      UId,
				FirstName: US.FirstName,
				LastName:  US.LastName,
				Email:     US.Email,
				Phone:     US.Phone,
			},
		})
		if err != nil {
			utils.WriteJSON(w, http.StatusNotFound, map[string]any{
				"error": "Error - User with that id doesn't exists.",
			})
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]any{
			"message": "Success - Your account have successful updated",
			"data":    userById.GetUser(),
		})
	}
}
func (app *UserResourceHandler) DeleteUserByIdHandler(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := userProto.NewUserServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		UUId := r.PathValue("uuid")
		findId, _ := client.FindUserById(ctx, &userProto.FindUserByIdRequest{
			UUId: UUId,
		})
		if UUId != findId.User.GetUUId() {
			utils.WriteJSON(w, http.StatusNotFound, map[string]any{
				"error": "Error - User with that id doesn't exists........",
			})
			return
		}
		_, err := client.DeleteUserById(ctx, &userProto.DeleteUserByIdRequest{
			UUId: UUId,
		})
		if err != nil {
			utils.WriteJSON(w, http.StatusNotFound, map[string]any{
				"error": "Error - User with that id doesn't exists.",
			})
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]any{
			"message": "This user have been deleted",
		})
	}
}
func (app *UserResourceHandler) FindEmailTokenVerificationHandler(conn *grpc.ClientConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := userProto.NewUserServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		UUId := r.URL.Query().Get("uuid")
		StrToken := r.URL.Query().Get("token")
		if UUId == "" || StrToken == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"error": "Error - That ID doesn't exist",
			})
			return
		}
		findToken, _ := client.FindEmailTokenVerification(ctx, &userProto.FindEmailTokenVerificationRequest{
			User: &userProto.User{
				UUId:       UUId,
				EmailToken: StrToken,
			},
		})
		if findToken.User.GetEmailToken() == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"error": "Error - Your accounts already activated.",
			})
			return
		}

		updateToken, err := client.UpdateEmailTokenVerification(ctx, &userProto.UpdateEmailTokenVerificationRequest{
			User: &userProto.User{
				EmailToken:        "",
				EmailVerified:     wrapperspb.Bool(true),
				AccountEnabled:    wrapperspb.Bool(true),
				AccountDisabled:   wrapperspb.Bool(false),
				AccountDisabledAt: "",
			},
		})
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]any{
				"error": "Error - an error verifying token",
			})
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]any{
			"data":    updateToken,
			"message": "Your account has been activated.",
		})
	}
}
