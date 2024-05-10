---
subcategory: "Manager"
layout: "nsxt"
page_title: "nsxt_manager_info"
description: A NSX-T manager information data source.
---

# nsxt_manager_info

This data source provides various information about the NSX-T manager.

## Example Usage

```hcl
data "nsxt_manager_info" "cluster" {}
```

## Attributes Reference

* `version` - The software version of the NSX-T manager node.
