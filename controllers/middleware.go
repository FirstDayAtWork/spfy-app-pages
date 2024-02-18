package controllers

import (
	"net/http"

	"github.com/FirstDayAtWork/mustracker/models"
)

// func AuthHandler(next http.Handler) http.Handler {
//  fn := func(w http.ResponseWriter, r *http.Request) {
//   switch r.URL.Path {
//   case AccountPath, LogoutPath:
//    log.Printf("Sensitive path %s hit, need auth")
//   default:
//    log.Printf("Not a sensitive path %s hit, no auth needed")
//   }
//  }
// }

// Revokes both auth and refresh cookikes
// func nullifyTokenCookies(w *http.ResponseWriter, r *http.Request) {
//  // Reset client's auth cookie
//  http.SetCookie(
//   *w,
//   &http.Cookie{
//    Name:     AuthTokenKey,
//    Value:    "",
//    Expires:  time.Now().Add(-1000 * time.Hour),
//    HttpOnly: true,
//   },
//  )
//  // Reset client's refresh cookie
//  http.SetCookie(
//   *w,
//   &http.Cookie{
//    Name:     RefreshTokenKey,
//    Value:    "",
//    Expires:  time.Now().Add(-1000 * time.Hour),
//    HttpOnly: true,
//   },
//  )

//  // Invalidate refresh cookie
//  refreshCookie, refreshErr := r.Cookie(RefreshTokenKey)
//  if errors.Is(refreshErr, http.ErrNoCookie) {
//   // No cookie == no problem
//   return
//  } else if refreshErr != nil {
//   log.Panic("panic: %+v", refreshErr)
//   http.Error(*w, http.StatusText(500), 500)
//  }
//  auth.RevokeRefreshToken(refreshCookie.Value)
// }

// Grants auth and refresh tokens
func setAuthAndRefreshCookies(w *http.ResponseWriter, authToken, refreshToken string) {
	http.SetCookie(
		*w,
		&http.Cookie{
			Name:     AuthTokenKey,
			Value:    authToken,
			HttpOnly: true,
		},
	)
	http.SetCookie(
		*w,
		&http.Cookie{
			Name:     RefreshTokenKey,
			Value:    refreshToken,
			HttpOnly: true,
		},
	)
}

// Fetches CSFR from form value or header
func grabCSFRFromRequest(r *http.Request) string {
	CSFR := r.FormValue(CSRFKey)
	if CSFR != models.EmptyString {
		return CSFR
	} else {
		return r.Header.Get(CSRFKey)
	}
}
