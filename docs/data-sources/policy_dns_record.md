---
subcategory: "DNS"
page_title: "NSXT: policy_dns_record"
description: A policy DNS Record data source.
---

# nsxt_policy_dns_record

This data source provides information about a policy DNS Record configured within an NSX project.

This data source is applicable to NSX Policy Manager and requires NSX 9.2.0 or higher.

## Example Usage

```hcl
data "nsxt_policy_dns_record" "test" {
  display_name = "www-record"
  zone_path    = nsxt_policy_dns_zone.example.path

  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
}
```

## Argument Reference

* `id` - (Optional) The ID of the DNS Record to retrieve.
* `display_name` - (Optional) The Display Name of the DNS Record to retrieve.
* `zone_path` - (Optional) Policy path of the parent DNS Zone, used to narrow the search.
* `context` - (Required) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
