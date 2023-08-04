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

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_gateway_qos_profile" "qosprof" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "qosprof"
}
```

## Argument Reference

* `id` - (Optional) The ID of GatewayQosProfile to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Gateway QoS Profile to retrieve.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
