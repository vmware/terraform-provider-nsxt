// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-nsxt/nsxt/metadata"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	tier0ipsecvpnservices "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s/ipsec_vpn_services"
	tier0localeservicesipsec "github.com/vmware/terraform-provider-nsxt/api/infra/tier_0s/locale_services/ipsec_vpn_services"
	tier1ipsecvpnservices "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s/ipsec_vpn_services"
	tier1localeservicesipsec "github.com/vmware/terraform-provider-nsxt/api/infra/tier_1s/locale_services/ipsec_vpn_services"
	utl "github.com/vmware/terraform-provider-nsxt/api/utl"
)

var cliTier0LocaleServiceIpsecVpnLocalEndpointsClient = tier0localeservicesipsec.NewLocalEndpointsClient
var cliTier0IpsecVpnLocalEndpointsClient = tier0ipsecvpnservices.NewLocalEndpointsClient
var cliTier1LocaleServiceIpsecVpnLocalEndpointsClient = tier1localeservicesipsec.NewLocalEndpointsClient
var cliTier1IpsecVpnLocalEndpointsClient = tier1ipsecvpnservices.NewLocalEndpointsClient

func getIPSecVpnLocalEndpointCommonSchema(isTransitGateway bool) map[string]*metadata.ExtendedSchema {
	schema := map[string]*metadata.ExtendedSchema{
		"nsx_id":       metadata.GetExtendedSchema(getNsxIDSchema()),
		"path":         metadata.GetExtendedSchema(getPathSchema()),
		"display_name": metadata.GetExtendedSchema(getDisplayNameSchema()),
		"description":  metadata.GetExtendedSchema(getDescriptionSchema()),
		"revision":     metadata.GetExtendedSchema(getRevisionSchema()),
		"tag":          metadata.GetExtendedSchema(getTagsSchema()),
		"context":      metadata.GetExtendedSchema(getContextSchema(false, false, false)),
		"local_address": {
			Schema: schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsIPAddress,
			},
		},
		"certificate_path": metadata.GetExtendedSchema(getPolicyPathSchema(false, false, "Policy path referencing site certificate")),
		"local_id": {
			Schema: schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		"trust_ca_paths": {
			Schema: schema.Schema{
				Type:     schema.TypeSet,
				Elem:     getElemPolicyPathSchema(),
				Optional: true,
			},
		},
		"trust_crl_paths": {
			Schema: schema.Schema{
				Type:     schema.TypeSet,
				Elem:     getElemPolicyPathSchema(),
				Optional: true,
			},
		},
	}

	if isTransitGateway {
		schema["parent_path"] = metadata.GetExtendedSchema(getPolicyPathSchema(true, true, "Policy path of the parent"))
	} else {
		schema["service_path"] = metadata.GetExtendedSchema(getPolicyPathSchema(true, true, "Policy path for IPSec VPN service"))
	}
	return schema
}

var IPSecVpnLocalEndpointSchema = getIPSecVpnLocalEndpointCommonSchema(false)

func resourceNsxtPolicyIPSecVpnLocalEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyIPSecVpnLocalEndpointCreate,
		Read:   resourceNsxtPolicyIPSecVpnLocalEndpointRead,
		Update: resourceNsxtPolicyIPSecVpnLocalEndpointUpdate,
		Delete: resourceNsxtPolicyIPSecVpnLocalEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: nsxtVPNServiceResourceImporter,
		},

		Schema: metadata.GetSchemaFromExtendedSchema(IPSecVpnLocalEndpointSchema),

		CustomizeDiff: validateIPv6LocaleServiceConflict,
	}
}

type localEndpointClient struct {
	isT0            bool
	gwID            string
	localeServiceID string
	serviceID       string
	sessionContext  utl.SessionContext
}

func getIPSecVpnLocalEndpointSessionContext(d *schema.ResourceData, m interface{}, servicePath string) utl.SessionContext {
	sessionContext := getSessionContext(d, m)
	if strings.HasPrefix(servicePath, "/orgs/") {
		if sessionContext.ProjectID == "" {
			sessionContext.ProjectID = extractProjectIDFromPolicyPath(servicePath)
		}
		if sessionContext.ProjectID != "" {
			sessionContext.ClientType = utl.Multitenancy
		}
	}
	return sessionContext
}

func newLocalEndpointClient(servicePath string, sessionContext utl.SessionContext) (*localEndpointClient, error) {
	isT0, gwID, localeServiceID, serviceID, err := parseIPSecVPNServicePolicyPath(servicePath)
	if err != nil {
		return nil, err
	}

	return &localEndpointClient{
		isT0:            isT0,
		gwID:            gwID,
		localeServiceID: localeServiceID,
		serviceID:       serviceID,
		sessionContext:  sessionContext,
	}, nil
}

func (c *localEndpointClient) Get(connector client.Connector, id string) (model.IPSecVpnLocalEndpoint, error) {
	if c.isT0 {
		if len(c.localeServiceID) > 0 {
			client := cliTier0LocaleServiceIpsecVpnLocalEndpointsClient(c.sessionContext, connector)
			if client == nil {
				return model.IPSecVpnLocalEndpoint{}, fmt.Errorf("unsupported client type")
			}
			return client.Get(c.gwID, c.localeServiceID, c.serviceID, id)
		}
		client := cliTier0IpsecVpnLocalEndpointsClient(c.sessionContext, connector)
		if client == nil {
			return model.IPSecVpnLocalEndpoint{}, fmt.Errorf("unsupported client type")
		}
		return client.Get(c.gwID, c.serviceID, id)

	}
	if len(c.localeServiceID) > 0 {
		if c.sessionContext.ClientType == utl.Multitenancy {
			return model.IPSecVpnLocalEndpoint{}, fmt.Errorf("project context is not supported for locale-service scoped IPSec VPN local endpoints")
		}
		client := cliTier1LocaleServiceIpsecVpnLocalEndpointsClient(c.sessionContext, connector)
		if client == nil {
			return model.IPSecVpnLocalEndpoint{}, fmt.Errorf("unsupported client type")
		}
		return client.Get(c.gwID, c.localeServiceID, c.serviceID, id)
	}
	client := cliTier1IpsecVpnLocalEndpointsClient(c.sessionContext, connector)
	if client == nil {
		return model.IPSecVpnLocalEndpoint{}, fmt.Errorf("unsupported client type")
	}
	return client.Get(c.gwID, c.serviceID, id)
}

// List retrieves a list of IPSecVpnLocalEndpoint objects for the specified gateway, service, and locale configuration.
// Note: We don't expect pagination to be relevant here.
func (c *localEndpointClient) List(connector client.Connector) ([]model.IPSecVpnLocalEndpoint, error) {
	boolFalse := false
	var cursor string
	var result model.IPSecVpnLocalEndpointListResult
	var err error
	if c.isT0 {
		if len(c.localeServiceID) > 0 {
			client := cliTier0LocaleServiceIpsecVpnLocalEndpointsClient(c.sessionContext, connector)
			if client == nil {
				return nil, fmt.Errorf("unsupported client type")
			}
			result, err = client.List(c.gwID, c.localeServiceID, c.serviceID, &cursor, &boolFalse, nil, nil, nil, nil)
		} else {
			client := cliTier0IpsecVpnLocalEndpointsClient(c.sessionContext, connector)
			if client == nil {
				return nil, fmt.Errorf("unsupported client type")
			}
			result, err = client.List(c.gwID, c.serviceID, &cursor, &boolFalse, nil, nil, nil, nil)
		}

	} else {
		if len(c.localeServiceID) > 0 {
			if c.sessionContext.ClientType == utl.Multitenancy {
				return nil, fmt.Errorf("project context is not supported for locale-service scoped IPSec VPN local endpoints")
			}
			client := cliTier1LocaleServiceIpsecVpnLocalEndpointsClient(c.sessionContext, connector)
			if client == nil {
				return nil, fmt.Errorf("unsupported client type")
			}
			result, err = client.List(c.gwID, c.localeServiceID, c.serviceID, &cursor, &boolFalse, nil, nil, nil, nil)
		} else {
			client := cliTier1IpsecVpnLocalEndpointsClient(c.sessionContext, connector)
			if client == nil {
				return nil, fmt.Errorf("unsupported client type")
			}
			result, err = client.List(c.gwID, c.serviceID, &cursor, &boolFalse, nil, nil, nil, nil)
		}
	}

	return result.Results, err
}

func (c *localEndpointClient) Patch(connector client.Connector, id string, obj model.IPSecVpnLocalEndpoint) error {
	if c.isT0 {
		if len(c.localeServiceID) > 0 {
			client := cliTier0LocaleServiceIpsecVpnLocalEndpointsClient(c.sessionContext, connector)
			if client == nil {
				return fmt.Errorf("unsupported client type")
			}
			return client.Patch(c.gwID, c.localeServiceID, c.serviceID, id, obj)
		}
		client := cliTier0IpsecVpnLocalEndpointsClient(c.sessionContext, connector)
		if client == nil {
			return fmt.Errorf("unsupported client type")
		}
		return client.Patch(c.gwID, c.serviceID, id, obj)

	}
	if len(c.localeServiceID) > 0 {
		if c.sessionContext.ClientType == utl.Multitenancy {
			return fmt.Errorf("project context is not supported for locale-service scoped IPSec VPN local endpoints")
		}
		client := cliTier1LocaleServiceIpsecVpnLocalEndpointsClient(c.sessionContext, connector)
		if client == nil {
			return fmt.Errorf("unsupported client type")
		}
		return client.Patch(c.gwID, c.localeServiceID, c.serviceID, id, obj)
	}
	client := cliTier1IpsecVpnLocalEndpointsClient(c.sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	return client.Patch(c.gwID, c.serviceID, id, obj)
}

func (c *localEndpointClient) Delete(connector client.Connector, id string) error {
	if c.isT0 {
		if len(c.localeServiceID) > 0 {
			client := cliTier0LocaleServiceIpsecVpnLocalEndpointsClient(c.sessionContext, connector)
			if client == nil {
				return fmt.Errorf("unsupported client type")
			}
			return client.Delete(c.gwID, c.localeServiceID, c.serviceID, id)
		}
		client := cliTier0IpsecVpnLocalEndpointsClient(c.sessionContext, connector)
		if client == nil {
			return fmt.Errorf("unsupported client type")
		}
		return client.Delete(c.gwID, c.serviceID, id)

	}
	if len(c.localeServiceID) > 0 {
		if c.sessionContext.ClientType == utl.Multitenancy {
			return fmt.Errorf("project context is not supported for locale-service scoped IPSec VPN local endpoints")
		}
		client := cliTier1LocaleServiceIpsecVpnLocalEndpointsClient(c.sessionContext, connector)
		if client == nil {
			return fmt.Errorf("unsupported client type")
		}
		return client.Delete(c.gwID, c.localeServiceID, c.serviceID, id)
	}
	client := cliTier1IpsecVpnLocalEndpointsClient(c.sessionContext, connector)
	if client == nil {
		return fmt.Errorf("unsupported client type")
	}
	return client.Delete(c.gwID, c.serviceID, id)
}

func resourceNsxtPolicyIPSecVpnLocalEndpointExistsOnService(id string, connector client.Connector, servicePath string) (bool, error) {
	// Used by acceptance tests and create flow; infer multitenancy from servicePath when applicable.
	sessionContext := utl.SessionContext{ClientType: utl.Local}
	if strings.HasPrefix(servicePath, "/orgs/") {
		projectID := extractProjectIDFromPolicyPath(servicePath)
		if projectID != "" {
			sessionContext.ClientType = utl.Multitenancy
			sessionContext.ProjectID = projectID
		}
	}
	client, err := newLocalEndpointClient(servicePath, sessionContext)
	if err != nil {
		return false, err
	}
	_, err = client.Get(connector, id)
	if err == nil {
		return true, nil
	}

	if isNotFoundError(err) {
		return false, nil
	}

	return false, logAPIError("Error retrieving resource", err)
}

func resourceNsxtPolicyIPSecVpnLocalEndpointExists(servicePath string) func(id string, connector client.Connector, isGlobalManager bool) (bool, error) {
	return func(id string, connector client.Connector, isGlobalManager bool) (bool, error) {
		return resourceNsxtPolicyIPSecVpnLocalEndpointExistsOnService(id, connector, servicePath)
	}
}

func ipSecVpnLocalEndpointInitStruct(d *schema.ResourceData) model.IPSecVpnLocalEndpoint {
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	tags := getPolicyTagsFromSchema(d)
	certificatePath := d.Get("certificate_path").(string)
	localAddress := d.Get("local_address").(string)
	localID := d.Get("local_id").(string)
	trustCaPaths := getStringListFromSchemaSet(d, "trust_ca_paths")
	trustCrlPaths := getStringListFromSchemaSet(d, "trust_crl_paths")

	obj := model.IPSecVpnLocalEndpoint{
		DisplayName:  &displayName,
		Description:  &description,
		Tags:         tags,
		LocalAddress: &localAddress,
	}

	if len(localID) > 0 {
		obj.LocalId = &localID
	}

	if len(certificatePath) > 0 {
		obj.CertificatePath = &certificatePath
	}

	if len(trustCaPaths) > 0 {
		obj.TrustCaPaths = trustCaPaths
	}

	if len(trustCrlPaths) > 0 {
		obj.TrustCrlPaths = trustCrlPaths
	}

	return obj
}

func resourceNsxtPolicyIPSecVpnLocalEndpointCreate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	servicePath := d.Get("service_path").(string)
	sessionContext := getIPSecVpnLocalEndpointSessionContext(d, m, servicePath)
	existsFunc := func(id string, connector client.Connector, isGlobalManager bool) (bool, error) {
		client, err := newLocalEndpointClient(servicePath, sessionContext)
		if err != nil {
			return false, err
		}
		_, err = client.Get(connector, id)
		if err == nil {
			return true, nil
		}
		if isNotFoundError(err) {
			return false, nil
		}
		return false, logAPIError("Error retrieving resource", err)
	}
	id, err := getOrGenerateID(d, m, existsFunc)
	if err != nil {
		return err
	}

	obj := ipSecVpnLocalEndpointInitStruct(d)

	log.Printf("[INFO] Creating IPSecVpnLocalEndpoint with ID %s", id)
	client, err := newLocalEndpointClient(servicePath, sessionContext)
	if err != nil {
		return handleCreateError("IPSecVpnLocalEndpoint", id, err)
	}
	err = client.Patch(connector, id, obj)
	if err != nil {
		return handleCreateError("IPSecVpnLocalEndpoint", id, err)
	}

	d.SetId(id)
	d.Set("nsx_id", id)

	return resourceNsxtPolicyIPSecVpnLocalEndpointRead(d, m)
}

func resourceNsxtPolicyIPSecVpnLocalEndpointRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IPSecVpnLocalEndpoint ID")
	}

	servicePath := d.Get("service_path").(string)
	sessionContext := getIPSecVpnLocalEndpointSessionContext(d, m, servicePath)
	client, err := newLocalEndpointClient(servicePath, sessionContext)
	if err != nil {
		return handleReadError(d, "IPSecVpnLocalEndpoint", id, err)
	}
	obj, err := client.Get(connector, id)
	if err != nil {
		return handleReadError(d, "IPSecVpnLocalEndpoint", id, err)
	}

	d.Set("display_name", obj.DisplayName)
	d.Set("description", obj.Description)
	setPolicyTagsInSchema(d, obj.Tags)
	d.Set("nsx_id", id)
	d.Set("path", obj.Path)
	d.Set("revision", obj.Revision)
	d.Set("certificate_path", obj.CertificatePath)
	d.Set("local_address", obj.LocalAddress)
	d.Set("local_id", obj.LocalId)
	d.Set("trust_ca_paths", obj.TrustCaPaths)
	d.Set("trust_crl_paths", obj.TrustCrlPaths)

	return nil
}

func resourceNsxtPolicyIPSecVpnLocalEndpointUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IPSecVpnLocalEndpoint ID")
	}

	servicePath := d.Get("service_path").(string)
	sessionContext := getIPSecVpnLocalEndpointSessionContext(d, m, servicePath)
	client, err := newLocalEndpointClient(servicePath, sessionContext)
	if err != nil {
		return handleUpdateError("IPSecVpnLocalEndpoint", id, err)
	}

	obj := ipSecVpnLocalEndpointInitStruct(d)
	revision := int64(d.Get("revision").(int))
	obj.Revision = &revision

	log.Printf("[INFO] Updating IPSecVpnLocalEndpoint with ID %s", id)
	err = client.Patch(connector, id, obj)
	if err != nil {
		return handleUpdateError("IPSecVpnLocalEndpoint", id, err)
	}

	return resourceNsxtPolicyIPSecVpnLocalEndpointRead(d, m)
}

func resourceNsxtPolicyIPSecVpnLocalEndpointDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	if id == "" {
		return fmt.Errorf("Error obtaining IPSecVpnLocalEndpoint ID")
	}

	servicePath := d.Get("service_path").(string)
	sessionContext := getIPSecVpnLocalEndpointSessionContext(d, m, servicePath)
	client, err := newLocalEndpointClient(servicePath, sessionContext)
	if err != nil {
		return handleUpdateError("IPSecVpnLocalEndpoint", id, err)
	}
	connector := getPolicyConnector(m)
	err = client.Delete(connector, id)

	if err != nil {
		return handleDeleteError("IPSecVpnLocalEndpoint", id, err)
	}

	return nil
}

func nsxtVPNServiceResourceImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	err := fmt.Errorf("VPN Local Endpoint Path expected, got %s", importID)
	// path format example: /infra/tier-1s/aaa/locale-services/default/ipsec-vpn-services/bbb/local-endpoints/ccc
	s := strings.Split(importID, "/")
	if len(s) < 8 {
		return []*schema.ResourceData{d}, err
	}

	endpointID := s[len(s)-1]
	d.SetId(endpointID)
	s = strings.Split(importID, "/local-endpoints/")
	if len(s) != 2 {
		return []*schema.ResourceData{d}, err
	}

	d.Set("service_path", s[0])

	return []*schema.ResourceData{d}, nil
}

func validateIPv6LocaleServiceConflict(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	return validateIPv6WithLocaleService(diff, "local_address", "validateIPv6LocaleServiceConflict")
}

func validateIPv6WithLocaleService(diff *schema.ResourceDiff, fieldName, logPrefix string) error {
	address := diff.Get(fieldName).(string)
	servicePath := diff.Get("service_path").(string)
	ip := net.ParseIP(address)
	isIPv6 := ip != nil && strings.Contains(address, ":")
	isLocaleService := strings.Contains(servicePath, "/locale-services/")

	if isIPv6 && isLocaleService {
		return fmt.Errorf("IPv6 addresses are not supported for VPN services configured with locale service path. Please use a VPN service configured with a gateway path")
	}
	return nil
}
