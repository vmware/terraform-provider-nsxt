---
layout: "nsxt"
page_title: "NSXT: mac_pool"
sidebar_current: "docs-nsxt-datasource-mac-pool"
description: A MAC pool data source.
---

# nsxt_mac_pool

This data source provides information about a MAC pool configured in NSX.

## Example Usage

```hcl
data "nsxt_mac_pool" "mac_pool" {
  display_name = "DefaultMacPool"
}
```

## Argument Reference

* `id` - (Optional) The ID of MAC pool to retrieve

* `display_name` - (Optional) The Display Name of the MAC pool to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the MAC pool.
