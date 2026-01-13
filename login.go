package tax1099

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

var ErrBadLogin = errors.New("bad login")

type loginRequest struct {
	Email    string `json:"login"`
	Password string `json:"password"`
	AppKey   string `json:"appKey"`
}

type loginResponse struct {
	SessionID          string `json:"sessionId,omitempty"`
	ValidationMessages []any  `json:"validationMessages,omitempty"`
}

func (t *tax1099Impl) Authorize(ctx context.Context, email, password, appKey string) error {
	slog.InfoContext(ctx, "Authorizing...")

	var res loginResponse
	if err := t.post(ctx, t.generateFullUrl(UrlMain, "login"), loginRequest{Email: email, Password: password, AppKey: appKey}, &res); err != nil {
		return err
	}

	if res.SessionID == "" {
		return ErrBadLogin
	}

	t.token = res.SessionID
	t.tokenExpiresAt = time.Now().Add(55 * time.Minute) // 5 minutes before the token expires

	slog.InfoContext(ctx, "...authorization complete")

	return nil
}
