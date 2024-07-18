---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: compute_manager"
description: A Compute Manager data source.
---

# nsxt_compute_manager

This data source provides information about a Compute Manager configured on NSX.
When only single Compute Manager is present on NSX, no need to specify any filter in the data source.

## Example Usage

```hcl
data "nsxt_compute_manager" "test_vcenter" {
  display_name = "some-compute-mgr"
}
```

## Argument Reference

* `id` - (Optional) The ID of Compute Manager to retrieve
* `display_name` - (Optional) The Display Name of the Compute Manager to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `server` - IP address or hostname of the resource.
