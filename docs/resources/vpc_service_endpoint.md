---
subcategory: "VPC"
page_title: "NSXT: nsxt_vpc_service_endpoint"
description: A resource to configure a VPC Service Endpoint.
---

# nsxt_vpc_service_endpoint

This resource provides a method for configuring a VPC Service Endpoint. A VPC service endpoint represents a service VM within a provider VPC that is exposed to authorized consumer tenants via a private link. The service endpoint IP identifies the VM providing the service.

This resource is applicable to NSX Policy Manager and is supported with NSX 9.2.0 onwards.

## Example Usage

```hcl
resource "nsxt_vpc_service_endpoint" "example" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
    vpc_id     = nsxt_vpc.provider_vpc.id
  }
  display_name        = "payment-service-ep"
  description         = "Service endpoint for the payment service VM"
  service_endpoint_ip = "192.168.1.100"
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
* `service_endpoint_ip` - (Required, ForceNew) IP address of the VM providing the service. Must be a valid IPv4 address. Immutable after creation.
* `service_endpoint_ip_type` - (Optional, ForceNew) Indicates what type of resource the IP is assigned to. Only `WORKLOAD` is supported. Defaults to `WORKLOAD`. Immutable after creation.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_vpc_service_endpoint.example PATH
```

The above command imports VpcServiceEndpoint named `example` with the policy path `PATH`.
