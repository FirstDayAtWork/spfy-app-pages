package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FirstDayAtWork/mustracker/models"
	"github.com/FirstDayAtWork/mustracker/utils"
	"github.com/golang-jwt/jwt"
)

type Authorizer struct {
	Secret     string
	Issuer     string
	Repository *models.Repository
}

func (auth *Authorizer) NewTokenClaims(uuid string, role int, tt models.TokenType, now *time.Time) (*models.TokenClaims, error) {
	var exp int
	switch tt {
	case models.AccessTokenType:
		exp = models.AccessTokenValidityHrs
	case models.RefreshTokenType:
		exp = models.RefreshTokenValidityHrs
	default:
		return nil, fmt.Errorf("%v token type is not supported", tt)
	}
	return &models.TokenClaims{
		UserUUID: uuid,
		Role:     role,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: now.Add(time.Hour * time.Duration(exp)).Unix(),
			Issuer:    auth.Issuer,
			Id:        utils.NewUUID(),
		},
	}, nil
}

func (auth *Authorizer) ClaimsToSignedString(clms *models.TokenClaims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, clms).SignedString([]byte(auth.Secret))
}

func (auth *Authorizer) JWTTokenToCookie(tkn string, tt models.TokenType) (*http.Cookie, error) {
	cookie := &http.Cookie{
		Value:    tkn,
		HttpOnly: true,
	}
	switch tt {
	case models.AccessTokenType:
		cookie.Name = models.AccessTokenCookieName
	case models.RefreshTokenType:
		cookie.Name = models.RefreshTokenCookieName
	default:
		return nil, fmt.Errorf("%v token type is not supported", tt)
	}
	if !utils.IsOkCookieSize(cookie) {
		return nil, errors.New("cookie size is too large")
	}
	return cookie, nil
}

func (auth *Authorizer) ParseTokenToClaims(tkn string) (*models.TokenClaims, error) {
	clms := &models.TokenClaims{}
	token, err := jwt.ParseWithClaims(tkn, clms,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(auth.Secret), nil
		},
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return clms, nil
}

func (auth *Authorizer) CheckAccessToken(
	r *http.Request,
	state *models.AuthCheckResult,
) {
	accessCookie, err := r.Cookie(models.AccessTokenCookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			state.ValidAccess = false
			return
		}
		state.Err = err
		return
	}
	accessClms, err := auth.ParseTokenToClaims(accessCookie.Value)
	if err != nil {
		state.Err = err
		return
	}
	state.AccessClms = accessClms
	state.ValidAccess = state.HasValidAccess()
}

func (auth *Authorizer) CheckRefreshToken(
	r *http.Request,
	state *models.AuthCheckResult,
) {
	refreshCookie, err := r.Cookie(models.RefreshTokenCookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			state.ValidRefresh = false
			return
		}
		state.Err = err
		return
	}
	refreshClms, err := auth.ParseTokenToClaims(refreshCookie.Value)
	if err != nil {
		state.Err = err
		return
	}
	state.RefreshClms = refreshClms
	refreshTknDB, err := auth.Repository.GetTokenByValue(refreshCookie.Value)
	if err != nil {
		state.Err = err
		return
	}
	log.Printf("got token data from db: %v\n", refreshTknDB)
	state.ValidRefresh = (refreshTknDB.ExpiresAt == refreshClms.ExpiresAt &&
		refreshTknDB.IsValid &&
		state.HasValidRefresh())
}

func (auth *Authorizer) CheckAuthorization(r *http.Request) *models.AuthCheckResult {
	// Get requirements from const
	perms, ok := restrictedPaths[r.URL.Path]
	res := &models.AuthCheckResult{}
	if ok {
		res.NeedsHandling = true
	}
	// Path is restricted, checking user access
	auth.CheckAccessToken(r, res)
	auth.CheckRefreshToken(r, res)
	if res.RefreshClms != nil {
		res.ValidRole = res.RefreshClms.Role >= perms.MinRoleNeeded
	}
	log.Println("checking if auth state has user uuid")
	if userUUID := res.GetUserUUID(); userUUID != models.EmptyString {
		log.Println("auth state has user uuid")
		user, err := auth.Repository.GetAccountDataByUUID(userUUID)
		if err == nil {
			res.User = user
		}
	}
	log.Printf("auth check pre return: %v\n", res)
	return res
}

func (auth *Authorizer) handleAuthStateError(
	authState *models.AuthCheckResult,
	res *models.AuthHandlingResult,
) *models.AuthHandlingResult {
	log.Printf("error conducting auth check: %v", authState.Err)
	res.Redirect = &models.Redirect{
		Path:   LoginPath,
		Status: http.StatusInternalServerError,
	}
	return res
}

func (auth *Authorizer) handleMissingRefreshResult(
	authState *models.AuthCheckResult,
	res *models.AuthHandlingResult,
) *models.AuthHandlingResult {
	log.Println("user has access, but no refresh")
	authState.AccessClms.ExpiresAt = -1 // TODO make this a method
	if auth.Repository.InvalidateUserTokens(authState.AccessClms.UserUUID) != nil {
		log.Printf(
			"error invalidating refresh token for user %s\n",
			authState.AccessClms.UserUUID,
		)
	}
	res.Redirect = &models.Redirect{
		Path:   LoginPath,
		Status: http.StatusFound,
	}
	return res
}

func (auth *Authorizer) handleMissingAccessResult(
	authState *models.AuthCheckResult,
	res *models.AuthHandlingResult,
	path string,
) *http.Cookie {
	log.Println("access token expired, but refresh is in place")
	now := time.Now()
	accessClms, err := auth.NewTokenClaims(
		authState.RefreshClms.UserUUID, authState.RefreshClms.Role,
		models.AccessTokenType, &now,
	)
	if err != nil {
		log.Printf("access token generation error: %+v\n", err)
		res.Redirect = &models.Redirect{
			Path:   LoginPath,
			Status: http.StatusInternalServerError,
		}
		return nil
	}
	accessTkn, err := auth.ClaimsToSignedString(accessClms)
	if err != nil {
		log.Printf("access token generation error: %+v\n", err)
		res.Redirect = &models.Redirect{
			Path:   LoginPath,
			Status: http.StatusInternalServerError,
		}
		return nil
	}
	// cookie creation
	accessCookie, err := auth.JWTTokenToCookie(accessTkn, models.AccessTokenType)
	if err != nil {
		log.Printf("access cookie generation error: %+v\n", err)
		res.Redirect = &models.Redirect{
			Path:   LoginPath,
			Status: http.StatusInternalServerError,
		}
		return nil
	}
	res.Redirect = &models.Redirect{
		Path:   path,
		Status: http.StatusTemporaryRedirect,
	}
	// Updating access token in auth state
	authState.AccessClms = accessClms
	return accessCookie
}

func (auth *Authorizer) handleInvalidAuthResult(
	authState *models.AuthCheckResult,
	res *models.AuthHandlingResult,
) *models.AuthHandlingResult {
	log.Println("access and refresh are both invalid")
	// revoke refresh in db if possible
	if userUUID := authState.GetUserUUID(); userUUID != models.EmptyString {
		if err := auth.Repository.InvalidateUserTokens(userUUID); err != nil {
			log.Printf("error revoking refresh token for %s\n", userUUID)
		}
	}
	res.Redirect = &models.Redirect{
		Path:   LoginPath,
		Status: http.StatusFound,
	}
	return res
}

func (auth *Authorizer) GrantNewToken(
	user *models.AccountData,
	tt models.TokenType,
	now *time.Time,
) (*http.Cookie, error) {
	clms, err := auth.NewTokenClaims(user.UUID, user.Role, tt, now)
	if err != nil {
		return nil, fmt.Errorf("(%d token type): claims creation failed: %v", tt, err)
	}
	tkn, err := auth.ClaimsToSignedString(clms)
	if err != nil {
		return nil, fmt.Errorf("(%d token type): claims-to-string conversion failed: %v", tt, err)
	}
	cookie, err := auth.JWTTokenToCookie(tkn, tt)
	if err != nil {
		return nil, fmt.Errorf("(%d token type): string-to-cookie conversion failed: %v", tt, err)
	}
	if tt == models.AccessTokenType {
		return cookie, err
	}
	// Refresh tokens require some extra work to be granted
	if err = auth.Repository.InvalidateUserTokens(user.UUID); err != nil {
		return nil, fmt.Errorf("(%d token type): old token invalidation failed: %v", tt, err)
	}
	if err = auth.Repository.StoreToken(tkn, clms); err != nil {
		return nil, fmt.Errorf("(%d token type): recording token to DB failed: %v", tt, err)
	}
	return cookie, err
}
