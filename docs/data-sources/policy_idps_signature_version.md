---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_idps_signature_version"
description: A data source to read IDPS signature version information.
---

# nsxt_policy_idps_signature_version

This data source provides information about IDPS signature versions in NSX Policy manager. Signature versions are system-managed resources that are automatically created when NSX downloads signature updates.

This data source is applicable to NSX Policy Manager (NSX version 4.2.0 onwards).

## Example Usage

```hcl
# Look up signature version by display name
data "nsxt_policy_idps_signature_version" "latest" {
  display_name = "2024.01.15"
}

# Look up signature version by ID
data "nsxt_policy_idps_signature_version" "specific" {
  id = "version-2024.01.15"
}

# Use in other resources
output "version_status" {
  value = data.nsxt_policy_idps_signature_version.latest.status
}
```

## Argument Reference

* `id` - (Optional) The ID of the signature version to retrieve.
* `display_name` - (Optional) The Display Name prefix of the signature version to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - Description of the resource.
* `path` - NSX policy path.
* `version_id` - Signature version identifier string.
* `version_name` - Human-readable version name.
* `change_log` - Version change log describing updates.
* `update_time` - Time when version was downloaded and saved (epoch milliseconds).
* `status` - Version status (`LATEST` or `OUTDATED`).
* `state` - Current state (`ACTIVE` or `NOTACTIVE`).
* `user_uploaded` - Whether this version was uploaded by a user.
* `auto_update` - Whether this version came via the auto-update mechanism.
* `sites` - List of sites mapped with this signature version.

**Note**: Auto-update and custom signature settings are global settings configured via the `nsxt_policy_idps_settings` data source, not per-version attributes.
