---
subcategory: "Realization"
layout: "nsxt"
page_title: "NSXT: policy_gateway_interface_realization"
description: A gateway interface realization information.
---

# nsxt_policy_gateway_interface_realization_info

This data source provides information about the realization of a Tier0/Tier1 gateway interface resource on NSX manager. This data source will wait until realization is determined as either success or error. It is recommended to use this data source if further configuration depends on resource realization.

This data source is applicable to NSX Policy Manager, NSX Global Manager and Multi-Tenancy.

## Example Usage

```hcl
data "nsxt_policy_tier1_gateway" "tier1_gw" {
  display_name = "tier1_gw"
}

data "nsxt_policy_gateway_interface_realization" "info" {
  gateway_path = data.nsxt_policy_tier1_gateway.tier1_gw.path
  display_name = "pepsi-it_t1-t1_lrp"
  timeout      = 60
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_tier1_gateway" "tier1_gw" {
  display_name = "tier1_gw"
}

data "nsxt_policy_gateway_interface_realization" "info" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  gateway_path = data.nsxt_policy_tier1_gateway.tier1_gw.path
  display_name = "pepsi-it_t1-t1_lrp"
  timeout      = 60
}
```

## Argument Reference

* `gateway_path` - (Required) The policy path of the resource.
* `id`  - (Optional) The ID of gateway interface.
* `display_name` - (Optional) The display name of the resource. If neither ID nor display name are set, the realized gateway interface with IP addresses will be retrieved.
* `delay` - (Optional) Delay (in seconds) before realization polling is started. Default is set to 1.
* `timeout` - (Optional) Timeout (in seconds) for realization polling. Default is set to 1200.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `state` - The realization state of the resource: "REALIZED", "UNKNOWN", "UNREALIZED" or "ERROR".
* `ip_address` - The IP addresses of the realized object.
* `mac_address` - The MAC addresses of the realized object.
