---
subcategory: "EVPN"
page_title: "NSXT: nsxt_policy_distributed_vxlan_connection"
description: A resource to configure a Distributed VXLAN Connection.
---

# nsxt_policy_distributed_vxlan_connection

This resource provides a method for the management of a Distributed VXLAN Connection.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_distributed_vxlan_connection" "test" {
  display_name          = "test"
  description           = "Terraform provisioned Distributed VXLAN Connection"
  l3_vni                = 5000
  route_controller_path = nsxt_policy_route_controller.rc.path
  route_distinguisher   = "65001:100"

  route_target {
    address_family = "L2VPN_EVPN"
    import_targets = ["65001:100"]
    export_targets = ["65001:200"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `l3_vni` - (Required) L3 VNI for overlay traffic.
* `connectivity_type` - (Optional) Connectivity type for the distributed VXLAN connection. Currently only `L3_EVPN` is supported. Default is `L3_EVPN`.
* `route_controller_path` - (Required) Policy path of the route controller.
* `route_distinguisher` - (Required) Route distinguisher for the VXLAN connection. Must be in ASN:nn or IP:nn format.
* `route_target` - (Required) Route target configuration. Only one block is supported.
    * `address_family` - (Optional) Address family for route targets. Currently only `L2VPN_EVPN` is supported. Default is `L2VPN_EVPN`.
    * `import_targets` - (Optional) List of import route targets in ASN:nn format.
    * `export_targets` - (Optional) List of export route targets in ASN:nn format.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_distributed_vxlan_connection.test PATH
```

The above command imports Distributed VXLAN Connection named `test` with the policy path `PATH`.
