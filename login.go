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
	const op = "tax1099.authorize"

	slog.InfoContext(ctx, "Authorizing...",
		slog.String("component", component),
		slog.String("op", op),
	)

	var res loginResponse
	if err := t.post(ctx, op, t.generateFullUrl(UrlMain, "login"), loginRequest{Email: email, Password: password, AppKey: appKey}, &res); err != nil {
		return err
	}

	if res.SessionID == "" {
		return ErrBadLogin
	}

	t.token = res.SessionID
	t.tokenExpiresAt = time.Now().Add(55 * time.Minute) // 5 minutes before the token expires

	slog.InfoContext(ctx, "...authorization complete",
		slog.String("component", component),
		slog.String("op", op),
	)

	return nil
}
