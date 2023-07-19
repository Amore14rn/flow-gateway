package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

const AuthService = "http://localhost:8001"

type Handler struct {
	router *gin.Engine
	secret string
	db     *sql.DB
}

func NewHandler(db *sql.DB, secret string) *Handler {
	return &Handler{
		router: gin.Default(),
		secret: secret,
		db:     db,
	}
}

func (h *Handler) RegisterRoutes() {
	authGroup := h.router.Group("/auth")
	{
		authGroup.POST("api/v1/secret/jwt/create/", h.Protected(createReverseProxy(AuthService)))
		authGroup.POST("api/v1/secret/jwt/refresh/", h.Protected(createReverseProxy(AuthService)))
		authGroup.POST("api/v1/secret/jwt/verify/", h.verifyJWTHandler(AuthService))

	}

	h.router.Any("/auth22/*path", h.auth22Handler)
}

func createReverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetURL, _ := url.Parse(strings.Replace(c.Request.RequestURI, "/auth", target, 1))

		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		c.Request.URL.Scheme = targetURL.Scheme
		c.Request.URL.Host = targetURL.Host
		c.Request.URL.Path = c.Param("path")

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func (h *Handler) Protected(handlerFunc func(c *gin.Context)) func(c *gin.Context) {
	return func(c *gin.Context) {

		// вызввать `VerifyJWT` и проверить результат
		isValidToken := h.verifyJWTHandler(c)

		if isValidToken {

			c.Header("X-USER-PERMISSIONs", "some-value")

			handlerFunc(c)
		} else {

			c.JSON(http.StatusUnauthorized, gin.H{
				"type": "client_error",
				"errors": []gin.H{
					{
						"code":   "token_not_valid",
						"detail": "Given token not valid for any token type",
						"attr":   "detail",
					},
				},
			})
		}
	}
}
