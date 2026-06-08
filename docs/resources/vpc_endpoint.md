---
subcategory: "VPC"
page_title: "NSXT: nsxt_vpc_endpoint"
description: A resource to configure a VPC Endpoint.
---

# nsxt_vpc_endpoint

This resource provides a method for configuring a VPC Endpoint. A VPC endpoint is created within a consumer VPC to connect to a shared VPC service endpoint. The endpoint obtains an IP address from a VPC IP address allocation in the same VPC. VMs in the consumer VPC use this IP as the destination to reach the service.

This resource is applicable to NSX Policy Manager and is supported with NSX 9.2.0 onwards.

## Example Usage

```hcl
resource "nsxt_vpc_endpoint" "example" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
    vpc_id     = nsxt_vpc.consumer_vpc.id
  }
  display_name         = "payment-ep"
  description          = "Client endpoint connecting to payment service"
  vpc_service_endpoint = nsxt_vpc_service_endpoint.example.path
  ip_allocation_path   = nsxt_vpc_ip_address_allocation.example.path
  tag {
    scope = "env"
    tag   = "prod"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `context` - (Required) The context which the object belongs to. Requires a VPC context.
    * `project_id` - (Required) The ID of the project which the object belongs to.
    * `vpc_id` - (Required) The ID of the VPC which the object belongs to.
* `vpc_service_endpoint` - (Required, ForceNew) Policy path to the VPC service endpoint being consumed. The referenced VPC service endpoint must belong to a different VPC within the same project or from another project. Immutable after creation.
* `ip_allocation_path` - (Required, ForceNew) Policy path to the VPC IP address allocation that supplies the client endpoint IP address. The allocation must reside in the same VPC as this VPC endpoint. Immutable after creation.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_vpc_endpoint.example PATH
```

The above command imports VpcEndpoint named `example` with the policy path `PATH`.
