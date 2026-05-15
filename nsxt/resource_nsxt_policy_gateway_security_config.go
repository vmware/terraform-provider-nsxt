// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_1s"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

// Package-level client functions for testability
var cliTier0SecurityConfigClient = func(connector vapiProtocolClient.Connector) tier_0s.SecurityConfigClient {
	return tier_0s.NewSecurityConfigClient(connector)
}

var cliTier1SecurityConfigClient = func(connector vapiProtocolClient.Connector) tier_1s.SecurityConfigClient {
	return tier_1s.NewSecurityConfigClient(connector)
}

func resourceNsxtPolicyGatewaySecurityConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceNsxtPolicyGatewaySecurityConfigCreate,
		Read:   resourceNsxtPolicyGatewaySecurityConfigRead,
		Update: resourceNsxtPolicyGatewaySecurityConfigUpdate,
		Delete: resourceNsxtPolicyGatewaySecurityConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: gatewaySecurityConfigImporter,
		},
		Description: "Manages security feature configuration on Tier-0 and Tier-1 gateways. Used to enable or disable North-South traffic security features such as IDPS, IDFW, Malware Prevention, and TLS Inspection.",

		Schema: map[string]*schema.Schema{
			"tier0_id": {
				Type:         schema.TypeString,
				Description:  "ID of the Tier-0 gateway. Exactly one of tier0_id or tier1_id must be specified. Changing this forces a new resource to be created.",
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"tier0_id", "tier1_id"},
			},
			"tier1_id": {
				Type:         schema.TypeString,
				Description:  "ID of the Tier-1 gateway. Exactly one of tier0_id or tier1_id must be specified. Changing this forces a new resource to be created.",
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"tier0_id", "tier1_id"},
			},
			"idps_enabled": {
				Type:        schema.TypeBool,
				Description: "Enable or disable Intrusion Detection and Prevention System (IDPS) on the gateway. Supported for both Tier-0 and Tier-1 gateways.",
				Optional:    true,
				Default:     false,
			},
			"idfw_enabled": {
				Type:        schema.TypeBool,
				Description: "Enable or disable Identity Firewall (IDFW) on the gateway. Supported for both Tier-0 and Tier-1 gateways.",
				Optional:    true,
				Default:     false,
			},
			"malware_prevention_enabled": {
				Type:        schema.TypeBool,
				Description: "Enable or disable Malware Prevention on the gateway. Supported for Tier-1 gateways only.",
				Optional:    true,
				Default:     false,
			},
			"tls_enabled": {
				Type:        schema.TypeBool,
				Description: "Enable or disable TLS (Transport Layer Security) Inspection on the gateway. Supported for Tier-1 gateways only.",
				Optional:    true,
				Default:     false,
			},
			"path": {
				Type:        schema.TypeString,
				Description: "NSX path of the gateway security configuration.",
				Computed:    true,
			},
			"revision": {
				Type:        schema.TypeInt,
				Description: "Revision number of the gateway security configuration.",
				Computed:    true,
			},
		},
	}
}

// gatewaySecurityConfigImporter parses import IDs in the format "tier0/<id>" or "tier1/<id>".
func gatewaySecurityConfigImporter(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()
	parts := strings.SplitN(importID, "/", 2)
	if len(parts) != 2 || parts[1] == "" {
		return nil, fmt.Errorf("invalid import ID %q: expected format 'tier0/<gateway-id>' or 'tier1/<gateway-id>'", importID)
	}
	gwType := parts[0]
	gwID := parts[1]
	switch gwType {
	case "tier0":
		d.Set("tier0_id", gwID)
	case "tier1":
		d.Set("tier1_id", gwID)
	default:
		return nil, fmt.Errorf("invalid gateway type %q in import ID: must be 'tier0' or 'tier1'", gwType)
	}
	d.SetId(importID)
	return []*schema.ResourceData{d}, nil
}

func gatewaySecurityConfigID(d *schema.ResourceData) string {
	if tier0ID := d.Get("tier0_id").(string); tier0ID != "" {
		return "tier0/" + tier0ID
	}
	return "tier1/" + d.Get("tier1_id").(string)
}

func resourceNsxtPolicyGatewaySecurityConfigCreate(d *schema.ResourceData, m interface{}) error {
	d.SetId(gatewaySecurityConfigID(d))
	return resourceNsxtPolicyGatewaySecurityConfigUpdate(d, m)
}

func resourceNsxtPolicyGatewaySecurityConfigRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining Gateway Security Config ID")
	}

	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid Gateway Security Config ID %q", id)
	}
	gwType := parts[0]
	gwID := parts[1]

	switch gwType {
	case "tier0":
		client := cliTier0SecurityConfigClient(connector)
		config, err := client.Get(gwID, nil, nil, nil, nil, nil, nil)
		if err != nil {
			return handleReadError(d, "Tier0 Gateway Security Config", gwID, err)
		}
		d.Set("tier0_id", gwID)
		return setTier0GatewaySecurityConfigInSchema(d, config, false)

	case "tier1":
		client := cliTier1SecurityConfigClient(connector)
		config, err := client.Get(gwID, nil, nil, nil, nil, nil, nil)
		if err != nil {
			return handleReadError(d, "Tier1 Gateway Security Config", gwID, err)
		}
		d.Set("tier1_id", gwID)
		return setTier1GatewaySecurityConfigInSchema(d, config, false)

	default:
		return fmt.Errorf("invalid gateway type %q in ID %q", gwType, id)
	}
}

func resourceNsxtPolicyGatewaySecurityConfigUpdate(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining Gateway Security Config ID")
	}

	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid Gateway Security Config ID %q", id)
	}
	gwType := parts[0]
	gwID := parts[1]

	switch gwType {
	case "tier0":
		client := cliTier0SecurityConfigClient(connector)
		config := buildTier0SecurityFeatures(d)
		log.Printf("[INFO] Updating Tier0 Gateway Security Config for gateway %s", gwID)
		_, err := client.Patch(gwID, config)
		if err != nil {
			return handleUpdateError("Tier0 Gateway Security Config", gwID, err)
		}

	case "tier1":
		client := cliTier1SecurityConfigClient(connector)
		config := buildTier1SecurityFeatures(d)
		log.Printf("[INFO] Updating Tier1 Gateway Security Config for gateway %s", gwID)
		_, err := client.Patch(gwID, config)
		if err != nil {
			return handleUpdateError("Tier1 Gateway Security Config", gwID, err)
		}

	default:
		return fmt.Errorf("invalid gateway type %q in ID %q", gwType, id)
	}

	return resourceNsxtPolicyGatewaySecurityConfigRead(d, m)
}

func resourceNsxtPolicyGatewaySecurityConfigDelete(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)
	id := d.Id()
	if id == "" {
		return fmt.Errorf("error obtaining Gateway Security Config ID")
	}

	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid Gateway Security Config ID %q", id)
	}
	gwType := parts[0]
	gwID := parts[1]

	// Gateway security config cannot be deleted; on destroy we disable all features.
	log.Printf("[INFO] Disabling all security features for %s gateway %s (DELETE operation)", gwType, gwID)

	switch gwType {
	case "tier0":
		client := cliTier0SecurityConfigClient(connector)
		config := disabledTier0SecurityFeatures()
		_, err := client.Patch(gwID, config)
		if err != nil {
			if isNotFoundError(err) {
				return nil
			}
			return handleDeleteError("Tier0 Gateway Security Config", gwID, err)
		}

	case "tier1":
		client := cliTier1SecurityConfigClient(connector)
		config := disabledTier1SecurityFeatures()
		_, err := client.Patch(gwID, config)
		if err != nil {
			if isNotFoundError(err) {
				return nil
			}
			return handleDeleteError("Tier1 Gateway Security Config", gwID, err)
		}

	default:
		return fmt.Errorf("invalid gateway type %q in ID %q", gwType, id)
	}

	return nil
}

func buildTier0SecurityFeatures(d *schema.ResourceData) model.Tier0SecurityFeatures {
	idpsEnabled := d.Get("idps_enabled").(bool)
	idfwEnabled := d.Get("idfw_enabled").(bool)

	idpsFeature := model.Tier0SecurityFeature_FEATURE_IDPS
	idfwFeature := model.Tier0SecurityFeature_FEATURE_IDFW

	// Only include features that are supported by Tier0 gateways
	// Based on API validation errors, Tier0 only supports IDFW and IDPS
	return model.Tier0SecurityFeatures{
		Features: []model.Tier0SecurityFeature{
			{Feature: &idpsFeature, Enable: &idpsEnabled},
			{Feature: &idfwFeature, Enable: &idfwEnabled},
		},
	}
}

func buildTier1SecurityFeatures(d *schema.ResourceData) model.SecurityFeatures {
	idpsEnabled := d.Get("idps_enabled").(bool)
	idfwEnabled := d.Get("idfw_enabled").(bool)
	malwareEnabled := d.Get("malware_prevention_enabled").(bool)
	tlsEnabled := d.Get("tls_enabled").(bool)

	idpsFeature := model.SecurityFeature_FEATURE_IDPS
	idfwFeature := model.SecurityFeature_FEATURE_IDFW
	malwareFeature := model.SecurityFeature_FEATURE_MALWAREPREVENTION
	tlsFeature := model.SecurityFeature_FEATURE_TLS

	// Only include features that are supported by Tier1 gateways
	// Based on API validation errors, Tier1 supports MALWAREPREVENTION, IDFW, IDPS, and TLS
	// but NOT GEOIP_MONITORING or GFW_MULTICAST
	return model.SecurityFeatures{
		Features: []model.SecurityFeature{
			{Feature: &idpsFeature, Enable: &idpsEnabled},
			{Feature: &idfwFeature, Enable: &idfwEnabled},
			{Feature: &malwareFeature, Enable: &malwareEnabled},
			{Feature: &tlsFeature, Enable: &tlsEnabled},
		},
	}
}

func disabledTier0SecurityFeatures() model.Tier0SecurityFeatures {
	disabled := false
	idpsFeature := model.Tier0SecurityFeature_FEATURE_IDPS
	idfwFeature := model.Tier0SecurityFeature_FEATURE_IDFW

	// Only disable features that are supported by Tier0 gateways (IDFW and IDPS)
	return model.Tier0SecurityFeatures{
		Features: []model.Tier0SecurityFeature{
			{Feature: &idpsFeature, Enable: &disabled},
			{Feature: &idfwFeature, Enable: &disabled},
		},
	}
}

func disabledTier1SecurityFeatures() model.SecurityFeatures {
	disabled := false
	idpsFeature := model.SecurityFeature_FEATURE_IDPS
	idfwFeature := model.SecurityFeature_FEATURE_IDFW
	malwareFeature := model.SecurityFeature_FEATURE_MALWAREPREVENTION
	tlsFeature := model.SecurityFeature_FEATURE_TLS

	// Only disable features that are supported by Tier1 gateways (MALWAREPREVENTION, IDFW, IDPS, TLS)
	return model.SecurityFeatures{
		Features: []model.SecurityFeature{
			{Feature: &idpsFeature, Enable: &disabled},
			{Feature: &idfwFeature, Enable: &disabled},
			{Feature: &malwareFeature, Enable: &disabled},
			{Feature: &tlsFeature, Enable: &disabled},
		},
	}
}

func setTier0GatewaySecurityConfigInSchema(d *schema.ResourceData, config model.Tier0SecurityFeatures, isDataSource bool) error {
	featureMap := make(map[string]bool)
	for _, f := range config.Features {
		if f.Feature != nil && f.Enable != nil {
			featureMap[*f.Feature] = *f.Enable
		}
	}

	if v, ok := featureMap[model.Tier0SecurityFeature_FEATURE_IDPS]; ok {
		d.Set("idps_enabled", v)
	}
	if v, ok := featureMap[model.Tier0SecurityFeature_FEATURE_IDFW]; ok {
		d.Set("idfw_enabled", v)
	}

	// The following features are not supported on Tier-0 gateways; always set to false to keep state consistent.
	d.Set("malware_prevention_enabled", false)
	d.Set("tls_enabled", false)

	if config.Path != nil {
		d.Set("path", *config.Path)
	}
	if !isDataSource && config.Revision != nil {
		d.Set("revision", int(*config.Revision))
	}

	return nil
}

func setTier1GatewaySecurityConfigInSchema(d *schema.ResourceData, config model.SecurityFeatures, isDataSource bool) error {
	featureMap := make(map[string]bool)
	for _, f := range config.Features {
		if f.Feature != nil && f.Enable != nil {
			featureMap[*f.Feature] = *f.Enable
		}
	}

	if v, ok := featureMap[model.SecurityFeature_FEATURE_IDPS]; ok {
		d.Set("idps_enabled", v)
	}
	if v, ok := featureMap[model.SecurityFeature_FEATURE_IDFW]; ok {
		d.Set("idfw_enabled", v)
	}
	if v, ok := featureMap[model.SecurityFeature_FEATURE_MALWAREPREVENTION]; ok {
		d.Set("malware_prevention_enabled", v)
	}
	if v, ok := featureMap[model.SecurityFeature_FEATURE_TLS]; ok {
		d.Set("tls_enabled", v)
	}

	if config.Path != nil {
		d.Set("path", *config.Path)
	}
	if !isDataSource && config.Revision != nil {
		d.Set("revision", int(*config.Revision))
	}

	return nil
}
