---
subcategory: "DNS"
page_title: "NSXT: policy_dns_rule"
description: A policy DNS Rule data source.
---

# nsxt_policy_dns_rule

This data source provides information about a policy DNS Rule attached to a Policy DNS Service.

This data source is applicable to NSX Policy Manager and requires NSX 9.2.0 or higher.

## Example Usage

```hcl
data "nsxt_policy_dns_rule" "test" {
  display_name = "my-dns-rule"
  parent_path  = nsxt_policy_dns_service.demo.path
}
```

## Argument Reference

* `id` - (Optional) The ID of the DNS Rule to retrieve.
* `display_name` - (Optional) The Display Name of the DNS Rule to retrieve.
* `parent_path` - (Required) Policy path of the parent DnsService.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
