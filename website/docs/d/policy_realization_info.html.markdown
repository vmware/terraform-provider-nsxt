---
layout: "nsxt"
page_title: "NSXT: policy_realization_info"
sidebar_current: "docs-nsxt-datasource-policy-realization-info"
description: A policy resource realization information.
---

# nsxt_policy_realization_info

This data source provides information about the realization of a policy resource on NSX manager. This data source will wait until realization is determined as either success or error. It is recommended to use this data source if further configuration depends on resource realization.

## Example Usage

```hcl
data "nsxt_policy_tier1_gateway" "tier1_gw" {
  display_name = "tier1_gw"
}

data "nsxt_policy_realization_info" "info" {
  path = data.nsxt_policy_tier1_gateway.tier1_gw.path
  entity_type = "RealizedLogicalRouter"
}
```

## Argument Reference

* `path` - (Required) The policy path of the resource.

* `entity_type` - (Optional) The entity type of realized resource. If not set, on of the realized resources of the policy resource will be retrieved.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `state` - The realization state of the resource: "REALIZED", "UNKNOWN", "UNREALIZED" or "ERROR".

* `realized_id` - The id of the realized object.
