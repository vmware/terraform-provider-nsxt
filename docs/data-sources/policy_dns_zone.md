---
subcategory: "DNS"
page_title: "NSXT: policy_dns_zone"
description: A policy DNS Zone data source.
---

# nsxt_policy_dns_zone

This data source provides information about a policy DNS Zone attached to a Policy DNS Service.

This data source is applicable to NSX Policy Manager and requires NSX 9.2.0 or higher.

## Example Usage

```hcl
data "nsxt_policy_dns_zone" "test" {
  display_name = "my-dns-zone"
  parent_path  = nsxt_policy_dns_service.demo.path
}
```

## Argument Reference

* `id` - (Optional) The ID of the DNS Zone to retrieve.
* `display_name` - (Optional) The Display Name of the DNS Zone to retrieve.
* `parent_path` - (Required) Policy path of the parent DnsService.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
