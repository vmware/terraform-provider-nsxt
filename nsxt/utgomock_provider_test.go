//go:build unittest

// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// To generate the mock for this test, run:
// mockgen -destination=mocks/nsx/LicensesClient.go -package=mocks -source=<local path>/vsphere-automation-sdk-go/services/nsxt-mp/nsx/LicensesClient.go LicensesClient

package nsxt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx"
	nsxModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-mp/nsx/model"
	"go.uber.org/mock/gomock"

	tf_api "github.com/vmware/terraform-provider-nsxt/api/utl"
	nsxmocks "github.com/vmware/terraform-provider-nsxt/mocks/nsx"
)

func providerTestData(overrides map[string]interface{}) map[string]interface{} {
	data := map[string]interface{}{
		"host":         "nsxmanager.example.com",
		"username":     "admin",
		"password":     "password",
		"session_auth": true,
	}
	for k, v := range overrides {
		data[k] = v
	}
	return data
}

func TestUnitNsxt_customHeaderProcessor(t *testing.T) {
	headers := map[string]string{"X-Custom": "value"}
	p := newCustomHeaderProcessor(&headers)
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	err := p.Process(req)
	require.NoError(t, err)
	assert.Equal(t, "value", req.Header.Get("X-Custom"))
}

func TestUnitNsxt_remoteAuthHeaderProcessor(t *testing.T) {
	p := newRemoteAuthHeaderProcessor()
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Authorization", "Basic dXNlcjpwYXNz")
	err := p.Process(req)
	require.NoError(t, err)
	assert.Equal(t, "Remote dXNlcjpwYXNz", req.Header.Get("Authorization"))
}

func TestUnitNsxt_logRequestProcessor(t *testing.T) {
	p := newLogRequestProcessor()
	req, _ := http.NewRequest("GET", "http://example.com/path", nil)
	err := p.Process(req)
	require.NoError(t, err)
}

func TestUnitNsxt_logResponseAcceptor(t *testing.T) {
	a := newLogResponseAcceptor()
	req, _ := http.NewRequest("GET", "http://example.com/path", nil)
	resp := &http.Response{StatusCode: 200, Body: http.NoBody, Request: req}
	a.Accept(resp)
}

func TestUnitNsxt_bearerAuthHeaderProcessor(t *testing.T) {
	p := newBearerAuthHeaderProcessor("my-token")
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	err := p.Process(req)
	require.NoError(t, err)
	assert.Equal(t, "Bearer my-token", req.Header.Get("Authorization"))
}

func TestUnitNsxt_sessionHeaderProcessor(t *testing.T) {
	p := newSessionHeaderProcessor("cookie-value", "xsrf-value")
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	err := p.Process(req)
	require.NoError(t, err)
	assert.Equal(t, "cookie-value", req.Header.Get("Cookie"))
	assert.Equal(t, "xsrf-value", req.Header.Get("X-XSRF-TOKEN"))
}

func TestUnitNsxt_isVMCCredentialSet(t *testing.T) {
	res := Provider()

	t.Run("vmc_token alone is sufficient", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{"vmc_token": "tok"}))
		assert.True(t, isVMCCredentialSet(d))
	})

	t.Run("client id and secret together are sufficient", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{
			"vmc_client_id": "id", "vmc_client_secret": "secret",
		}))
		assert.True(t, isVMCCredentialSet(d))
	})

	t.Run("neither set returns false", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(nil))
		assert.False(t, isVMCCredentialSet(d))
	})

	t.Run("client id without secret returns false", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{"vmc_client_id": "id"}))
		assert.False(t, isVMCCredentialSet(d))
	})
}

func TestUnitNsxt_getVmcAuthInfo(t *testing.T) {
	res := Provider()

	t.Run("explicit auth host is preserved", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{"vmc_auth_host": "custom.host"}))
		info := getVmcAuthInfo(d)
		assert.Equal(t, "custom.host", info.authHost)
	})

	t.Run("token auth defaults to api-tokens authorize host", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{"vmc_token": "tok"}))
		info := getVmcAuthInfo(d)
		assert.Contains(t, info.authHost, "api-tokens/authorize")
	})

	t.Run("oauth app defaults to authorize host", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{
			"vmc_client_id": "id", "vmc_client_secret": "secret",
		}))
		info := getVmcAuthInfo(d)
		assert.Contains(t, info.authHost, "/authorize")
		assert.NotContains(t, info.authHost, "api-tokens")
	})

	t.Run("no credentials leaves auth host empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(nil))
		info := getVmcAuthInfo(d)
		assert.Equal(t, "", info.authHost)
	})
}

func TestUnitNsxt_vmcAuthInfoIsZero(t *testing.T) {
	assert.True(t, (&vmcAuthInfo{}).IsZero())
	assert.False(t, (&vmcAuthInfo{accessToken: "t"}).IsZero())
	assert.False(t, (&vmcAuthInfo{clientID: "i", clientSecret: "s"}).IsZero())
}

func TestUnitNsxt_getAPIToken(t *testing.T) {
	t.Run("no credentials errors without a network call", func(t *testing.T) {
		info := &vmcAuthInfo{authHost: "example.com"}
		_, err := info.getAPIToken()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid VMC auth input")
	})

	t.Run("successful token exchange via refresh token", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"abc123"}`))
		}))
		defer server.Close()

		restore := redirectDefaultClientTo(server)
		defer restore()

		info := &vmcAuthInfo{authHost: "vmc.example.com", accessToken: "refresh-token"}
		token, err := info.getAPIToken()
		require.NoError(t, err)
		assert.Equal(t, "abc123", token)
	})

	t.Run("successful token exchange via oauth app", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			assert.True(t, ok)
			assert.Equal(t, "client-id", user)
			assert.Equal(t, "client-secret", pass)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"xyz789"}`))
		}))
		defer server.Close()

		restore := redirectDefaultClientTo(server)
		defer restore()

		info := &vmcAuthInfo{authHost: "vmc.example.com", clientID: "client-id", clientSecret: "client-secret"}
		token, err := info.getAPIToken()
		require.NoError(t, err)
		assert.Equal(t, "xyz789", token)
	})

	t.Run("non-200 response is an error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("unauthorized"))
		}))
		defer server.Close()

		restore := redirectDefaultClientTo(server)
		defer restore()

		info := &vmcAuthInfo{authHost: "vmc.example.com", accessToken: "refresh-token"}
		_, err := info.getAPIToken()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "401")
	})
}

// redirectDefaultClientTo makes http.DefaultClient route all requests to the given
// test server regardless of the original request URL, so getAPIToken's hardcoded
// "https://"+authHost URL can be exercised against an httptest.Server without TLS.
func redirectDefaultClientTo(server *httptest.Server) func() {
	original := http.DefaultClient.Transport
	http.DefaultClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		redirected := req.Clone(req.Context())
		redirected.URL.Scheme = "http"
		redirected.URL.Host = server.Listener.Addr().String()
		return http.DefaultTransport.RoundTrip(redirected)
	})
	return func() { http.DefaultClient.Transport = original }
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func generateTestCertKeyPEM(t *testing.T) (certPEM, keyPEM []byte) {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "test"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
	}
	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	require.NoError(t, err)
	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	return certPEM, keyPEM
}

func TestUnitNsxt_getConnectorTLSConfig(t *testing.T) {
	res := Provider()
	certPEM, keyPEM := generateTestCertKeyPEM(t)

	t.Run("insecure flag is propagated", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{"allow_unverified_ssl": true}))
		cfg, err := getConnectorTLSConfig(d)
		require.NoError(t, err)
		assert.True(t, cfg.InsecureSkipVerify)
	})

	t.Run("cert file without key file errors", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{"client_auth_cert_file": "/tmp/cert.pem"}))
		_, err := getConnectorTLSConfig(d)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "key file")
	})

	t.Run("cert string without key string errors", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{"client_auth_cert": string(certPEM)}))
		_, err := getConnectorTLSConfig(d)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "key")
	})

	t.Run("valid cert/key files load a client certificate", func(t *testing.T) {
		dir := t.TempDir()
		certFile := filepath.Join(dir, "cert.pem")
		keyFile := filepath.Join(dir, "key.pem")
		require.NoError(t, os.WriteFile(certFile, certPEM, 0o600))
		require.NoError(t, os.WriteFile(keyFile, keyPEM, 0o600))

		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{
			"client_auth_cert_file": certFile, "client_auth_key_file": keyFile,
		}))
		cfg, err := getConnectorTLSConfig(d)
		require.NoError(t, err)
		require.NotNil(t, cfg.GetClientCertificate)
	})

	t.Run("valid cert/key strings load a client certificate", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{
			"client_auth_cert": string(certPEM), "client_auth_key": string(keyPEM),
		}))
		cfg, err := getConnectorTLSConfig(d)
		require.NoError(t, err)
		require.NotNil(t, cfg.GetClientCertificate)
	})

	t.Run("invalid cert/key pair errors", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{
			"client_auth_cert": "not-a-cert", "client_auth_key": "not-a-key",
		}))
		_, err := getConnectorTLSConfig(d)
		require.Error(t, err)
	})

	t.Run("ca file is loaded into RootCAs", func(t *testing.T) {
		dir := t.TempDir()
		caFile := filepath.Join(dir, "ca.pem")
		require.NoError(t, os.WriteFile(caFile, certPEM, 0o600))

		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{"ca_file": caFile}))
		cfg, err := getConnectorTLSConfig(d)
		require.NoError(t, err)
		require.NotNil(t, cfg.RootCAs)
	})

	t.Run("missing ca file errors", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{"ca_file": "/no/such/file.pem"}))
		_, err := getConnectorTLSConfig(d)
		require.Error(t, err)
	})

	t.Run("ca string is loaded into RootCAs", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{"ca": string(certPEM)}))
		cfg, err := getConnectorTLSConfig(d)
		require.NoError(t, err)
		require.NotNil(t, cfg.RootCAs)
	})
}

func TestUnitNsxt_initCommonConfig(t *testing.T) {
	res := Provider()

	t.Run("defaults are applied when unset", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(nil))
		cfg := initCommonConfig(d)
		assert.Equal(t, "admin", cfg.Username)
		assert.Equal(t, "password", cfg.Password)
		assert.NotEmpty(t, cfg.RetryStatusCodes)
	})

	t.Run("retry_max_delay lower than min is bumped up", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{
			"retry_min_delay": 500, "retry_max_delay": 100,
		}))
		cfg := initCommonConfig(d)
		assert.Equal(t, 501, cfg.MaxRetryInterval)
	})

	t.Run("explicit retry_on_status_codes are preserved", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{
			"retry_on_status_codes": []interface{}{500, 503},
		}))
		cfg := initCommonConfig(d)
		assert.Equal(t, []int{500, 503}, cfg.RetryStatusCodes)
	})
}

func TestUnitNsxt_getGlobalPolicyEnforcementPointPath(t *testing.T) {
	clients := nsxtClients{PolicyEnforcementPoint: "default"}
	site := "/global-infra/sites/site1"
	assert.Equal(t, "/global-infra/sites/site1/enforcement-points/default", getGlobalPolicyEnforcementPointPath(clients, &site))
}

func TestUnitNsxt_getSessionContextFromParentPath(t *testing.T) {
	t.Run("multitenancy project path", func(t *testing.T) {
		ctx := testAccGetSessionContextFromParentPath(nsxtClients{}, "/orgs/default/projects/proj1/groups/g1")
		assert.EqualValues(t, tf_api.Multitenancy, ctx.ClientType)
		assert.Equal(t, "proj1", ctx.ProjectID)
	})

	t.Run("VPC path", func(t *testing.T) {
		ctx := testAccGetSessionContextFromParentPath(nsxtClients{}, "/orgs/default/projects/proj1/vpcs/vpc1/subnets/s1")
		assert.EqualValues(t, tf_api.VPC, ctx.ClientType)
		assert.Equal(t, "vpc1", ctx.VPCID)
	})

	t.Run("global manager path falls back to Global", func(t *testing.T) {
		ctx := testAccGetSessionContextFromParentPath(nsxtClients{PolicyGlobalManager: true}, "/infra/domains/default/groups/g1")
		assert.EqualValues(t, tf_api.Global, ctx.ClientType)
	})

	t.Run("local manager path falls back to Local", func(t *testing.T) {
		ctx := testAccGetSessionContextFromParentPath(nsxtClients{}, "/infra/domains/default/groups/g1")
		assert.EqualValues(t, tf_api.Local, ctx.ClientType)
	})
}

func TestMockNsxtGetLicenses(t *testing.T) {
	t.Run("filters out default vShield Endpoint license", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK := nsxmocks.NewMockLicensesClient(ctrl)
		original := cliLicensesClient
		cliLicensesClient = func(_ client.Connector) nsx.LicensesClient { return mockSDK }
		defer func() { cliLicensesClient = original }()

		defaultDesc := "NSX for vShield Endpoint"
		key1 := "key-1"
		key2 := "key-2"
		mockSDK.EXPECT().List().Return(nsxModel.LicensesListResult{
			Results: []nsxModel.License{
				{LicenseKey: &key1, Description: &defaultDesc},
				{LicenseKey: &key2},
			},
		}, nil)

		licenses, err := getLicenses(nil)
		require.NoError(t, err)
		assert.Equal(t, []string{"key-2"}, licenses)
	})

	t.Run("List error is wrapped", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK := nsxmocks.NewMockLicensesClient(ctrl)
		original := cliLicensesClient
		cliLicensesClient = func(_ client.Connector) nsx.LicensesClient { return mockSDK }
		defer func() { cliLicensesClient = original }()

		mockSDK.EXPECT().List().Return(nsxModel.LicensesListResult{}, errors.New("boom"))

		_, err := getLicenses(nil)
		require.Error(t, err)
	})
}

func TestMockNsxtApplyLicense(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSDK := nsxmocks.NewMockLicensesClient(ctrl)
	original := cliLicensesClient
	cliLicensesClient = func(_ client.Connector) nsx.LicensesClient { return mockSDK }
	defer func() { cliLicensesClient = original }()

	mockSDK.EXPECT().Create(gomock.Any()).Return(nsxModel.License{}, nil)

	err := applyLicense(nil, "new-key")
	require.NoError(t, err)
}

func TestMockNsxtConfigureLicenses(t *testing.T) {
	t.Run("no intent licenses is a no-op", func(t *testing.T) {
		err := configureLicenses(nil, nil)
		require.NoError(t, err)
	})

	t.Run("applies only missing licenses", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSDK := nsxmocks.NewMockLicensesClient(ctrl)
		original := cliLicensesClient
		cliLicensesClient = func(_ client.Connector) nsx.LicensesClient { return mockSDK }
		defer func() { cliLicensesClient = original }()

		existingKey := "existing-key"
		mockSDK.EXPECT().List().Return(nsxModel.LicensesListResult{
			Results: []nsxModel.License{{LicenseKey: &existingKey}},
		}, nil)
		mockSDK.EXPECT().Create(gomock.Any()).Return(nsxModel.License{}, nil)

		err := configureLicenses(nil, []string{"existing-key", "new-key"})
		require.NoError(t, err)
	})
}

func TestUnitNsxt_configureNsxtClient(t *testing.T) {
	res := Provider()

	t.Run("on_demand_connection skips client setup", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{"on_demand_connection": true}))
		clients := &nsxtClients{}
		err := configureNsxtClient(d, clients)
		require.NoError(t, err)
		assert.Nil(t, clients.NsxtClientConfig)
	})

	t.Run("vmc basic auth mode skips client setup", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{"vmc_auth_mode": "Basic"}))
		clients := &nsxtClients{}
		err := configureNsxtClient(d, clients)
		require.NoError(t, err)
		assert.Nil(t, clients.NsxtClientConfig)
	})

	t.Run("vmc credentials set skips client setup", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{"vmc_token": "tok"}))
		clients := &nsxtClients{}
		err := configureNsxtClient(d, clients)
		require.NoError(t, err)
		assert.Nil(t, clients.NsxtClientConfig)
	})

	t.Run("cert file without key file errors", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{"client_auth_cert_file": "/tmp/cert.pem"}))
		err := configureNsxtClient(d, &nsxtClients{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "key file")
	})

	t.Run("cert string without key string errors", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{"client_auth_cert": "cert-data"}))
		err := configureNsxtClient(d, &nsxtClients{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "key")
	})

	t.Run("missing username errors when credentials required", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{"username": ""}))
		err := configureNsxtClient(d, &nsxtClients{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "username")
	})

	t.Run("missing password errors when credentials required", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{"password": ""}))
		err := configureNsxtClient(d, &nsxtClients{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "password")
	})

	t.Run("missing host errors", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{"host": ""}))
		err := configureNsxtClient(d, &nsxtClients{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "host")
	})

	t.Run("valid config with session_auth disabled succeeds without network access", func(t *testing.T) {
		// session_auth=false skips api.NewAPIClient's session-cookie fetch, which
		// would otherwise attempt a real HTTP call to the (fake) host.
		d := schema.TestResourceDataRaw(t, res.Schema, providerTestData(map[string]interface{}{"session_auth": false}))
		clients := &nsxtClients{}
		err := configureNsxtClient(d, clients)
		require.NoError(t, err)
		require.NotNil(t, clients.NsxtClientConfig)
		assert.Equal(t, "nsxmanager.example.com", clients.NsxtClientConfig.Host)
		assert.NotNil(t, clients.NsxtClient)
	})
}
