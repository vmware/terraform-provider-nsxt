---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_lb_cookie_persistence_profile"
description: |-
  Provides a resource to configure lb cookie persistence profile on NSX-T manager
---

# nsxt_lb_cookie_persistence_profile

Provides a resource to configure lb cookie persistence profile on NSX-T manager

~> **NOTE:** This resource requires NSX version 2.3 or higher.

## Example Usage

```hcl
resource "nsxt_lb_cookie_persistence_profile" "lb_cookie_persistence_profile" {
  description        = "lb_cookie_persistence_profile provisioned by Terraform"
  display_name       = "lb_cookie_persistence_profile"
  cookie_name        = "my_cookie"
  persistence_shared = "false"
  cookie_fallback    = "false"
  cookie_garble      = "false"
  cookie_mode        = "INSERT"

  insert_mode_params {
    cookie_domain      = ".example2.com"
    cookie_path        = "/subfolder"
    cookie_expiry_type = "SESSION_COOKIE_TIME"
    max_idle_time      = "1000"
    max_life_time      = "2000"
  }

  tag {
    scope = "color"
    tag   = "red"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional) The display name of this resource. Defaults to ID if not set.
* `description` - (Optional) Description of this resource.
* `cookie_mode` - (Optional) The cookie persistence mode. Accepted values: PREFIX, REWRITE and INSERT which is the default.
* `cookie_name` - (Required) cookie name.
* `persistence_shared` - (Optional) A boolean flag which reflects whether the cookie persistence is private or shared. When false (which is the default value), the cookie persistence is private to each virtual server and is qualified by the pool. If set to true, in cookie insert mode, cookie persistence could be shared across multiple virtual servers that are bound to the same pools.
* `cookie_fallback` - (Optional) A boolean flag which reflects whether once the server points by this cookie is down, a new server is selected, or the requests will be rejected.
* `cookie_garble` - (Optional) A boolean flag which reflects whether the cookie value (server IP and port) would be encrypted or in plain text.
* `insert_mode_params` - (Optional) Additional parameters for the INSERT cookie mode:
  * `cookie_domain` - (Optional) HTTP cookie domain (for INSERT mode only).
  * `cookie_path` - (Optional) HTTP cookie path (for INSERT mode only).
  * `cookie_expiry_type` - (Optional) Type of cookie expiration timing (for INSERT mode only). Accepted values: SESSION_COOKIE_TIME for session cookie time setting and PERSISTENCE_COOKIE_TIME for persistence cookie time setting.
  * `max_idle_time` - (Required if cookie_expiry_type is set) Maximum interval the cookie is valid for from the last time it was seen in a request.
  * `max_life_time` - (Required for INSERT mode with SESSION_COOKIE_TIME expiration) Maximum interval the cookie is valid for from the first time the cookie was seen in a request.
* `tag` - (Optional) A list of scope + tag pairs to associate with this lb cookie persistence profile.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the lb cookie persistence profile.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.


## Importing

An existing lb cookie persistence profile can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_lb_cookie_persistence_profile.lb_cookie_persistence_profile UUID
```

The above would import the lb cookie persistence profile named `lb_cookie_persistence_profile` with the nsx id `UUID`
