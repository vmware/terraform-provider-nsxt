---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_idps_signature_diff"
description: A data source to read the difference between published and unpublished custom signatures.
---

# nsxt_policy_idps_signature_diff

This data source returns the difference between the current (unpublished) custom signatures and the last published state for a given custom signature version. Use it to see newly added, deleted, and existing signatures before or after publishing.

Applicable to NSX Policy Manager (NSX 4.2.0 onwards).

## Example Usage

```hcl
data "nsxt_policy_idps_signature_diff" "diff" {
  signature_version_id = "default"
}

output "newly_added" {
  value = data.nsxt_policy_idps_signature_diff.diff.newly_added_signatures
}

output "deleted" {
  value = data.nsxt_policy_idps_signature_diff.diff.deleted_signatures
}
```

## Argument Reference

- `signature_version_id` - (Required) ID of the custom signature version (e.g. `default`).

## Attributes Reference

- `id` - Same as `signature_version_id`.
- `newly_added_signatures` - List of custom signature IDs that are newly added (unpublished).
- `deleted_signatures` - List of custom signature IDs that are deleted (removed from published).
- `existing_signatures` - List of custom signature IDs that exist in both published and current state.
