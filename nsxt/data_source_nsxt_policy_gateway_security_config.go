// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package nsxt

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vapiProtocolClient "github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_0s"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_1s"
)

// Package-level client functions for testability
var cliTier0SecurityConfigClientDS = func(connector vapiProtocolClient.Connector) tier_0s.SecurityConfigClient {
	return tier_0s.NewSecurityConfigClient(connector)
}

var cliTier1SecurityConfigClientDS = func(connector vapiProtocolClient.Connector) tier_1s.SecurityConfigClient {
	return tier_1s.NewSecurityConfigClient(connector)
}

func dataSourceNsxtPolicyGatewaySecurityConfig() *schema.Resource {
	return &schema.Resource{
		Read:        dataSourceNsxtPolicyGatewaySecurityConfigRead,
		Description: "Data source to read security feature configuration on Tier-0 and Tier-1 gateways. Retrieves the status of North-South traffic security features such as IDPS, IDFW, Malware Prevention, and TLS Inspection.",

		Schema: map[string]*schema.Schema{
			"tier0_id": {
				Type:         schema.TypeString,
				Description:  "ID of the Tier-0 gateway. Exactly one of tier0_id or tier1_id must be specified.",
				Optional:     true,
				ExactlyOneOf: []string{"tier0_id", "tier1_id"},
			},
			"tier1_id": {
				Type:         schema.TypeString,
				Description:  "ID of the Tier-1 gateway. Exactly one of tier0_id or tier1_id must be specified.",
				Optional:     true,
				ExactlyOneOf: []string{"tier0_id", "tier1_id"},
			},
			"idps_enabled": {
				Type:        schema.TypeBool,
				Description: "Whether Intrusion Detection and Prevention System (IDPS) is enabled on the gateway.",
				Computed:    true,
			},
			"idfw_enabled": {
				Type:        schema.TypeBool,
				Description: "Whether Identity Firewall (IDFW) is enabled on the gateway.",
				Computed:    true,
			},
			"malware_prevention_enabled": {
				Type:        schema.TypeBool,
				Description: "Whether Malware Prevention is enabled on the gateway. Only supported on Tier-1 gateways.",
				Computed:    true,
			},
			"tls_enabled": {
				Type:        schema.TypeBool,
				Description: "Whether TLS (Transport Layer Security) Inspection is enabled on the gateway. Only supported on Tier-1 gateways.",
				Computed:    true,
			},
			"path": {
				Type:        schema.TypeString,
				Description: "NSX path of the gateway security configuration.",
				Computed:    true,
			},
		},
	}
}

func dataSourceNsxtPolicyGatewaySecurityConfigRead(d *schema.ResourceData, m interface{}) error {
	connector := getPolicyConnector(m)

	tier0ID := d.Get("tier0_id").(string)
	tier1ID := d.Get("tier1_id").(string)

	if tier0ID != "" {
		client := cliTier0SecurityConfigClientDS(connector)
		config, err := client.Get(tier0ID, nil, nil, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("error reading Tier0 Gateway Security Config for gateway %s: %v", tier0ID, err)
		}
		d.SetId("tier0/" + tier0ID)
		return setTier0GatewaySecurityConfigInSchema(d, config, true)
	}

	if tier1ID != "" {
		client := cliTier1SecurityConfigClientDS(connector)
		config, err := client.Get(tier1ID, nil, nil, nil, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("error reading Tier1 Gateway Security Config for gateway %s: %v", tier1ID, err)
		}
		d.SetId("tier1/" + tier1ID)
		return setTier1GatewaySecurityConfigInSchema(d, config, true)
	}

	return fmt.Errorf("one of tier0_id or tier1_id must be specified")
}
