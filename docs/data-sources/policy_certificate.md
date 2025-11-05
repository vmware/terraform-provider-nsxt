---
subcategory: "Certificates"
page_title: "NSXT: policy_certificate"
description: Policy Certificate data source.
---

# nsxt_policy_certificate

This data source provides information about Service Certificate configured on NSX Policy.

This data source is applicable to NSX Global Manager, and NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_certificate" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "certificate1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Certificate to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Certificate to retrieve.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
