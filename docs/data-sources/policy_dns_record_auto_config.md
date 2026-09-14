---
subcategory: "DNS"
page_title: "NSXT: policy_dns_record_auto_config"
description: A policy DNS Record Auto Config data source.
---

# nsxt_policy_dns_record_auto_config

This data source provides information about a policy DNS Record Auto Config (`DnsAutoRecordConfig`) configured within an NSX project.

This data source is applicable to NSX Policy Manager and requires NSX 9.2.0 or higher.

## Example Usage

```hcl
data "nsxt_policy_dns_record_auto_config" "test" {
  display_name = "auto-config"

  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
}
```

## Argument Reference

* `id` - (Optional) The ID of the DNS Record Auto Config to retrieve.
* `display_name` - (Optional) The Display Name of the DNS Record Auto Config to retrieve.
* `context` - (Required) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
