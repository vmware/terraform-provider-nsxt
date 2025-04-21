---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_lb_generic_persistence_profile"
description: A resource to configure a Load Balancer Generic Persistence Profile.
---

# nsxt_policy_lb_generic_persistence_profile

This resource provides a method for the management of a LB Generic Persistence Profile.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_lb_generic_persistence_profile" "test" {
  display_name       = "test"
  description        = "Terraform provisioned profile"
  persistence_shared = false
  timeout            = 1800
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `ha_persistence_mirroring_enabled` - (Optional) If enabled, persistence entries are synchronized to HA peer.
* `timeout` - (Optional) Expiration time once all connections are complete. Default is 300.
* `persistence_shared` - (Optional) If enabled, all virtual servers with this profile will share the same persistence mechanism.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_lb_generic_persistence_profile.test PATH
```

The above command imports LB Generic Persistence Profile named `test` with the NSX path `PATH`.
