package handler

import (
	"context"
	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"lms-system-internship/config"
	"net/http"
)

type RegisterRequest struct {
	Username string   `json:"username" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	Roles    []string `json:"roles"`
}

func RegisterUser(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	client := gocloak.NewClient(config.GetKeycloakBaseURL())

	token, err := client.LoginAdmin(ctx, config.GetKeycloakAdmin(), config.GetKeycloakPassword(), config.GetKeycloakRealm())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login as admin"})
		return
	}

	user := gocloak.User{
		Username:      gocloak.StringP(req.Username),
		Email:         gocloak.StringP(req.Email),
		Enabled:       gocloak.BoolP(true),
		EmailVerified: gocloak.BoolP(true),
	}

	userID, err := client.CreateUser(ctx, token.AccessToken, config.GetKeycloakRealm(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Set password
	err = client.SetPassword(ctx, token.AccessToken, userID, config.GetKeycloakRealm(), req.Password, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set password"})
		return
	}

	// Assign roles
	if len(req.Roles) > 0 {
		rolesToAssign := []*gocloak.Role{}
		for _, roleName := range req.Roles {
			role, err := client.GetRealmRole(ctx, token.AccessToken, config.GetKeycloakRealm(), roleName)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role: " + roleName})
				return
			}
			rolesToAssign = append(rolesToAssign, role)
		}

		var realmRoles []gocloak.Role
		for _, r := range rolesToAssign {
			if r != nil {
				realmRoles = append(realmRoles, *r)
			}
		}

		err = client.AddRealmRoleToUser(ctx, token.AccessToken, config.GetKeycloakRealm(), userID, realmRoles)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign roles"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created", "user_id": userID})
}
