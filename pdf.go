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
	FormID              uint       `json:"formId,omitempty"`
	FormType            string     `json:"formType"`
	Status              FormStatus `json:"status,omitempty"`
	ClientPayerID       string     `json:"clientPayerId,omitempty"`
	PayerTin            string     `json:"payerTin,omitempty"`
	TaxYear             string     `json:"taxYear,omitempty"`
	DisregardedEntity   string     `json:"disregardedEntity,omitempty"`
	CardReferenceID     string     `json:"cardReferenceId,omitempty"`
	IsAllCopies         bool       `json:"isAllCopies,omitempty"`
	IsPayerCopyOnly     bool       `json:"isPayerCopyOnly,omitempty"`
	IsRecipientCopyOnly bool       `json:"isRecipientCopyOnly,omitempty"`
	IsStateCopyOnly     bool       `json:"isStateCopyOnly,omitempty"`
	UnMaskPDF           bool       `json:"unMaskPDF,omitempty"`
}

// DownloadFilledForm downloads a filled 1099 form PDF based on the provided criteria.
// Either FormID or the combination of PayerTin and TaxYear must be provided.
func (t *tax1099Impl) DownloadFilledForm(ctx context.Context, payload DownloadFormRequest) ([]byte, error) {
	if payload.FormID > 0 {
		if payload.PayerTin != "" || payload.TaxYear != "" {
			return nil, fmt.Errorf("formId cannot be combined with payerTin or taxYear")
		}
	} else if payload.PayerTin == "" || payload.TaxYear == "" {
		return nil, fmt.Errorf("formId or payerTin with taxYear must be provided")
	}

	if payload.FormType == "" {
		return nil, fmt.Errorf("formType is required")
	}

	if payload.Status != "" && payload.Status != FormStatusNotSubmitted && payload.Status != FormStatusSubmitted {
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
