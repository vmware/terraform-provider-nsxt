---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_idps_signature_version"
description: A resource to reference IDPS signature versions in Terraform state.
---

# nsxt_policy_idps_signature_version

This resource provides a method for referencing IDPS signature versions in Terraform state.

**Important**: IDPS signature versions are **system-managed resources** in NSX-T. They are automatically created by NSX when signature updates are downloaded. This resource provides a way to:

- Reference existing signature versions in Terraform
- Activate a specific signature version (change state to ACTIVE)
- Track signature version information in Terraform state

**Note**: Auto-update and custom signature settings are configured globally via the `nsxt_policy_idps_settings` resource, not per-version.

This resource is applicable to NSX Policy Manager (NSX version 4.2.0 onwards).

## Example Usage

```hcl
# Reference an existing signature version and make it active
resource "nsxt_policy_idps_signature_version" "latest" {
  nsx_id      = "version-2024.01.15" # Must reference existing version
  description = "Latest signature version"
  state       = "ACTIVE"

  tag {
    scope = "environment"
    tag   = "production"
  }
}

# Configure auto-update globally (separate resource)
resource "nsxt_policy_idps_settings" "global" {
  auto_update_signatures   = true
  enable_custom_signatures = false
  oversubscription         = "BYPASSED"
}
```

## Argument Reference

The following arguments are supported:

- `nsx_id` - (Required) NSX ID of an existing signature version. Signature versions are system-managed and cannot be created directly. Use this to reference an existing version.
- `description` - (Optional) Description for tracking purposes in Terraform.
- `state` - (Optional) Version state. Set to `ACTIVE` to activate this version. Possible values:
    - `ACTIVE` - Make this version active (signatures will be used in IDS profiles)
    - `NOTACTIVE` - Version is available but not active (read-only, cannot be set directly)
- `tag` - (Optional) List of scope + tag pairs for organizing resources in Terraform.

**Note**: The following settings are configured globally via `nsxt_policy_idps_settings`, not per-version:

- Auto-update enablement → Use `nsxt_policy_idps_settings.auto_update_signatures`
- Custom signature enablement → Use `nsxt_policy_idps_settings.enable_custom_signatures`
- Update frequency → Managed internally by NSX-T

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

- `id` - ID of the signature version.
- `display_name` - Display name of the signature version (assigned by NSX).
- `revision` - Current revision number as seen by NSX-T API server (useful for debugging).
- `path` - NSX policy path.
- `version_id` - Version identifier string.
- `version_name` - Human-readable version name.
- `change_log` - Version change log describing updates.
- `update_time` - Time when version was downloaded and saved (epoch milliseconds).
- `status` - Version status. Possible values:
    - `LATEST` - This is the most recent version available
    - `OUTDATED` - A newer version is available
- `user_uploaded` - Whether this version was uploaded by a user (vs auto-downloaded).
- `auto_update` - Whether this version came via the auto-update mechanism.
- `sites` - List of sites mapped with this signature version.

## Importing

An existing signature version can be [imported][docs-import] into this resource:

```shell
terraform import nsxt_policy_idps_signature_version.latest ID
```

The above imports the signature version with NSX ID `ID`.

[docs-import]: https://developer.hashicorp.com/terraform/cli/import
