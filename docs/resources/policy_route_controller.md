---
subcategory: "Gateways and Routing"
page_title: "NSXT: nsxt_policy_route_controller"
description: A resource to configure a Route Controller.
---

# nsxt_policy_route_controller

This resource provides a method for the management of a Route Controller.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_route_controller" "test" {
  display_name = "test"
  description  = "Terraform provisioned Route Controller"
  ha_mode      = "ACTIVE_STANDBY"

  bgp_config {
    ecmp                   = true
    graceful_restart_mode  = "HELPER_ONLY"
    graceful_restart_timer = 180
    local_as_num           = "65000"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `ha_mode` - (Optional) High-availability mode for route controller. Currently only `ACTIVE_STANDBY` is supported. This value is computed if not specified.
* `virtual_network_appliance_cluster_path` - (Optional) Policy path for the virtual network appliance cluster.
* `bgp_config` - (Optional) BGP routing configuration block. The following arguments are supported:
    * `ecmp` - (Optional) Flag to enable ECMP. Defaults to `true`.
    * `local_as_num` - (Required) BGP AS number in ASPLAIN/ASDOT format.
    * `multipath_relax` - (Optional) Flag to enable BGP multipath relax option.
    * `graceful_restart_mode` - (Optional) BGP graceful restart configuration mode. One of `DISABLE`, `GR_AND_HELPER`, `HELPER_ONLY`. Defaults to `HELPER_ONLY`.
    * `graceful_restart_timer` - (Optional) BGP graceful restart timer in seconds. Defaults to `180`.
    * `peer_route_convergence_timer` - (Optional) Extra time in seconds the router waits before sending UP notification after the peer session is established.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.
* `bgp_config`:
    * `path` - The NSX path of the BGP configuration.
    * `revision` - Indicates current revision number of the BGP configuration object.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_route_controller.test PATH
```

The above command imports Route Controller named `test` with the policy path `PATH`.
