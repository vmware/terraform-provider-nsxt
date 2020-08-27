---
subcategory: "Policy - Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: policy_gateway_qos_profile"
description: Policy GatewayQosProfile data source.
---

# nsxt_policy_gateway_qos_profile

This data source provides information about policy GatewayQosProfile configured on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_gateway_qos_profile" "test" {
  display_name = "gateway-qos-profile1"
}
```

## Argument Reference

* `id` - (Optional) The ID of GatewayQosProfile to retrieve.

* `display_name` - (Optional) The Display Name prefix of the GatewayQosProfile to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.

* `path` - The NSX path of the policy resource.
