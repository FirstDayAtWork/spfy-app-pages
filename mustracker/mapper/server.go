package mapper

import (
	"io"
	"mustracker/entity"
	"net/http"
)

func RegistrationRequestToAccountData(r *http.Request) (*entity.AccountData, error) {
	bts, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO log this
		return nil, err
	}
	rd := &entity.AccountData{}
	if err = rd.Unmarshal(bts); err != nil {
		return nil, err
	}

	return rd, nil
}
