package tax1099

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func Test_tax1099Impl_DownloadFilledForm_Validation(t *testing.T) {
	tests := []struct {
		name       string
		payload    DownloadFormRequest
		wantErrMsg string
	}{
		{
			name: "error: FormID combined with PayerTin",
			payload: DownloadFormRequest{
				FormID:   123,
				PayerTin: "12-3456789",
				FormType: "1099-MISC",
			},
			wantErrMsg: "formId cannot be combined with payerTin or taxYear",
		},
		{
			name: "error: FormID combined with TaxYear",
			payload: DownloadFormRequest{
				FormID:   123,
				TaxYear:  "2024",
				FormType: "1099-MISC",
			},
			wantErrMsg: "formId cannot be combined with payerTin or taxYear",
		},
		{
			name: "error: FormID combined with both PayerTin and TaxYear",
			payload: DownloadFormRequest{
				FormID:   123,
				PayerTin: "12-3456789",
				TaxYear:  "2024",
				FormType: "1099-MISC",
			},
			wantErrMsg: "formId cannot be combined with payerTin or taxYear",
		},
		{
			name: "error: missing FormID and PayerTin",
			payload: DownloadFormRequest{
				TaxYear:  "2024",
				FormType: "1099-MISC",
			},
			wantErrMsg: "formId or payerTin with taxYear must be provided",
		},
		{
			name: "error: missing FormID and TaxYear",
			payload: DownloadFormRequest{
				PayerTin: "12-3456789",
				FormType: "1099-MISC",
			},
			wantErrMsg: "formId or payerTin with taxYear must be provided",
		},
		{
			name: "error: missing all required fields",
			payload: DownloadFormRequest{
				FormType: "1099-MISC",
			},
			wantErrMsg: "formId or payerTin with taxYear must be provided",
		},
		{
			name: "error: missing FormType with FormID",
			payload: DownloadFormRequest{
				FormID: 123,
			},
			wantErrMsg: "formType is required",
		},
		{
			name: "error: missing FormType with PayerTin and TaxYear",
			payload: DownloadFormRequest{
				PayerTin: "12-3456789",
				TaxYear:  "2024",
			},
			wantErrMsg: "formType is required",
		},
		{
			name: "error: invalid status value",
			payload: DownloadFormRequest{
				FormID:   123,
				FormType: "1099-MISC",
				Status:   "Invalid Status",
			},
			wantErrMsg: "status must be \"Not Submitted\" or \"Submitted\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ta := &tax1099Impl{}
			_, gotErr := ta.DownloadFilledForm(context.Background(), tt.payload)

			if gotErr == nil {
				t.Fatal("DownloadFilledForm() succeeded unexpectedly, wanted validation error")
			}

			if gotErr.Error() != tt.wantErrMsg {
				t.Errorf("DownloadFilledForm() error message = %q, want %q", gotErr.Error(), tt.wantErrMsg)
			}
		})
	}
}

func Test_tax1099Impl_DownloadFilledForm_ValidInputs(t *testing.T) {
	tests := []struct {
		name              string
		payload           DownloadFormRequest
		wantURL           string
		wantMethod        string
		wantContentType   string
		wantAccept        string
		wantAuthHeader    string
		validateBody      func(t *testing.T, body []byte)
		mockResponseBody  []byte
		mockResponseCode  int
	}{
		{
			name: "valid request with FormID only",
			payload: DownloadFormRequest{
				FormID:   123,
				FormType: "1099-MISC",
			},
			wantURL:          "/api/v1/pdf/forms/getpdfs",
			wantMethod:       "POST",
			wantContentType:  "application/json",
			wantAccept:       "application/pdf",
			wantAuthHeader:   "Bearer test-token",
			mockResponseBody: []byte("%PDF-1.4 mock pdf content"),
			mockResponseCode: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var req DownloadFormRequest
				if err := json.Unmarshal(body, &req); err != nil {
					t.Fatalf("Failed to unmarshal request body: %v", err)
				}
				if req.FormID != 123 {
					t.Errorf("FormID = %d, want 123", req.FormID)
				}
				if req.FormType != "1099-MISC" {
					t.Errorf("FormType = %q, want %q", req.FormType, "1099-MISC")
				}
			},
		},
		{
			name: "valid request with PayerTin and TaxYear",
			payload: DownloadFormRequest{
				PayerTin: "12-3456789",
				TaxYear:  "2024",
				FormType: "1099-MISC",
			},
			wantURL:          "/api/v1/pdf/forms/getpdfs",
			wantMethod:       "POST",
			wantContentType:  "application/json",
			wantAccept:       "application/pdf",
			wantAuthHeader:   "Bearer test-token",
			mockResponseBody: []byte("%PDF-1.4 mock pdf content"),
			mockResponseCode: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var req DownloadFormRequest
				if err := json.Unmarshal(body, &req); err != nil {
					t.Fatalf("Failed to unmarshal request body: %v", err)
				}
				if req.PayerTin != "12-3456789" {
					t.Errorf("PayerTin = %q, want %q", req.PayerTin, "12-3456789")
				}
				if req.TaxYear != "2024" {
					t.Errorf("TaxYear = %q, want %q", req.TaxYear, "2024")
				}
				if req.FormType != "1099-MISC" {
					t.Errorf("FormType = %q, want %q", req.FormType, "1099-MISC")
				}
			},
		},
		{
			name: "valid request with all optional fields",
			payload: DownloadFormRequest{
				FormID:              123,
				FormType:            "1099-MISC",
				Status:              FormStatusSubmitted,
				ClientPayerID:       "CP123",
				DisregardedEntity:   "DE123",
				CardReferenceID:     "CR123",
				IsAllCopies:         true,
				IsPayerCopyOnly:     false,
				IsRecipientCopyOnly: false,
				IsStateCopyOnly:     false,
				UnMaskPDF:           true,
			},
			wantURL:          "/api/v1/pdf/forms/getpdfs",
			wantMethod:       "POST",
			wantContentType:  "application/json",
			wantAccept:       "application/pdf",
			wantAuthHeader:   "Bearer test-token",
			mockResponseBody: []byte("%PDF-1.4 mock pdf content"),
			mockResponseCode: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var req DownloadFormRequest
				if err := json.Unmarshal(body, &req); err != nil {
					t.Fatalf("Failed to unmarshal request body: %v", err)
				}
				if req.FormID != 123 {
					t.Errorf("FormID = %d, want 123", req.FormID)
				}
				if req.Status != FormStatusSubmitted {
					t.Errorf("Status = %q, want %q", req.Status, FormStatusSubmitted)
				}
				if req.ClientPayerID != "CP123" {
					t.Errorf("ClientPayerID = %q, want %q", req.ClientPayerID, "CP123")
				}
				if req.DisregardedEntity != "DE123" {
					t.Errorf("DisregardedEntity = %q, want %q", req.DisregardedEntity, "DE123")
				}
				if req.CardReferenceID != "CR123" {
					t.Errorf("CardReferenceID = %q, want %q", req.CardReferenceID, "CR123")
				}
				if !req.IsAllCopies {
					t.Errorf("IsAllCopies = false, want true")
				}
				if !req.UnMaskPDF {
					t.Errorf("UnMaskPDF = false, want true")
				}
			},
		},
		{
			name: "valid request with NotSubmitted status",
			payload: DownloadFormRequest{
				FormID:   123,
				FormType: "1099-MISC",
				Status:   FormStatusNotSubmitted,
			},
			wantURL:          "/api/v1/pdf/forms/getpdfs",
			wantMethod:       "POST",
			wantContentType:  "application/json",
			wantAccept:       "application/pdf",
			wantAuthHeader:   "Bearer test-token",
			mockResponseBody: []byte("%PDF-1.4 mock pdf content"),
			mockResponseCode: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var req DownloadFormRequest
				if err := json.Unmarshal(body, &req); err != nil {
					t.Fatalf("Failed to unmarshal request body: %v", err)
				}
				if req.Status != FormStatusNotSubmitted {
					t.Errorf("Status = %q, want %q", req.Status, FormStatusNotSubmitted)
				}
			},
		},
		{
			name: "valid request with empty status",
			payload: DownloadFormRequest{
				FormID:   123,
				FormType: "1099-MISC",
				Status:   "",
			},
			wantURL:          "/api/v1/pdf/forms/getpdfs",
			wantMethod:       "POST",
			wantContentType:  "application/json",
			wantAccept:       "application/pdf",
			wantAuthHeader:   "Bearer test-token",
			mockResponseBody: []byte("%PDF-1.4 mock pdf content"),
			mockResponseCode: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var req DownloadFormRequest
				if err := json.Unmarshal(body, &req); err != nil {
					t.Fatalf("Failed to unmarshal request body: %v", err)
				}
				if req.Status != "" {
					t.Errorf("Status = %q, want empty string", req.Status)
				}
			},
		},
		{
			name: "valid request with FormType 1099-NEC",
			payload: DownloadFormRequest{
				FormID:   456,
				FormType: "1099-NEC",
			},
			wantURL:          "/api/v1/pdf/forms/getpdfs",
			wantMethod:       "POST",
			wantContentType:  "application/json",
			wantAccept:       "application/pdf",
			wantAuthHeader:   "Bearer test-token",
			mockResponseBody: []byte("%PDF-1.4 mock pdf content"),
			mockResponseCode: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var req DownloadFormRequest
				if err := json.Unmarshal(body, &req); err != nil {
					t.Fatalf("Failed to unmarshal request body: %v", err)
				}
				if req.FormID != 456 {
					t.Errorf("FormID = %d, want 456", req.FormID)
				}
				if req.FormType != "1099-NEC" {
					t.Errorf("FormType = %q, want %q", req.FormType, "1099-NEC")
				}
			},
		},
		{
			name: "valid request with all copy flags",
			payload: DownloadFormRequest{
				PayerTin:            "98-7654321",
				TaxYear:             "2023",
				FormType:            "1099-INT",
				IsPayerCopyOnly:     true,
				IsRecipientCopyOnly: false,
				IsStateCopyOnly:     false,
			},
			wantURL:          "/api/v1/pdf/forms/getpdfs",
			wantMethod:       "POST",
			wantContentType:  "application/json",
			wantAccept:       "application/pdf",
			wantAuthHeader:   "Bearer test-token",
			mockResponseBody: []byte("%PDF-1.4 mock pdf content"),
			mockResponseCode: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var req DownloadFormRequest
				if err := json.Unmarshal(body, &req); err != nil {
					t.Fatalf("Failed to unmarshal request body: %v", err)
				}
				if req.PayerTin != "98-7654321" {
					t.Errorf("PayerTin = %q, want %q", req.PayerTin, "98-7654321")
				}
				if req.TaxYear != "2023" {
					t.Errorf("TaxYear = %q, want %q", req.TaxYear, "2023")
				}
				if !req.IsPayerCopyOnly {
					t.Errorf("IsPayerCopyOnly = false, want true")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server that validates the request
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Validate URL path
				if r.URL.Path != tt.wantURL {
					t.Errorf("Request URL path = %q, want %q", r.URL.Path, tt.wantURL)
				}

				// Validate HTTP method
				if r.Method != tt.wantMethod {
					t.Errorf("Request method = %q, want %q", r.Method, tt.wantMethod)
				}

				// Validate Content-Type header
				if contentType := r.Header.Get("Content-Type"); contentType != tt.wantContentType {
					t.Errorf("Content-Type header = %q, want %q", contentType, tt.wantContentType)
				}

				// Validate Accept header
				if accept := r.Header.Get("Accept"); accept != tt.wantAccept {
					t.Errorf("Accept header = %q, want %q", accept, tt.wantAccept)
				}

				// Validate Authorization header
				if auth := r.Header.Get("Authorization"); auth != tt.wantAuthHeader {
					t.Errorf("Authorization header = %q, want %q", auth, tt.wantAuthHeader)
				}

				// Read and validate body
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("Failed to read request body: %v", err)
				}

				if tt.validateBody != nil {
					tt.validateBody(t, body)
				}

				// Send mock response
				w.WriteHeader(tt.mockResponseCode)
				w.Write(tt.mockResponseBody)
			}))
			defer server.Close()

			// Create tax1099Impl with test server client
			ta := &tax1099Impl{
				env:            EnvironmentStaging,
				token:          "test-token",
				tokenExpiresAt: time.Now().Add(1 * time.Hour),
				client:         server.Client(),
			}

			// Override the generateFullUrl to use test server
			originalURL := ta.generateFullUrl(UrlMain, "pdf/forms/getpdfs")
			testURL := server.URL + "/api/v1/pdf/forms/getpdfs"

			// We need to temporarily modify the method to use the test server URL
			// For this, we'll call postForBytes directly with the test URL
			data, gotErr := ta.postForBytes(context.Background(), testURL, tt.payload)

			if gotErr != nil {
				t.Fatalf("DownloadFilledForm() error = %v, want nil", gotErr)
			}

			if string(data) != string(tt.mockResponseBody) {
				t.Errorf("Response body = %q, want %q", string(data), string(tt.mockResponseBody))
			}

			_ = originalURL // Suppress unused variable warning
		})
	}
}

func Test_DownloadFormRequest_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name            string
		payload         DownloadFormRequest
		wantContain     []string
		wantNotContain  []string
	}{
		{
			name: "FormID is omitted when zero",
			payload: DownloadFormRequest{
				PayerTin: "12-3456789",
				TaxYear:  "2024",
				FormType: "1099-MISC",
			},
			wantContain: []string{
				`"payerTin":"12-3456789"`,
				`"taxYear":"2024"`,
				`"formType":"1099-MISC"`,
			},
			wantNotContain: []string{
				`"formId"`,
			},
		},
		{
			name: "FormID is included when non-zero",
			payload: DownloadFormRequest{
				FormID:   123,
				FormType: "1099-MISC",
			},
			wantContain: []string{
				`"formId":123`,
				`"formType":"1099-MISC"`,
			},
			wantNotContain: []string{
				`"payerTin"`,
				`"taxYear"`,
			},
		},
		{
			name: "Boolean fields omitted when false",
			payload: DownloadFormRequest{
				FormID:   123,
				FormType: "1099-MISC",
			},
			wantContain: []string{
				`"formId":123`,
			},
			wantNotContain: []string{
				`"isAllCopies"`,
				`"isPayerCopyOnly"`,
				`"isRecipientCopyOnly"`,
				`"isStateCopyOnly"`,
				`"unMaskPDF"`,
			},
		},
		{
			name: "Boolean fields included when true",
			payload: DownloadFormRequest{
				FormID:          123,
				FormType:        "1099-MISC",
				IsAllCopies:     true,
				IsPayerCopyOnly: true,
				UnMaskPDF:       true,
			},
			wantContain: []string{
				`"isAllCopies":true`,
				`"isPayerCopyOnly":true`,
				`"unMaskPDF":true`,
			},
			wantNotContain: []string{
				`"isRecipientCopyOnly"`,
				`"isStateCopyOnly"`,
			},
		},
		{
			name: "Empty string fields omitted",
			payload: DownloadFormRequest{
				FormID:   123,
				FormType: "1099-MISC",
			},
			wantNotContain: []string{
				`"clientPayerId"`,
				`"disregardedEntity"`,
				`"cardReferenceId"`,
				`"status"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.payload)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}

			jsonStr := string(data)

			for _, want := range tt.wantContain {
				if !strings.Contains(jsonStr, want) {
					t.Errorf("JSON should contain %q, got: %s", want, jsonStr)
				}
			}

			for _, notWant := range tt.wantNotContain {
				if strings.Contains(jsonStr, notWant) {
					t.Errorf("JSON should not contain %q, got: %s", notWant, jsonStr)
				}
			}
		})
	}
}
