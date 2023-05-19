---
subcategory: "Load Balancer"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_lb_virtual_server"
description: A resource to configure a Load Balancer Virtual Server.
---

# nsxt_policy_lb_virtual_server

This resource provides a method for the management of a Load Balancer Virtual Server.

This resource is applicable to NSX Policy Manager.

Note that the 'rule' section has been added at a later date. In order to preserve backward compatibility for users that have created rules manually, existing ("live") rules will not be changed if there is no 'rule' section present in the resource definition. If you want to delete manually created rules from a managed resource, you might have to initially add a 'rule' section and subsequentially delete it again. 

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

  rule {
    display_name   = "rule_test"
    match_strategy = "ALL"
    phase          = "HTTP_REQUEST_REWRITE"

    action {
      http_request_header_delete {
        header_name = "X-something"
      }
    }
    condition {
      http_request_body {
        body_value = "xyz"
        match_type = "REGEX"
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `application_profile_path` - (Required) Application profile path for this virtual server. Note that this also differentiates between Layer 4 TCP/UDP and Layer 7 HTTP virtual servers. 
* `access_log_enabled` - (Optional) If set, all connections/requests sent to the virtual server are logged to access log.
* `default_pool_member_ports` - (Optional) Default pool member ports to use when member port is not defined on the pool.
* `enabled` - (Optional) Flag to enable this Virtual Server.
* `ip_address` - (Required) Virtual Server IP address.
* `ports` - (Required) Virtual Server Ports.
* `persistence_profile_path` - (Optional) Path to persistence profile allowing related client connections to be sent to the same backend server.
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

* `rule` - (Optional) Specifies one or more rules to manipulate traffic passing through HTTP or HTTPS virtual server.
  * `display_name` - (Optional) Display name of the rule.
  * `match_strategy` - (Optional) Match strategy for determining match of multiple conditions, one of `ALL`, `ANY`. Default is `ANY`.
  * `phase` - (Optional) Load balancer processing phase, one of `HTTP_REQUEST_REWRITE`, `HTTP_FORWARDING`, `HTTP_RESPONSE_REWRITE`, `HTTP_ACCESS` or `TRANSPORT`. Default is `HTTP_FORWARDING`.

  * `action` - (Required) A list of actions to be executed at specified phase when load balancer rule matches.

    * `connection_drop` - (Optional) Action to drop the connections. (There is no argument to this action)

    * `http_redirect_acion` - (Optional) Action to redirect HTTP request message toto a new URL.
      * `redirect_status` - (Required) HTTP response status code.
      * `redirect_url` - (Required) The URL that the HTTP request is redirected to.

    * `http_reject` - (Optional) Action to reject HTTP request message.
      * `reply_message` - (Optional) Response message.
      * `reply_status` - (Required) HTTP response status code.

    * `http_request_header_delete` - (Optional) Action to delete header fields of HTTP request messages.
      * `header_name`- (Required) Name of a header field of HTTP request message.

    * `http_request_header_rewrite` - (Optional) Action to rewrite header fields of matched HTTP request messages to specified new values.
      * `header_name` - (Required) Name of HTTP request header.
      * `header_value`  - (Required) Value of HTTP request header.

    * `http_request_uri_rewrite` - (Optional) Action to rewrite URIs in matched HTTP request messages.
      * `uri` - (Required) URI of HTTP request.
      * `uri_arguments` - (Optional) URI arguments.

    * `http_response_header_delete` - (Optional) Action to delete header fields of HTTP response messages.
      * `header_name` - (Required) Name of a header field of HTTP response messages.

    * `http_response_header_rewrite` - (Optional) Action to rewrite header fields of matched HTTP request message to specified new values.
      * `header_name` - (Required) Name of HTTP request header.
      * `header_value` - (Required) Value of HTTp request header.

    * `jwt_auth` - (Optional) Action to control access to backend server resources using JSON Web Token (JWT) authentication.
      * `pass_jwt_to_pool` - (Optional) Whether to pass JWT to backend server or remove it, Boolean, Default `false`.
      * `realm` - (Optional) JWT realm.
      * `tokens` - (Optional) List of JWT tokens.
      * `key` - (Optional) Key to verify signature of JWT token, specify exactly one of the arguments.
        * `certificate_path` - (Optional) Use certficate to verify signature of JWT token.
        * `public_key_content` - (Optional) Use public key to verify signature of JWT token.
        * `symmetric_key` - (Optional) Use symmetric key to verify signature of JWT token, this argument indicates presence only, the value is discarded.

    * `select_pool` - (Optional) Action used to select a pool for matched HTTP request messages.
      * `pool_id` - (Required) Path of load balancer pool.

    * `ssl_mode_selection` - (Optional) Action to select SSL mode.
      * `ssl_mode` - (Required) Type of SSL mode, one of `SSL_PASSTHROUGH`, `SSL_END_TO_END` or `SSL_OFFLOAD`.

    * `variable_assignment` - (Optional) Action to create new variable and assign value to it.
      * `variable_name` - (Required) Name of the variable to be assigned.
      * `variable_value` - (Required) Value of variable.

    * `variable_persistence_learn` - (Optional) Action to learn the value of variable from the HTTP response.
      * `persistence_profile_path` - (Optional) Path to nsxt_policy_persistence_profile.
      * `variable_hash_enabled` - (Optional) Whether to enable a hash operation for variable value, Boolean, Default `false`.
      * `variable_name` - (Required) Variable name.

    * `variable_persistence_on` - (Optional) Action to inspect the variable of HTTP request.
      * `persistence_profile_path` - (Optional) Path to nsxt_policy_persistence_profile.
      * `variable_hash_enabled` - (Optional) Whether to enable a hash operation for variable value, Boolean, Default `false`.
      * `variable_name` - (Required) Variable name.

  * `condition` - (Optional) A list of match conditions used to match application traffic.

    * `http_request_body` - (Optional) Condition to match the message body of an HTTP request.
      * `body_value` - (Required) HTTP request body.
      * `case_sensitive` - (Optional) A case sensitive flag for HTTP body comparison, Boolean, Default `true`.
      * `inverse` - (Optional) A flag to indicate whether to reverse the match result of this condition, Boolean, Default `false`
      * `match_type` - (Optional) Match type of HTTP body, one of `REGEX`, `STARTS_WITH`, `ENDS_WITH`. `EQUALS` or `CONTAINS`. Default is `REGEX`.

    * `http_request_cookie` - (Optional) Condition to match HTTP request messages by cookie.
      * `cookie_name` - (Required) Name of cookie.
      * `cookie_value` - (Required) Value of cookie.
      * `case_sensitive` - (Optional) A case sensitive flag for cookie comparison, Boolean, Default `true`.
      * `inverse` - (Optional) A flag to indicate whether to reverse the match result of this condition, Boolean, Default `false`
      * `match_type` - (Optional) Match type of cookie, one of `REGEX`, `STARTS_WITH`, `ENDS_WITH`. `EQUALS` or `CONTAINS`. Default is `REGEX`.

    * `http_request_header` - (Optional) Condition to match HTTP request messages by HTTP header fields.
      * `header_name` - (Required) Name of HTTP header.
      * `header_value` - (Required) Value of HTTP header.
      * `case_sensitive` - (Optional) A case sensitive flag for HTTP header comparison, Boolean, Default `true`.
      * `inverse` - (Optional) A flag to indicate whether to reverse the match result of this condition, Boolean, Default `false`
      * `match_type` - (Optional) Match type of HTTP header, one of `REGEX`, `STARTS_WITH`, `ENDS_WITH`. `EQUALS` or `CONTAINS`. Default is `REGEX`.

    * `http_request_method` - (Optional) Condition to match method of HTTP requests.
      * `method` - (Required) Type of HTTP request method, one of `GET`, `OPTIONS`, `POST`, `HEAD` or `PUT`.
      * `inverse` - (Optional) A flag to indicate whether to reverse the match result of this condition, Boolean, Default `false`

    * `http_request_uri_arguments` - (Optional) Condition to match URI arguments.
      * `uri_arguments` - (Required) URI arguments.
      * `case_sensitive` - (Optional) A case sensitive flag for HTTP uri arguments comparison, Boolean, Default `true`.
      * `inverse` - (Optional) A flag to indicate whether to reverse the match result of this condition, Boolean, Default `false`
      * `match_type` - (Optional) Match type of HTTP uri arguments, one of `REGEX`, `STARTS_WITH`, `ENDS_WITH`. `EQUALS` or `CONTAINS`. Default is `REGEX`.

    * `http_request_uri` - (Optional) Condition to match URIs of HTTP requests messages.
      * `uri` - (Required) A string used to identify resource.
      * `case_sensitive` - (Optional) A case sensitive flag for HTTP uri comparison, Boolean, Default `true`.
      * `inverse` - (Optional) A flag to indicate whether to reverse the match result of this condition, Boolean, Default `false`
      * `match_type` - (Optional) Match type of HTTP uri, one of `REGEX`, `STARTS_WITH`, `ENDS_WITH`. `EQUALS` or `CONTAINS`. Default is `REGEX`.

    * `http_request_version` - (Optional) Condition to match the HTTP protocol version of the HTTP request messages.
      * `version` (Required) HTTP version, one of `HTTP_VERSION_1_0` or `HTTP_VERSION_1_1`.
      * `inverse` - (Optional) A flag to indicate whether to reverse the match result of this condition, Boolean, Default `false`

    * `http_response_header` - (Optional) Condition to match HTTP response messages from backend servers by HTTP header fields.
      * `header_name` - (Required) Name of HTTP header field.
      * `header_value` - (Required) Value of HTTP header field.
      * `case_sensitive` - (Optional) A case sensitive flag for HTTP header comparison, Boolean, Default `true`.
      * `inverse` - (Optional) A flag to indicate whether to reverse the match result of this condition, Boolean, Default `false`
      * `match_type` - (Optional) Match type of HTTP header, one of `REGEX`, `STARTS_WITH`, `ENDS_WITH`. `EQUALS` or `CONTAINS`. Default is `REGEX`.

    * `http_ssl` - (Optional) Condition to match SSL handshake and SSL connection.
      * `client_certificate_issuer_dn` - (Optional) The issuer DN match condition of the client certificate.
        * `issuer_dn` - (Required) Value of issuer DN.
        * `case_sensitive` - (Optional) A case sensitive flag for issuer DN comparison, Boolean, Default `true`.
        * `match_type` - (Optional) Match type of issuer DN, one of `REGEX`, `STARTS_WITH`, `ENDS_WITH`. `EQUALS` or `CONTAINS`. Default is `REGEX`.
      * `client_certificate_subject_dn` - (Optional) The subject DN match condition of the client certificate.
        * `subject_dn` - (Required) Value of subject DN.
        * `case_sensitive` - (Optional) A case sensitive flag for subject DN comparison, Boolean, Default `true`.
        * `match_type` - (Optional) Match type of subject DN, one of `REGEX`, `STARTS_WITH`, `ENDS_WITH`. `EQUALS` or `CONTAINS`. Default is `REGEX`.
      * `client_support_ssl_ciphers` - (Optional) List of ciphers supported by client (see documentation for possible values).
      * `inverse` - (Optional) A flag to indicate whether to reverse the match result of this condition, Boolean, Default `false`
      * `session_reused` - (Optional) The type of SSL session reused, one of `IGNORE`, `REUSED` or `NEW`. Default is `IGNORE`.
      * `used_protocol` - (Optional) Protocol of an established SSL connection, one of `SSL_V2`, `SSL_V3`, `TLS_V1`, `TLS_V1_1` or `TLS_V1_2`.
      * `used_ssl_cipher` - (Optional) Cypher used for an established SSL connection (see documentation for possible values).

    * `ip_header` - (Optional) Condition to match IP header fields.
      * `group_path` - (Optional) Grouping object path.
      * `source_address`  - (Optional) Source IP address, range or subnet of HTTP message.
      * `inverse` - (Optional) A flag to indicate whether to reverse the match result of this condition, Boolean, Default `false`

    * `ssl_sni` - (Optional) Condition to match SSL SNI in client hello.
      * `sni` - (Required) The server name indication.
      * `case_sensitive` - (Optional) A case sensitive flag for SNI comparison, Boolean, Default `true`.
      * `inverse` - (Optional) A flag to indicate whether to reverse the match result of this condition, Boolean, Default `false`
      * `match_type` - (Optional) Match type of SNI, one of `REGEX`, `STARTS_WITH`, `ENDS_WITH`. `EQUALS` or `CONTAINS`. Default is `REGEX`.

    * `tcp_header` - (Optional) Condition to match TCP header fields.
      * `source_port` - (Required) TCP source port or port range of HTTP message.
      * `inverse` - (Optional) A flag to indicate whether to reverse the match result of this condition, Boolean, Default `false`

    * `variable` - (Optional) Condition to match variable's name and value.
      * `variable_name` - (Required) Name of the variable to be matched.
      * `variable_value` - (Required) Value of the variable to be matched.
      * `case_sensitive` - (Optional) A case sensitive flag for variable comparison, Boolean, Default `true`.
      * `inverse` - (Optional) A flag to indicate whether to reverse the match result of this condition, Boolean, Default `false`
      * `match_type` - (Optional) Match type of variable, one of `REGEX`, `STARTS_WITH`, `ENDS_WITH`. `EQUALS` or `CONTAINS`. Default is `REGEX`.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing Virtual Server can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_lb_virtual_server.test ID
```

The above command imports Load Balancer Virtual Server named `test` with the NSX Load Balancer Virtual Server ID `ID`.
