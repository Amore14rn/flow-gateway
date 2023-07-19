package handlers

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (h *Handler) verifyJWTHandler(c *gin.Context) bool {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []gin.H{
			{"code": "not_authenticated", "detail": "Учетные данные не были предоставлены."},
		}})
		return false
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверка алгоритма подписи токена
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неправильный алгоритм подписи токена: %v", token.Header["alg"])
		}

		// Возвращение секретного ключа
		return []byte(h.secret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": []gin.H{
			{"code": "invalid_token", "detail": "Неверный JWT"},
		}})
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": []gin.H{
			{"code": "invalid_token", "detail": "Неверный JWT"},
		}})
		return false
	}

	c.JSON(http.StatusOK, gin.H{"message": "JWT верифицирован", "claims": claims})
	return true
}
