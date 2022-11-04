---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_lb_http_forwarding_rule"
description: |-
  Provides a resource to configure lb rule on NSX-T manager
---

# nsxt_lb_http_forwarding_rule

Provides a resource to configure lb http forwarding rule on NSX-T manager. This rule will be executed when HTTP request message is forwarded by load balancer.

~> **NOTE:** This resource requires NSX version 2.3 or higher.

## Example Usages
This example represents a superset of all possible action and conditions (and thus doesn't make much sense).
More specific examples are provided below.

```hcl
resource "nsxt_lb_http_forwarding_rule" "lb_rule" {
  description    = "lb_rule provisioned by Terraform"
  display_name   = "lb_rule"
  match_strategy = "ANY"

  tag {
    scope = "color"
    tag   = "red"
  }

  body_condition {
    value          = "XXX"
    match_type     = "CONTAINS"
    case_sensitive = false
  }

  header_condition {
    name       = "header1"
    value      = "bad"
    match_type = "EQUALS"
    inverse    = true
  }

  cookie_condition {
    name           = "name"
    value          = "cookie1"
    match_type     = "STARTS_WITH"
    case_sensitive = true
  }

  cookie_condition {
    name           = "name"
    value          = "cookie2"
    match_type     = "STARTS_WITH"
    case_sensitive = true
  }

  method_condition {
    method = "HEAD"
  }

  version_condition {
    version = "HTTP_VERSION_1_0"
    inverse = true
  }

  uri_condition {
    uri        = "/index.html"
    match_type = "EQUALS"
  }

  ip_condition {
    source_address = "1.1.1.1"
  }

  tcp_condition {
    source_port = 7887
  }

  http_reject_action {
    reply_status  = "500"
    reply_message = "rejected"
  }

  http_redirect_action {
    redirect_status = "200"
    redirect_url    = "/abc.com"
  }

  select_pool_action {
    pool_id = nsxt_lb_pool.pool.id
  }
}
```

The following rule will match if header X-FORWARDED-FOR does not start with "192.168", request method is GET and URI contains "books":

```hcl
resource "nsxt_lb_http_forwarding_rule" "lb_rule1" {
  match_strategy = "ALL"

  header_condition {
    name       = "X-FORWARDED-FOR"
    value      = "192.168"
    match_type = "STARTS_WITH"
    inverse    = true
  }

  method_condition {
    method = "GET"
  }

  uri_condition {
    uri        = "books"
    match_type = "CONTAINS"
  }

  http_reject_action {
    reply_status  = "500"
    reply_message = "rejected"
  }
}
```


The following rule will match if header X-TEST contains "apples" or "pears", regardless of the case:

```hcl
resource "nsxt_lb_http_forwarding_rule" "lb_rule1" {
  match_strategy = "ANY"

  header_condition {
    name           = "X-TEST"
    value          = "apples"
    match_type     = "CONTAINS"
    case_sensitive = false
  }

  header_condition {
    name           = "X-TEST"
    value          = "pears"
    match_type     = "CONTAINS"
    case_sensitive = false
  }

  select_pool_action {
    pool_id = nsxt_lb_pool.pool.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `tag` - (Optional) A list of scope + tag pairs to associate with this lb rule.
* `match_strategy` - (Required) Strategy to define how load balancer rule is considered a match when multiple match conditions are specified in one rule. If set to ALL, then load balancer rule is considered a match only if all the conditions match. If set to ANY, then load balancer rule is considered a match if any one of the conditions match.

* `body_condition` - (Optional) Set of match conditions used to match http request body:
  * `value` - (Required) The value to look for in the body.
  * `match_type` - (Required) Defines how value field is used to match the body of HTTP requests. Accepted values are STARTS_WITH, ENDS_WITH, CONTAINS, EQUALS, REGEX.
  * `case_sensitive` - (Optional) If true, case is significant in the match. Default is true.
  * `inverse` - (Optional) A flag to indicate whether reverse the match result of this condition. Default is false.

* `header_condition` - (Optional) Set of match conditions used to match http request header:
  * `name` - (Required) The name of HTTP header to match.
  * `value` - (Required) The value of HTTP header to match.
  * `match_type` - (Required) Defines how value field is used to match the header value of HTTP requests. Accepted values are STARTS_WITH, ENDS_WITH, CONTAINS, EQUALS, REGEX. Header name field does not support match types.
  * `case_sensitive` - (Optional) If true, case is significant in the match. Default is true.
  * `inverse` - (Optional) A flag to indicate whether reverse the match result of this condition. Default is false.

* `cookie_condition` - (Optional) Set of match conditions used to match http request cookie:
  * `name` - (Required) The name of cookie to match.
  * `value` - (Required) The value of cookie to match.
  * `match_type` - (Required) Defines how value field is used to match the cookie. Accepted values are STARTS_WITH, ENDS_WITH, CONTAINS, EQUALS, REGEX.
  * `case_sensitive` - (Optional) If true, case is significant in the match. Default is true.
  * `inverse` - (Optional) A flag to indicate whether reverse the match result of this condition. Default is false.

* `method_condition` - (Optional) Set of match conditions used to match http request method:
  * `method` - (Required) One of GET, HEAD, POST, PUT, OPTIONS.
  * `inverse` - (Optional) A flag to indicate whether reverse the match result of this condition. Default is false.

* `version_condition` - (Optional) Match condition used to match http version of the request:
  * `version` - (Required) One of HTTP_VERSION_1_0, HTTP_VERSION_1_1.
  * `inverse` - (Optional) A flag to indicate whether reverse the match result of this condition. Default is false.

* `ip_condition` - (Optional) Set of match conditions used to match IP header values of HTTP request:
  * `source_address` - (Required) The value source IP address to match.
  * `inverse` - (Optional) A flag to indicate whether reverse the match result of this condition. Default is false.

* `uri_condition` - (Optional) Set of match conditions used to match http request URI:
  * `uri` - (Required) The value of URI to match.
  * `match_type` - (Required) Defines how value field is used to match the URI. Accepted values are STARTS_WITH, ENDS_WITH, CONTAINS, EQUALS, REGEX.
  * `case_sensitive` - (Optional) If true, case is significant in the match. Default is true.
  * `inverse` - (Optional) A flag to indicate whether reverse the match result of this condition. Default is false.

* `http_reject_action` - (At least one action is required) Set of http reject actions to be executed when load balancer rule matches:
  * `reply_status` - (Required) The HTTP reply status.
  * `reply_message` - (Required) The HTTP reply message.

* `http_redirect_action` - (At least one action is required) Set of http redirect actions to be executed when load balancer rule matches:
  * `redirect_status` - (Required) The HTTP reply status.
  * `redirect_url` - (Required) The URL to redirect to.

* `select_pool_action` - (At least one action is required) Set of pool selection actions to be executed when load balancer rule matches:
  * `pool_id` - (Required) The loadbalancer pool the request will be forwarded to.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the lb rule.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing lb rule can be [imported][docs-import] into this resource, via the following command:
  }
}

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_lb_http_forwarding_rule.lb_rule UUID
```

The above would import the lb rule named `lb_rule` with the nsx id `UUID`
