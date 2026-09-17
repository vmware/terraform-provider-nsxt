---
subcategory: "VPC"
page_title: "NSXT: policy_project_centralized_network_attachment"
description: A policy Project Centralized Network Attachment data source.
---

# nsxt_policy_project_centralized_network_attachment

This data source provides information about a Centralized Network Attachment (CNA) configured within a project.

This data source is applicable to NSX Policy Manager and is supported with NSX 9.0.0 onwards.

## Example Usage

```hcl
data "nsxt_policy_project_centralized_network_attachment" "test" {
  display_name = "my-cna"

  context {
    project_id = nsxt_policy_project.demoproj.id
  }
}
```

## Argument Reference

* `id` - (Optional) The ID of the Centralized Network Attachment to retrieve.
* `display_name` - (Optional) The Display Name of the Centralized Network Attachment to retrieve.
* `context` - (Required) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
