---
subcategory: "VPC"
page_title: "NSXT: nsxt_vpc"
description: A resource to configure a VPC under Project.
---

# nsxt_vpc

This resource provides a method for the management of a VPC.

This resource is applicable to NSX Policy Manager and is supported with NSX 9.0.0 onwards.

## Example Usage

```hcl
resource "nsxt_vpc" "test" {
  context {
    project_id = nsxt_policy_project.test.id
  }

  display_name        = "test-vpc"
  description         = "Terraform provisioned VPC"
  private_ips         = ["10.1.0.0/16"]
  short_id            = "vpc-test"
  vpc_service_profile = nsxt_vpc_service_profile.test.path
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `private_ips` - (Optional) IP CIDRs to manage private IPv4 subnets.
* `vpc_service_profile` - (Optional) The path of the configuration profile of the VPC services.
* `load_balancer_vpc_endpoint` - (Optional) Configuration for Load Balancer Endpoint
    * `enabled` - (Optional) Flag to indicate whether support for load balancing is needed. Setting this flag to `true` causes allocation of private IPs from the private block associated with this VPC for the use of the load balancer.
* `ip_address_type` - (Optional) This defines the IP address type that will be allocated for subnets.
* `short_id` - (Optional) Defaults to id if id is less than equal to 8 characters or defaults to random generated id if not set. Can not be updated once VPC is created.
* `quotas` - (Optional) List of policy paths for quota resources that are applicable to this VPC.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_vpc.test PATH
```

The above command imports VOC named `test` with the NSX path `PATH`.
