package handler

import (
	"context"
	"net/http"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"lms-system-internship/config"
)

type UpdateUserRequest struct {
	Username        string `json:"username,omitempty"`
	Email           string `json:"email,omitempty"`
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password,omitempty"`
}

func UpdateUser(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	client := gocloak.NewClient(config.GetKeycloakBaseURL())

	tokenRaw, exists := c.Get("access_token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token is missing"})
		return
	}
	token := tokenRaw.(string)

	userInfo, err := client.GetUserInfo(ctx, token, config.GetKeycloakRealm())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	userID := userInfo.Sub

	user := gocloak.User{
		ID: gocloak.StringP(userID),
	}

	if req.Username != "" {
		user.Username = gocloak.StringP(req.Username)
	}
	if req.Email != "" {
		user.Email = gocloak.StringP(req.Email)
	}

	err = client.UpdateUser(ctx, token, config.GetKeycloakRealm(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user data"})
		return
	}

	client.Update
	if req.NewPassword != "" {
		err := client.UpdatePassword(ctx, gocloak.UpdatePasswordRequest{
			Username:    gocloak.StringP(userInfo.PreferredUsername),
			OldPassword: gocloak.StringP(req.CurrentPassword),
			NewPassword: gocloak.StringP(req.NewPassword),
		}, config.KeycloakRealm)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update password"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "User data updated successfully"})
}
