---
subcategory: "VPN"
layout: "nsxt"
page_title: "NSXT: policy_ipsec_vpn_local_endpoint"
description: Policy IPSec VPN Local Endpoint data source.
---

# nsxt_policy_ipsec_vpn_local_endpoint

This data source provides information about IPSec VPN policy local endpoint configured on NSX.

This data source is applicable to NSX Policy Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_ipsec_vpn_local_endpoint" "test" {
  display_name = "local_endpoint1"
}
```

## VMC Example Usage

```hcl
data "nsxt_policy_ipsec_vpn_local_endpoint" "test" {
  id = "Public-IP1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Local Endpoint to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Local Endpoint to retrieve.
* `service_path` - (Optional) Service Path for this Local Endpoint.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
