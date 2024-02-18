package utils

import (
	"net/http"

	"github.com/FirstDayAtWork/mustracker/models"
	"github.com/google/uuid"
)

func NewUUID() string {
	return uuid.NewString()
}

func IsOkCookieSize(c *http.Cookie) bool {
	if c == nil {
		return false
	}
	return len(c.String()) < models.MaxCookieSizeBytes
}
