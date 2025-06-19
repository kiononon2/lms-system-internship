package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

//

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

//

// LoginHandler godoc
// @Summary Аутентификация пользователя
// @Description Получить JWT токены по имени пользователя и паролю
// @Tags auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Данные для входа"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	form := url.Values{}
	form.Set("grant_type", "password")
	form.Set("client_id", "backend-client")
	form.Set("username", req.Username)
	form.Set("password", req.Password)

	keycloakTokenURL := "http://keycloak:8080/realms/lms/protocol/openid-connect/token"

	resp, err := http.Post(
		keycloakTokenURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Keycloak"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		c.JSON(resp.StatusCode, gin.H{"error": "Login failed", "details": string(body)})
		return
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token response"})
		return
	}

	c.JSON(http.StatusOK, tokenResp)
}

//

// RefreshTokenHandler godoc
// @Summary Обновление токена
// @Description Получить новые access и refresh токены по refresh_token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh body RefreshRequest true "Refresh токен"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/refresh [post]
func RefreshTokenHandler(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", "backend-client")
	form.Set("refresh_token", req.RefreshToken)

	keycloakTokenURL := "http://keycloak:8080/realms/lms/protocol/openid-connect/token"

	resp, err := http.Post(
		keycloakTokenURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Keycloak"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		c.JSON(resp.StatusCode, gin.H{"error": "Token refresh failed", "details": string(body)})
		return
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token response"})
		return
	}

	c.JSON(http.StatusOK, tokenResp)
}
