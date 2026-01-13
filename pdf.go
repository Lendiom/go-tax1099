package tax1099

import (
	"context"
	"fmt"
	"log/slog"
)

type FormStatus string

const (
	FormStatusNotSubmitted FormStatus = "Not Submitted"
	FormStatusSubmitted    FormStatus = "Submitted"
)

type DownloadFormRequest struct {
	FormID   uint       `json:"formId"`
	FormType string     `json:"formType"`
	Status   FormStatus `json:"status"`
}

func (t *tax1099Impl) DownloadFilledForm(ctx context.Context, payload DownloadFormRequest) ([]byte, error) {
	if payload.FormID == 0 {
		return nil, fmt.Errorf("formId must be greater than zero")
	}

	if payload.Status != FormStatusNotSubmitted && payload.Status != FormStatusSubmitted {
		return nil, fmt.Errorf("status must be %q or %q", FormStatusNotSubmitted, FormStatusSubmitted)
	}

	slog.InfoContext(ctx, "Downloading filled form PDF...")

	data, err := t.postForBytes(ctx, t.generateFullUrl(UrlMain, "pdf/forms/getpdfs"), payload)
	if err != nil {
		return nil, err
	}

	slog.InfoContext(ctx, "...filled form PDF downloaded")

	return data, nil
}
