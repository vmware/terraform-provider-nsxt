package nsxt

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
	api "github.com/vmware/go-vmware-nsxt"
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

func constructMockProviderClient() nsxtClients {
	commonConfig := commonProviderConfig{
		RemoteAuth:             false,
		ToleratePartialSuccess: false,
		MaxRetries:             2,
		MinRetryInterval:       0,
		MaxRetryInterval:       0,
		RetryStatusCodes:       []int{404, 400},
		Username:               "username",
		Password:               "password",
	}

	nsxtClient := nsxtClients{
		CommonConfig: commonConfig,
	}
	nsxtClient.NsxtClientConfig = &api.Configuration{
		BasePath:   "/api/v1",
		Scheme:     "https",
		UserAgent:  "terraform-provider-nsxt",
		UserName:   "username",
		Password:   "password",
		RemoteAuth: true,
		Insecure:   true,
	}

	return nsxtClient
}

func newMockProviderClient(response string) (nsxtClients, *mockRoundTripper) {
	mockRoundTripper := constructMockResponse(response)
	mockProviderClient := constructMockProviderClient()
	mockProviderClient.NsxtClientConfig.HTTPClient = &http.Client{
		Transport: mockRoundTripper,
	}
	mockProviderClient.PolicyHTTPClient = &http.Client{
		Transport: mockRoundTripper,
	}
	return mockProviderClient, mockRoundTripper
}

func getGjsonString(input, jqpath string) string {
	return gjson.Get(input, jqpath).String()
}
