---
subcategory: "Gateways and Routing"
page_title: "NSXT: nsxt_policy_route_controller_interface"
description: A resource to configure a Route Controller Interface.
---

# nsxt_policy_route_controller_interface

This resource provides a method for the management of a Route Controller Interface.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_route_controller_interface" "test" {
  display_name        = "test"
  description         = "Terraform provisioned Route Controller Interface"
  parent_path         = nsxt_policy_route_controller.rc.path
  urpf_mode           = "STRICT"
  floating_ip_subnets = ["192.168.200.1/24"]

  interface_address {
    subnets                        = ["192.168.100.1/24"]
    portgroup_id                   = "6d7b6a6d-8c5b-4d6c-9f9c-0e5b6a6d8c5b:dvportgroup-11"
    virtual_network_appliance_path = nsxt_policy_virtual_network_appliance_cluster.vna.path
  }

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `parent_path` - (Required) Policy path of the parent Route Controller.
* `urpf_mode` - (Optional) Unicast Reverse Path Forwarding mode. One of `NONE`, `STRICT`. Defaults to `STRICT`.
* `mtu` - (Optional) Maximum transmission unit specifies the size of the largest packet that a network protocol can transmit. Minimum value is `64`.
* `floating_ip_subnets` - (Required) List of floating IP subnets in CIDR notation (IP/prefix-length) for this interface.
* `interface_address` - (Optional) List of interface address configurations. The following arguments are supported:
    * `subnets` - (Required) List of IP addresses and network prefixes in CIDR notation for this interface address.
    * `portgroup_id` - (Required) DV port group identifier discovered from vCenter.
    * `virtual_network_appliance_path` - (Optional) Policy path for the virtual network appliance.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing Route Controller Interface can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```shell
terraform import nsxt_policy_route_controller_interface.test <path>
```

The above command imports the Route Controller Interface named `test` using the policy path as import ID.
The import ID should be the full policy path of the interface, for example:
`/infra/route-controllers/<controller-id>/interfaces/<interface-id>`
