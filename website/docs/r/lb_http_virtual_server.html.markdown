---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_lb_http_virtual_server"
description: |-
  Provides a resource to configure lb http virtual server on NSX-T manager
---

# nsxt_lb_http_virtual_server

Provides a resource to configure lb http or https virtual server on NSX-T manager

## Example Usage

```hcl
resource "nsxt_lb_http_application_profile" "http_xff" {
  x_forwarded_for = "INSERT"
}

resource "nsxt_lb_cookie_persistence_profile" "session_persistence" {
  cookie_name = "SESSION"
}

resource "nsxt_lb_pool" "pool1" {
  algorithm = "LEAST_CONNECTION"
  member {
    ip_address = "3.0.0.1"
    port       = "443"
  }
  member {
    ip_address = "3.0.0.2"
    port       = "443"
  }
}

resource "nsxt_lb_pool" "sorry_pool" {
  member {
    ip_address = "3.0.0.15"
    port       = "443"
  }
}

resource "nsxt_lb_http_request_rewrite_rule" "redirect_post" {
  match_strategy = "ALL"
  method_condition {
    method = "POST"
  }

  uri_rewrite_action {
    uri = "/sorry_page.html"
  }
}

resource "nsxt_lb_client_ssl_profile" "ssl1" {
  prefer_server_ciphers = true
}

resource "nsxt_lb_server_ssl_profile" "ssl1" {
  session_cache_enabled = false
}

resource "nsxt_lb_http_virtual_server" "lb_virtual_server" {
  description                = "lb_virtual_server provisioned by terraform"
  display_name               = "virtual server 1"
  access_log_enabled         = true
  application_profile_id     = nsxt_lb_http_application_profile.http_xff.id
  enabled                    = true
  ip_address                 = "10.0.0.2"
  port                       = "443"
  default_pool_member_port   = "8888"
  max_concurrent_connections = 50
  max_new_connection_rate    = 20
  persistence_profile_id     = nsxt_lb_cookie_persistence_profile.session_persistence.id
  pool_id                    = nsxt_lb_pool.pool1.id
  sorry_pool_id              = nsxt_lb_pool.sorry_pool.id
  rule_ids                   = [nsxt_lb_http_request_rewrite_rule.redirect_post.id]

  client_ssl {
    client_ssl_profile_id   = nsxt_lb_client_ssl_profile.ssl1.id
    default_certificate_id  = data.nsxt_certificate.cert1.id
    certificate_chain_depth = 2
    client_auth             = true
    ca_ids                  = [data.nsxt_certificate.ca.id]
    crl_ids                 = ["fa27d79e-13cd-4dd1-a088-4f2fb2cc91f9"]
    sni_certificate_ids     = [data.nsxt_certificate.sni.id]
  }

  server_ssl {
    server_ssl_profile_id   = nsxt_lb_server_ssl_profile.ssl1.id
    client_certificate_id   = data.nsxt_certificate.client.id
    certificate_chain_depth = 2
    server_auth             = true
    ca_ids                  = [data.nsxt_certificate.server_ca.id]
    crl_ids                 = ["fa27d79e-13cd-4dd1-a088-4f2fb2cc91f9"]
  }

  tag {
    scope = "color"
    tag   = "green"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `enabled` - (Optional) Whether the virtual server is enabled. Default is true.
* `ip_address` - (Required) Virtual server IP address.
* `port` - (Required) Virtual server port.
* `tag` - (Optional) A list of scope + tag pairs to associate with this lb http virtual server.
* `access_log_enabled` - (Optional) Whether access log is enabled. Default is false.
* `application_profile_id` - (Required) The application profile defines the application protocol characteristics.
* `default_pool_member_port` - (Optional) Default pool member port.
* `max_concurrent_connections` - (Optional) To ensure one virtual server does not over consume resources, affecting other applications hosted on the same LBS, connections to a virtual server can be capped. If it is not specified, it means that connections are unlimited.
* `max_new_connection_rate` - (Optional) To ensure one virtual server does not over consume resources, connections to a member can be rate limited. If it is not specified, it means that connection rate is unlimited.
* `persistence_profile_id` - (Optional) Persistence profile is used to allow related client connections to be sent to the same backend server.
* `pool_id` - (Optional) Pool of backend servers. Server pool consists of one or more servers, also referred to as pool members, that are similarly configured and are running the same application.
* `sorry_pool_id` - (Optional) When load balancer can not select a backend server to serve the request in default pool or pool in rules, the request would be served by sorry server pool.
* `rule_ids` - (Optional) List of load balancer rules that provide customization of load balancing behavior using match/action rules.
* `client_ssl` - (Optional) Client side SSL customization.
  * `client_ssl_profile_id` - (Required) Id of client SSL profile that defines reusable properties.
  * `default_certificate_id` - (Required) Id of certificate that will be used if the server does not host     multiple hostnames on the same IP address or if the client does not support SNI extension.
  * `certificate_chain_depth` - (Optional) Allowed depth of certificate chain. Default is 3.
  * `client_auth` - (Optional) Whether client authentication is mandatory. Default is false.
  * `ca_ids` - (Optional) List of CA certificate ids for client authentication.
  * `crl_ids` - (Optional) List of CRL certificate ids for client authentication.
  * `sni_certificate_ids` - (Optional) List of certificates to serve different hostnames.

* `server_ssl` - (Optional) Server side SSL customization.
  * `server_ssl_profile_id` - (Required) Id of server SSL profile that defines reusable properties.
  * `server_auth` - (Optional) Whether server authentication is needed. Default is False. If true, ca_ids should be provided.
  * `certificate_chain_depth` - (Optional) Allowed depth of certificate chain. Default is 3.
  * `client_certificate_id` - (Optional) Whether server authentication is required. Default is false.
  * `ca_ids` - (Optional) List of CA certificate ids for server authentication.
  * `crl_ids` - (Optional) List of CRL certificate ids for server authentication.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the lb http virtual server.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing lb http virtual server can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_lb_http_virtual_server.lb_http_virtual_server UUID
```

The above would import the lb http virtual server named `lb_http_virtual_server` with the nsx id `UUID`
