/* Copyright © 2024 Broadcom, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api "github.com/vmware/go-vmware-nsxt"

	"github.com/vmware/terraform-provider-nsxt/nsxt/customtypes"
)

type nsxtProvider struct {
	version string
}

type nsxtProviderModel struct {
	AllowUnverifiedSsl     types.Bool   `tfsdk:"allow_unverified_ssl"`
	Username               types.String `tfsdk:"username"`
	Password               types.String `tfsdk:"password"`
	RemoteAuth             types.Bool   `tfsdk:"remote_auth"`
	SessionAuth            types.Bool   `tfsdk:"session_auth"`
	Host                   types.String `tfsdk:"host"`
	ClientAuthCertFile     types.String `tfsdk:"client_auth_cert_file"`
	ClientAuthKeyFile      types.String `tfsdk:"client_auth_key_file"`
	CAFile                 types.String `tfsdk:"ca_file"`
	MaxRetries             types.Int64  `tfsdk:"max_retries"`
	RetryMinDelay          types.Int64  `tfsdk:"retry_min_delay"`
	RetryMaxDelay          types.Int64  `tfsdk:"retry_max_delay"`
	RetryOnStatusCodes     types.List   `tfsdk:"retry_on_status_codes"`
	ToleratePartialSuccess types.Bool   `tfsdk:"tolerate_partial_success"`
	VMCAuthHost            types.String `tfsdk:"vmc_auth_host"`
	VMCToken               types.String `tfsdk:"vmc_token"`
	VMCClientId            types.String `tfsdk:"vmc_client_id"`
	VMCClientSecret        types.String `tfsdk:"vmc_client_secret"`
	VMCAuthMode            types.String `tfsdk:"vmc_auth_mode"`
	EnforcementPoint       types.String `tfsdk:"enforcement_point"`
	GlobalManager          types.Bool   `tfsdk:"global_manager"`
	LicenseKeys            types.List   `tfsdk:"license_keys"`
	ClientAuthCert         types.String `tfsdk:"client_auth_cert"`
	ClientAuthKey          types.String `tfsdk:"client_auth_key"`
	CA                     types.String `tfsdk:"ca"`
	OnDemandConnection     types.Bool   `tfsdk:"on_demand_connection"`
}

func (n nsxtProvider) Metadata(ctx context.Context, request provider.MetadataRequest, response *provider.MetadataResponse) {
	response.TypeName = "nsxt"
	//TODO: implant provider version here
	response.Version = n.version
}

func (n nsxtProvider) Schema(ctx context.Context, request provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"allow_unverified_ssl": schema.BoolAttribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"remote_auth": schema.BoolAttribute{
				Optional: true,
			},
			"session_auth": schema.BoolAttribute{
				Optional: true,
			},
			"host": schema.StringAttribute{
				Optional:    true,
				Description: "The hostname or IP address of the NSX manager.",
			},
			"client_auth_cert_file": schema.StringAttribute{
				Optional: true,
			},
			"client_auth_key_file": schema.StringAttribute{
				Optional: true,
			},
			"ca_file": schema.StringAttribute{
				Optional: true,
			},
			"max_retries": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum number of HTTP client retries",
			},
			"retry_min_delay": schema.Int64Attribute{
				Optional:    true,
				Description: "Minimum delay in milliseconds between retries of a request",
			},
			"retry_max_delay": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum delay in milliseconds between retries of a request",
			},
			"retry_on_status_codes": schema.ListAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
				Description: "HTTP replies status codes to retry on",
			},
			"tolerate_partial_success": schema.BoolAttribute{
				Optional:    true,
				Description: "Treat partial success status as success",
			},
			"vmc_auth_host": schema.StringAttribute{
				Optional:    true,
				Description: "URL for VMC authorization service (CSP)",
			},
			"vmc_token": schema.StringAttribute{
				Optional:    true,
				Description: "Long-living API token for VMC authorization",
			},
			"vmc_client_id": schema.StringAttribute{
				Optional:    true,
				Description: "ID of OAuth App associated with the VMC organization",
			},
			"vmc_client_secret": schema.StringAttribute{
				Optional:    true,
				Description: "Secret of OAuth App associated with the VMC organization",
			},
			"vmc_auth_mode": schema.StringAttribute{
				Optional:    true,
				Validators:  []validator.String{stringvalidator.OneOf("Default", "Bearer", "Basic")},
				Description: "Mode for VMC authorization",
			},
			"enforcement_point": schema.StringAttribute{
				Optional:    true,
				Description: "Enforcement Point for NSXT Policy",
			},
			"global_manager": schema.BoolAttribute{
				Optional:    true,
				Description: "Is this a policy global manager endpoint",
			},
			"license_keys": schema.ListAttribute{
				Optional:    true,
				Description: "license keys",
				ElementType: customtypes.NSXLicenseType{},
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.Expressions{path.MatchRoot("vmc_token")}...),
				},
			},
			"client_auth_cert": schema.StringAttribute{
				Description: "Client certificate passed as string",
				Optional:    true,
			},
			"client_auth_key": schema.StringAttribute{
				Description: "Client certificate key passed as string",
				Optional:    true,
			},
			"ca": schema.StringAttribute{
				Description: "CA certificate passed as string",
				Optional:    true,
			},
			"on_demand_connection": schema.BoolAttribute{
				Optional:    true,
				Description: "Avoid initializing NSX connection on startup",
			},
		},
	}
}

func getBoolFromConfigOrEnv(cfgVal types.Bool, varName string, defVal bool, diag diag.Diagnostics) bool {
	if !cfgVal.IsNull() {
		return cfgVal.ValueBool()
	}
	strVal := os.Getenv(varName)
	if strVal == "" {
		return defVal
	}
	boolVal, err := strconv.ParseBool(strVal)
	if err != nil {
		diag.AddError(fmt.Sprintf("Error parsing boolean value %s from environment variable %s", strVal, varName), err.Error())
		return false
	}
	return boolVal
}

func getInt64FromConfigOrEnv(cfgVal types.Int64, varName string, defVal int64, diag diag.Diagnostics) int64 {
	if !cfgVal.IsNull() {
		return cfgVal.ValueInt64()
	}
	strVal := os.Getenv(varName)
	if strVal == "" {
		return defVal
	}
	intVal, err := strconv.ParseInt(strVal, 10, 64)
	if err != nil {
		diag.AddError(fmt.Sprintf("Error parsing int64 value %s from environment variable %s", strVal, varName), err.Error())
		return 0
	}
	return intVal
}

func getStringFromConfigOrEnv(cfgVal types.String, varName string, defVal string) string {
	if !cfgVal.IsNull() {
		return cfgVal.ValueString()
	}
	strVal := os.Getenv(varName)
	if strVal == "" {
		return defVal
	}
	return strVal
}

func initFrameworkCommonConfig(ctx context.Context, config nsxtProviderModel, response *provider.ConfigureResponse) commonProviderConfig {
	remoteAuth := getBoolFromConfigOrEnv(config.RemoteAuth, "NSXT_REMOTE_AUTH", false, response.Diagnostics)
	toleratePartialSuccess := getBoolFromConfigOrEnv(config.ToleratePartialSuccess, "NSXT_TOLERATE_PARTIAL_SUCCESS", false, response.Diagnostics)
	maxRetries := int(getInt64FromConfigOrEnv(config.MaxRetries, "NSXT_MAX_RETRIES", 4, response.Diagnostics))
	retryMinDelay := int(getInt64FromConfigOrEnv(config.RetryMinDelay, "NSXT_RETRY_MIN_DELAY", 0, response.Diagnostics))
	retryMaxDelay := int(getInt64FromConfigOrEnv(config.RetryMaxDelay, "NSXT_RETRY_MAX_DELAY", 500, response.Diagnostics))
	username := getStringFromConfigOrEnv(config.Username, "NSXT_USERNAME", "")
	password := getStringFromConfigOrEnv(config.Password, "NSXT_PASSWORD", "")

	retryStatuses := getIntListFromFrameworkList(ctx, config.RetryOnStatusCodes, response.Diagnostics)
	if len(retryStatuses) == 0 {
		// Set to the defaults if empty
		retryStatuses = append(retryStatuses, defaultRetryOnStatusCodes...)
	}

	licenses := getStringListFromFrameworkList(ctx, config.LicenseKeys, response.Diagnostics)

	commonConfig := commonProviderConfig{
		RemoteAuth:             remoteAuth,
		ToleratePartialSuccess: toleratePartialSuccess,
		MaxRetries:             maxRetries,
		MinRetryInterval:       retryMinDelay,
		MaxRetryInterval:       retryMaxDelay,
		RetryStatusCodes:       retryStatuses,
		Username:               username,
		Password:               password,
		LicenseKeys:            licenses,
	}

	return commonConfig
}

func configureFrameworkNsxtClient(config nsxtProviderModel, clients *nsxtClients, response *provider.ConfigureResponse) error {
	onDemandConn := getBoolFromConfigOrEnv(config.OnDemandConnection, "NSXT_ON_DEMAND_CONNECTION", false, response.Diagnostics)
	clientAuthCertFile := getStringFromConfigOrEnv(config.ClientAuthCertFile, "NSXT_CLIENT_AUTH_CERT_FILE", "")
	clientAuthKeyFile := getStringFromConfigOrEnv(config.ClientAuthKeyFile, "NSXT_CLIENT_AUTH_Key_FILE", "")
	clientAuthCert := getStringFromConfigOrEnv(config.ClientAuthCert, "NSXT_CLIENT_AUTH_CERT", "")
	clientAuthKey := getStringFromConfigOrEnv(config.ClientAuthKey, "NSXT_CLIENT_AUTH_KEY", "")
	vmcToken := getStringFromConfigOrEnv(config.VMCToken, "NSXT_VMC_TOKEN", "")
	var vmcAuthMode string
	if !config.VMCAuthMode.IsNull() {
		vmcAuthMode = config.VMCAuthMode.ValueString()
	}

	if onDemandConn {
		// On demand connection option is not supported with old SDK
		return nil
	}

	if (len(vmcToken) > 0) || (vmcAuthMode == "Basic") {
		// VMC can operate without token with basic auth, however MP API is not
		// available for cloud admin user
		return nil
	}

	needCreds := true
	if len(clientAuthCertFile) > 0 {
		if len(clientAuthKeyFile) == 0 {
			return fmt.Errorf("please provide key file for client certificate")
		}
		needCreds = false
	}

	if len(clientAuthCert) > 0 {
		if len(clientAuthKey) == 0 {
			return fmt.Errorf("please provide key for client certificate")
		}
		// only supported for policy resources
		needCreds = false
	}

	insecure := getBoolFromConfigOrEnv(config.AllowUnverifiedSsl, "NSXT_ALLOW_UNVERIFIED_SSL", false, response.Diagnostics)
	username := getStringFromConfigOrEnv(config.Username, "NSXT_USERNAME", "")
	password := getStringFromConfigOrEnv(config.Password, "NSXT_PASSWORD", "")

	if needCreds {
		if username == "" {
			return fmt.Errorf("username must be provided")
		}

		if password == "" {
			return fmt.Errorf("password must be provided")
		}
	}

	host := getStringFromConfigOrEnv(config.Host, "NSXT_HOST", "")
	// Remove schema
	host = strings.TrimPrefix(host, "https://")

	if host == "" {
		return fmt.Errorf("host must be provided")
	}

	caFile := getStringFromConfigOrEnv(config.CAFile, "NSXT_CA_FILE", "")
	caString := getStringFromConfigOrEnv(config.CA, "NSXT_CA", "")
	sessionAuth := getBoolFromConfigOrEnv(config.SessionAuth, "NSXT_SESSION_AUTH", false, response.Diagnostics)
	skipSessionAuth := !sessionAuth

	retriesConfig := api.ClientRetriesConfiguration{
		MaxRetries:      clients.CommonConfig.MaxRetries,
		RetryMinDelay:   clients.CommonConfig.MinRetryInterval,
		RetryMaxDelay:   clients.CommonConfig.MaxRetryInterval,
		RetryOnStatuses: clients.CommonConfig.RetryStatusCodes,
	}

	clients.NsxtClientConfig = &api.Configuration{
		BasePath:             "/api/v1",
		Host:                 host,
		Scheme:               "https",
		UserAgent:            "terraform-provider-nsxt",
		UserName:             username,
		Password:             password,
		RemoteAuth:           clients.CommonConfig.RemoteAuth,
		ClientAuthCertFile:   clientAuthCertFile,
		ClientAuthKeyFile:    clientAuthKeyFile,
		CAFile:               caFile,
		ClientAuthCertString: clientAuthCert,
		ClientAuthKeyString:  clientAuthKey,
		CAString:             caString,
		Insecure:             insecure,
		RetriesConfiguration: retriesConfig,
		SkipSessionAuth:      skipSessionAuth,
	}

	nsxClient, err := api.NewAPIClient(clients.NsxtClientConfig)
	if err != nil {
		return err
	}

	clients.NsxtClient = nsxClient

	return nil
}

func getFrameworkConnectorTLSConfig(config nsxtProviderModel, response *provider.ConfigureResponse) (*tls.Config, error) {

	insecure := getBoolFromConfigOrEnv(config.AllowUnverifiedSsl, "NSXT_ALLOW_UNVERIFIED_SSL", false, response.Diagnostics)
	clientAuthCertFile := getStringFromConfigOrEnv(config.ClientAuthCertFile, "NSXT_CLIENT_AUTH_CERT_FILE", "")
	clientAuthKeyFile := getStringFromConfigOrEnv(config.ClientAuthKeyFile, "NSXT_CLIENT_AUTH_Key_FILE", "")
	caFile := getStringFromConfigOrEnv(config.CAFile, "NSXT_CA_FILE", "")
	clientAuthCert := getStringFromConfigOrEnv(config.ClientAuthCert, "NSXT_CLIENT_AUTH_CERT", "")
	clientAuthKey := getStringFromConfigOrEnv(config.ClientAuthKey, "NSXT_CLIENT_AUTH_KEY", "")
	caCert := getStringFromConfigOrEnv(config.CA, "NSXT_CA", "")
	tlsConfig := tls.Config{InsecureSkipVerify: insecure}

	if len(clientAuthCertFile) > 0 {

		// cert and key are passed via filesystem
		if len(clientAuthKeyFile) == 0 {
			return nil, fmt.Errorf("please provide key file for client certificate")
		}

		cert, err := tls.LoadX509KeyPair(clientAuthCertFile, clientAuthKeyFile)

		if err != nil {
			return nil, fmt.Errorf("failed to load client cert/key pair: %v", err)
		}

		tlsConfig.GetClientCertificate = func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
			return &cert, nil
		}
	}

	if len(clientAuthCert) > 0 {
		// cert and key are passed as strings
		if len(clientAuthKey) == 0 {
			return nil, fmt.Errorf("please provide key for client certificate")
		}

		cert, err := tls.X509KeyPair([]byte(clientAuthCert), []byte(clientAuthKey))

		if err != nil {
			return nil, fmt.Errorf("failed to load client cert/key pair: %v", err)
		}

		tlsConfig.GetClientCertificate = func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
			return &cert, nil
		}
	}

	if len(caFile) > 0 {
		caCert, err := ioutil.ReadFile(caFile)
		if err != nil {
			return nil, err
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig.RootCAs = caCertPool
	}

	if len(caCert) > 0 {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCert))

		tlsConfig.RootCAs = caCertPool
	}

	return &tlsConfig, nil
}

func getFrameworkVmcAuthInfo(config nsxtProviderModel) *vmcAuthInfo {
	vmcInfo := vmcAuthInfo{
		authHost:     getStringFromConfigOrEnv(config.VMCAuthHost, "NSXT_VMC_AUTH_HOST", ""),
		authMode:     getStringFromConfigOrEnv(config.VMCAuthMode, "NSXT_VMC_AUTH_MODE", ""),
		accessToken:  getStringFromConfigOrEnv(config.VMCToken, "NSXT_VMC_TOKEN", ""),
		clientID:     getStringFromConfigOrEnv(config.VMCClientId, "NSXT_VMC_CLIENT_ID", ""),
		clientSecret: getStringFromConfigOrEnv(config.VMCClientSecret, "NSXT_VMC_CLIENT_SECRET", ""),
	}
	if len(vmcInfo.authHost) > 0 {
		return &vmcInfo
	}

	// Fill in default auth host + url based on auth method
	if len(vmcInfo.accessToken) > 0 {
		vmcInfo.authHost = "console.cloud.vmware.com/csp/gateway/am/api/auth/api-tokens/authorize"
	} else if len(vmcInfo.clientSecret) > 0 && len(vmcInfo.clientID) > 0 {
		vmcInfo.authHost = "console.cloud.vmware.com/csp/gateway/am/api/auth/authorize"
	}
	return &vmcInfo
}

func isFrameworkVMCCredentialSet(config nsxtProviderModel) bool {
	// Refresh token
	vmcToken := getStringFromConfigOrEnv(config.VMCToken, "NSXT_VMC_TOKEN", "")
	if len(vmcToken) > 0 {
		return true
	}

	// Oauth app
	vmcClientID := getStringFromConfigOrEnv(config.VMCClientId, "NSXT_VMC_CLIENT_ID", "")
	vmcClientSecret := getStringFromConfigOrEnv(config.VMCClientSecret, "NSXT_VMC_CLIENT_SECRET", "")
	if len(vmcClientSecret) > 0 && len(vmcClientID) > 0 {
		return true
	}

	return false
}

func configureFrameworkPolicyConnectorData(ctx context.Context, config nsxtProviderModel, clients *nsxtClients, response *provider.ConfigureResponse) error {
	onDemandConn := getBoolFromConfigOrEnv(config.OnDemandConnection, "NSXT_ON_DEMAND_CONNECTION", false, response.Diagnostics)
	host := getStringFromConfigOrEnv(config.Host, "NSXT_HOST", "")
	username := getStringFromConfigOrEnv(config.Username, "NSXT_USERNAME", "")
	password := getStringFromConfigOrEnv(config.Password, "NSXT_PASSWORD", "")
	// vmcToken := getStringFromConfigOrEnv(config.VMCToken, "NSXT_VMC_TOKEN", "")
	// vmcAuthHost := getStringFromConfigOrEnv(config.VMCAuthHost, "NSXT_VMC_AUTH_HOST", "")
	clientAuthCertFile := getStringFromConfigOrEnv(config.ClientAuthCertFile, "NSXT_CLIENT_AUTH_CERT_FILE", "")
	clientAuthCert := getStringFromConfigOrEnv(config.ClientAuthCert, "NSXT_CLIENT_AUTH_CERT", "")
	clientAuthDefined := (len(clientAuthCertFile) > 0) || (len(clientAuthCert) > 0)
	policyEnforcementPoint := getStringFromConfigOrEnv(config.EnforcementPoint, "NSXT_POLICY_ENFORCEMENT_POINT", "default")
	policyGlobalManager := getBoolFromConfigOrEnv(config.GlobalManager, "NSXT_GLOBAL_MANAGER", false, response.Diagnostics)
	// vmcAuthMode := getStringFromConfigOrEnv(config.VMCAuthMode, "NSXT_VMC_AUTH_MODE", "Default")
	vmcInfo := getFrameworkVmcAuthInfo(config)

	isVMC := false
	if (vmcInfo.authMode == "Basic") || isFrameworkVMCCredentialSet(config) {
		isVMC = true
		if onDemandConn {
			return fmt.Errorf("on demand connection option is not supported with VMC")
		}
	}

	if host == "" {
		return fmt.Errorf("host must be provided")
	}

	if !strings.HasPrefix(host, "https://") {
		host = fmt.Sprintf("https://%s", host)
	}

	securityContextNeeded := true
	if clientAuthDefined && !clients.CommonConfig.RemoteAuth {
		securityContextNeeded = false
	}
	if securityContextNeeded {
		securityCtx, err := getConfiguredSecurityContext(clients, vmcInfo, username, password)
		if err != nil {
			return err
		}
		clients.PolicySecurityContext = securityCtx
	}

	tlsConfig, err := getFrameworkConnectorTLSConfig(config, response)
	if err != nil {
		return err
	}

	tr := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: tlsConfig,
	}

	httpClient := http.Client{Transport: tr}
	clients.PolicyHTTPClient = &httpClient
	clients.Host = host
	clients.PolicyEnforcementPoint = policyEnforcementPoint
	clients.PolicyGlobalManager = policyGlobalManager

	if onDemandConn {
		// version init will happen on demand
		return nil
	}

	if !isVMC {
		err = configureLicenses(getStandalonePolicyConnector(*clients, true), clients.CommonConfig.LicenseKeys)
		if err != nil {
			return err
		}
	}

	err = initNSXVersion(getStandalonePolicyConnector(*clients, true))
	if err != nil && isVMC {
		// In case version API does not work for VMC, we workaround by testing version-specific APIs
		// TODO - remove this when /node/version API works for all auth methods on VMC
		initNSXVersionVMC(*clients)
		return nil
	}
	return err
}

func (n nsxtProvider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	var config nsxtProviderModel
	diags := request.Config.Get(ctx, &config)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	commonConfig := initFrameworkCommonConfig(ctx, config, response)
	clients := nsxtClients{
		CommonConfig: commonConfig,
	}

	err := configureFrameworkNsxtClient(config, &clients, response)
	if err != nil {
		response.Diagnostics.AddError("Error configuring framework NSXT client", err.Error())
	}

	err = configureFrameworkPolicyConnectorData(ctx, config, &clients, response)
	if err != nil {
		response.Diagnostics.AddError("Error configuring framework policy connector data", err.Error())
	}

	response.ResourceData = clients
}

func getIntListFromFrameworkList(ctx context.Context, values types.List, diag diag.Diagnostics) []int {
	var intValues []int
	if !values.IsNull() && !values.IsUnknown() {
		elems := make([]types.Int64, 0, len(values.Elements()))
		intValues = make([]int, len(values.Elements()))
		diags := values.ElementsAs(ctx, &elems, false)
		diag.Append(diags...)
		if !diag.HasError() {
			for i, v := range elems {
				intValues[i] = int(v.ValueInt64())
			}
		}
	}
	return intValues
}

func getStringListFromFrameworkList(ctx context.Context, values types.List, diag diag.Diagnostics) []string {
	var stringValues []string
	if !values.IsNull() && !values.IsUnknown() {
		elems := make([]types.String, 0, len(values.Elements()))
		stringValues = make([]string, len(values.Elements()))
		diags := values.ElementsAs(ctx, &elems, false)
		diag.Append(diags...)
		if !diag.HasError() {
			for i, v := range elems {
				stringValues[i] = v.ValueString()
			}
		}
	}
	return stringValues
}

func (n nsxtProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (n nsxtProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPolicyIPDiscoveryProfileResource,
		NewPolicyQOSProfileResource,
	}
}

func NewFrameworkProvider() provider.Provider {
	return &nsxtProvider{}
}
