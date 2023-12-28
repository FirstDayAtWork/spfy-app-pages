package simpleclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mustracker/entity"
	"net/http"
)

func Get() {
	r, err := http.Get("http://localhost:2228/test_get")
	if err != nil {
		fmt.Printf("Get got an error: %+v", err)
		return
	}
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Reading Get resp to string got an error: %v\n", err)
		return
	}

	fmt.Printf("Get response: %s\n", string(data))
}

func Post() {
	bodyData := entity.SimpleJson{
		Username: "makarik",
		Email:    "agent228@gmail.com",
		Password: "unhashed!",
	}

	jsonBody, err := json.Marshal(bodyData)
	if err != nil {
		fmt.Printf("Marshalling error: %v\n", err)
		return
	}
	r, err := http.Post(
		"http://localhost:2228/test_post",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		fmt.Printf("Post got an error: %+v", err)
		return
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Reading Get resp to string got an error: %v\n", err)
		return
	}
	fmt.Printf("Post response: %s\n", string(data))

}
