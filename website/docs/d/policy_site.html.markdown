---
layout: "nsxt"
page_title: "NSXT: policy_site"
sidebar_current: "docs-nsxt-datasource-policy-site"
description: Policy Site data source.
---

# nsxt_policy_site

This data source provides information about Site (or Region) configured in NSX Global Manager. This data source is applicable for Global Manager provider only.

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

* `path` - The NSX path of the policy resource.
