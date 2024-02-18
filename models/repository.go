package models

import (
	"log"

	"gorm.io/gorm"
)

// Repository interacts with db
type Repository struct {
	DB *gorm.DB
}

// CreateAccountData addes a new row to a table with accounts.
func (r *Repository) CreateAccountData(ad *AccountData) error {
	res := r.DB.Create(ad)
	if res.Error != nil {
		// Logging?
		return res.Error
	}
	return nil
}

// GetAccountDataByUserame fetches account data with a given username.
func (r *Repository) GetAccountDataByUserame(username string) (*AccountData, error) {
	// Create a dummy struct for query filters
	resultData := &AccountData{}
	res := r.DB.Where(
		&AccountData{Username: username},
	).First(resultData)

	if res.Error != nil {
		return nil, res.Error
	}
	return resultData, nil
}

// Store token creates a token entry in DB
func (r *Repository) StoreToken(token string, clms *TokenClaims) error {
	res := r.DB.Create(
		&AccessToken{
			IsValid:       true,
			Token:         token,
			GrantedToUUID: clms.UserUUID,
			ExpiresAt:     clms.ExpiresAt,
			JIT:           clms.Id,
		},
	)
	if res.Error != nil {
		log.Printf("error storing token: %+v", res.Error)
		return res.Error
	}
	return nil
}

// InvalidateToken sets a token's IsValid to false
func (r *Repository) InvalidateToken(jit string) error {
	token, err := r.GetTokenByValue(jit)
	if err != nil {
		return err
	}
	if res := r.DB.Model(&token).Updates(AccessToken{IsValid: false}); res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *Repository) GetTokenByValue(token string) (*AccessToken, error) {
	log.Printf("fetching token by value %s\n", token)
	return r.getTokenByCondition(&AccessToken{Token: token})
}

func (r *Repository) getTokenByCondition(conditions *AccessToken) (*AccessToken, error) {
	rows, err := r.queryDataByCondition(conditions)
	if err != nil {
		return nil, err
	}
	resultData := &AccessToken{}
	if res := rows.First(resultData); res.Error != nil {
		return nil, err
	}
	log.Printf("result data after fetching first row: %v\n", resultData)
	return resultData, nil
}

func (r *Repository) queryDataByCondition(conditions interface{}) (*gorm.DB, error) {
	res := r.DB.Where(conditions)
	if res.Error != nil {
		return nil, res.Error
	}
	log.Printf("fetched rows from db: %v\n", res)
	return res, nil
}

// InvalidateUserTokens sets all tokens of a user to false
func (r *Repository) InvalidateUserTokens(uuid string) error {
	tokenRows, err := r.queryDataByCondition(&AccessToken{GrantedToUUID: uuid})
	if err != nil {
		return err
	}
	// Doing a select because GORM does not write empty values when using structs
	return tokenRows.Select(isValidColumnName).
		Updates(&AccessToken{IsValid: false}).
		Error
}
