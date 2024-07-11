---
subcategory: "VPC"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_vpc"
description: A resource to configure a VPC.
---

# nsxt_policy_vpc

This resource provides a method for the management of a VPC.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_vpc" "test" {
  context {
    project_id = "test_proj"
  }

  display_name = "test"
  description  = "Terraform provisioned VPC"

  external_ipv4_blocks = [data.nsxt_policy_project.test.external_ipv4_blocks.0]

  private_ips = ["192.168.55.2"]
  short_id    = "vpc-tf"

  load_balancer_vpc_endpoint {
    enabled = false
  }

  vpc_dns_forwarder {
    enabled     = true
    listener_ip = "20.32.1.15"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `external_ipv4_blocks` - (Optional) IP block used for allocating CIDR blocks for public subnets. IP block must be subset of Project IPv4 blocks. This field is mutually exclusive with `vpc_connectivity_profile`, `vpc_service_profile`.
* `private_ips` - (Optional) The user is expected to use this field to manage private IPv4 subnets. `private_ipv4_blocks` computed field will be updated internally upon the update of this attribute. For each private IPv4 subnet specified through private_ips field, a private IP block will be created.
* `vpc_service_profile` - (Optional) The path of the configuration profile of the VPC services.
* `load_balancer_vpc_endpoint` - (Optional) Configuration for Load Balancer Endpoint
  * `enabled` - (Optional) Flag to indicate whether support for load balancing is needed. Setting this flag to `true` causes allocation of private IPs from the private block associated with this VPC for the use of the load balancer.
* `ip_address_type` - (Optional) This defines the IP address type that will be allocated for subnets.
* `vpc_dns_forwarder` - (Optional) DNS Forwarder Configuration
  * `enabled` - (Optional) Flag to indicate whether the DNS forwarder is enabled.
  * `listener_ip` - (Optional) The IP on which the DNS forwarder listens. This is allocated from the external IP block.
* `subnet_profiles` - (Optional) Patsh for Subnet Profiles
  * `ip_discovery` - (Optional) Policy path for IP Discovery profile
  * `qos` - (Optional) Policy path for Segment Qos Profile
  * `segment_security` - (Optional) Policy path for Segment Security profile
  * `spoof_guard` - (Optional) Policy path for SpoofGuard profile
  * `mac_discovery` - (Optional) Policy path for Mac Discovery Profile
* `private_ipv4_blocks` - (Optional) IP block used for allocating CIDR blocks for private subnets. IP block must be defined by the Project admin.
* `short_id` - (Optional) Defaults to id if id is less than equal to 8 characters or defaults to random generated id if not set. Can not be updated once VPC is created.
* `vpc_connectivity_profile` - (Optional) The path of the configuration profile for the VPC transit gateway attachment.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.
* `private_ipv4_blocks` - (Optional) Policy paths of automatically created private IPv4 blocks, based on `private_ips` specified for the VPC. 

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_vpc.test PATH
```

The above command imports Vpc named `test` with the NSX path `PATH`.
