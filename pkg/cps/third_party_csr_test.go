package cps

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetChangeThirdPartyCSR(t *testing.T) {
	tests := map[string]struct {
		params           GetChangeRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *ThirdPartyCSRResponse
		withError        func(*testing.T, error)
	}{
		"200 OK": {
			params: GetChangeRequest{
				EnrollmentID: 1,
				ChangeID:     2,
			},
			responseStatus: http.StatusOK,
			responseBody: `
{
    "csrs": [
		{
			"csr": "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
			"keyAlgorithm": "RSA"
		},
		{
			"csr": "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
			"keyAlgorithm": "ECDSA"
		}
	]
}`,
			expectedPath: "/cps/v2/enrollments/1/changes/2/input/info/third-party-csr",
			expectedResponse: &ThirdPartyCSRResponse{
				CSRs: []CertSigningRequest{
					{CSR: "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----", KeyAlgorithm: "RSA"},
					{CSR: "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----", KeyAlgorithm: "ECDSA"},
				},
			},
		},
		"500 internal server error": {
			params: GetChangeRequest{
				EnrollmentID: 1,
				ChangeID:     2,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
  "type": "internal_error",
  "title": "Internal Server Error",
  "detail": "Error making request",
  "status": 500
}`,
			expectedPath: "/cps/v2/enrollments/1/changes/2/input/info/third-party-csr",
			withError: func(t *testing.T, err error) {
				want := &Error{
					Type:       "internal_error",
					Title:      "Internal Server Error",
					Detail:     "Error making request",
					StatusCode: http.StatusInternalServerError,
				}
				assert.True(t, errors.Is(err, want), "want: %s; got: %s", want, err)
			},
		},
		"validation error": {
			params: GetChangeRequest{},
			withError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrStructValidation), "want: %s; got: %s", ErrStructValidation, err)
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "application/vnd.akamai.cps.csr.v2+json", r.Header.Get("Accept"))
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetChangeThirdPartyCSR(context.Background(), test.params)
			if test.withError != nil {
				test.withError(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}
