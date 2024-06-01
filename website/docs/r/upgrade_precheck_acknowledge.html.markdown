---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_upgrade_precheck_acknowledge"
description: A resource to acknowledge failed NSXT upgrade prechecks.
---

# nsxt_upgrade_precheck_acknowledge

This resource provides a method for acknowledging the failed prechecks
for NSXT upgrade.

## Example Usage

```hcl
resource "nsxt_upgrade_precheck_acknowledge" "test" {
  precheck_ids   = ["backupOperationCheck", "pUBCheck"]
  target_version = nsxt_upgrade_prepare.test.target_version
}
```

## Argument Reference

The following arguments are supported:

* `precheck_ids` - (Required) List of ids of failed prechecks user wants to acknowledge.
* `target_version` - (Required) Target system version

## Attributes Reference

* `precheck_warnings` - (Computed) State of all precheck warnings.
  * `id` - ID of the precheck warning.
  * `message` - Message of the precheck warning.
  * `is_acknowledged` - Boolean value which identifies if the precheck warning has been acknowledged.

## Importing

Importing is not supported for this resource.
