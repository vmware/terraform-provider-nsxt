---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_vtep_ha_host_switch_profile"
description: A resource to configure a VTEP HA host switch profile in NSX Policy.
---

# nsxt_policy_vtep_ha_host_switch_profile

This resource provides a method for the management of a VTEP HA host switch profile which can be used within NSX Policy.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_vtep_ha_host_switch_profile" "vtep_ha_host_switch_profile" {
  description  = "VTEP host switch profile provisioned by Terraform"
  display_name = "vtep_ha_host_switch_profile"

  auto_recovery              = true
  auto_recovery_initial_wait = 3003
  auto_recovery_max_backoff  = 80000
  enabled                    = true
  failover_timeout           = 40

  tag {
    scope = "color"
    tag   = "blue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `auto_recovery` - (Optional) Enabled status of autonomous recovery option. Default is true.
* `auto_recovery_initial_wait` - (Optional) Start time of autonomous recovery (in seconds). Default is 300.
* `auto_recovery_max_backoff` - (Optional) Maximum backoff time for autonomous recovery (in seconds). Default is 86400.
* `enabled` - (Optional) Enabled status of VTEP High Availability feature. Default is false.
* `failover_timeout` - (Optional) VTEP High Availability failover timeout (in seconds). Default is 5.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful
  for debugging.
* `path` - The NSX path of the policy resource.
* `realized_id` - Realized ID for the profile. For reference in fabric resources (such as `transport_node`), `realized_id` should be used rather than `id`.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_vtep_ha_host_switch_profile.vtep_ha_host_switch_profile UUID
```

The above command imports VTEP HA host switch profile named `vtep_ha_host_switch_profile` with the NSX ID `UUID`.

```
terraform import nsxt_policy_vtep_ha_host_switch_profile.vtep_ha_host_switch_profile POLICY_PATH
```

The above command imports the VTEP HA host switch profile named `vtep_ha_host_switch_profile` with policy
path `POLICY_PATH`.
