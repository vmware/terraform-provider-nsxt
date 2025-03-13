---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: policy_tags"
description: Tags data source.
---

# nsxt_policy_tags

This data source provides the list of tags with a particular scope.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_tags" "tags" {
  scope = "dev"
}
```

## Argument Reference

* `scope` - (Required) The scope of the tags to retrieve. Supports starts with, ends with, equals, and contains filters.
Use * as a suffix for "starts with" and a prefix for "ends with". For "contains," wrap the value with *. Use * alone to fetch all tags, irrespective scope.

## Attributes Reference

The following attributes are exported:

* `items` - List of tags with the scope given in the argument.
