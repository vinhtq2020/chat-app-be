package usecase

import (
	"go-service/internal/auth/domain"
	user_info_domain "go-service/internal/user/domain"
	"go-service/pkg/jwt"
	"go-service/pkg/model"
	"go-service/pkg/uuid"
	"go-service/pkg/validate"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	repository         domain.AuthRepository
	userInfoRepository user_info_domain.UserRepository
	Validator          domain.UserLoginDataValidator
}

func NewAuthUsecase(repository domain.AuthRepository, Validator domain.UserLoginDataValidator, userInfoRepository user_info_domain.UserRepository,
) *AuthUsecase {
	return &AuthUsecase{
		repository:         repository,
		Validator:          Validator,
		userInfoRepository: userInfoRepository,
	}
}

// authen with email, phone or pin send code to authenticate
func AuthRegister(e *gin.Context) (int64, error) {
	panic("")
}

func (u *AuthUsecase) LoginWithGoogle(e *gin.Context, email string) ([]validate.ErrorMsg, int64, *jwt.TokenData, error) {
	errMsg, err := u.Validator.ValidateEmailGoogle(e, email)
	if err != nil || len(errMsg) > 0 {
		return errMsg, 0, nil, err
	}
	res, err := u.repository.GetUserLoginData(e, email)
	if err != nil {
		return nil, 0, nil, err
	}
	// case account not already created
	if res == nil {
		id := uuid.Pseudo_uuid()
		username := strings.Split(email, "@")[0]
		user := domain.UserLoginData{
			Id:       id,
			Username: &username,
			Email:    &email,
			Provider: model.GoogleProvider(),
			Status:   model.StatusActive(),
		}
		res, err := u.repository.Register(e, user)
		if err != nil {
			return nil, res, nil, err
		}
		//
		token := jwt.GenerateTokens(id, username, "my-secret-key")
		return nil, res, &token, err
	} else {
		token := jwt.GenerateTokens(res.Id, *res.Username, "my-secret-key")
		return nil, 1, &token, err
	}

}

func (u *AuthUsecase) Register(e *gin.Context, userLoginData domain.UserLoginData) ([]validate.ErrorMsg, int64, error) {
	listErr, err := u.Validator.ValidateLogin(e, userLoginData)
	if err != nil {
		return nil, -1, err
	}
	if len(listErr) > 0 {
		return listErr, -1, err
	}
	id := uuid.Pseudo_uuid()
	userLoginData.Id = id
	hash, err := bcrypt.GenerateFromPassword([]byte(*userLoginData.Password), 10)
	if err != nil {
		return nil, -1, err
	}
	userLoginData.PasswordHash = hash
	res, err := u.repository.InTransaction(e, func() (int64, error) {
		res, err := u.repository.Register(e, userLoginData)
		if err != nil {
			return res, err
		}
		res, err = u.userInfoRepository.Create(e, user_info_domain.User{
			Id:       id,
			UserName: userLoginData.Username,
			Email:    userLoginData.Email,
		})

		return res, err
	})

	if err != nil {
		return nil, res, err
	}
	return nil, res, err
}
