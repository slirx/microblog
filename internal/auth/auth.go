package auth

import (
	"time"
)

// LoginRequest represents fields of login request.
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// LoginResponse represents fields of login response.
type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

// Config represents configuration fields for auth service.
type Config struct {
	Secret          string            // JWT secret
	InternalSecrets map[string]string // JWT secrets for microservices (service-to-service communication)
}

// Auth represents fields for columns in auth table.
type Auth struct {
	ID        int
	UserID    int
	Email     string
	Login     string
	Password  string
	CreatedAt time.Time
}

// InternalLoginRequest represents fields of login request.
type InternalLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// InternalLoginResponse represents fields of login response.
type InternalLoginResponse struct {
	AccessToken string `json:"access_token"`
}

// InternalAuth represents fields for columns in internal_auth table.
type InternalAuth struct {
	ID          int
	ServiceName string
	Password    string
	CreatedAt   time.Time
}

// AdminLoginRequest represents fields of login request.
type AdminLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// AdminLoginResponse represents fields of login response.
type AdminLoginResponse struct {
	AccessToken string `json:"access_token"`
}
