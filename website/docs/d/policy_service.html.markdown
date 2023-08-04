---
subcategory: "Firewall"
layout: "nsxt"
page_title: "NSXT: policy_service"
description: A policy service data source.
---

# nsxt_policy_service

This data source provides information about policy services configured in NSX.

## Example Usage

```hcl
data "nsxt_policy_service" "dns_service" {
  display_name = "DNS"
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_service" "demodnssvc" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "demodnssvc"
}
```

## Argument Reference

* `id` - (Optional) The ID of service to retrieve.
* `display_name` - (Optional) The Display Name prefix of the service to retrieve.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
