---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_virtual_network_appliance"
description: A resource to configure a Virtual Network Appliance.
---

# nsxt_policy_virtual_network_appliance

This resource provides a method for the management of an individual Virtual Network Appliance (VNA) within a VNA Cluster under an NSX enforcement point.

A Virtual Network Appliance is a VM deployed by NSX inside a VNA Cluster. It provides stateful networking services such as load balancing. This resource allows you to define deployment parameters (vSphere placement, management network, credentials) and update the appliance's display attributes.

This resource is supported with NSX 9.1.1 onwards.

## Example Usage

```hcl
data "nsxt_policy_edge_transport_node" "edge1" {
  display_name = "edge-node-1"
}

data "nsxt_policy_transport_zone" "overlay_tz" {
  display_name = "tz-overlay"
}

resource "nsxt_policy_virtual_network_appliance_cluster" "example" {
  display_name          = "vna-cluster-1"
  appliance_form_factor = "MEDIUM"
  service_type          = "VPC_SERVICES"

  member {
    edge_transport_node_path = data.nsxt_policy_edge_transport_node.edge1.path
  }

  advanced_configuration {
    overlay_transport_zone_path = data.nsxt_policy_transport_zone.overlay_tz.path
  }
}

resource "nsxt_policy_virtual_network_appliance" "example" {
  display_name = "vna-1"
  description  = "VNA node 1"
  cluster_path = nsxt_policy_virtual_network_appliance_cluster.example.path
  hostname     = "vna-1.example.com"

  credentials {
    cli_password  = var.vna_cli_password
    root_password = var.vna_root_password
  }

  management_interface {
    network_id = "dvportgroup-19"

    ip_assignment {
      dhcp_v4 = true
    }
  }

  vm_deployment_config {
    compute_manager_id          = "vc-compute-manager-uuid"
    cluster_or_resource_pool_id = "domain-c1"
    datastore_id                = "datastore-1"
  }

  tag {
    scope = "environment"
    tag   = "production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_path` - (Required, ForceNew) Policy path of the parent `nsxt_policy_virtual_network_appliance_cluster`.
* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional, ForceNew) The NSX ID of this resource. If omitted, a random UUID is generated.
* `hostname` - (Optional) Desired hostname or FQDN for the Virtual Network Appliance VM.
* `failure_domain_path` - (Optional) Policy path of the failure domain for auto-placement of Tier-1 gateways and DHCP servers.
* `credentials` - (Optional) Credentials for the VNA node. Passwords are write-only (used at deployment time only; not returned by NSX GET responses). This block supports:
    * `cli_password` - (Required) Admin CLI password. Sensitive.
    * `root_password` - (Required) Root user password. Sensitive.
    * `audit_password` - (Optional) Audit user password. Sensitive.
    * `cli_username` - (Computed) CLI admin username returned by NSX.
    * `audit_username` - (Computed) Audit username returned by NSX.
* `management_interface` - (Optional) Management interface configuration. This block supports:
    * `network_id` - (Required) Distributed portgroup identifier for the management vNIC.
    * `ip_assignment` - (Required) IP assignment specification. Exactly one assignment type must be specified per entry. Supported types for VNA management: `dhcp_v4`, `static_ipv4`, `static_ipv6`. See `nsxt_policy_edge_transport_node` for full sub-argument documentation.
* `vm_deployment_config` - (Optional) vSphere deployment settings for the VNA VM. This block supports:
    * `compute_manager_id` - (Optional) NSX compute manager (vCenter) UUID.
    * `cluster_or_resource_pool_id` - (Optional) vSphere cluster or resource pool identifier.
    * `datastore_id` - (Optional) vSphere datastore identifier.
    * `reservation_info` - (Optional) Resource reservation settings. This block supports:
        * `cpu_reservation_in_mhz` - (Optional) CPU reservation in MHz.
        * `cpu_reservation_in_shares` - (Optional) CPU reservation priority shares. One of `EXTRA_HIGH_PRIORITY`, `HIGH_PRIORITY` (default), `NORMAL_PRIORITY`, `LOW_PRIORITY`.
        * `memory_reservation_percentage` - (Optional) Memory reservation percentage (0-100, default 100).

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - NSX ID of the Virtual Network Appliance.
* `path` - Policy path of the Virtual Network Appliance.
* `revision` - Indicates current revision of the object as seen by NSX-T API.

## Importing

An existing Virtual Network Appliance can be imported using its full policy path:

```shell
terraform import nsxt_policy_virtual_network_appliance.example \
  /infra/sites/default/enforcement-points/default/virtual-network-appliance-clusters/<cluster-id>/virtual-network-appliances/<vna-id>
```
