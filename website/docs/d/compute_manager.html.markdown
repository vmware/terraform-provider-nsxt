---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: compute_manager"
description: A Compute Manager data source.
---

# nsxt_compute_manager

This data source provides information about a Compute Manager configured on NSX.

## Example Usage

```hcl
data "nsxt_compute_manager" "test_vcenter" {
  display_name = "test-vcenter"
}
```

## Argument Reference

* `id` - (Optional) The ID of Compute Manager to retrieve
* `display_name` - (Optional) The Display Name of the Compute Manager to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the Compute Manager.
* `server` - IP address or hostname of compute manager.
