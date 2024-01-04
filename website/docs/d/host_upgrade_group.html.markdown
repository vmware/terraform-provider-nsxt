---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: host_upgrade_group"
description: A Host Upgrade Group data source.
---

# nsxt_host_upgrade_group

This data source provides information about a Host Upgrade Group configured on NSX.

## Example Usage

```hcl
data "nsxt_host_upgrade_group" "test" {
  upgrade_prepare_id = nsxt_upgrade_prepare.test.id
  display_name       = "hostupgradegroup1"
}
```

## Argument Reference

* `upgrade_prepare_id` - (Required) ID of corresponding `nsxt_upgrade_prepare` resource.
* `id` - (Optional) The ID of Host Upgrade Group to retrieve
* `display_name` - (Optional) The Display Name of the Host Upgrade Group to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
