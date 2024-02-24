---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_lb_https_monitor_profile"
description: A resource to configure a LBHttpsMonitorProfile.
---

# nsxt_policy_lb_https_monitor_profile

This resource provides a method for the management of a LBHttpsMonitorProfile.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_lb_https_monitor_profile" "test" {
  display_name          = "test"
  description           = "Terraform provisioned LBHttpsMonitorProfile"
  request_body          = "test"
  request_method        = "POST"
  request_url           = "test"
  request_version       = "HTTP_VERSION_1_1"
  response_body         = "test"
  response_status_codes = [200]
  fall_count            = 2
  interval              = 2
  monitor_port          = 8080
  rise_count            = 2
  timeout               = 2
  server_ssl {
    certificate_chain_depth = 3
    server_auth             = "IGNORE"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `request_body` - (Optional) String to send as part of HTTP health check request body. Valid only for certain HTTP methods like POST.
* `request_header` - (Optional) Array of HTTP request headers.
  * `header_name` - (Optional) Name of HTTP request header
  * `header_value` - (Optional) Value of HTTP request header
* `request_method` - (Optional) The health check method for HTTP monitor type.
* `request_url` - (Optional) For HTTPS active healthchecks, the HTTPS request url sent can be customized and can include query parameters.
* `request_version` - (Optional) HTTP request version. Possible values are: `HTTP_VERSION_1_0`, `HTTP_VERSION_1_1`.
* `response_body` - (Optional) If HTTP response body match string (regular expressions not supported) is specified then the healthcheck HTTP response body is matched against the specified string and server is considered healthy only if there is a match. If the response body string is not specified, HTTP healthcheck is considered successful if the HTTP response status code is 2xx, but it can be configured to accept other status codes as successful.
* `response_status_codes` - (Optional) The HTTP response status code should be a valid HTTP status code.
* `server_ssl` - (Optional) 
  * `certificate_chain_depth` - (Optional) Authentication depth is used to set the verification depth in the server certificates chain.
  * `client_certificate_path` - (Optional) To support client authentication (load balancer acting as a client authenticating to the backend server), client certificate can be specified in the server-side SSL profile binding
  * `server_auth` - (Optional) Server authentication mode. Possible values are: `REQUIRED`, `IGNORE`, `AUTO_APPLY`.
  * `server_auth_ca_paths` - (Optional) If server auth type is REQUIRED, server certificate must be signed by one of the trusted Certificate Authorities (CAs), also referred to as root CAs, whose self signed certificates are specified.
  * `server_auth_crl_paths` - (Optional) A Certificate Revocation List (CRL) can be specified in the server-side SSL profile binding to disallow compromised server certificates.
  * `ssl_profile_path` - (Optional) Server SSL profile defines reusable, application-independent server side SSL properties.
* `fall_count` - (Optional) Mark member status DOWN if the healtcheck fails consecutively for fall_count times.
* `interval` - (Optional) Active healthchecks are initiated periodically, at a configurable interval (in seconds), to each member of the Group.
* `monitor_port` - (Optional) Typically, monitors perform healthchecks to Group members using the member IP address and pool_port. However, in some cases, customers prefer to run healthchecks against a different port than the pool member port which handles actual application traffic. In such cases, the port to run healthchecks against can be specified in the monitor_port value.
* `rise_count` - (Optional) Bring a DOWN member UP if rise_count successive healthchecks succeed.
* `timeout` - (Optional) Timeout specified in seconds. After a healthcheck is initiated, if it does not complete within a certain period, then also the healthcheck is considered to be unsuccessful. Completing a healthcheck within timeout means establishing a connection (TCP or SSL), if applicable, sending the request and receiving the response, all within the configured timeout.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_lb_https_monitor_profile.test UUID
```

The above command imports LBHttpsMonitorProfile named `test` with the NSX ID `UUID`.
