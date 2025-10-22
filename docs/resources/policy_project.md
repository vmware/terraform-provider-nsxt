---
subcategory: "Multitenancy"
page_title: "NSXT: nsxt_policy_project"
description: A resource to configure a Project.
---

# nsxt_policy_project

This resource provides a method for the management of a Project.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_project" "test" {
  display_name        = "test"
  description         = "Terraform provisioned Project"
  short_id            = "test"
  tier0_gateway_paths = ["/infra/tier-0s/test"]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `short_id` - (Optional) Defaults to id if id is less than equal to 8 characters or defaults to random generated id if not set.
* `activate_default_dfw_rules` - (Optional) By default, Project is created with default distributed firewall rules, this boolean flag allows to deactivate those default rules. If not set, the default rules are enabled. Available since NSX 4.2.0.
* `site_info` - (Optional) Information related to sites applicable for given Project. For on-prem deployment, only 1 is allowed.
    * `edge_cluster_paths` - (Optional) The edge cluster on which the networking elements for the Org will be created.
    * `site_path` - (Optional) This represents the path of the site which is managed by Global Manager. For the local manager, if set, this needs to point to 'default'.
* `tier0_gateway_paths` - (Optional) The tier 0 has to be pre-created before Project is created. The tier 0 typically provides connectivity to external world. List of sites for Project has to be subset of sites where the tier 0 spans.
* `external_ipv4_blocks` - (Optional) IP blocks used for allocating CIDR blocks for public subnets. These can be consumed by all the VPCs under this project. Available since NSX 4.1.1.
* `tgw_external_connections` - (Optional) Transit gateway connection objects available to the project. Gateway connection and distributed VLAN connection object path will be allowed. Available since NSX 9.0.0.
* `default_security_profile`- (Optional) Default security profile properties for project.
    * `north_south_firewall` - (Required) North South firewall configuration.
          * `enabled` - (Required) This flag indicates whether north-south firewall (Gateway Firewall) is enabled. If set to false, then gateway firewall policies will not be enforced on the VPCs associated with this configuration.
* `vc_folder` - (Optional) Flag to specify whether the DVPGs created for project segments are grouped under a folder on the VC. Defaults to `true`. Available since NSX 9.0.0.
* `quotas` - (Optional) List of policy paths for quota resources that are applicable to this project. Available since NSX 9.0.0.
* `default_span_path` - (Optional) Policy path of the Cluster based default Span object of type NetworkSpan. Available since NSX 9.1.0.
* `non_default_span_paths` - (Optional) List of non default policy paths of the Span objects of type NetworkSpan.
* `zone_external_ids` - (Optional) An array of Zone object's external IDs.
* `id_suffix` - (Optional) A suffix appended to the IDs of the project's default objects. This attribute is supported only in NSX 9.1.0 and later, and cannot be modified after creation.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_project.test UUID
```

The above command imports Project named `test` with the NSX ID `UUID`.
