---
subcategory: "Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_ospf_config"
description: A resource to configure OSPF Settings of Tier0 Gateway on NSX Policy Manager.
---

# nsxt_policy_ospf_config

This resource provides a method for the management of OSPF for T0 Gateway on default locale service. A single resource should be specified per T0 Gateway. Edge Cluster is expected to be configured on the Gateway.

This resource is applicable to NSX Policy Manager only.
This resource is supported with NSX 3.1.1 onwards.

~> **NOTE:** NSX does not support deleting OSPF config on gateway, therefore this resource will update NSX object, but never delete it. To undo OSPF configuration, please disable it within the resource.

## Example Usage

```hcl
resource "nsxt_policy_ospf_config" "test" {
  gateway_path = nsxt_policy_tier0_gateway.gw1.path

  enabled               = true
  ecmp                  = true
  default_originate     = false
  graceful_restart_mode = "HELPER_ONLY"

  summary_address {
    prefix    = "20.1.0.0/24"
    advertise = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `gateway_path` - (Required) Policy Path for Tier0 Gateway.
* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `ecmp` - (Optional) A boolean flag to enable/disable ECMP. Default is `true`.
* `enabled` - (Optional) A boolean flag to enable/disable OSPF. Default is `true`.
* `default_originate` - (Optional) A boolean flag to configure advertisement of default route into OSPF domain. Default is `false`.
* `graceful_restart_mode` - (Optional) Graceful Restart Mode, one of `HELPER_ONLY` or `DISABLED`. Defaut is `HELPER_ONLY`.
* `summary_address`- (Optional) Repeatable block to define addresses to summarize or filter external routes.
  * `prefix` - (Required) OSPF Summary address in CIDR format.
  * `advertise` - (Optional) A boolean flag to configure advertisement of external routes into the OSPF domain. Default is `true`.
* `tag` - (Optional) A list of scope + tag pairs to associate with this Tier-0 gateway's OSPF configuration.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.
* `gateway_id` - Tier0 Gateway ID on which OSPF is configured.
* `locale_service_id` - Tier0 Gateway Locale Service ID on which OSPF is configured.

## Importing

Importing the resource is not supported - creating the resource would update it to desired state on backend.
