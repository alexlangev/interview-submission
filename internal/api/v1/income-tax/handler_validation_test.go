package v1

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalculate_ValidationErrors(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		wantStatus int
		wantSubstr string
	}{
		{
			name:       "missing both params",
			query:      "",
			wantStatus: http.StatusBadRequest,
			wantSubstr: "missing query params",
		},
		{
			name:       "invalid year non-numeric",
			query:      "year=abcd&salary=50000",
			wantStatus: http.StatusBadRequest,
			wantSubstr: "invalid or unsupported year",
		},
		{
			name:       "unsupported year",
			query:      "year=1990&salary=50000",
			wantStatus: http.StatusBadRequest,
			wantSubstr: "invalid or unsupported year",
		},
		{
			name:       "invalid salary format",
			query:      "year=2020&salary=abc",
			wantStatus: http.StatusBadRequest,
			wantSubstr: "invalid salary format",
		},
		{
			name:       "negative salary",
			query:      "year=2020&salary=-1000",
			wantStatus: http.StatusBadRequest,
			wantSubstr: "salary must be non-negative",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/calculate?"+testCase.query, nil)
			rr := httptest.NewRecorder()

			handler := Calculate(nil) // nil doesn't matter here
			handler.ServeHTTP(rr, req)

			if rr.Code != testCase.wantStatus {
				t.Fatalf("status = %d, want %d; body: %q", rr.Code, testCase.wantStatus, rr.Body.String())
			}
			if !strings.Contains(rr.Body.String(), testCase.wantSubstr) {
				t.Fatalf("body %q does not contain %q", rr.Body.String(), testCase.wantSubstr)
			}
		})
	}
}
