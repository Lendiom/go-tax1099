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

const component = "go-tax1099"

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

func New(ctx context.Context, env Environment, username, password, appKey string, timeout time.Duration) (Tax1099, error) {
	c := &http.Client{}
	c.Timeout = timeout

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

func (t *tax1099Impl) post(ctx context.Context, op, url string, payload, returnValue interface{}) error {
	// Re-authorize if the token has expired, but only if the URL is not the login URL
	if time.Now().After(t.tokenExpiresAt) && !strings.Contains(url, "/login") {
		if err := t.Authorize(ctx, t.username, t.password, t.appKey); err != nil {
			return fmt.Errorf("failed to re-authorize: %v", err)
		}
	}

	slog.InfoContext(ctx, "Tax1099 POST",
		slog.String("component", component),
		slog.String("op", op),
		slog.String("url", url),
	)

	body, err := json.Marshal(payload)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to marshal payload",
			slog.String("component", component),
			slog.String("op", op),
			slog.Any("payload", payload),
			slog.Any("error", err),
		)
		return err
	}

	// Payload bodies for tax forms include TINs, recipient names, addresses,
	// and dollar amounts. Keep the verbatim body at debug-only so production
	// INFO output stays free of PII; consumers that need the raw payload can
	// turn the package's slog level up to debug for a single call.
	slog.DebugContext(ctx, "Payload",
		slog.String("component", component),
		slog.String("op", op),
		slog.String("body", string(body)),
	)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create request",
			slog.String("component", component),
			slog.String("op", op),
			slog.Any("error", err),
		)
		return err
	}

	if len(body) != 0 {
		req.Header.Add("Content-Type", "application/json")
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to make request",
			slog.String("component", component),
			slog.String("op", op),
			slog.Any("error", err),
		)
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read response body",
			slog.String("component", component),
			slog.String("op", op),
			slog.Any("error", err),
		)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		slog.ErrorContext(ctx, "tax1099 request returned non-200",
			slog.String("component", component),
			slog.String("op", op),
			slog.String("url", url),
			slog.Int("status_code", resp.StatusCode),
			slog.String("body", string(data)),
		)
		return fmt.Errorf("status code %d return from %s with body: %s", resp.StatusCode, url, data)
	}

	if returnValue == nil {
		return nil
	}

	return json.Unmarshal(data, returnValue)
}

func (t *tax1099Impl) postForBytes(ctx context.Context, op, url string, payload interface{}) ([]byte, error) {
	// Re-authorize if the token has expired, but only if the URL is not the login URL
	if time.Now().After(t.tokenExpiresAt) && !strings.Contains(url, "/login") {
		if err := t.Authorize(ctx, t.username, t.password, t.appKey); err != nil {
			return nil, fmt.Errorf("failed to re-authorize: %v", err)
		}
	}

	slog.InfoContext(ctx, "Tax1099 POST",
		slog.String("component", component),
		slog.String("op", op),
		slog.String("url", url),
	)

	body, err := json.Marshal(payload)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to marshal payload",
			slog.String("component", component),
			slog.String("op", op),
			slog.Any("payload", payload),
			slog.Any("error", err),
		)
		return nil, err
	}

	// See post() for why this is debug-only.
	slog.DebugContext(ctx, "Payload",
		slog.String("component", component),
		slog.String("op", op),
		slog.String("body", string(body)),
	)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create request",
			slog.String("component", component),
			slog.String("op", op),
			slog.Any("error", err),
		)
		return nil, err
	}

	if len(body) != 0 {
		req.Header.Add("Content-Type", "application/json")
	}

	req.Header.Add("Accept", "application/pdf")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to make request",
			slog.String("component", component),
			slog.String("op", op),
			slog.Any("error", err),
		)
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read response body",
			slog.String("component", component),
			slog.String("op", op),
			slog.Any("error", err),
		)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		slog.ErrorContext(ctx, "tax1099 request returned non-200",
			slog.String("component", component),
			slog.String("op", op),
			slog.String("url", url),
			slog.Int("status_code", resp.StatusCode),
			slog.String("body", string(data)),
		)
		return nil, fmt.Errorf("status code %d return from %s with body: %s", resp.StatusCode, url, data)
	}

	// The provider can return a 200 with a JSON error envelope; without this check
	// those bytes would be passed along as if they were the requested PDF.
	if !bytes.HasPrefix(data, []byte("%PDF")) {
		slog.ErrorContext(ctx, "tax1099 request returned a non-PDF body",
			slog.String("component", component),
			slog.String("op", op),
			slog.String("url", url),
			slog.String("body", truncateForError(data)),
		)
		return nil, fmt.Errorf("response from %s is not a PDF, body: %s", url, truncateForError(data))
	}

	return data, nil
}

// truncateForError limits a response body to a readable length for error messages.
func truncateForError(data []byte) string {
	const maxLen = 200
	if len(data) <= maxLen {
		return string(data)
	}

	return string(data[:maxLen]) + "..."
}
