package usecase

import (
	"context"
	"errors"
	"go-service/internal/auth/domain"
	user_info_domain "go-service/internal/user/domain"
	"go-service/pkg/jwt"
	"go-service/pkg/model"
	"go-service/pkg/uuid"
	"go-service/pkg/validate"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUsecase struct {
	accountRepository      domain.AccountRepository
	userRepository         user_info_domain.UserRepository
	refreshTokenRepository domain.RefreshTokenRepository
	Validator              domain.AccountValidator
	secretKey              string
}

func NewAuthService(repository domain.AccountRepository,
	Validator domain.AccountValidator,
	userInfoRepository user_info_domain.UserRepository,
	refreshTokenRepository domain.RefreshTokenRepository,
	secretKey string,
) *AuthUsecase {
	return &AuthUsecase{
		accountRepository:      repository,
		Validator:              Validator,
		userRepository:         userInfoRepository,
		refreshTokenRepository: refreshTokenRepository,
		secretKey:              secretKey,
	}
}

// authen with email, phone or pin send code to authenticate
func AuthRegister(ctx context.Context) (int64, error) {
	panic("")
}

func (u *AuthUsecase) LoginWithGoogle(ctx context.Context, email string) ([]validate.ErrorMsg, int64, *jwt.TokenData, error) {
	errMsg, err := u.Validator.ValidateEmailGoogle(ctx, email)
	if err != nil || len(errMsg) > 0 {
		return errMsg, 0, nil, err
	}
	res, err := u.accountRepository.Load(ctx, email)
	if err != nil {
		return nil, 0, nil, err
	}
	// case account not already created
	if res == nil {
		id := uuid.Pseudo_uuid()
		username := strings.Split(email, "@")[0]
		user := domain.Account{
			Id:       id,
			Username: &username,
			Email:    &email,
			Provider: model.GoogleProvider(),
			Status:   model.StatusActive(),
		}
		res, err := u.accountRepository.Insert(ctx, user)
		if err != nil {
			return nil, res, nil, err
		}
		//
		token := jwt.GenerateTokens(id, u.secretKey, jwt.AccessTokenDuration, jwt.RefreshTokenDuration)
		return nil, res, &token, err
	} else {
		token := jwt.GenerateTokens(res.Id, u.secretKey, jwt.AccessTokenDuration, jwt.RefreshTokenDuration)
		return nil, 1, &token, err
	}

}

func (u *AuthUsecase) Login(ctx context.Context, email string, password string, browser string, ipAdress string, deviceId string) ([]validate.ErrorMsg, *jwt.TokenData, error) {
	errs, err := u.Validator.ValidateLogin(ctx, email)
	if err != nil || len(errs) > 0 {
		return errs, nil, err
	}

	// check user
	exist, err := u.accountRepository.Exist(ctx, email)
	if err != nil {
		return nil, nil, err
	}

	if exist == 0 {
		return nil, nil, nil
	}

	userInfo, err := u.accountRepository.Load(ctx, email)
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

	token := jwt.GenerateTokens(userInfo.Id, u.secretKey, jwt.AccessTokenDuration, jwt.RefreshTokenDuration)
	_, err = u.refreshTokenRepository.InTransaction(ctx, func(db *gorm.DB) (int64, error) {
		res, err := u.refreshTokenRepository.Delete(ctx, token.UserId, ipAdress, deviceId, browser)
		if err != nil {
			return res, err
		}

		_, err = u.refreshTokenRepository.Insert(ctx, domain.RefreshToken{
			UserId:    userInfo.Id,
			Token:     token.RefreshToken,
			Expiry:    jwt.RefreshTokenDuration,
			IPAddress: ipAdress,
			DeviceId:  deviceId,
			Browser:   browser,
			CreatedAt: time.Now(),
		})
		if err != nil {
			return -1, err
		}
		return 1, err
	})

	if err != nil {
		return nil, nil, err
	}

	return nil, &token, nil
}

func (u *AuthUsecase) Register(ctx context.Context, userLoginData domain.Account) ([]validate.ErrorMsg, int64, error) {
	listErr, err := u.Validator.ValidateRegister(ctx, userLoginData)
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

	res, err := u.accountRepository.InTransaction(ctx, func(tx *gorm.DB) (int64, error) {
		res, err := u.accountRepository.Insert(ctx, userLoginData)
		if err != nil {
			return res, err
		}

		res, err = u.userRepository.Create(ctx, user_info_domain.User{
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

func (u *AuthUsecase) Logout(ctx context.Context, userId string, browser string, ipAdress string, deviceId string) (int64, error) {
	return u.refreshTokenRepository.Delete(ctx, userId, ipAdress, deviceId, browser)
}

func (u *AuthUsecase) RefreshToken(ctx context.Context, userId string, browser string, ipAddress string, deviceId string) (int64, string, error) {
	oldToken, err := u.refreshTokenRepository.Load(ctx, browser, ipAddress, deviceId)
	if err != nil || oldToken == nil {
		return 0, "", err
	}

	// check is expiry
	if oldToken.CreatedAt.Add(oldToken.Expiry).Before(time.Now()) {
		return -2, "", nil
	}

	newToken, _ := jwt.GenerateAccessToken(userId, u.secretKey, jwt.AccessTokenDuration)
	return 1, newToken, nil
}
