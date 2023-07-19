package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (h *Handler) refreshJWTHandler(c *gin.Context) {
	// Обработка логики обновления JWT
	oldTokenString := c.GetHeader("Authorization")
	if oldTokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Отсутствует заголовок 'Authorization'"})
		return
	}

	oldToken, err := jwt.Parse(oldTokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверка алгоритма подписи старого токена
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неправильный алгоритм подписи токена: %v", token.Header["alg"])
		}

		// Возвращение секретного ключа
		return []byte(h.secret), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный старый JWT"})
		return
	}

	// Проверка валидности старого токена
	if _, ok := oldToken.Claims.(jwt.Claims); !ok || !oldToken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный старый JWT"})
		return
	}

	// Создание нового токена на основе старого токена с обновленными данными
	claims := oldToken.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Обновление времени истечения срока действия

	newToken := jwt.NewWithClaims(oldToken.Method, claims)

	// Подписывание нового токена с использованием секретного ключа
	newTokenString, err := newToken.SignedString([]byte(h.secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать новый JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"new_token": newTokenString})
}
