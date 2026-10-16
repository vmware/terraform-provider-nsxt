---
subcategory: "VPC"
page_title: "NSXT: nsxt_vpc_lb_http_application_profile"
description: A resource to configure a VPC LBHttpApplicationProfile.
---

# nsxt_vpc_lb_http_application_profile

This resource provides a method for the management of a VPC LBHttpApplicationProfile.

This resource is applicable to NSX Policy Manager and is supported with NSX 9.0.2 onwards.

## Example Usage

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_vpc" "demovpc" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "vpc1"
}

resource "nsxt_vpc_lb_http_application_profile" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
    vpc_id     = data.nsxt_vpc.demovpc.id
  }
  display_name           = "test"
  description            = "Terraform provisioned LBHttpProfile"
  http_redirect_to       = "http://www.google.com"
  http_redirect_to_https = false
  idle_timeout           = 15
  request_body_size      = 256
  request_header_size    = 1024
  response_buffering     = true
  response_header_size   = 4096
  response_timeout       = 100
  server_keep_alive      = true
  x_forwarded_for        = "REPLACE"

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
* `http_redirect_to` - (Optional) If a website is temporarily down or has moved, incoming requests for that virtual server can be temporarily redirected to a URL.
* `http_redirect_to_https` - (Optional) A boolean flag which reflects whether the client will automatically be redirected to use SSL.
* `idle_timeout` - (Optional) Timeout in seconds to specify how long an HTTP application can remain idle.
* `request_body_size` - (Optional) Maximum request body size in bytes (Unlimited if not specified).
* `request_header_size` - (Optional) Maximum request header size in bytes. Requests with larger header size will be processed as best effort whereas a request with header below this specified size is guaranteed to be processed.
* `response_buffering` - (Optional) A boolean flag indicating whether the response received by LB from the backend will be saved into the buffers. When buffering is disabled, the response is passed to a client synchronously, immediately as it is received. When buffering is enabled, LB receives a response from the backend server as soon as possible, saving it into the buffers.
* `response_header_size` - (Optional) Maximum request header size in bytes. Requests with larger header size will be processed as best effort whereas a request with header below this specified size is guaranteed to be processed.
* `response_timeout` - (Optional) Number of seconds waiting for the server response before the connection is closed.
* `server_keep_alive` - (Optional) A boolean flag indicating whether the backend connection will be kept alive for client connection. If server_keep_alive is true, it means the backend connection will keep alive for the client connection. Every client connection is tied 1:1 with the corresponding server-side connection. If server_keep_alive is false, it means the backend connection won't keep alive for the client connection.
* `x_forwarded_for` - (Optional) When X-Forwarded-For is configured, X-Forwarded-Proto and X-Forwarded-Port information is added automatically to the request header. Possible values are:`INSERT`, `REPLACE`

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_vpc_lb_http_application_profile.test PATH
```

The above command imports the VPC LB http application profile named `test` with the policy path `PATH`, which must be a VPC path of the format `/orgs/[org]/projects/[project]/vpcs/[vpc]/vpc-lb-app-profiles/[id]`. The project and VPC context is derived from the path.
