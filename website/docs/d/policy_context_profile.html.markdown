---
layout: "nsxt"
page_title: "NSXT: policy_context_profile"
sidebar_current: "docs-nsxt-datasource-policy-context-profile"
description: Policy Context Profile Profile data source.
---

# nsxt_policy_context_profile

This data source provides information about Policy Context Profile configured in NSX.

**NOTE:** This data source is experimental.

## Example Usage

```hcl
data "nsxt_policy_context_profile" "diameter" {
  display_name = "DIAMETER"
}
```

## Argument Reference

* `id` - (Optional) The ID of Profile to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Profile to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
