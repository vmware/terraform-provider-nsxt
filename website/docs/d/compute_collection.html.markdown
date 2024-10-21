---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: compute_collection"
description: A Compute Collection data source.
---

# nsxt_compute_collection

This data source provides information about a Compute Collection configured on NSX.

## Example Usage

```hcl
data "nsxt_compute_collection" "test_cluster" {
  display_name = "Compute_Cluster"
}
```

## Argument Reference

* `id` - (Optional) The ID of Compute Collection to retrieve.
* `display_name` - (Optional) The Display Name of the Compute Collection to retrieve.
* `origin_type` - (Optional) ComputeCollection type, f.e. `VC_Cluster`. Here the Compute Manager type prefix would help in differentiating similar named Compute Collection types from different Compute Managers.
* `origin_id` - (Optional) Id of the compute manager from where this Compute Collection was discovered.
* `cm_local_id` - (Optional) Local Id of the compute collection in the Compute Manager.

## Attributes Reference

All of the arguments above are computed unless specified.
