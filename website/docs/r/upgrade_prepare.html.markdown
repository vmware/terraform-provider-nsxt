---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_upgrade_prepare"
description: A resource to prepare for the upgrade of NSXT cluster.
---

# nsxt_upgrade_prepare

This resource provides a method for preparing the upgrade of NSXT cluster,
it will upload upgrade bundle to NSX and run upgrade prechecks.
If there are errors in precheck result, the creation of resource will fail.
Precheck warnings will be listed in the `failed_prechecks` field of terraform state,
user can use the nsxt_upgrade_precheck_acknowledge resource to acknowledge these
warnings.
When using this resource to prepare for NSXT upgrade, the username and password
for NSXT must be provided in the nsxt provider config.

## Example Usage

```hcl
resource "nsxt_upgrade_prepare" "test" {
  upgrade_bundle_url    = "http://url-to-upgrade-bundle.mub"
  precheck_bundle_url   = "http://url-to-precheck-bundle.pub"
  accept_user_agreement = true
  bundle_upload_timeout = 1600
  uc_upgrade_timeout    = 600
  precheck_timeout      = 1600
  version               = "4.1.2"
}
```

## Argument Reference

The following arguments are supported:

* `upgrade_bundle_url` - (Required) The url to download the Manager Upgrade bundle.
* `version` - (Optional) Target upgrade version for NSX, format is x.x.x..., should include at least 3 digits, example: 4.1.2
* `precheck_bundle_url` - (Optional) The url to download the Precheck bundle.
* `accept_user_agreement` - (Required) This field must be set to true otherwise upgrade will not proceed.
* `bundle_upload_timeout` - (Optional) Timeout for uploading upgrade bundle in seconds. Default is `3600`.
* `uc_upgrade_timeout` - (Optional) Timeout for upgrading upgrade coordinator in seconds. Default is `1800`.
* `precheck_timeout` - (Optional) Timeout for executing pre-upgrade checks in seconds. Default is `3600`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `failed_prechecks` - (Computed) Failed prechecks from running pre-upgrade check.
  * `id` - ID of the failed precheck.
  * `message` - Message of the failed precheck.
  * `type` - Type of the failed precheck, possible values are `WARNING` and `FAILURE`.
  * `needs_ack` - Boolean value which identifies if acknowledgement is required for the precheck.
  * `needs_resolve` - Boolean value identifies if resolution is required for the precheck.
  * `acked` - Boolean value which identifies if precheck has been acknowledged.
  * `resolution_status` - The resolution status of precheck failure.
* `target_version` - Target system version

## Importing

Importing is not supported for this resource.
