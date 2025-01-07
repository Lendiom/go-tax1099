package tax1099

import (
	"errors"
	"log"
	"time"
)

var ErrBadLogin = errors.New("bad login")

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	AppKey   string `json:"appKey"`
}

type loginResponse struct {
	SessionID          string `json:"sessionId,omitempty"`
	ValidationMessages []any  `json:"validationMessages,omitempty"`
}

func (t *tax1099Impl) Authorize(email, password, appKey string) error {
	log.Println("Authorizing...")

	var res loginResponse
	if err := t.post("/login", loginRequest{Email: email, Password: password, AppKey: appKey}, &res); err != nil {
		return err
	}

	if res.SessionID == "" {
		return ErrBadLogin
	}

	t.token = res.SessionID
	t.tokenExpiresAt = time.Now().Add(55 * time.Minute) // 5 minutes before the token expires

	log.Println("...authorization complete")

	return nil
}
