package models

import "time"

const (
	emailRegex            = `\b[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}\b`
	maxPasswordLen        = 72
	EmptyString           = ""
	UserRole              = "user"
	AdminRole             = "admin"
	RefreshTokenValidTime = time.Hour * 72
	AuthTokenValidTime    = time.Minute * 15
	CSFRInputLen          = 32
)
