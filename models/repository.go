package models

import "gorm.io/gorm"

// Repository interacts with db
type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateAccountData(ad *AccountData) error {
	res := r.DB.Create(ad)
	if res.Error != nil {
		// Logging?
		return res.Error
	}
	return nil
}

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
