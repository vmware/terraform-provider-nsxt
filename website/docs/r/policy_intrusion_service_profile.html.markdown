---
subcategory: "Firewall"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_intrusion_service_profile"
description: A resource to configure Intrusion Service Profile.
---

# nsxt_policy_intrusion_service_profile

This resource provides a method for the management of Intrusion Service (IDS) Profile.

This resource is applicable to NSX Policy Manager and VMC (NSX version 3.1.0 and up).

## Example Usage

```hcl
resource "nsxt_policy_intrusion_service_profile" "profile1" {
  display_name = "test"
  description  = "Terraform provisioned Profile"
  severities   = ["HIGH", "CRITICAL"]

  criteria {
    attack_types      = ["trojan-activity", "successful-admin", "attempted-admin", "attempted-user"]
    attack_targets    = ["Web_Server", "SQL_Server", "Email_Server", "Client_and_Server", "Server"]
    products_affected = ["Windows_XP_Vista_7_8_10_Server_32_64_Bit", "Web_Server_Applications", "Linux"]
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

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_intrusion_service_profile" "profile1" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "test"
  description  = "Terraform provisioned Profile"
  severities   = ["HIGH", "CRITICAL"]

  criteria {
    attack_types      = ["trojan-activity", "successful-admin", "attempted-admin", "attempted-user"]
    attack_targets    = ["Web_Server", "SQL_Server", "Email_Server", "Client_and_Server", "Server"]
    products_affected = ["Windows_XP_Vista_7_8_10_Server_32_64_Bit", "Web_Server_Applications", "Linux"]
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
* `tag` - (Optional) A list of scope + tag pairs to associate with this profile.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `context` - (Optional) The context which the object belongs to
  * `project_id` - (Required) The ID of the project which the object belongs to
* `severities` - (Required) List of profile severities, supported values are `LOW`, `MEDIUM`, `HIGH`, 'CRITICAL`.
* `criteria` - (Required) Filtering criteria for the IDS Profile.
  * `attack_types` - (Optional) List of supported attack types.
  * `attack_targets` - (Optional) List of supported attack targets. Please refer to example above to ensure correct formatting - in some versions, UI shows a different format than NSX expects.
  * `cvss` - (Optional) List of CVSS (Common Vulnerability Scoring System) ranges. Supported values are `NONE`, `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`.
  * `products_affected` - (Optional) List of supported products that are affected. Please refer to example above to ensure correct formatting - in some versions, UI shows a different format than NSX expects.
* `overridden_signature` - (Optional) List of signatures that has been overridden this profile.
  * `signature_id` - (Required) Id for the existing signature that profile wishes to override.
  * `action` - (Optional) Overridden action, one of `ALERT`, `DROP`, `REJECT`. Default is `ALERT`.
  * `enabled` - (Optional) Flag to enable/disable this signature.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the Profile.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX policy path of the resource.

## Importing

An existing profile can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_intrusion_service_profile.profile1 ID
```
The above command imports the profile named `profile1` with the NSX Policy ID `ID`.

```
terraform import nsxt_policy_intrusion_service_profile.profile1 POLICY_PATH
```
The above command imports the profile named `profile1` with the NSX policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.
