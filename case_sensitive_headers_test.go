package traefik_middleware_case_sensitive_headers

import (
	"net/http"
	"testing"
)

func TestHeaderConfig(t *testing.T) {
	headers := &http.Header{}

	headers.Add("X-To-Remove-Header", "X-To-Remove-Header-Value")
	headers.Add("Authorization", "token")
	headers.Add("X-Client-Cert", "certificate-data")

	rewriteHeaders(headers, &headerConfig{
		AddHeaders: []*addHeaderConfig{
			createAddHeaderConfig("X-To-Add-Header-1", "X-To-Add-Header-1-Value"),
			createAddHeaderConfig("X-To-Add-Header-2", "X-To-Add-Header-2-Value"),
		},
		RemoveHeaders: createRemoveHeaderConfig([]string{"X-To-Remove-Header"}),
		ModifyHeaders: []*modifyHeaderConfig{
			createModifyHeaderConfig("Authorization", "X-Auth", "Bearer ", ";", true, true),
			createModifyHeaderConfig("X-Client-Cert", "SSL_CLIENT_CERT", "-----BEGIN CERTIFICATE-----", "-----END CERTIFICATE-----", false, true),
		}})

	assertHeader(t, headers, "X-To-Add-Header-1", "X-To-Add-Header-1-Value")
	assertHeader(t, headers, "X-To-Add-Header-2", "X-To-Add-Header-2-Value")

	assertHeaderIsAbsent(t, headers, "X-To-Remove-Header")

	assertHeaderIsAbsent(t, headers, "Authorization")

	assertHeader(t, headers, "X-Client-Cert", "certificate-data")
	assertHeader(t, headers, "SSL_CLIENT_CERT", "-----BEGIN CERTIFICATE-----certificate-data-----END CERTIFICATE-----")
}

func assertHeader(t *testing.T, headers *http.Header, headerName, headerValue string) {
	t.Helper()
	if ((*headers)[headerName])[0] != headerValue {
		t.Errorf("\ninvalid header value for header:\n %s: %s", headerName, ((*headers)[headerName])[0])
	}
}

func assertHeaderIsAbsent(t *testing.T, headers *http.Header, headerName string) {
	t.Helper()
	if ((*headers)[headerName]) != nil {
		t.Errorf("\nheader is present:\n%s: %s", headerName, ((*headers)[headerName])[0])
	}
}
