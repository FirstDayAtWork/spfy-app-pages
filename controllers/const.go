package controllers

import "github.com/FirstDayAtWork/mustracker/models"

const (
	RegisterPath    = "/register"
	LoginPath       = "/login"
	HomePath        = "/home"
	AboutPath       = "/about"
	DonatePath      = "/donate"
	LogoutPath      = "/logout"
	AccountPath     = "/account"
	AuthTokenKey    = "AuthToken"
	RefreshTokenKey = "RefreshToken"
)

type accessReqs struct {
	MinRoleNeeded int
}

// Make it account for role as well
var restrictedPaths = map[string]accessReqs{
	AccountPath: {MinRoleNeeded: models.UserRole},
	LogoutPath:  {MinRoleNeeded: models.UserRole},
}
