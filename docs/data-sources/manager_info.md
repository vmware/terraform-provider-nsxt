---
subcategory: "Manager"
page_title: "nsxt_manager_info"
description: A NSX manager information data source.
---

# nsxt_manager_info

This data source provides various information about the NSX manager.

## Example Usage

```hcl
data "nsxt_manager_info" "cluster" {}
```

## Attributes Reference

* `version` - The software version of the NSX manager node.
