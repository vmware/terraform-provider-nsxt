---
subcategory: "EVPN"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_evpn_tenant"
description: A resource to configure EVPN Tenant on NSX Policy manager.
---

# nsxt_policy_evpn_tenant

This resource provides a method for the management of EVPN Tenant.

This resource is applicable to NSX Policy Manager.
This resource is supported with NSX 3.1.0 onwards.

## Example Usage

```hcl
resource "nsxt_policy_evpn_tenant" "tenant1" {
  display_name = "tenant1"
  description  = "terraform provisioned tenant"

  transport_zone_path = data.nsxt_policy_transport_zone.overlay1.path
  vni_pool_path       = data.nsxt_policy_vni_pool.pool1.path

  mapping {
    vnis  = "75012-75015"
    vlans = "112-115"
  }

  mapping {
    vnis  = "75003"
    vlans = "113"
  }

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
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the policy resource.
* `transport_zone_path` - (Required) Policy path for Overlay Transport Zone.
* `vni_pool_path` - (Required) Policy path for existing VNI pool.
* `mapping` - (Required) List of VLAN - VNI mappings for this tenant.
  * `vlans` - (Required) Single VLAN Id or range of VLAN Ids.
  * `vnis` - (Required) Single VNI or range of VNIs. Please note that the range should match the range of vlans exactly.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing EVPN Tenant can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_evpn_tenant.tenant1 ID
```

The above command imports EVPN Tenant named `tenant1` with the NSX Policy ID `ID`.
