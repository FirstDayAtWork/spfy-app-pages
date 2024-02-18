package models

import (
	"time"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

const (
	isValidColumnName = "is_valid"
)

type AccessToken struct {
	JIT           string    `gorm:"column:jit"`
	Token         string    `gorm:"column:token"`
	IsValid       bool      `gorm:"column:is_valid"`
	GrantedToUUID string    `gorm:"column:granted_to_uuid"`
	ExpiresAt     int64     `gorm:"column:expires_at"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func MigrateAccessTokenData(db *gorm.DB) error {
	return db.AutoMigrate(&AccessToken{})
}

type TokenClaims struct {
	UserUUID string `json:"user_uuid,omitempty"`
	Role     int    `json:"role,omitempty"`
	*jwt.StandardClaims
}

type TokenType uint

const (
	RefreshTokenType        TokenType = 0
	AccessTokenType         TokenType = 1
	RefreshTokenValidityHrs           = 24
	AccessTokenValidityHrs            = 1
	RefreshTokenCookieName            = "refresh_token"
	AccessTokenCookieName             = "access_token"
	LoggedInCookieName                = "logged_in"
	TrueString                        = "true"
)

// AuthCheckResult represents results of user authorization check
type AuthCheckResult struct {
	ValidAccess   bool
	ValidRefresh  bool
	ValidRole     bool
	Err           error
	NeedsHandling bool
	AccessClms    *TokenClaims
	RefreshClms   *TokenClaims
}

// Redirect represents a setup for redirect used by AuthHandlingResult
type Redirect struct {
	Path   string
	Status int
}

// AuthHandlingResult represents results of user authorization handling
type AuthHandlingResult struct {
	Redirect     *Redirect
	IsAuthorized bool
}
