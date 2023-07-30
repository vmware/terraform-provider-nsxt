---
subcategory: "Realization"
layout: "nsxt"
page_title: "NSXT: policy_realization_info"
description: A policy resource realization information.
---

# nsxt_policy_realization_info

This data source provides information about the realization of a policy resource on NSX manager. This data source will wait until realization is determined as either success or error. It is recommended to use this data source if further configuration depends on resource realization.

This data source is applicable to NSX Policy Manager and NSX Global Manager.

## Example Usage

```hcl
data "nsxt_policy_tier1_gateway" "tier1_gw" {
  display_name = "tier1_gw"
}

data "nsxt_policy_realization_info" "info" {
  path        = data.nsxt_policy_tier1_gateway.tier1_gw.path
  entity_type = "RealizedLogicalRouter"
  timeout     = 60
}
```

## Global Manager Example

```hcl
data "nsxt_policy_tier1_gateway" "tier1_gw" {
  display_name = "tier1_gw"
}

data "nsxt_policy_site" "site" {
  display_name = "Paris"
}

data "nsxt_policy_realization_info" "info" {
  path        = data.nsxt_policy_tier1_gateway.tier1_gw.path
  entity_type = "RealizedLogicalRouter"
  site_path   = data.nsxt_policy_site.site.path
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

data "nsxt_policy_realization_info" "info" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  path        = data.nsxt_policy_tier1_gateway.tier1_gw.path
  entity_type = "RealizedLogicalRouter"
  timeout     = 60
}
```

## Argument Reference

* `path` - (Required) The policy path of the resource.
* `entity_type` - (Optional) The entity type of realized resource. If not set, on of the realized resources of the policy resource will be retrieved.
* `site_path` - (Optional) The path of the site which the resource belongs to, this configuration is required for global manager only. `path` field of the existing `nsxt_policy_site` can be used here.
* `delay` - (Optional) Delay (in seconds) before realization polling is started. Default is set to 1.
* `timeout` - (Optional) Timeout (in seconds) for realization polling. Default is set to 1200.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `state` - The realization state of the resource: "REALIZED", "UNKNOWN", "UNREALIZED" or "ERROR".
* `realized_id` - The id of the realized object.
