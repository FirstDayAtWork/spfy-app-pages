package mapper

import (
	"io"
	"mustracker/entity"
	"net/http"
)

func RegistrationRequestToRegistrationData(r *http.Request) (*entity.RegistrationData, error) {
	bts, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO log this
		return nil, err
	}
	rd := &entity.RegistrationData{}
	if err = rd.Unmarshal(bts); err != nil {
		return nil, err
	}

	return rd, nil
}
