---
subcategory: "EVPN"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_evpn_config"
description: A resource to configure EVPN Settings of Tier0 Gateway.
---

# nsxt_policy_evpn_config

This resource provides a method to configure EVPN on T0 Gateway. A single resource should be configured per Gateway.

This resource is applicable to NSX Policy Manager only.
This resource is supported with NSX 3.1.0 onwards.

## Example Usage

```hcl
resource "nsxt_policy_evpn_config" "test" {
  display_name = "evpn"
  gateway_path = nsxt_policy_tier0_gateway.gw1.path

  mode          = "INLINE"
  vni_pool_path = data.nsxt_policy_vni_pool.pool1.path
}
```

## Argument Reference

The following arguments are supported:

  * `display_name` - (Required) Display name for the resource.
  * `description` - (Optional) Description for the resource.
  * `gateway_path` - (Required) Policy Path for Tier0 Gateway to configure EVPN on.
  * `mode` - (Required) EVPN Mode, one of `INLINE` or `ROUTE_SERVER`. In `ROUTE_SERVER` mode, edge nodes participate in the BGP EVPN control plane route exchanges only and do not participate in data forwarding.
  * `vni_pool_path` - (Optional) Path of VNI pool to use. This setting is only applicable (and required) with `INLINE` mode.
  * `evpn_tenant_path` - (Optional) Policy path for EVPN tenant. Relevant for `ROUTE_SERVER` mode.
  * `tag` - (Optional) A list of scope + tag pairs to associate with this resource.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing EVPN Config can be [imported][docs-import] into this resource, via the following command:

 [docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_evpn_config.config1 gwPath
```

The above command imports EVPN Config named `config1` for NSX Policy Tier0 Gateway with full Policy Path `gwPath`.

~> **NOTE:** Note that import parameter is non-standard here. Please make sure you use full policy path for the gateway, such as `/infra/tier-0s/mygateway`
