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
* `centralized_config` - (Optional) Singleton block for high-availability and edge cluster for centralized connectivity (gateway connections, VPN). Sent in the same H-API transaction as the transit gateway (as a child object, like security policy rules). Available since NSX 9.1.0.
    * `ha_mode` - (Optional) High-availability mode. Values: `ACTIVE_ACTIVE`, `ACTIVE_STANDBY`. Default is `ACTIVE_ACTIVE`.
    * `edge_cluster_paths` - (Optional) Policy paths of edge clusters. At most one item. Must be authorized for the project.
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
