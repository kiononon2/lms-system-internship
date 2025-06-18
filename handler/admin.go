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

type UpdateUserRequest struct {
	Email     string `json:"email" binding:"omitempty,email"`
	FirstName string `json:"first_name" binding:"omitempty"`
	LastName  string `json:"last_name" binding:"omitempty"`
	Password  string `json:"password" binding:"omitempty,min=6"`
}

type UpdateUserRolesRequest struct {
	UserID   string   `json:"user_id" binding:"required"`
	NewRoles []string `json:"new_roles" binding:"required"` // Роли, которые нужно оставить
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

func UpdateUserProfile(c *gin.Context) {
	usernameRaw, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No username in context"})
		return
	}
	username := usernameRaw.(string)

	var req UpdateUserRequest
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

	users, err := client.GetUsers(ctx, token.AccessToken, config.GetKeycloakRealm(), gocloak.GetUsersParams{
		Username: &username,
	})
	if err != nil || len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user := users[0]

	if req.Email != "" {
		user.Email = &req.Email
	}
	if req.FirstName != "" {
		user.FirstName = &req.FirstName
	}
	if req.LastName != "" {
		user.LastName = &req.LastName
	}

	err = client.UpdateUser(ctx, token.AccessToken, config.GetKeycloakRealm(), *user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	if req.Password != "" {
		err := client.SetPassword(ctx, token.AccessToken, *user.ID, config.GetKeycloakRealm(), req.Password, false)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set password"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "User profile updated"})
}

func UpdateUserRolesHandler(c *gin.Context) {
	var req UpdateUserRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	client := gocloak.NewClient(config.GetKeycloakBaseURL())

	token, err := client.LoginAdmin(ctx, config.GetKeycloakAdmin(), config.GetKeycloakPassword(), config.GetKeycloakRealm())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin login failed"})
		return
	}

	// Получаем все текущие роли пользователя
	currentRoles, err := client.GetRealmRolesByUserID(ctx, token.AccessToken, config.GetKeycloakRealm(), req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get current roles"})
		return
	}

	// Получаем список ролей, которые админ хочет оставить
	var newRolesObjs []gocloak.Role
	for _, roleName := range req.NewRoles {
		role, err := client.GetRealmRole(ctx, token.AccessToken, config.GetKeycloakRealm(), roleName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role: " + roleName})
			return
		}
		newRolesObjs = append(newRolesObjs, *role)
	}

	// Определяем роли, которые нужно удалить
	var rolesToRemove []gocloak.Role
	newRolesMap := map[string]bool{}
	for _, r := range req.NewRoles {
		newRolesMap[r] = true
	}
	for _, r := range currentRoles {
		if r == nil || r.Name == nil {
			continue
		}
		if _, ok := newRolesMap[*r.Name]; !ok {
			rolesToRemove = append(rolesToRemove, *r)
		}
	}

	// Удаляем роли, которых нет в списке новых
	if len(rolesToRemove) > 0 {
		err = client.DeleteRealmRoleFromUser(ctx, token.AccessToken, config.GetKeycloakRealm(), req.UserID, rolesToRemove)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove old roles"})
			return
		}
	}

	// Добавляем новые роли
	if len(newRolesObjs) > 0 {
		err = client.AddRealmRoleToUser(ctx, token.AccessToken, config.GetKeycloakRealm(), req.UserID, newRolesObjs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add new roles"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "User roles updated successfully"})
}
