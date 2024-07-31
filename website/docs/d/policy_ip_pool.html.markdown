---
subcategory: "IPAM"
layout: "nsxt"
page_title: "NSXT: policy_ip_pool"
description: Policy IP Pool Config data source.
---

# nsxt_policy_ip_pool

This data source provides information about policy IP Pools configured on NSX.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_ip_pool" "test" {
  display_name = "ippool1"
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_ip_pool" "demopool" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "demopool"
}
```

## Argument Reference

* `id` - (Optional) The ID of IP Pool Config to retrieve.
* `display_name` - (Optional) The Display Name prefix of the IP Pool Config to retrieve.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
* `realized_id` - The id of realized pool object. This id should be used in `nsxt_edge_transport_node` resource.
