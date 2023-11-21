---
subcategory: "Firewall"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_firewall_exclude_list_member"
description: A resource to configure a member within the firewall exclude list.
---

# nsxt_policy_firewall_exclude_list_member

This resource provides a method for the management of a firewall exclude list members.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_group" "testgroup" {
  display_name = "Terraform created test group"
  description  = "Group for testing purposes"
}

resource "nsxt_policy_firewall_exclude_list_member" "excludetestgroup" {
  member = nsxt_policy_group.testgroup.path
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_group" "testgroup" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "Terraform created test group"
  description  = "Group for testing purposes"
}

resource "nsxt_policy_firewall_exclude_list_member" "excludetestgroup" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  member = nsxt_policy_group.testgroup.path
}
```

## Argument Reference

The following arguments are supported:

* `member` - Exclusion list member policy path

## Importing

An existing exclude list member can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_firewall_exclude_list_member.test POLICY_PATH
```
The above command imports the exclude list member named `test` with policy path `POLICY_PATH`.
