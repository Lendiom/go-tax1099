package tax1099

import (
	"context"
	"log/slog"
	"time"
)

type Submit1098sRequest struct {
	TaxYear         string     `json:"taxYear"`
	FormName        string     `json:"formName"`
	ScheduledDate   time.Time  `json:"scheduledDate"`
	IsCorrected     bool       `json:"isCorrected"`
	CouponCode      string     `json:"couponCode"`
	CardReferenceID string     `json:"cardReferenceId"`
	Items           []Item1098 `json:"items"`
}

type Submit1098sResponse struct {
	TraceIdentifier        string `json:"traceIdentifier,omitempty"`
	Message                string `json:"message,omitempty"`
	StatusCode             int    `json:"statusCode,omitempty"`
	OriginalStatusCode     int    `json:"originalStatusCode,omitempty"`
	IsError                bool   `json:"isError,omitempty"`
	ReferenceIDs           []int  `json:"referenceIds,omitempty"`
	PaymentResponseMessage string `json:"paymentResponseMessage,omitempty"`
	TotalCount             int    `json:"totalCount,omitempty"`
}

func (t *tax1099Impl) Submit1098s(ctx context.Context, payload Submit1098sRequest) (Submit1098sResponse, error) {
	slog.InfoContext(ctx, "Submitting the 1098 forms...")

	var res Submit1098sResponse
	if err := t.post(ctx, t.generateFullUrl(UrlPayment, "payment/forms/import/submit/1098"), payload, &res); err != nil {
		return res, err
	}

	slog.InfoContext(ctx, "...1098 forms submitted", "response", res)

	return res, nil
}
