---
subcategory: "DHCP"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_dhcp_v4_static_binding"
description: A resource to configure IPv4 DHCP Static Binding.
---

# nsxt_policy_dhcp_v4_static_binding

This resource provides a method for the management of IPv4 DHCP Static Binding.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC (NSX version 3.0.0 and up)

## Example Usage

```hcl
resource "nsxt_policy_dhcp_v4_static_binding" "test" {
  segment_path    = nsxt_policy_segment.test.path
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

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_dhcp_v4_static_binding" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  segment_path    = nsxt_policy_segment.test.path
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

* `segment_path` - (Required) Policy path for segment to configure this binding on.
* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `context` - (Optional) The context which the object belongs to
  * `project_id` - (Required) The ID of the project which the object belongs to
* `ip_address` - (Required) The IPv4 address must belong to the subnet, if any, configured on Segment.
* `mac_address` - (Required) MAC address of the host.
* `gateway_address` - (Optional) Gateway IPv4 Address. When not specified, gateway address is auto-assigned from segment configuration.
* `hostname` - (Optional) Hostname to assign to the host.
* `lease_time` - (Optional) Lease time, in seconds. Defaults to 86400.
* `dhcp_option_121` - (Optional) DHCP classless static routes.
  * `network` - (Required) Destination in cidr format.
  * `next_hop` - (Required) IP address of next hop.
* `dhcp_generic_option` - (Optional) Generic DHCP options.
  * `dhcp_generic_option` - (Optional) Generic DHCP option number. Please note not all options are supported by the platform.
  * `values` - (Required) List of DHCP option values.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_dhcp_v4_static_binding.test [GW-ID]/SEG-ID/ID
```
The above command imports DHCP V4 static binding named `test` with the NSX ID `ID` on segment `SEG-ID`.
For fixed segments (VMC), `GW-ID` needs to be specified. Otherwise, `GW-ID` should be omitted.

```
terraform import nsxt_policy_dhcp_v4_static_binding.test POLICY_PATH
```
The above command imports DHCP V4 static binding named `test` with the NSX policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.
