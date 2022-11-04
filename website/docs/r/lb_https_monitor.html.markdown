---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_lb_https_monitor"
description: |-
  Provides a resource to configure lb https monitor on NSX-T manager
---

# nsxt_lb_https_monitor

Provides a resource to configure lb https monitor on NSX-T manager

## Example Usage

```hcl

data "nsxt_certificate" "client" {
  display_name = "client-1"
}

data "nsxt_certificate" "CA" {
  display_name = "ca-1"
}

resource "nsxt_lb_https_monitor" "lb_https_monitor" {
  description             = "lb_https_monitor provisioned by Terraform"
  display_name            = "lb_https_monitor"
  fall_count              = 2
  interval                = 5
  monitor_port            = 8080
  rise_count              = 5
  timeout                 = 10
  certificate_chain_depth = 2
  ciphers                 = ["TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256", "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384"]
  client_certificate_id   = data.nsxt_certificate.client.id
  protocols               = ["TLS_V1_2"]
  request_body            = "ping"
  request_method          = "HEAD"
  request_url             = "/index.html"
  request_version         = "HTTP_VERSION_1_1"
  response_body           = "pong"
  response_status_codes   = [200, 304]
  server_auth             = "REQUIRED"
  server_auth_ca_ids      = data.nsxt_certificate.CA.id
  server_auth_crl_ids     = ["78ba3814-bfe1-45e5-89d3-46862bed7896"]

  request_header {
    name  = "X-healthcheck"
    value = "NSX"
  }

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
* `tag` - (Optional) A list of scope + tag pairs to associate with this lb https monitor.
* `fall_count` - (Optional) Number of consecutive checks that must fail before marking it down.
* `interval` - (Optional) The frequency at which the system issues the monitor check (in seconds).
* `monitor_port` - (Optional) If the monitor port is specified, it would override pool member port setting for healthcheck. A port range is not supported.
* `rise_count` - (Optional) Number of consecutive checks that must pass before marking it up.
* `timeout` - (Optional) Number of seconds the target has to respond to the monitor request.
* `certificate_chain_depth` - (Optional) Authentication depth is used to set the verification depth in the server certificates chain.
* `ciphers` - (Optional) List of supported SSL ciphers.
* `client_certificate_id` - (Optional) Client certificate can be specified to support client authentication.
* `protocols` - (Optional) SSL versions TLS1.1 and TLS1.2 are supported and enabled by default. SSLv2, SSLv3, and TLS1.0 are supported, but disabled by default.
* `request_body` - (Optional) String to send as HTTP health check request body. Valid only for certain HTTP methods like POST.
* `request_header` - (Optional) HTTP request headers.
* `request_method` - (Optional) Health check method for HTTP monitor type. Valid values are GET, HEAD, PUT, POST and OPTIONS.
* `request_url` - (Optional) URL used for HTTP monitor.
* `request_version` - (Optional) HTTP request version. Valid values are HTTP_VERSION_1_0 and HTTP_VERSION_1_1.
* `response_body` - (Optional) If response body is specified, healthcheck HTTP response body is matched against the specified string and server is considered healthy only if there is a match (regular expressions not supported). If response body string is not specified, HTTP healthcheck is considered successful if the HTTP response status code is among configured values.
* `response_status_codes` - (Optional) HTTP response status code should be a valid HTTP status code.
* `server_auth` - (Optional) Server authentication mode - REQUIRED or IGNORE.
* `server_auth_ca_ids` - (Optional) If server auth type is REQUIRED, server certificate must be signed by one of the trusted Certificate Authorities (CAs), also referred to as root CAs, whose self signed certificates are specified.
* `server_auth_crl_ids` - (Optional) A Certificate Revocation List (CRL) can be specified in the server-side SSL profile binding to disallow compromised server certificates.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the lb_https_monitor.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `is_secure` - This flag is set to true when all the ciphers and protocols are secure. It is set to false when one of the ciphers or protocols is insecure.


## Importing

An existing lb https monitor can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_lb_https_monitor.lb_https_monitor UUID
```

The above would import the lb https monitor named `lb_https_monitor` with the nsx id `UUID`
