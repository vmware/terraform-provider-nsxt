---
subcategory: "DNS"
page_title: "NSXT: policy_dns_service"
description: A policy DNS Service data source.
---

# nsxt_policy_dns_service

This data source provides information about a Policy DNS Service (`DnsService`) configured within an NSX project.

This data source is applicable to NSX Policy Manager and requires NSX 9.2.0 or higher.

## Example Usage

```hcl
data "nsxt_policy_dns_service" "test" {
  display_name = "my-dns-service"

  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
}
```

## Argument Reference

* `id` - (Optional) The ID of the DNS Service to retrieve.
* `display_name` - (Optional) The Display Name of the DNS Service to retrieve.
* `context` - (Required) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
