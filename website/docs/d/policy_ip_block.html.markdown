---
subcategory: "IPAM"
layout: "nsxt"
page_title: "NSXT: policy_ip_block"
description: Policy IP Block data source.
---

# nsxt_policy_ip_block

This data source provides information about Policy IP Blocks configured on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_ip_block" "test" {
  display_name = "ipblock1"
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_ip_block" "demoblock" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "demoblock"
}
```

## Argument Reference

* `id` - (Optional) The ID of IP Block to retrieve.
* `display_name` - (Optional) The Display Name prefix of the IP Block to retrieve.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
