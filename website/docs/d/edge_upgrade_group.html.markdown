---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: edge_upgrade_group"
description: An Edge Upgrade Group data source.
---

# nsxt_edge_upgrade_group

This data source provides information about an Edge Upgrade Group configured on NSX.

## Example Usage

```hcl
data "nsxt_edge_upgrade_group" "test" {
  upgrade_prepare_id = nsxt_upgrade_prepare.test.id
  display_name       = "edgeupgradegroup1"
}
```

## Argument Reference

* `upgrade_prepare_id` - (Required) ID of corresponding `nsxt_upgrade_prepare` resource.
* `id` - (Optional) The ID of Edge Upgrade Group to retrieve
* `display_name` - (Optional) The Display Name of the Edge Upgrade Group to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
