---
subcategory: "Beta"
page_title: "NSXT: nsxt_policy_lb_cookie_persistence_profile"
description: A resource to configure a Load Balancer Cookie Persistence Profile.
---

# nsxt_policy_lb_cookie_persistence_profile

This resource provides a method for the management of a LB Cookie Persistence Profile.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_lb_cookie_persistence_profile" "test" {
  display_name       = "test"
  description        = "Terraform provisioned profile"
  persistence_shared = false

  cookie_name     = "NSXT-LB"
  cookie_mode     = "INSERT"
  cookie_fallback = true
  cookie_garble   = true
  cookie_domain   = "cookie.domain.org"
  cookie_path     = "cookie/path"
  cookie_httponly = true
  cookie_secure   = true
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `cookie_mode` - (Optional) One of `INSERT`, `PREFIX` or `REWRITE`. Default is `INSERT`.
* `cookie_name` - (Optional) Cookie name, default is `NSXLB`
* `cookie_domain` - (Optional) HTTP cookie domain. Only relevant for `INSERT` mode.
* `cookie_fallback` - (Optional) If true, once the cookie points to a server that is down (i.e. admin state DISABLED or healthcheck state is DOWN), then a new server is selected by default to handle that request. If fallback is false, it will cause the request to be rejected.
* `cookie_garble` - (Optional) If enabled, cookie value (server IP and port) would be encrypted.
* `cookie_http_only` - (Optional) Prevents a script running in the browser from accessing the cookie. Only relevant for `INSERT` mode.
* `cookie_path` - (Optional) HTTP cookie path. Only relevant for `INSERT` mode.
* `cookie_secure` - (Optional) If enabled, this cookie will only be sent over HTTPS. Only relevant for `INSERT` mode.
* `session_cookie_time` - (Optional) Session cookie time preferences
  * `max_idle` - (Optional) Maximum interval the cookie is valid for from the last time it was seen in a request
  * `max_life` - (Optional) Maximum interval the cookie is valid for from the first time it was seen in a request
* `persistence_cookie_time` - (Optional) Persistence cookie time preferences
  * `max_idle` - (Optional) Maximum interval the cookie is valid for from the last time it was seen in a request
* `persistence_shared` - (Optional) If enabled, all virtual servers with this profile will share the same persistence mechanism.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_lb_cookie_persistence_profile.test PATH
```

The above command imports LBCookiePersistenceProfile named `test` with the NSX path `PATH`.
