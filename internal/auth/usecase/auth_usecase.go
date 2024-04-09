package usecase

import (
	"errors"
	"go-service/internal/auth/domain"
	user_info_domain "go-service/internal/user/domain"
	"go-service/pkg/jwt"
	"go-service/pkg/model"
	"go-service/pkg/uuid"
	"go-service/pkg/validate"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	accessTokenDuration  = time.Minute
	refreshTokenDuration = time.Hour
)

type AuthUsecase struct {
	repository         domain.AuthRepository
	userInfoRepository user_info_domain.UserRepository
	Validator          domain.UserLoginDataValidator
	secretKey          string
}

func NewAuthUsecase(repository domain.AuthRepository, Validator domain.UserLoginDataValidator, userInfoRepository user_info_domain.UserRepository, secretKey string,
) *AuthUsecase {
	return &AuthUsecase{
		repository:         repository,
		Validator:          Validator,
		userInfoRepository: userInfoRepository,
		secretKey:          secretKey,
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
		token := jwt.GenerateTokens(id, username, u.secretKey, accessTokenDuration, refreshTokenDuration)
		return nil, res, &token, err
	} else {
		token := jwt.GenerateTokens(res.Id, *res.Username, u.secretKey, accessTokenDuration, refreshTokenDuration)
		return nil, 1, &token, err
	}

}

func (u *AuthUsecase) Login(e *gin.Context, email string, password string, browser string, ipAdress string, deviceId string) ([]validate.ErrorMsg, *jwt.TokenData, error) {
	errs, err := u.Validator.ValidateLogin(e, email)
	if err != nil || len(errs) > 0 {
		return errs, nil, err
	}

	exist, err := u.repository.Exist(e, email)
	if err != nil {
		return nil, nil, err
	}

	if exist == 0 {
		return nil, nil, nil
	}

	userInfo, err := u.repository.GetUserLoginData(e, email)
	if err != nil {
		return nil, nil, err
	}

	if userInfo.PasswordHash == nil {
		return nil, nil, errors.New("account haven't already had password")
	}
	err = bcrypt.CompareHashAndPassword([]byte(*userInfo.PasswordHash), []byte(password))
	if err != nil {
		return nil, nil, err
	}

	token := jwt.GenerateTokens(userInfo.Id, *userInfo.Email, u.secretKey, accessTokenDuration, refreshTokenDuration)

	_, err = u.repository.AddRefreshToken(e, domain.RefreshToken{
		UserId:    userInfo.Id,
		Token:     token.RefreshToken,
		Expiry:    refreshTokenDuration,
		IPAddress: ipAdress,
		DeviceId:  deviceId,
		Browser:   browser,
		CreatedAt: time.Now(),
	})

	if err != nil {
		return nil, nil, err
	}

	return nil, &token, nil
}

func (u *AuthUsecase) Register(e *gin.Context, userLoginData domain.UserLoginData) ([]validate.ErrorMsg, int64, error) {
	listErr, err := u.Validator.ValidateRegister(e, userLoginData)
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

	hashString := string(hash[:])
	userLoginData.PasswordHash = &hashString

	currentTime := time.Now()
	userLoginData.CreatedAt = &currentTime
	userLoginData.UpdatedAt = &currentTime
	userLoginData.Version = 1

	res, err := u.repository.InTransaction(e, func(tx *gorm.DB) (int64, error) {
		res, err := u.repository.Register(e, userLoginData)
		if err != nil {
			return res, err
		}

		res, err = u.userInfoRepository.Create(e, user_info_domain.User{
			Id:       id,
			UserName: userLoginData.Username,
			Version:  1,
		})

		return res, err
	})

	if err != nil {
		return nil, -1, err
	}
	return nil, res, nil
}

func (u *AuthUsecase) Logout(e *gin.Context, userId string, browser string, ipAdress string, deviceId string) (int64, error) {
	return u.repository.RemoveRefreshToken(e, userId, ipAdress, deviceId, browser)
}
