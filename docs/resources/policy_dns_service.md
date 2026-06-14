---
subcategory: "DNS"
page_title: "NSXT: nsxt_policy_dns_service"
description: A resource to configure a Policy DNS Service.
---

# nsxt_policy_dns_service

This resource provides a method for the management of a Policy DNS Service (aDNS service) within an NSX project.

A PolicyDnsService is deployed on VNS clusters and provides authoritative DNS resolution for VPCs. Listener IPs must reference valid `IpAddressAllocation` objects within the same project.

This resource is applicable to NSX Policy Manager and requires NSX 9.2.0 or higher.

## Example Usage

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_dns_service" "example" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }

  display_name = "my-dns-service"
  description  = "Terraform provisioned DNS Service"

  allocated_listener_ips = [
    nsxt_project_ip_address_allocation.listener.path,
  ]

  vns_clusters = [
    data.nsxt_policy_vns_cluster.cluster1.path,
  ]

  forwarder_config {
    cache_size       = 1000
    upstream_servers = ["8.8.8.8", "8.8.4.4"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `context` - (Required) The context which the object belongs to.
    * `project_id` - (Required) The ID of the project which the object belongs to.
* `allocated_listener_ips` - (Required) Policy paths to `IpAddressAllocation` objects providing IPv4 listener IP addresses for this DNS service. At least one entry is required.
* `vns_clusters` - (Required) Policy paths to VNS clusters on which this DNS service is deployed. At least one entry is required.
* `transit_gateway` - (Optional) Policy path to the transit gateway providing north-south connectivity.
* `forwarder_config` - (Optional) Forwarder and cache settings. When present, enables recursive resolution by forwarding unmatched queries to the configured upstream servers.
    * `cache_size` - (Optional) Number of DNS cache entries (100-100000).
    * `upstream_servers` - (Optional) Upstream DNS server IP addresses for catch-all recursive resolution. Maximum 3 entries.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server.
* `path` - The NSX policy path of the resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_dns_service.example PATH
```

The above command imports Policy DNS Service named `example` using the NSX policy path `PATH`.
The expected path format is `/orgs/[org]/projects/[project]/dns-services/[dns-service]`.
