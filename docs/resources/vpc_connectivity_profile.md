---
subcategory: "VPC"
page_title: "NSXT: nsxt_vpc_connectivity_profile"
description: A resource to configure Connectivity Profile for VPC.
---

# nsxt_vpc_connectivity_profile

This resource provides a method for the management of a Connectivity Profile for VPC.

This resource is applicable to NSX Policy Manager and is supported with NSX 9.0.0 onwards.

## Example Usage

```hcl
resource "nsxt_vpc_connectivity_profile" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }

  display_name         = "dev-connectivity-profile"
  description          = "Connectivity profile for development VPCs"
  transit_gateway_path = nsxt_policy_transit_gateway.test.path
  external_ip_blocks = [nsxt_policy_ip_block.test.path]

  service_gateway {
    enable             = true
    edge_cluster_paths = [data.nsxt_policy_edge_cluster.edge_cluster.path]
    nat_config {
      enable_default_snat = true
    }
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
    * `edge_cluster_paths` - (Optional) Array of edge cluster or service cluster path for VPC service instance realization. Edge / Service cluster must be associated with the Project. Edge cluster path shall be specified when VPC is connected to centralised transit gateway and service cluster path shall be specified when VPC is connected to distributed transit gateway Changing edge or service cluster is allowed but relocation of vpc services to new cluster is a disruptive operation and affects traffic.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_vpc_connectivity_profile.test PATH
```

The above command imports VpcConnectivityProfile named `test` with the NSX policy path `PATH`.
