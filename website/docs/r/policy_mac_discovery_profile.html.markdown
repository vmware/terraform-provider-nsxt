---
subcategory: "Segments"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_mac_discovery_profile"
description: A resource to configure a MAC Discovery Profile.
---

# nsxt_policy_mac_discovery_profile

This resource provides a method for the management of a MAC Discovery Profile.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_mac_discovery_profile" "test" {
  display_name                     = "test"
  description                      = "Terraform provisioned"
  mac_change_enabled               = true
  mac_learning_enabled             = true
  mac_limit                        = 4096
  mac_limit_policy                 = "ALLOW"
  remote_overlay_mac_limit         = 2048
  unknown_unicast_flooding_enabled = true
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_mac_discovery_profile" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name                     = "test"
  description                      = "Terraform provisioned"
  mac_change_enabled               = true
  mac_learning_enabled             = true
  mac_limit                        = 4096
  mac_limit_policy                 = "ALLOW"
  remote_overlay_mac_limit         = 2048
  unknown_unicast_flooding_enabled = true
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
* `mac_change_enabled` - (Optional) MAC address change feature.
* `mac_learning_enabled` - (Optional) MAC learning feature.
* `mac_limit` - (Optional) The maximum number of MAC addresses that can be learned on this port.
* `mac_limit_policy` - (Optional) The policy after MAC Limit is exceeded. Possible values are `ALLOW` and `DROP`, with default being `ALLOW`.
* `remote_overlay_mac_limit` - (Optional) The maximum number of MAC addresses learned on an overlay Logical Switch.
* `unknown_unicast_flooding_enabled` - (Optional) Allowing flooding for unlearned MAC for ingress traffic.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_mac_discovery_profile.test UUID
```
The above command imports MAC Discovery Profile named `test` with ID `UUID`.

```
terraform import nsxt_policy_mac_discovery_profile.test POLICY_PATH
```
The above command imports MAC Discovery Profile named `test` with policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.
