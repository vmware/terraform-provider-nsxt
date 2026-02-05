---
subcategory: "IDPS"
page_title: "NSXT: nsxt_policy_idps_signature_version"
description: A resource to manage IDPS signature versions and auto-update settings.
---

# nsxt_policy_idps_signature_version

This resource provides a method for the management of IDPS signature versions. It allows configuration of signature auto-update settings and custom signature enablement.

This resource is applicable to NSX Policy Manager and VMC (NSX version 3.1.0 onwards).

## Example Usage

```hcl
# Basic configuration
resource "nsxt_policy_idps_signature_version" "latest" {
  display_name             = "latest-signatures"
  description              = "Latest signature version with auto-update"
  version                  = "2024.01.15"
  auto_update_enabled      = true
  custom_signature_enabled = true
  update_frequency         = "DAILY"

  tag {
    scope = "environment"
    tag   = "production"
  }
}

# Multi-tenancy
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_idps_signature_version" "latest_mt" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name        = "latest-signatures"
  version             = "2024.01.15"
  auto_update_enabled = true
  update_frequency    = "WEEKLY"
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name.
* `description` - (Optional) Description.
* `version` - (Required) Signature version identifier. Cannot be changed after creation.
* `auto_update_enabled` - (Optional) Enable automatic updates. Default: `true`.
* `custom_signature_enabled` - (Optional) Enable custom signature support. Default: `false`.
* `update_frequency` - (Optional) Update frequency: `HOURLY`, `DAILY`, `WEEKLY`. Default: `DAILY`.
* `tag` - (Optional) List of scope + tag pairs.
* `nsx_id` - (Optional) NSX ID. If set, this ID will be used to create the resource.
* `context` - (Optional) Context for multi-tenancy.
    * `project_id` - (Required) Project ID.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the signature version.
* `revision` - Current revision number as seen by NSX-T API server (useful for debugging).
* `path` - NSX policy path.
* `signature_count` - Number of signatures in this version.
* `last_update_time` - Last update timestamp.
* `status` - Version status.

## Importing

An existing signature version can be [imported][docs-import] into this resource:

```shell
terraform import nsxt_policy_idps_signature_version.latest ID
```

The above imports the signature version with NSX ID `ID`.

[docs-import]: https://www.terraform.io/docs/import/

## Notes

* Manages signature version containers for IDPS.
* The `version` field is immutable after creation.
* Auto-update settings control signature refresh frequency.
* Custom signature enablement allows user-defined detection rules.
* `signature_count` and `status` are read-only, system-updated attributes.
