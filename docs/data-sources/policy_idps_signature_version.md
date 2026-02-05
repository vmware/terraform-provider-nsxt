---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_idps_signature_version"
description: A data source to read IDPS signature version information.
---

# nsxt_policy_idps_signature_version

This data source provides information about an IDPS signature version configured in NSX Policy manager.

This data source is applicable to NSX Policy Manager and VMC (NSX version 4.2.0 onwards).

## Example Usage

```hcl
# Basic usage
data "nsxt_policy_idps_signature_version" "latest" {
  display_name = "latest-signatures"
}
```

## Argument Reference

* `id` - (Optional) The ID of the signature version to retrieve.
* `display_name` - (Optional) The Display Name prefix of the signature version to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - Description of the resource.
* `path` - NSX policy path.
* `version` - Signature version identifier.
* `auto_update_enabled` - Whether automatic updates are enabled.
* `custom_signature_enabled` - Whether custom signatures are enabled.
* `update_frequency` - Automatic update frequency.
* `signature_count` - Number of signatures in this version.
* `last_update_time` - Last update timestamp.
* `status` - Version status.
