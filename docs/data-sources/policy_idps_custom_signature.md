---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_idps_custom_signature"
description: A data source to read a single IDPS custom signature by ID.
---

# nsxt_policy_idps_custom_signature

This data source reads one IDPS custom signature by ID. You can specify the ID as a composite value `signature_version_id/signature_id` or as a standalone signature ID together with `signature_version_id`.

Applicable to NSX Policy Manager (NSX 4.2.0 onwards).

## Example Usage

```hcl
# By composite ID (from resource)
data "nsxt_policy_idps_custom_signature" "sig" {
  id = nsxt_policy_idps_custom_signature.malware.id
}

# By signature ID and version
data "nsxt_policy_idps_custom_signature" "by_id" {
  id                   = "5000001"
  signature_version_id = "default"
}

output "validation_status" {
  value = data.nsxt_policy_idps_custom_signature.sig.validation_status
}
```

## Argument Reference

- `id` - (Required) ID of the custom signature. Use format `signature_version_id/signature_id` (e.g. `default/5000001`), or the signature ID only when `signature_version_id` is set.
- `signature_version_id` - (Optional) Custom signature version ID (e.g. `default`). Required when `id` is only the signature ID.

## Attributes Reference

In addition to arguments above, the following attributes are exported:

- `path` - Absolute path of the custom signature.
- `display_name` - Display name.
- `revision` - Revision of the custom signature.
- `signature_id` - System-assigned signature ID.
- `original_signature_id` - Original signature ID from the rule.
- `original_signature` - Original raw signature (Snort rule).
- `validation_status` - One of: `VALID`, `INVALID`, `PENDING`, `WARNING`.
- `validation_message` - Validation message or error if invalid.
- `severity` - Severity of the signature.
- `signature_severity` - Vendor-assigned signature severity.
- `name` - Signature name (from msg in rule).
- `signature_revision` - Signature revision from the rule.
- `enable` - Whether the signature is enabled.
