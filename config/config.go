package config

import "os"

func GetKeycloakBaseURL() string  { return os.Getenv("KEYCLOAK_BASE_URL") }
func GetKeycloakRealm() string    { return os.Getenv("KEYCLOAK_REALM") }
func GetKeycloakAdmin() string    { return os.Getenv("KEYCLOAK_ADMIN") }
func GetKeycloakPassword() string { return os.Getenv("KEYCLOAK_PASSWORD") }
