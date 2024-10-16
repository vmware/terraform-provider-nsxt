---
subcategory: "VPC"
layout: "nsxt"
page_title: "NSXT: nsxt_vpc_connectivity_profile"
description: A resource to configure Connectivity Profile for VPC.
---

# nsxt_vpc_connectivity_profile

This resource provides a method for the management of a Connectivity Profile for VPC.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_vpc_connectivity_profile" "test" {
  context {
    project_id = "dev"
  }

  display_name = "test"
  description  = "Terraform provisioned profile"

  transit_gateway_path = nsxt_policy_transit_gateway.gw1.path
  service_gateway {
    nat_config {
      enable_default_snat = true
    }
    enable = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `transit_gateway_path` - (Required) Transit Gateway path.
* `private_tgw_ip_blocks` - (Optional) Policy path of Private IP block
* `external_ip_blocks` - (Optional) Policy path of External IP block
* `service_gateway` - (Optional) Service Gateway configuration
  * `nat_config` - (Optional) NAT configuration
    * `enable_default_snat` - (Optional) Auto configured SNAT for private subnet.
  * `qos_config` - (Optional) None
    * `ingress_qos_profile_path` - (Optional) Policy path to gateway QoS profile in ingress direction.
    * `egress_qos_profile_path` - (Optional) Policy path to gateway QoS profile in egress direction.
  * `enable` - (Optional) Status of the VPC attachment SR.
  * `edge_cluster_paths` - (Optional) List of edge cluster paths for VPC attachment SR realization. If edge cluster is not specified transit gateway's edge cluster will be used.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_vpc_connectivity_profile.test PATH
```

The above command imports VpcConnectivityProfile named `test` with the NSX policy path `PATH`.
