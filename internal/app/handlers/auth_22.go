package handlers

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	pg "gitlab.com/mildd/flow-gateway/internal/app/commands/secret"
)

func (h *Handler) auth22Handler(c *gin.Context) {
	// Handle auth22 logic

	// Проверка наличия заголовка 'X-User-Permissions'
	if c.GetHeader("X-User-Permissions") != "" {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []gin.H{
			{"code": "not_authenticated", "detail": "Учетные данные не были предоставлены."},
		}})
		return
	}

	// Проверка JWT-токена
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"errors": []gin.H{
			{"code": "not_authenticated", "detail": "Учетные данные не были предоставлены."},
		}})
		return
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
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": []gin.H{
			{"code": "invalid_token", "detail": "Неверный JWT"},
		}})
		return
	}

	// Получение user_id из payload токена
	userID, ok := claims["user_id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": []gin.H{
			{"code": "invalid_token", "detail": "Неверный JWT"},
		}})
		return
	}

	// Запрос в базу данных для получения перемещений по пользователю
	// permissions :=  pg.GetSecretById(userID)
	permissions := pg.GetHandler(userID)

	// Добавление перемещений в заголовок 'X-User-Permissions'
	c.Header("X-User-Permissions", permissions)

	// Проксирование запроса в CreateReverseProxy
	createReverseProxy(AuthService)(c)
}
