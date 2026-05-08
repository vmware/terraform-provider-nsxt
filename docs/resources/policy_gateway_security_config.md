---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_gateway_security_config"
description: A resource to configure security features on Tier-0 and Tier-1 gateways for North-South traffic inspection.
---

# nsxt_policy_gateway_security_config

This resource manages security feature configuration on NSX-T Tier-0 and Tier-1 gateways. It is used to activate or deactivate North-South traffic security features including IDPS (Intrusion Detection and Prevention System), IDFW (Identity Firewall), Malware Prevention, and TLS Inspection.

**Note:** Feature availability varies by gateway type. Tier-0 gateways support only IDPS and IDFW features, while Tier-1 gateways support IDPS, IDFW, Malware Prevention, and TLS features. See the Feature Availability table below for details.

This corresponds to the **Security > IDPS > Settings > Shared > Activate Gateways for North-South Traffic** configuration in the NSX-T UI.

This resource manages existing security configuration that is automatically created by NSX when a gateway is provisioned. It cannot create new gateway security configurations or delete them — on destroy, it disables all configured features.

## Example Usage

### Enable IDPS on a Tier-1 Gateway

```hcl
resource "nsxt_policy_gateway_security_config" "tier1_idps" {
  tier1_id     = "test"
  idps_enabled = true
}
```

### Enable IDPS on a Tier-0 Gateway

```hcl
resource "nsxt_policy_gateway_security_config" "tier0_idps" {
  tier0_id     = "tier0-gw"
  idps_enabled = true
}
```

### Enable Multiple Security Features on a Tier-1 Gateway

```hcl
resource "nsxt_policy_gateway_security_config" "tier1_security" {
  tier1_id                   = "tier1-gw"
  idps_enabled               = true
  idfw_enabled               = true
  malware_prevention_enabled = true
  tls_enabled                = true
}
```

### Enable Security Features on a Tier-0 Gateway Using Data Source

```hcl
data "nsxt_policy_tier0_gateway" "tier0" {
  display_name = "DefaultT0Gateway"
}

resource "nsxt_policy_gateway_security_config" "tier0_security" {
  tier0_id     = data.nsxt_policy_tier0_gateway.tier0.id
  idps_enabled = true
  idfw_enabled = true
}
```

### Enable IDPS on Tier-1 Gateway (IDPS Prerequisite)

```hcl
data "nsxt_policy_tier1_gateway" "tier1" {
  display_name = "Tier1-GW-01"
}

# Activate IDPS on the gateway for North-South traffic inspection
resource "nsxt_policy_gateway_security_config" "tier1_idps" {
  tier1_id     = data.nsxt_policy_tier1_gateway.tier1.id
  idps_enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `tier0_id` - (Optional) The ID of the Tier-0 gateway. Exactly one of `tier0_id` or `tier1_id` must be specified. Changing this forces a new resource to be created.
* `tier1_id` - (Optional) The ID of the Tier-1 gateway. Exactly one of `tier0_id` or `tier1_id` must be specified. Changing this forces a new resource to be created.
* `idps_enabled` - (Optional) Enable or disable Intrusion Detection and Prevention System (IDPS) on the gateway. Supported for both Tier-0 and Tier-1 gateways. Defaults to `false`.
* `idfw_enabled` - (Optional) Enable or disable Identity Firewall (IDFW) on the gateway. Supported for both Tier-0 and Tier-1 gateways. Defaults to `false`.
* `malware_prevention_enabled` - (Optional) Enable or disable Malware Prevention on the gateway. **Supported for Tier-1 gateways only.** Setting this on a Tier-0 gateway will be ignored (always `false`). Defaults to `false`.
* `tls_enabled` - (Optional) Enable or disable TLS (Transport Layer Security) Inspection on the gateway. **Supported for Tier-1 gateways only.** Setting this on a Tier-0 gateway will be ignored (always `false`). Defaults to `false`.

## Attribute Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - The resource ID in the format `tier0/<gateway-id>` or `tier1/<gateway-id>`.
* `path` - The NSX path of the gateway security configuration.
* `revision` - The revision number of the gateway security configuration.

## Importing

An existing gateway security configuration can be [imported][docs-import] into this resource using the following command:

[docs-import]: https://www.terraform.io/cli/import

```shell
terraform import nsxt_policy_gateway_security_config.example tier1/<gateway-id>
```

Or for a Tier-0 gateway:

```shell
terraform import nsxt_policy_gateway_security_config.example tier0/<gateway-id>
```

For example:

```shell
terraform import nsxt_policy_gateway_security_config.tier1_idps tier1/test
terraform import nsxt_policy_gateway_security_config.tier0_idps tier0/tier0-gw
```

## Important Notes

### Gateway Security Config Lifecycle

* **Auto-created**: Gateway security configuration is automatically created by NSX when a gateway is provisioned. This resource does not create new configurations.
* **Cannot delete**: Gateway security configurations cannot be deleted via the NSX API. When this resource is destroyed, it disables all configured security features but the configuration itself remains in NSX.
* **Update only**: This resource manages the security feature settings of an existing gateway configuration.

### Feature Availability by Gateway Type

| Feature | Tier-0 | Tier-1 | Notes |
| --- | --- | --- | --- |
| `idps_enabled` | ✅ Yes | ✅ Yes | Fully supported on both gateway types |
| `idfw_enabled` | ✅ Yes | ✅ Yes | Fully supported on both gateway types |
| `malware_prevention_enabled` | ❌ No | ✅ Yes | Only supported on Tier-1 gateways |
| `tls_enabled` | ❌ No | ✅ Yes | Only supported on Tier-1 gateways |

**Important:** Unsupported features will be silently ignored and always return `false` in the state, even if specified in the configuration.

### NSX-T UI Mapping

This resource corresponds to the following NSX-T UI path:

`Security > IDPS > Settings > Shared > Activate Gateways for North-South Traffic`

### IDPS Prerequisites

For IDPS to function on a gateway:

1. The gateway must be activated for North-South traffic inspection using this resource.
2. An IDPS Gateway Policy (`nsxt_policy_intrusion_service_gateway_policy`) must be configured and applied.
3. An IDPS Profile must reference the appropriate signatures.

## See Also

* [nsxt_policy_tier0_gateway](../data-sources/policy_tier0_gateway.md) data source
* [nsxt_policy_tier1_gateway](../data-sources/policy_tier1_gateway.md) data source
* [nsxt_policy_gateway_security_config](../data-sources/policy_gateway_security_config.md) data source
* [nsxt_policy_intrusion_service_gateway_policy](./policy_intrusion_service_gateway_policy.md) resource
