---
subcategory: "EVPN"
layout: "nsxt"
page_title: "NSXT: policy_vni_pool"
description: Policy VNI Pool Config data source.
---

# nsxt_policy_vni_pool

This data source provides information about policy VNI Pools configured in NSX.

This data source is applicable to NSX Policy Manager.

This data source is supported with NSX 3.0.0 onwards.

## Example Usage

```hcl
data "nsxt_policy_vni_pool" "test" {
  display_name = "vnipool1"
}
```

## Argument Reference

* `id` - (Optional) The ID of VNI Pool Config to retrieve.
* `display_name` - (Optional) The Display Name prefix of the VNI Pool Config to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
* `start` - The start range of VNI Pool.
* `end` - The end range of VNI Pool.
