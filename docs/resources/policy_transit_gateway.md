---
subcategory: "VPC"
page_title: "NSXT: nsxt_policy_transit_gateway"
description: A resource to configure a Transit Gateway.
---

# nsxt_policy_transit_gateway

This resource provides a method for the management of a Transit Gateway.

This resource is applicable to NSX Policy Manager.

In NSX 9.0, users cannot create TGWs. Terraform users can import and update the default TGW.
In NSX 9.1 onwards, users can create and delete TGWs

## Example Usage

```hcl

resource "nsxt_policy_transit_gateway" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }

  display_name = "test-transit-gateway"
  description  = "Development transit gateway for VPC connectivity"

  high_availability_config {
    ha_mode            = "ACTIVE_STANDBY"
    edge_cluster_paths = [data.nsxt_policy_edge_cluster.edge_cluster.path]
  }
}

resource "nsxt_policy_transit_gateway" "testsample" {
  context {
    project_id = "dev"
  }

  display_name    = "test"
  description     = "Terraform provisioned TransitGateway"
  transit_subnets = ["10.203.4.0/24"]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `transit_subnets` - (Optional) Array of IPV4 CIDRs for internal VPC attachment networks.
* `high_availability_config` - (Optional) Transit Gateway high availability config centralized transit gateway. Available since NSX 9.1.0.
    * `ha_mode` - (Optional) High-availability Mode for Transit Gateway. Accepted values are: "ACTIVE_ACTIVE", "ACTIVE_STANDBY". Default is "ACTIVE_ACTIVE".
    * `edge_cluster_paths` - (Required) The Edge cluster should be authorized to be used in the transit gateway. A single edge cluster will be supported when the transit gateway is created from the local NSX manager.
* `span` - (Optional) Span configuration. Note that one of `cluster_based_span` and `zone_based_span` is required. Available since NSX 9.1.0.
    * `cluster_based_span` - (Optional) Span based on vSphere Clusters.
        * `span_path` - (Required) Policy path of the network span object.
    * `zone_based_span` - (Optional) Span based on zones.
        * `zone_external_ids` - (Optional) An array of Zone object's external IDs.
        * `use_all_zones` - (Optional) Flag to indicate that TransitGateway is associated with all project zones.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_transit_gateway.test PATH
```

The above command imports Transit Gateway named `test` with the NSX path `PATH`.
