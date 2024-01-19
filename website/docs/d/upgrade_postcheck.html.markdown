---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: edge_upgrade_group"
description: An Edge Upgrade Group data source.
---

# nsxt_upgrade_postcheck

This data source provides information about Upgrade post-check results.

## Example Usage

```hcl
data "nsxt_upgrade_postcheck" "pc" {
  upgrade_run_id = nsxt_upgrade_run.test.id
  type           = "EDGE"
}
```

## Argument Reference

* `upgrade_run_id` - (Required) The ID of corresponding `nsxt_upgrade_run` resource.
* `type` - (Required) Component type of upgrade post-checks. Supported type: EDGE, HOST.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `failed_group` - (Computed) The description of the resource.
    * `id` - ID of upgrade unit group.
    * `status` - Status of execution of upgrade post-checks.
    * `details` - Details about current execution of upgrade post-checks.
    * `failure_count` - Total count of generated failures or warnings in last execution of upgrade post-checks      },
