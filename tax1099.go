package tax1099

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type Tax1099 interface {
	Authorize(ctx context.Context, email, password, appKey string) error
	Validate1098(ctx context.Context, payload Submit1098Request) (Submit1098Response, error)
	Import1098(ctx context.Context, payload Submit1098Request) (Submit1098Response, error)
	Submit1098s(ctx context.Context, payload Submit1098sRequest) (Submit1098sResponse, error)
	DownloadFilledForm(ctx context.Context, payload DownloadFormRequest) ([]byte, error)
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

func New(ctx context.Context, env Environment, username, password, appKey string) (Tax1099, error) {
	c := &http.Client{}
	c.Timeout = 90 * time.Second

	tximpl := &tax1099Impl{
		env:      env,
		username: username,
		password: password,
		appKey:   appKey,
		client:   c,
	}

	return tximpl, tximpl.Authorize(ctx, username, password, appKey)
}

func (t *tax1099Impl) isProduction() bool {
	return t.env == EnvironmentProduction
}

func (t *tax1099Impl) generateFullUrl(urlType UrlType, endpoint string) string {
	var baseUrl string

	switch urlType {
	case UrlMain:
		baseUrl = "https://tax1099api.1099cloud.com/api/v1"

		if t.isProduction() {
			baseUrl = "https://app.tax1099.com/api/v1"
		}
	case Url1098:
		baseUrl = "https://apiforms.1099cloud.com/api/v1"

		if t.isProduction() {
			baseUrl = "https://form1098.tax1099.com/api/v1"
		}
	case UrlPayment:
		baseUrl = "https://apipayment.1099cloud.com/api/v1"

		if t.isProduction() {
			baseUrl = "https://apipayment.tax1099.com/api/v1"
		}
	}

	return fmt.Sprintf("%s/%s", baseUrl, endpoint)
}

func (t *tax1099Impl) post(ctx context.Context, url string, payload, returnValue interface{}) error {
	// Re-authorize if the token has expired, but only if the URL is not the login URL
	if time.Now().After(t.tokenExpiresAt) && !strings.Contains(url, "/login") {
		if err := t.Authorize(ctx, t.username, t.password, t.appKey); err != nil {
			return fmt.Errorf("failed to re-authorize: %v", err)
		}
	}

	slog.InfoContext(ctx, "Tax1099 POST", "url", url)

	body, err := json.Marshal(payload)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to marshal payload", "payload", payload, "error", err)
		return err
	}

	slog.InfoContext(ctx, "Payload", "body", string(body))

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create request", "error", err)
		return err
	}

	if len(body) != 0 {
		req.Header.Add("Content-Type", "application/json")
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to make request", "error", err)
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read response body", "error", err)
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

func (t *tax1099Impl) postForBytes(ctx context.Context, url string, payload interface{}) ([]byte, error) {
	// Re-authorize if the token has expired, but only if the URL is not the login URL
	if time.Now().After(t.tokenExpiresAt) && !strings.Contains(url, "/login") {
		if err := t.Authorize(ctx, t.username, t.password, t.appKey); err != nil {
			return nil, fmt.Errorf("failed to re-authorize: %v", err)
		}
	}

	slog.InfoContext(ctx, "Tax1099 POST", "url", url)

	body, err := json.Marshal(payload)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to marshal payload", "payload", payload, "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Payload", "body", string(body))

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create request", "error", err)
		return nil, err
	}

	if len(body) != 0 {
		req.Header.Add("Content-Type", "application/json")
	}

	req.Header.Add("Accept", "application/pdf")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to make request", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read response body", "error", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d return from %s with body: %s", resp.StatusCode, url, data)
	}

	return data, nil
}
