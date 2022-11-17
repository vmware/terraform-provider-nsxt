---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_lb_http_application_profile"
description: |-
  Provides a resource to configure LB HTTP application profile on NSX-T manager
---

# nsxt_lb_http_application_profile

Provides a resource to configure Load Balancer HTTP application profile on NSX-T manager

~> **NOTE:** This resource requires NSX version 2.3 or higher.

## Example Usage

```hcl
resource "nsxt_lb_http_application_profile" "lb_http_application_profile" {
  description            = "lb_http_application_profile provisioned by Terraform"
  display_name           = "lb_http_application_profile"
  http_redirect_to       = "http://www.example.com"
  http_redirect_to_https = "false"
  idle_timeout           = "15"
  request_body_size      = "100"
  request_header_size    = "1024"
  response_timeout       = "60"
  x_forwarded_for        = "INSERT"
  ntlm                   = "true"

  tag {
    scope = "color"
    tag   = "red"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `http_redirect_to` - (Optional) A URL that incoming requests for that virtual server can be temporarily redirected to, If a website is temporarily down or has moved. When set, http_redirect_to_https should be false.
* `http_redirect_to_https` - (Optional) A boolean flag which reflects whether the client will automatically be redirected to use SSL. When true, the http_redirect_to should not be specified.
* `idle_timeout` - (Optional) Timeout in seconds to specify how long an HTTP application can remain idle. Defaults to 15 seconds.
* `ntlm` - (Optional) A boolean flag which reflects whether NTLM challenge/response methodology will be used over HTTP. Can be set to true only if http_redirect_to_https is false.
* `request_body_size` - (Optional) Maximum request body size in bytes. If it is not specified, it means that request body size is unlimited.
* `request_header_size` - (Optional) Maximum request header size in bytes. Requests with larger header size will be processed as best effort whereas a request with header below this specified size is guaranteed to be processed. Defaults to 1024 bytes.
* `response_timeout` - (Optional) Number of seconds waiting for the server response before the connection is closed. Defaults to 60 seconds.
* `x_forwarded_for` - (Optional) When this value is set, the x_forwarded_for header in the incoming request will be inserted or replaced. Supported values are "INSERT" and "REPLACE".
* `tag` - (Optional) A list of scope + tag pairs to associate with this lb http profile.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the lb http application profile.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing lb http profile can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_lb_http_application_profile.lb_http_application_profile UUID
```

The above would import the LB HTTP application profile named `lb_http_application_profile` with the nsx id `UUID`
