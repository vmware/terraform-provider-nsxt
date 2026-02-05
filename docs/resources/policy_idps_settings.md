---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_idps_settings"
description: A resource to configure global IDPS settings including syslog, oversubscription, signature auto-update, and custom signature enablement.
---

# nsxt_policy_idps_settings

This resource provides a method for the management of global IDPS (Intrusion Detection and Prevention System) settings. It allows configuration of:

- **Signature Auto Update** - Automatically download and install latest IDPS signatures
- **Syslog Settings** - Forward IDPS events to syslog servers
- **Oversubscription Settings** - Control behavior when IDPS engine capacity is exceeded
- **Custom Signature Enablement** - Enable/disable custom user-defined signatures globally

**⚠️ Singleton Resource**: This resource represents global IDPS settings and should only be declared **once per NSX-T environment**. Multiple instances in your Terraform configuration will be treated as the same resource, potentially causing configuration conflicts.

This resource is applicable to NSX Policy Manager and VMC (NSX version 4.2.0 onwards).

## Example Usage

```hcl
resource "nsxt_policy_idps_settings" "global" {
  description = "Global IDPS configuration for production environment"

  # Enable automatic signature updates
  auto_update_signatures = true

  # Enable syslog for IDPS events
  enable_syslog = true

  # Set oversubscription behavior
  oversubscription = "BYPASSED"

  # Enable custom signatures (requires custom signature version)
  enable_custom_signatures    = true
  custom_signature_version_id = "version-1"

  tag {
    scope = "environment"
    tag   = "production"
  }
}
```

## Example Usage - With Custom Signature Version

```hcl
# Create a custom signature version first
resource "nsxt_policy_idps_signature_version" "custom" {
  display_name = "Custom Signatures v1"
  description  = "Custom IDPS signatures for organization-specific threats"
}

# Configure global IDPS settings to use custom signatures
resource "nsxt_policy_idps_settings" "global" {
  description = "IDPS settings with custom signatures enabled"

  auto_update_signatures      = true
  enable_syslog               = true
  oversubscription            = "BYPASSED"
  enable_custom_signatures    = true
  custom_signature_version_id = nsxt_policy_idps_signature_version.custom.nsx_id
}
```

## Argument Reference

The following arguments are supported:

- `display_name` - (Computed) Display name of the resource. This is read-only for this singleton resource.
- `description` - (Optional) Description of the resource.
- `tag` - (Optional) A list of scope + tag pairs to associate with this resource.

### IDPS Settings

- `auto_update_signatures` - (Optional) Enable automatic update of IDS/IPS signatures. When enabled, NSX-T will automatically download and install the latest signature updates from VMware. Default is `false`.

- `enable_syslog` - (Optional) Enable sending IDS/IPS events to syslog server. When enabled, IDPS events will be forwarded to configured syslog servers (configured separately in NSX-T). Default is `false`.

- `oversubscription` - (Optional) Action to take when IDPS engine is oversubscribed (processing capacity exceeded). Supported values:
    - `BYPASSED` - Traffic bypasses IDPS inspection when oversubscribed (maintains connectivity)
    - `DROPPED` - Traffic is dropped when oversubscribed (maintains security)
  
  Default is `BYPASSED`.

- `enable_custom_signatures` - (Optional) Enable custom signatures globally. When enabled, user-defined custom signatures in the specified version will be active in IDPS policies. Requires `custom_signature_version_id` to be set. Default is `false`.

- `custom_signature_version_id` - (Optional) The custom signature version ID to use when enabling custom signatures. This should reference an existing custom signature version resource (e.g., `nsxt_policy_idps_signature_version.custom.nsx_id`). Required when `enable_custom_signatures` is `true`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

- `id` - ID of the IDPS settings (always `intrusion-services` for this singleton resource).
- `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
- `path` - The NSX policy path of the resource.

## Importing

An existing IDPS settings configuration can be [imported][docs-import] into this resource, via the following command:

```shell
terraform import nsxt_policy_idps_settings.global intrusion-services
```

The above command imports the IDPS settings named `global` with the NSX ID `intrusion-services`.

[docs-import]: https://www.terraform.io/docs/import/

## Notes

- The `DELETE` operation resets settings to defaults rather than deleting the resource (as it's a singleton resource)
