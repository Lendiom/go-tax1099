package tax1099

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Tax1099 interface {
	Authorize(email, password, appKey string) error
	Validate1098(payload Submit1098Request) (Submit1098Response, error)
	Import1098(payload Submit1098Request) (Submit1098Response, error)
	Submit1098s(payload Submit1098sRequest) (Submit1098sResponse, error)
}

type tax1099Impl struct {
	env            Environment
	username       string
	password       string
	appKey         string
	token          string
	tokenExpiresAt time.Time

	client *http.Client
}

func New(env Environment, username, password, appKey string) (Tax1099, error) {
	c := &http.Client{}
	c.Timeout = 60 * time.Second

	tximpl := &tax1099Impl{
		env:      env,
		username: username,
		password: password,
		appKey:   appKey,
		client:   c,
	}

	return tximpl, tximpl.Authorize(username, password, appKey)
}

func (t *tax1099Impl) generateFullUrl(urlType UrlType, endpoint string) string {
	var baseUrl string

	switch urlType {
	case UrlMain:
		baseUrl = "https://tax1099api.1099cloud.com/api/v1"

		if t.env == EnvironmentProduction {
			baseUrl = "https://app.tax1099.com/api/v1"
		}
	case Url1098:
		baseUrl = "https://apiform1098.1099cloud.com/api/v1"

		if t.env == EnvironmentProduction {
			baseUrl = "https://form1098.tax1099.com/api/v1"
		}
	case UrlPayment:
		baseUrl = "https://apipayment.1099cloud.com/api/v1"

		if t.env == EnvironmentProduction {
			baseUrl = "https://apipayment.tax1099.com/api/v1"
		}
	}

	return fmt.Sprintf("%s/%s", baseUrl, endpoint)
}

func (t *tax1099Impl) post(url string, payload, returnValue interface{}) error {
	// Re-authorize if the token has expired, but only if the URL is not the login URL
	if time.Now().After(t.tokenExpiresAt) && !strings.Contains(url, "/login") {
		if err := t.Authorize(t.username, t.password, t.appKey); err != nil {
			return fmt.Errorf("failed to re-authorize: %v", err)
		}
	}

	log.Printf("Tax1099 POST %s", url)

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal the payload: %+v", payload)
		return err
	}

	log.Printf("Payload: %s", body)

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
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
		return fmt.Errorf("status code %d return from %s with body: %s", resp.StatusCode, url, data)
	}

	if returnValue == nil {
		return nil
	}

	return json.Unmarshal(data, returnValue)
}
