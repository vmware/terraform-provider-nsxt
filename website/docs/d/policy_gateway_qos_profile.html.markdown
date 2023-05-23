---
subcategory: "Gateways and Routing"
layout: "nsxt"
page_title: "NSXT: policy_gateway_qos_profile"
description: Policy Gateway QoS Profile data source.
---

# nsxt_policy_gateway_qos_profile

This data source provides information about policy Gateway Quality of Service Profile configured on NSX.

This data source is applicable to NSX Global Manager, and NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_gateway_qos_profile" "test" {
  display_name = "gateway-qos-profile1"
}
```

## Argument Reference

* `id` - (Optional) The ID of GatewayQosProfile to retrieve.

* `display_name` - (Optional) The Display Name prefix of the Gateway QoS Profile to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.

* `path` - The NSX path of the policy resource.
