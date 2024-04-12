---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: policy_site"
description: Policy Site data source.
---

# nsxt_policy_site

This data source provides information about Site (or Location) configured on NSX Global Manager.

This data source is applicable to NSX Global Manager only.

## Example Usage

```hcl
data "nsxt_policy_site" "paris" {
  display_name = "Paris"
}
```

## Argument Reference

* `id` - (Optional) The ID of Site to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Site to retrieve.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource. This attribute can serve as `site_path` field of `nsxt_policy_transport_zone` data source.
