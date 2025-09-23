package nsxt

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	emptyBody     = ""
	serviceWithId = `
{
  "resource_type": "Service",
  "description": "My HTTP",
  "id": "foo",
  "display_name": "foo",
  "path": "/infra/services/my-http",
  "parent_path": "/infra/services/my-http",
  "relative_path": "my-http",
  "service_entries": [
      {
          "resource_type": "L4PortSetServiceEntry",
          "id": "MyHttpEntry",
          "display_name": "MyHttpEntry",
          "path": "/infra/services/my-http/service-entries/MyHttpEntry",
          "parent_path": "/infra/services/my-http",
          "relative_path": "MyHttpEntry",
          "destination_ports": [
              "8080"
          ],
          "l4_protocol": "TCP",
          "_create_user": "admin",
          "_create_time": 1517310677617,
          "_last_modified_user": "admin",
          "_last_modified_time": 1517310677617,
          "_system_owned": false,
          "_protection": "NOT_PROTECTED",
          "_revision": 0
      }
  ],
  "_create_user": "admin",
  "_create_time": 1517310677604,
  "_last_modified_user": "admin",
  "_last_modified_time": 1517310677604,
  "_system_owned": false,
  "_protection": "NOT_PROTECTED",
  "_revision": 0
}
`

	nsxtVersion = `
{
    "node_version": "9.1.0",
    "product_version": "9.1.0"
}
`
)

type mockRoundTripper struct {
	cache []string
	fn    func(req *http.Request) *http.Response
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if (req.Method == http.MethodPatch || req.Method == http.MethodPost || req.Method == http.MethodPut) && req.Body != nil {
		body, _ := io.ReadAll(req.Body)
		m.cache = append(m.cache, string(body))
	}
	return m.fn(req), nil
}

func (m *mockRoundTripper) ClearCache() {
	m.cache = []string{}
}

func constructMockResponse(response string) *mockRoundTripper {

	return &mockRoundTripper{
		cache: []string{},
		fn: func(req *http.Request) *http.Response {
			if req.URL.Path == "/api/v1/node/version" {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(nsxtVersion)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}
			} else if response != "" {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(response)),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}
			} else {
				return &http.Response{
					StatusCode: http.StatusOK,
					// Body:       ioutil.NopCloser(bytes.NewBufferString(serviceWithId)),
					Header: http.Header{"Content-Type": []string{"application/json"}},
				}
			}
		},
	}
}

func newMockProviderClient(response string) (nsxtClients, *mockRoundTripper) {
	mockRoundTripper := constructMockResponse(response)
	mockProviderClient := ConstructMockProviderClient()
	mockProviderClient.NsxtClientConfig.HTTPClient = &http.Client{
		Transport: mockRoundTripper,
	}
	mockProviderClient.PolicyHTTPClient = &http.Client{
		Transport: mockRoundTripper,
	}
	return mockProviderClient, mockRoundTripper
}
