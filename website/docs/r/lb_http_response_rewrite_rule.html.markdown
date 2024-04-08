---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_lb_http_response_rewrite_rule"
description: |-
  Provides a resource to configure lb rule on NSX-T manager
---

# nsxt_lb_http_response_rewrite_rule

Provides a resource to configure lb http response rewrite rule on NSX-T manager. This rule will be executed when HTTP response message is received by load balancer.

~> **NOTE:** This resource requires NSX version 2.3 or higher.

## Example Usages
This example represents a superset of all possible conditions (and thus doesn't make much sense).
More specific examples are provided below.

```hcl
resource "nsxt_lb_http_response_rewrite_rule" "lb_rule" {
  description    = "lb_rule provisioned by Terraform"
  display_name   = "lb_rule"
  match_strategy = "ALL"

  tag {
    scope = "color"
    tag   = "blue"
  }

  request_header_condition {
    name       = "header1"
    value      = "bad"
    match_type = "EQUALS"
    inverse    = true
  }

  response_header_condition {
    name       = "header1"
    value      = "good"
    match_type = "EQUALS"
    inverse    = false
  }

  cookie_condition {
    name           = "name1"
    value          = "cookie1"
    match_type     = "STARTS_WITH"
    case_sensitive = true
  }

  cookie_condition {
    name           = "name2"
    value          = "cookie2"
    match_type     = "STARTS_WITH"
    case_sensitive = true
  }

  method_condition {
    method = "HEAD"
  }

  version_condition {
    version = "HTTP_VERSION_1_1"
    inverse = true
  }

  uri_condition {
    uri        = "/index.html"
    match_type = "EQUALS"
  }

  uri_arguments_condition {
    uri_arguments = "delete"
    match_type    = "CONTAINS"
    inverse       = true
  }

  ip_condition {
    source_address = "1.1.1.1"
  }

  tcp_condition {
    source_port = 7887
  }

  header_rewrite_action {
    name  = "header1"
    value = "even better"
  }
}
```

The following rule will match if request header X-FORWARDED-FOR does not start with "192.168", request method is GET and response content is json:

```hcl
resource "nsxt_lb_http_response_rewrite_rule" "lb_rule1" {
  match_strategy = "ALL"

  request_header_condition {
    name       = "X-FORWARDED-FOR"
    value      = "192.168"
    match_type = "STARTS_WITH"
    inverse    = true
  }

  response_header_condition {
    name       = "Content-Type"
    value      = "/json"
    match_type = "CONTAINS"
    inverse    = false
  }

  method_condition {
    method = "GET"
  }

  header_rewrite_action {
    name  = "header1"
    value = "value2"
  }
}
```


The following rule will match if response header X-TEST contains "apples" or "pears", regardless of the case:

```hcl
resource "nsxt_lb_http_response_rewrite_rule" "lb_rule1" {
  match_strategy = "ANY"

  response_header_condition {
    name           = "X-TEST"
    value          = "apples"
    match_type     = "CONTAINS"
    case_sensitive = false
  }

  response_header_condition {
    name           = "X-TEST"
    value          = "pears"
    match_type     = "CONTAINS"
    case_sensitive = false
  }

  header_rewrite_action {
    name  = "header1"
    value = "value2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of this resource.
* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `tag` - (Optional) A list of scope + tag pairs to associate with this lb rule.
* `match_strategy` - (Required) Strategy to define how load balancer rule is considered a match when multiple match conditions are specified in one rule. If set to ALL, then load balancer rule is considered a match only if all the conditions match. If set to ANY, then load balancer rule is considered a match if any one of the conditions match.

* `request_header_condition` - (Optional) Set of match conditions used to match http request header:
  * `name` - (Required) The name of HTTP header to match.
  * `value` - (Required) The value of HTTP header to match.
  * `match_type` - (Required) Defines how value field is used to match the header value of HTTP request. Accepted values are STARTS_WITH, ENDS_WITH, CONTAINS, EQUALS, REGEX. Header name field does not support match types.
  * `case_sensitive` - (Optional) If true, case is significant in the match. Default is true.
  * `inverse` - (Optional) A flag to indicate whether reverse the match result of this condition. Default is false.

* `response_header_condition` - (Optional) Set of match conditions used to match http response header:
  * `name` - (Required) The name of HTTP header to match.
  * `value` - (Required) The value of HTTP header to match.
  * `match_type` - (Required) Defines how value field is used to match the header value of HTTP response. Accepted values are STARTS_WITH, ENDS_WITH, CONTAINS, EQUALS, REGEX. Header name field does not support match types.
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

* `uri_condition` - (Optional) Set of match conditions used to match http request URI:
  * `uri` - (Required) The value of URI to match.
  * `match_type` - (Required) Defines how value field is used to match the URI. Accepted values are STARTS_WITH, ENDS_WITH, CONTAINS, EQUALS, REGEX.
  * `case_sensitive` - (Optional) If true, case is significant in the match. Default is true.
  * `inverse` - (Optional) A flag to indicate whether reverse the match result of this condition. Default is false.

* `uri_arguments_condition` - (Optional) Set of match conditions used to match http request URI arguments (query string):
  * `uri_arguments` - (Required) Query string of URI, typically contains key value pairs.
  * `match_type` - (Required) Defines how value field is used to match the URI. Accepted values are STARTS_WITH, ENDS_WITH, CONTAINS, EQUALS, REGEX.
  * `case_sensitive` - (Optional) If true, case is significant in the match. Default is true.
  * `inverse` - (Optional) A flag to indicate whether reverse the match result of this condition. Default is false.

* `ip_condition` - (Optional) Set of match conditions used to match IP header values of HTTP message:
  * `source_address` - (Required) The value source IP address to match.
  * `inverse` - (Optional) A flag to indicate whether reverse the match result of this condition. Default is false.

* `header_rewrite_action` - (Required) Set of header rewrite actions to be executed on the outgoing response when load balancer rule matches:
  * `name` - (Required) The name of HTTP header to be rewritten.
  * `value` - (Required) The new value of HTTP header.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the lb rule.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing lb rule can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_lb_http_response_rewrite_rule.lb_rule UUID
```

The above would import the lb rule named `lb_rule` with the nsx id `UUID`
