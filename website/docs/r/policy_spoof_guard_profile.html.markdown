---
subcategory: "Segments"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_spoof_guard_profile"
description: A resource to configure SpoofGuard Profile.
---

# nsxt_policy_spoof_guard_profile

This resource provides a method for the management of SpoofGuard Profile.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_spoof_guard_profile" "test" {
  display_name              = "test"
  description               = "Terraform provisioned SpoofGuardProfile"
  address_binding_allowlist = true
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_spoof_guard_profile" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name              = "test"
  description               = "Terraform provisioned SpoofGuardProfile"
  address_binding_allowlist = true
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to
* `address_binding_allowlist` - (Optional) If true, enable the SpoofGuard, which only allows IPs listed in address bindings.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_spoof_guard_profile.test UUID
```
The above command imports SpoofGuard Profile named `test` with the NSX ID `UUID`.

```
terraform import nsxt_policy_spoof_guard_profile.test POLICY_PATH
```
The above command imports SpoofGuard Profile named `test` with policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.
