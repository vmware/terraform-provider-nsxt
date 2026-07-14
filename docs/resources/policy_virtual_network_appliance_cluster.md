---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_virtual_network_appliance_cluster"
description: A resource to configure a Virtual Network Appliance Cluster.
---

# nsxt_policy_virtual_network_appliance_cluster

This resource provides a method for the management of a Virtual Network Appliance (VNA) Cluster under an NSX enforcement point.

A VNA Cluster is a logical grouping of Virtual Network Appliances that share common configuration such as form factor, service type, and advanced settings. Cluster members (VNA appliances and their associated edge transport nodes) are managed by NSX and are exposed as read-only attributes on the `nsxt_policy_virtual_network_appliance_cluster` data source.

This resource is supported with NSX 9.2.0 onwards.

## Example Usage

```hcl
data "nsxt_policy_transport_zone" "overlay_tz" {
  display_name = "overlay-tz"
}

resource "nsxt_policy_virtual_network_appliance_cluster" "example" {
  display_name          = "vna-cluster-1"
  description           = "VNA cluster for VPC services"
  appliance_form_factor = "MEDIUM"
  service_type          = "VPC_SERVICES"

  advanced_configuration {
    core_allocation_profile     = "L4LBSERVICE"
    overlay_transport_zone_path = data.nsxt_policy_transport_zone.overlay_tz.path
  }

  tag {
    scope = "environment"
    tag   = "production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `site_path` - (Optional) The path of the site this cluster belongs to. `path` field of the existing `nsxt_policy_site` can be used here. Defaults to default site path.
* `enforcement_point` - (Optional) The ID of enforcement point under given `site_path` to manage the cluster. Defaults to default enforcement point.
* `appliance_form_factor` - (Optional, Computed) The form factor of the virtual network appliances in the cluster. Supported values: `SMALL`, `MEDIUM`, `LARGE`, `XLARGE`. When the form factor is updated, new VNAs deployed at the cluster level will use the updated value. Use the VNA redeploy API to apply the new form factor to existing appliances.
* `appliance_type` - (Optional, Computed) The virtual network appliance type of the cluster.
* `service_type` - (Optional, Computed, ForceNew) The service type for the cluster. Supported values: `VPC_SERVICES`, `ROUTE_CONTROLLER`. When set to `ROUTE_CONTROLLER`, the cluster is exclusively for route controller and cannot be used to connect VPC workloads. Defaults to `VPC_SERVICES` and cannot be modified after creation.
* `password_managed_by_vcf` - (Optional, Computed) When set to `true`, enables VCF password management for all virtual network appliances in the cluster.
* `advanced_configuration` - (Optional) Advanced configuration for virtual network appliances in the cluster. This block supports:
    * `core_allocation_profile` - (Optional, Computed) Core allocation profile for VNAs in the cluster. Defines core allocation for new or redeployed VNAs. Supported values: `L4FORWARDING`, `L7LBSERVICE`, `L4LBSERVICE`. For `VPC_SERVICES` clusters, NSX defaults to `L4LBSERVICE` when not set. A manual reboot is required for a profile change to take effect on existing appliances.
    * `high_availability_profile` - (Optional, Computed) Path to the high availability profile. If not specified, NSX assigns a default profile.
    * `overlay_transport_zone_path` - (Optional) An overlay transport zone path associated with the VNA host switch and TEP.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `path` - The NSX policy path of the resource.
* `revision` - Indicates current revision number of the object as seen by the NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing Virtual Network Appliance Cluster can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_virtual_network_appliance_cluster.example POLICY_PATH
```

The above command imports the Virtual Network Appliance Cluster named `example` with the NSX policy path `POLICY_PATH`.

Note: The policy path format is `/infra/sites/<site-id>/enforcement-points/<ep-id>/virtual-network-appliance-clusters/<cluster-id>`.
