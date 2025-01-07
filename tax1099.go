package tax1099

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Tax1099 interface {
	Authorize(email, password, appKey string) error
}

type tax1099Impl struct {
	baseAPI        string
	username       string
	password       string
	appKey         string
	token          string
	tokenExpiresAt time.Time

	client *http.Client
}

func New(baseAPI, username, password, appKey string) (Tax1099, error) {
	c := &http.Client{}
	c.Timeout = 60 * time.Second

	tximpl := &tax1099Impl{
		baseAPI:  baseAPI,
		username: username,
		password: password,
		appKey:   appKey,
		client:   c,
	}

	return tximpl, tximpl.Authorize(username, password, appKey)
}

func (t *tax1099Impl) generateFullUrl(endpoint string) string {
	return fmt.Sprintf("%s/%s", t.baseAPI, endpoint)
}

func (t *tax1099Impl) post(endpoint string, payload, returnValue interface{}) error {
	fullUrl := t.generateFullUrl(endpoint)

	// Re-authorize if the token has expired
	if time.Now().After(t.tokenExpiresAt) {
		if err := t.Authorize(t.username, t.password, t.appKey); err != nil {
			return fmt.Errorf("failed to re-authorize: %v", err)
		}
	}

	log.Printf("Tax1099 POST %s: %+v", fullUrl, payload)

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal the payload: %+v", payload)
		return err
	}

	req, err := http.NewRequest("POST", fullUrl, bytes.NewReader(body))
	if err != nil {
		log.Printf("failed to make the request: %v", err)
		return err
	}

	if len(body) != 0 {
		req.Header.Add("Content-Type", "application/json")
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("failed to make the request: %v", err)
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read the response body: %v", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code %d return from %s with body: %s", resp.StatusCode, fullUrl, data)
	}

	if returnValue == nil {
		return nil
	}

	return json.Unmarshal(data, returnValue)
}
