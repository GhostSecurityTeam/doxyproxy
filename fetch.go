package doxyproxy

import (
	"encoding/json"
	"net/http"

	"github.com/getlantern/errors"
)

//Response holds a response from the API server
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

var (
	//ErrAPIResponseSuccessFalse is returned when the response was false
	ErrAPIResponseSuccessFalse = errors.New("api response success false")
	//ErrClientNotFound is returned when the client was not found on the proxy server
	ErrClientNotFound = errors.New("client not found")
)

//Fetch will query the proxy API for the real IP
func (entry IPEntry) Fetch() (string, error) {

	req, err := http.NewRequest("GET", entry.proxy.API+entry.ID+"/ip", nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", entry.proxy.Key)

	resp, err := entry.proxy.Client.Do(req)
	if err != nil {
		return "", err
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	if response.Success == false {
		if resp.StatusCode == 404 {
			return "", ErrClientNotFound
		}

		return "", ErrAPIResponseSuccessFalse
	}

	return response.Message, nil
}
