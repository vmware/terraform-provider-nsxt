---
subcategory: "Policy - Firewall"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_intrusion_service_profile"
description: A resource to configure Intrusion Service Profile.
---

# nsxt_policy_intrusion_service_profile

This resource provides a method for the management of Intrusion Service (IDS) Profile.

This resource is applicable to NSX Policy Manager and VMC (NSX version 3.0.0 and up).

## Example Usage

```hcl
resource "nsxt_policy_intrusion_service_profile" "profile1" {
  display_name = "test"
  description  = "Terraform provisioned Profile"
  severities   = ["HIGH", "CRITICAL"]

  criteria {
    attack_types      = ["trojan-activity", "successful-admin"]
    products_affected = ["Linux"]
  }

  overridden_signature {
    signature_id = "2026323"
    action       = "REJECT"
  }

  overridden_signature {
    signature_id = "2026324"
    action       = "REJECT"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this policy.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `severities` - (Required) List of profile severities, supported values are `LOW`, `MEDIUM`, `HIGH`, 'CRITICAL`.
* `criteria` - (Required) Filtering criteria for the IDS Profile.
  * `attack_types` - (Optional) List of supported attack types.
  * `attack_targets` - (Optional) List of supported attack targets.
  * `cvss` - (Optional) List of CVSS (Common Vulnerability Scoring System) ranges. Supported values are `NONE`, `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`.
  * `products_affected` - (Optional) List of supported products that are affected.
* `overridden_signature` - (Optional) List of signatures that has been overridden this profile.
  * `signature_id` - (Required) Id for the existing signature that profile wishes to override.
  * `action` - (Optional) Overridden action, one of `ALERT`, `DROP`, `REJECT`. Default is `ALERT`.
  * `enabled` - (Optional) Flag to enable/disable this signature.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the Security Policy.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing security policy can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_policy_intrusion_service_profile.profile1 ID
```

The above command imports the policy named `profile1` with the NSX Policy ID `ID`.
