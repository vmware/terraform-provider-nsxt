---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_idps_settings"
description: A data source to read global IDPS settings information.
---

# nsxt_policy_idps_settings

This data source provides information about global IDPS settings configured in NSX Policy manager, including signature auto-update, syslog, oversubscription, and custom signature enablement settings.

**ℹ️ Singleton Resource**: This data source reads global IDPS settings which exist as a single instance per NSX-T environment.

This data source is applicable to NSX Policy Manager 4.2.0 onwards.

## Example Usage

```hcl
data "nsxt_policy_idps_settings" "global" {
}

output "auto_update_enabled" {
  value = data.nsxt_policy_idps_settings.global.auto_update_signatures
}

output "syslog_enabled" {
  value = data.nsxt_policy_idps_settings.global.enable_syslog
}

output "oversubscription_action" {
  value = data.nsxt_policy_idps_settings.global.oversubscription
}
```

## Example Usage - With Custom Signature Version

```hcl
# Read IDPS settings with custom signature configuration
data "nsxt_policy_idps_settings" "global" {
  custom_signature_version_id = "version-1"
}

output "custom_signatures_enabled" {
  value = data.nsxt_policy_idps_settings.global.enable_custom_signatures
}
```

## Argument Reference

* `custom_signature_version_id` - (Optional) The custom signature version ID to query for custom signature enablement status. If provided, the data source will also retrieve the custom signature settings for this version.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - The ID of the IDPS settings (always `intrusion-services` for this singleton resource).
* `display_name` - The display name of the resource.
* `description` - The description of the resource.
* `path` - The NSX policy path of the resource.

### IDPS Settings Attributes

* `auto_update_signatures` - Whether automatic update of IDS/IPS signatures is enabled. When `true`, NSX-T automatically downloads and installs the latest signature updates.

* `enable_syslog` - Whether sending IDS/IPS events to syslog server is enabled. When `true`, IDPS events are forwarded to configured syslog servers (configured separately in NSX-T).

* `oversubscription` - Action to take when IDPS engine is oversubscribed (processing capacity exceeded). Possible values:
    * `BYPASSED` - Traffic bypasses IDPS inspection when oversubscribed (maintains connectivity)
    * `DROPPED` - Traffic is dropped when oversubscribed (maintains security)

* `enable_custom_signatures` - Whether custom signatures are enabled globally. When `true`, user-defined custom signatures in the configured version are active. This attribute is only populated when `custom_signature_version_id` is provided.
