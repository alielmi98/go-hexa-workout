package auth

import (
	"time"

	"github.com/alielmi98/go-hexa-workout/constants"
	"github.com/alielmi98/go-hexa-workout/internal/user/adapter/http/dto"
	"github.com/alielmi98/go-hexa-workout/internal/user/entity"
	"github.com/alielmi98/go-hexa-workout/pkg/config"
	"github.com/alielmi98/go-hexa-workout/pkg/service_errors"
	"github.com/golang-jwt/jwt"
)

type JwtProvider struct {
	cfg *config.Config
}

func NewJwtProvider(cfg *config.Config) *JwtProvider {
	return &JwtProvider{
		cfg: cfg,
	}
}

func (s *JwtProvider) GenerateToken(token *entity.TokenPayload) (*dto.TokenDetail, error) {
	td := &dto.TokenDetail{}
	td.AccessTokenExpireTime = time.Now().Add(s.cfg.JWT.AccessTokenExpireDuration * time.Minute).Unix()
	td.RefreshTokenExpireTime = time.Now().Add(s.cfg.JWT.RefreshTokenExpireDuration * time.Minute).Unix()

	atc := jwt.MapClaims{}

	atc[constants.UserIdKey] = token.UserId
	atc[constants.FirstNameKey] = token.FirstName
	atc[constants.LastNameKey] = token.LastName
	atc[constants.UsernameKey] = token.Username
	atc[constants.EmailKey] = token.Email
	atc[constants.MobileNumberKey] = token.MobileNumber
	atc[constants.ExpireTimeKey] = td.AccessTokenExpireTime

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atc)

	var err error
	td.AccessToken, err = at.SignedString([]byte(s.cfg.JWT.Secret))

	if err != nil {
		return nil, err
	}

	rtc := jwt.MapClaims{}
	rtc[constants.UserIdKey] = token.UserId
	rtc[constants.FirstNameKey] = token.FirstName
	rtc[constants.LastNameKey] = token.LastName
	rtc[constants.UsernameKey] = token.Username
	rtc[constants.EmailKey] = token.Email
	rtc[constants.MobileNumberKey] = token.MobileNumber
	rtc[constants.ExpireTimeKey] = td.RefreshTokenExpireTime

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtc)

	td.RefreshToken, err = rt.SignedString([]byte(s.cfg.JWT.RefreshSecret))

	if err != nil {
		return nil, err
	}

	return td, nil
}

func (s *JwtProvider) VerifyToken(token string) (*jwt.Token, error) {
	at, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, &service_errors.ServiceError{EndUserMessage: service_errors.UnExpectedError}
		}
		return []byte(s.cfg.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	return at, nil
}

func (s *JwtProvider) GetClaims(token string) (claimMap map[string]interface{}, err error) {
	claimMap = map[string]interface{}{}

	verifyToken, err := s.VerifyToken(token)
	if err != nil {
		return nil, err
	}
	claims, ok := verifyToken.Claims.(jwt.MapClaims)
	if ok && verifyToken.Valid {
		for k, v := range claims {
			claimMap[k] = v
		}
		return claimMap, nil
	}
	return nil, &service_errors.ServiceError{EndUserMessage: service_errors.ClaimsNotFound}
}
func (s *JwtProvider) RefreshToken(refreshToken string) (*dto.TokenDetail, error) {
	claims, err := s.GetClaims(refreshToken)
	if err != nil {
		return nil, err
	}

	tokenDto := entity.TokenPayload{
		UserId:       int(claims[constants.UserIdKey].(float64)),
		FirstName:    claims[constants.FirstNameKey].(string),
		LastName:     claims[constants.LastNameKey].(string),
		Username:     claims[constants.UsernameKey].(string),
		MobileNumber: claims[constants.MobileNumberKey].(string),
		Email:        claims[constants.EmailKey].(string),
	}
	newTokenDetail, err := s.GenerateToken(&tokenDto)
	if err != nil {
		return nil, err
	}

	return newTokenDetail, nil
}
