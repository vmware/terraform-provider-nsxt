---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_gateway_security_config"
description: A data source to read security feature configuration on Tier-0 and Tier-1 gateways.
---

# nsxt_policy_gateway_security_config

This data source provides information about the security feature configuration on NSX-T Tier-0 and Tier-1 gateways. It retrieves the current status of North-South traffic security features including IDPS, IDFW, Malware Prevention, and TLS Inspection.

**Note:** Feature availability varies by gateway type. Tier-0 gateways support only IDPS and IDFW features, while Tier-1 gateways support IDPS, IDFW, Malware Prevention, and TLS features.

## Example Usage

### Read Security Config for a Tier-1 Gateway

```hcl
data "nsxt_policy_gateway_security_config" "tier1_security" {
  tier1_id = "test"
}

output "idps_enabled" {
  value = data.nsxt_policy_gateway_security_config.tier1_security.idps_enabled
}
```

### Read Security Config for a Tier-0 Gateway

```hcl
data "nsxt_policy_tier0_gateway" "tier0" {
  display_name = "DefaultT0Gateway"
}

data "nsxt_policy_gateway_security_config" "tier0_security" {
  tier0_id = data.nsxt_policy_tier0_gateway.tier0.id
}

output "tier0_idps_enabled" {
  value = data.nsxt_policy_gateway_security_config.tier0_security.idps_enabled
}
```

### Check Security Feature Status Before Creating Policy

```hcl
data "nsxt_policy_tier1_gateway" "tier1" {
  display_name = "Tier1-GW-01"
}

data "nsxt_policy_gateway_security_config" "tier1_security" {
  tier1_id = data.nsxt_policy_tier1_gateway.tier1.id
}

output "security_summary" {
  value = {
    idps_enabled               = data.nsxt_policy_gateway_security_config.tier1_security.idps_enabled
    idfw_enabled               = data.nsxt_policy_gateway_security_config.tier1_security.idfw_enabled
    malware_prevention_enabled = data.nsxt_policy_gateway_security_config.tier1_security.malware_prevention_enabled
    tls_enabled                = data.nsxt_policy_gateway_security_config.tier1_security.tls_enabled
  }
}
```

## Argument Reference

The following arguments are supported:

* `tier0_id` - (Optional) The ID of the Tier-0 gateway. Exactly one of `tier0_id` or `tier1_id` must be specified.
* `tier1_id` - (Optional) The ID of the Tier-1 gateway. Exactly one of `tier0_id` or `tier1_id` must be specified.

## Attribute Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - The data source ID in the format `tier0/<gateway-id>` or `tier1/<gateway-id>`.
* `idps_enabled` - Whether Intrusion Detection and Prevention System (IDPS) is enabled on the gateway. Supported on both Tier-0 and Tier-1 gateways.
* `idfw_enabled` - Whether Identity Firewall (IDFW) is enabled on the gateway. Supported on both Tier-0 and Tier-1 gateways.
* `malware_prevention_enabled` - Whether Malware Prevention is enabled on the gateway. Only supported on Tier-1 gateways (always `false` for Tier-0).
* `tls_enabled` - Whether TLS (Transport Layer Security) Inspection is enabled on the gateway. Only supported on Tier-1 gateways (always `false` for Tier-0).
* `path` - The NSX path of the gateway security configuration.

## See Also

* [nsxt_policy_tier0_gateway](./policy_tier0_gateway.md) data source
* [nsxt_policy_tier1_gateway](./policy_tier1_gateway.md) data source
* [nsxt_policy_gateway_security_config](../resources/policy_gateway_security_config.md) resource
