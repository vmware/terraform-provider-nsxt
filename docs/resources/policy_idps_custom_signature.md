---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_idps_custom_signature"
description: A resource to create, read, update, and delete IDPS custom signatures (vDefend Signature Management).
---

# nsxt_policy_idps_custom_signature

This resource manages IDPS custom signatures in NSX-T. Custom signatures allow you to define detection rules using Snort syntax. They are created under a custom signature version (e.g. `default`) and can be validated and optionally published so they become active in IDS/IPS policies.

Applicable to NSX Policy Manager (NSX 4.2.0 onwards).

## Example Usage

```hcl
resource "nsxt_policy_idps_custom_signature" "malware" {
  signature_version_id = "default"
  signature            = "alert tcp any any -> any any (msg:\"Custom Malware\"; flow:to_server,established; content:\"malware_string\"; nocase; metadata:signature_severity Critical; classtype:trojan-activity; sid:5000001; rev:1;)"
  description          = "Custom malware detection rule"
  enable               = true
  publish              = false
}

# Publish so the signature becomes active in IDS/IPS
resource "nsxt_policy_idps_custom_signature" "published" {
  signature_version_id = "default"
  signature            = "alert tcp any any -> any 80 (msg:\"Web Attack\"; content:\"attack\"; nocase; sid:5000002; rev:1;)"
  publish              = true
}
```

## Argument Reference

The following arguments are supported:

- `signature_version_id` - (Required) ID of the custom signature version (e.g. `default`) to which this signature belongs. Force new.
- `signature` - (Required) Raw custom signature rule in Snort syntax.
- `description` - (Optional) Description of the custom signature.
- `enable` - (Optional) Whether the signature is enabled. Default is `true`.
- `publish` - (Optional) If `true`, validate and publish the signature so it becomes active in IDS/IPS. Default is `false` (signature remains in preview).

## Attributes Reference

In addition to arguments above, the following attributes are exported:

- `id` - Composite ID in the form `signature_version_id/signature_id`.
- `path` - Absolute path of the custom signature.
- `display_name` - Display name (derived from signature).
- `revision` - Revision of the custom signature.
- `signature_id` - System-assigned signature ID (e.g. SID in Snort).
- `original_signature_id` - Original signature ID from the rule (sid in Snort).
- `original_signature` - Original raw signature as stored.
- `validation_status` - One of: `VALID`, `INVALID`, `PENDING`, `WARNING`.
- `validation_message` - Validation message or error if invalid.
- `severity` - Severity of the signature.
- `signature_severity` - Vendor-assigned signature severity.
- `name` - Signature name (from msg in rule).
- `signature_revision` - Signature revision from the rule.

## Behavior and lifecycle

The following sections describe how the resource behaves for create, read, update, and delete. This ensures predictable results and explains fallbacks and edge cases.

### Create

1. The raw signature is sent to the Custom Signature Versions API with action `ADD_CUSTOM_SIGNATURES`. The signature is added to the version in **preview** (unpublished) state.
2. The provider lists custom signatures (including preview) and finds the newly added one by matching the exact `signature` (rule text). If multiple entries match (e.g. same rule from a previous run), the **last** match in the list is used so the newly created signature is stored.
3. The resource ID is set to `signature_version_id/signature_id`, where `signature_id` is the last path segment of the signature (e.g. from the path returned by the API).
4. If `publish = true`, the provider runs VALIDATE (no modifications) then PUBLISH so the signature becomes active in IDS/IPS.

### Read

1. The provider calls Get by `signature_version_id` and `signature_id`. Published signatures are returned by Get.
2. If Get returns 404 (e.g. signature is preview-only and the API does not support Get for preview), the provider falls back to listing custom signatures (preview and published) and finds the one whose ID or path matches the stored `signature_id`.
3. Preview signatures may report `enable = false` until published; the provider preserves the configured `enable` value in state when the API returns false to avoid unnecessary plan drift.

### Update

1. If only `publish` or other non-signature attributes change, the provider may only run VALIDATE and PUBLISH when `publish = true`.
2. If the raw `signature` (rule text) changes, the provider sends a VALIDATE request with `ModifiedSignatures` containing the new rule and the signature ID, then optionally PUBLISH when `publish = true`.

### Delete

Deletion uses the NSX API pattern: stage the removal with VALIDATE, then apply it with PUBLISH. There is no single-signature DELETE endpoint. This flow works for both **preview** (unpublished) and **published** custom signatures.

**Primary path (per-signature delete):**

1. Resolve the signature to its full path (from Get or from list by ID).
2. Call VALIDATE with `DeletedSignatures` set to that path (normalized with a leading `/`).
3. If VALIDATE succeeds, call PUBLISH with an empty payload (revision only) so the staged deletion is applied.
4. If the signature was only in preview, it may still appear in the list after PUBLISH. In that case the provider may call CANCEL to remove all unpublished signatures in that version so the resource is fully removed.
5. If the signature still exists after the first PUBLISH (e.g. published signature or eventual consistency), the provider retries VALIDATE+PUBLISH with a fresh revision and tries alternative path forms (preview-style and signatures-style) so the deletion is applied.

**Fallback when path is rejected:**

- Some NSX versions or configurations reject the path format in `DeletedSignatures` (e.g. error "path is invalid" or code 500012). The provider first retries VALIDATE with a relative path form `signatures-preview/<signature_id>`.
- If VALIDATE still fails, the provider falls back to **CANCEL**. CANCEL reverts the entire custom signature version to the last published state and removes **all** unpublished custom signatures in that version (per vDefend: "A signature that is validated but not yet published can be cancelled"). Use this only when there are no other unpublished signatures in the same version that you need to keep.

**After delete:**

- The resource is removed from Terraform state. The backend may take a few seconds to reflect the removal (eventual consistency). Acceptance tests retry the destroy check for a short period before asserting the signature is gone.

### Summary of guarantees

|Operation|What is guaranteed|
|---------|------------------|
|Create|Signature exists in the version (preview or published if `publish = true`); state has composite ID and matches the created signature when multiple matches exist (last match).|
|Read|State is populated from Get or, on 404, from list by ID; drift on `enable` for preview is minimized.|
|Update|Rule changes are sent as ModifiedSignatures + VALIDATE; PUBLISH is used when `publish = true`.|
|Delete|Terraform state is cleared. Backend removal is via VALIDATE + PUBLISH (for both preview and published signatures; retries with alternative path forms if needed), or CANCEL when path is invalid (removes all unpublished in the version).|

## Import

Import is supported using the composite ID:

```shell
terraform import nsxt_policy_idps_custom_signature.malware default/5000001
```

Format: `signature_version_id/signature_id`.
