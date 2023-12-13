---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: upgrade_prepare_ready"
description: A data source to check if NSX Manager is ready for upgrade.
---

# nsxt_upgrade_prepare_ready

This data source provides information about the upgrade readiness of NSX Manager.
The data source will be successfully created only if there are no errors in precheck and
all warnings have been acknowledged.
If there is any nsxt_upgrade_precheck_acknowledge resource, then nsxt_upgrade_prepare_ready resource
will only be successfully created if all nsxt_upgrade_precheck_acknowledge resources are
successfully created.
When created together with nsxt_upgrade_precheck_acknowledge resources, the dependency must be specified
via depends_on in nsxt_upgrade_prepare_ready data source's template.

## Example Usage

```hcl
data "nsxt_upgrade_prepare_ready" "test" {
  depends_on         = ["nsxt_upgrade_precheck_acknowledge.test"]
  upgrade_prepare_id = nsxt_upgrade_prepare.test.id
}
```

## Argument Reference

* `upgrade_prepare_id` - (Required) ID of corresponding `nsxt_upgrade_prepare` resource.
* `id` - (Optional) The ID the upgrade_prepare_status data source.
