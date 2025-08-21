---
subcategory: "Gateways and Routing"
page_title: "NSXT: policy_tier1_gateway"
description: A policy Tier-1 gateway data source.
---

# nsxt_policy_tier1_gateways

This data source provides list of Tier-1s configured on NSX.

This data source is applicable to NSX Policy Manager, NSX Global Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_tier1_gateways" "all" {}

output "nsxt_policy_tier1_gateways_result" {
  value = data.nsxt_policy_tier1_gateways.all.items
}
```

## Example Usage - Using Regex

```hcl
data "nsxt_policy_tier1_gateways" "all" {
  display_name = ".*"
}

output "nsxt_policy_tier1_gateways_result" {
  value = data.nsxt_policy_tier1_gateways.all.items
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "all" {}

data "nsxt_policy_tier1_gateways" "all" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
}
```

## Example Usage - Global infra

```hcl
data "nsxt_policy_tier1_gateways" "all" {
  context {
    from_global = true
  }
}
```

## Argument Reference

* `display_name` - (Optional) Display name for the Tier1. Supports regular expressions.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Optional) The ID of the project which the object belongs to
    * `from_global` - (Optional) Set to True if the data source will need to search Tier-1 gateway created in a global manager instance (/global-infra)

## Attributes Reference

* `items` - Map of IDs by Display Name.
