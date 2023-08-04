---
subcategory: "DHCP"
layout: "nsxt"
page_title: "NSXT: policy_dhcp_server"
description: A policy DHCP server data source.
---

# nsxt_policy_dhcp_server

This data source provides information about policy DHCP servers configured on NSX.

This data source is applicable to NSX Policy Manager, NSX Global Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_dhcp_server" "test" {
  display_name = "dhcp-2"
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_dhcp_server" "demodhcp" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "demodhcp"
}
```

## Argument Reference

* `id` - (Optional) The ID of DHCP Server to retrieve. If ID is specified, no additional argument should be configured.
* `display_name` - (Optional) The Display Name prefix of DHCP server to retrieve.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
