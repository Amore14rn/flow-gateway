package handlers

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const secret_key = "django-insecure-test"

func (h *Handler) createJWTHandler(c *gin.Context) {
	// Обработка логики создания JWT
	// В этом примере предполагается, что у тебя есть данные пользователя для создания токена
	// Ты можешь использовать пакет jwt-go для создания токена с необходимыми данными

	// Пример создания токена с некоторыми данными пользователя
	claims := jwt.MapClaims{
		"username": "example_user",
		"role":     "admin",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Получение секретного ключа из базы данных
	secret := secret_key

	// Подписывание токена с использованием секретного ключа
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
