package port

import (
	"github.com/alielmi98/go-hexa-workout/internal/user/adapter/http/dto"
	"github.com/alielmi98/go-hexa-workout/internal/user/entity"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type TokenProvider interface {
	GenerateToken(token *entity.TokenPayload) (*dto.TokenDetail, error)
	VerifyToken(token string) (*jwt.Token, error)
	GetClaims(token string) (map[string]interface{}, error)
	RefreshToken(c *gin.Context) (*dto.TokenDetail, error)
}
