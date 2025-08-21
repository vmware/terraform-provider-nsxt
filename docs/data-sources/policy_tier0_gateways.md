---
subcategory: "Gateways and Routing"
page_title: "NSXT: policy_tier0_gateways"
description: A policy Tier-0 gateways data source.
---

# nsxt_policy_tier0_gateways

This data source provides list of Tier-0s configured on NSX.

This data source is applicable to NSX Policy Manager, NSX Global Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_tier0_gateways" "all" {}

output "nsxt_policy_tier0_gateways_result" {
  value = data.nsxt_policy_tier0_gateways.all.items
}
```

## Example Usage - Using Regex

```hcl
data "nsxt_policy_tier0_gateways" "all" {
  display_name = ".*"
}

output "nsxt_policy_tier0_gateways_result" {
  value = data.nsxt_policy_tier0_gateways.all.items
}
```

## Example Usage - Global infra

```hcl
data "nsxt_policy_tier0_gateways" "all" {
  context {
    from_global = true
  }
}
```

## Argument Reference

* `display_name` - (Optional) Display name for the Tier0. Supports regular expressions.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Optional) The ID of the project which the object belongs to
    * `from_global` - (Optional) Set to True if the data source will need to search Tier-0 gateway created in a global manager instance (/global-infra)

## Attributes Reference

* `items` - Map of IDs by Display Name.
