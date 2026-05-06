---
subcategory: "Upgrade"
page_title: "NSXT: upgrade_postcheck"
description: A data source to retrieve upgrade post-check results.
---

# nsxt_upgrade_postcheck

This data source waits for upgrade post-checks to complete and reports any failures.

The data source polls post-check execution status until all groups have finished running
checks (i.e. none remain in `IN_PROGRESS` state). On success, no output is expected:
`failed_group` will be empty when all post-checks passed. If any group fails or has
warnings, it will be listed in `failed_group`.

Post-checks are triggered by the `nsxt_upgrade_run` resource when `post_upgrade_check`
is enabled in `edge_upgrade_setting` or `host_upgrade_setting`.

## Example Usage

```hcl
data "nsxt_upgrade_postcheck" "pc" {
  upgrade_run_id = nsxt_upgrade_run.test.id
  type           = "EDGE"
}
```

## Argument Reference

* `upgrade_run_id` - (Required) The ID of corresponding `nsxt_upgrade_run` resource.
* `type` - (Required) Component type of upgrade post-checks. Supported values: `EDGE`, `HOST`.
* `timeout` - (Optional) Post-check status checks timeout in seconds. Default: 1800 seconds.
* `interval` - (Optional) Interval to check post-check status in seconds. Default: 30 seconds.
* `delay` - (Optional) Initial delay before starting post-check status checks in seconds. Default: 30 seconds.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `failed_group` - Upgrade unit groups that failed post-checks. This list will be empty when all post-checks passed successfully.
    * `id` - ID of upgrade unit group.
    * `status` - Status of execution of upgrade post-checks.
    * `details` - Details about current execution of upgrade post-checks.
    * `failure_count` - Total count of generated failures or warnings in last execution of upgrade post-checks.
