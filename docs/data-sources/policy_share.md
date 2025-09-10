---
subcategory: "Multitenancy"
page_title: "NSXT: nsxt_policy_share"
description: A Share data source.
---

# nsxt_policy_share

This data source provides information about a Share configured on NSX.
This data source is applicable to NSX Policy Manager.
Note: The env variable NSXT_READ_RETRY_TIMEOUT_SECONDS can be set to add a retry while reading the resources.

## Example Usage

```hcl
data "nsxt_policy_share" "share" {
  display_name = "share1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Share to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Share to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
