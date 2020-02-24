---
layout: "nsxt"
page_title: "NSXT: nsxt_policy_lb_virtual_server"
sidebar_current: "docs-nsxt-resource-policy-lb-virtual-server"
description: A resource to configure a Load Balancer Virtual Server.
---

# nsxt_policy_lb_virtual_server

This resource provides a method for the management of a Load Balancer Virtual Server.
 
## Example Usage

```hcl
resource "nsxt_policy_lb_virtual_server" "test" {
  display_name               = "test"
  description                = "Terraform provisioned Virtual Server"
  access_log_enabled         = true
  application_profile_path   = data.nsxt_policy_lb_app_profile.tcp.path
  enabled                    = true
  ip_address                 = "10.10.10.21"
  ports                      = ["443"]
  default_pool_member_ports  = ["80"]
  service_path               = nsxt_policy_lb_service.app1.path
  max_concurrent_connections = 6
  max_new_connection_rate    = 20
  pool_path                  = nsxt_policy_lb_pool.pool1.path
  sorry_pool_path            = nsxt_policy_lb_pool.sorry_pool.path

  client_ssl {
    client_auth              = "REQUIRED"
    default_certificate_path = data.nsxt_policy_certificate.cert1.path
    ca_paths                 = [data.nsxt_policy_certificate.lb_ca.path]
    certificate_chain_depth  = 3
    ssl_profile_path         = data.nsxt_policy_lb_client_ssl_profile.lb_profile.path
  }

  server_ssl {
    server_auth             = "REQUIRED"
    client_certificate_path = data.nsxt_policy_certificate.client_ca.path
    certificate_chain_depth = 3
    ssl_profile_path        = data.nsxt_policy_lb_server_ssl_profile.lb_profile.path
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `application_profile_path` - (Required) Application profile path for this virtual server.
* `access_log_enabled` - (Optional) If set, all connections/requests sent to the virtual server are logged to access log.
* `default_pool_member_ports` - (Optional) Default pool member ports to use when member port is not defined on the pool.
* `enabled` - (Optional) Flag to enable this Virtual Server.
* `ip_address` - (Required) Virtual Server IP address.
* `ports` - (Required) Virtual Server Ports.
* `persistence_profile_path` - (Optional) Path to persistance profile allowing related client connections to be sent to the same backend server.
* `service_path` - (Optional) Virtual Server can be associated with Load Balancer Service.
* `max_concurrent_connections` - (Optional) To ensure one virtual server does not over consume resources, connections to Virtual Server can be capped.
* `max_new_connection_rate` - (Optional) To ensure one virtual server does not over consume resources, connections to a member can be rate limited.
* `pool_path` - (Optional)Path for Load Balancer Pool.
* `sorry_pool_path` - (Optional) When load balancer can not select server in pool, the request would be served by sorry pool

* `server_ssl` - (Optional)
  * `server_auth` - (Optional) Server Authentication Mode, one of `REQUIRED`, `IGNORE`, `AUTO_APPLY`. Default is `AUTO_APPLY`.
  * `certificate_chain_depth` - (Optional) Allowed certificate chain depth.
  * `client_certificate_path` - (Optional) Client certificat path for client authentication against the server.
  * `ca_paths` - (Optional) If server auth type is REQUIRED, client certificate must be signed by one Certificate Authorities provided here.
  * `crl_paths` - (Optional) Certificate Revocation Lists can be specified to disallow compromised certificate.
  * `ssl_profile_path` - (Optional) Server SSL profile path.

* `client_ssl` - (Optional)
  * `client_auth` - (Optional) Client Authentication Mode, one of `REQUIRED`, `IGNORE`. Default is `IGNORE`.
  * `certificate_chain_depth` - (Optional) Allowed certificate chain depth.
  * `ca_paths` - (Optional) If client auth type is REQUIRED, client certificate must be signed by one Certificate Authorities provided here.
  * `crl_paths` - (Optional) Certificate Revocation Lists can be specified to disallow compromised client certificate.
  * `default_certificate_path` - (Optional) Default Certificate Path. Must be specified if client_auth is set to `REQUIRED`.
  * `sni_paths` - (Optional) This setting allows multiple certificates(for different hostnames) to be bound to the same virtual server.
  * `ssl_profile_path` - (Optional) Client SSL profile path.
* `log_significant_event_only` - (Optional) If true, significant events are logged in access log. This flag is supported since NSX 3.0.0.
* `access_list_control` - (Optional) Specifies the access list control to define how to filter client connections.
  * `action` - (Required) Action for connections matching the grouping object.
  * `group_path` - (Required) The path of grouping object which defines the IP addresses or ranges to match client IP.
  * `enabled` - (Optional) Indicates whether to enable access list control option. Default is true.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing Virtual Server can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_policy_lb_virtual_server.test ID
```

The above command imports Load Balancer Virtual Server named `test` with the NSX Load Balancer Virtual Server ID `ID`.
