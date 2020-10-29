---
subcategory: "Policy - DHCP"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_dhcp_v4_static_binding_config"
description: A resource to configure IPv4 DHCP Static Binding.
---

# nsxt_policy_dhcp_v4_static_binding_config

This resource provides a method for the management of IPv4 DHCP Static Binding.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_dhcp_v4_static_binding_config" "test" {
  display_name    = "test"
  description     = "Terraform provisioned static binding"
  gateway_address = "10.0.2.1"
  hostname        = "host1"
  ip_address      = "10.0.2.167"
  lease_time      = 6400
  mac_address     = "10:ff:22:11:cc:02"

  dhcp_generic_option {
    code   = 43
    values = ["cat"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `segment_path` - (Required) Policy path for segment to configure this binding on.
* `ip_address` - (Required) The IPv4 address must belong to the subnet, if any, configured on Segment.
* `mac_address` - (Required) MAC address of the host.
* `gateway_address` - (Optional) Gateway IPv4 Address. When not specified, gateway address is auto-assigned from segment configuration.
* `hostname` - (Optional) Hostname to assign to the host.
* `lease_time` - (Optional) Lease time, in seconds. Defaults to 86400.
* `dhcp_option_121` - (Optional) DHCP classless static routes.
  * `network` - (Required) Destination in cidr format.
  * `next_hop` - (Required) IP address of next hop.
* `dhcp_generic_option` - (Optional) Generic DHCP options.
  * `dhcp_generic_option` - (Optional) Generic DHCP options.
  * `values` - (Required) List of DHCP option values.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_policy_dhcp_v4_static_binding_config.test UUID
```

The above command imports DhcpV4StaticBindingConfig named `test` with the NSX DhcpV4StaticBindingConfig ID `UUID`.
