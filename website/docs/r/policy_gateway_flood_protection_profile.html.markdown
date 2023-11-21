---
subcategory: "Security"
layout: "nsxt"
page_title: "NSXT: policy_gateway_flood_protection_profile"
description: A resource to configure Policy Gateway Flood Protection Profile on NSX Policy manager.
---

# nsxt_policy_gateway_flood_protection_profile

This resource provides a method for the management of a Gateway Flood Protection Profile.

This resource is applicable to NSX Global Manager and NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_gateway_flood_protection_profile" "test" {
  display_name             = "test"
  description              = "test"
  icmp_active_flow_limit   = 3
  other_active_conn_limit  = 3
  tcp_half_open_conn_limit = 3
  udp_active_flow_limit    = 3
  nat_active_conn_limit    = 3

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_gateway_flood_protection_profile" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name             = "test"
  description              = "test"
  icmp_active_flow_limit   = 3
  other_active_conn_limit  = 3
  tcp_half_open_conn_limit = 3
  udp_active_flow_limit    = 3
  nat_active_conn_limit    = 3

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the policy resource.
* `icmp_active_flow_limit` - (Optional) Active ICMP connections limit. If this field is empty, firewall will not set a limit to active ICMP connections. Minimum: 1, Maximum: 1000000.
* `other_active_conn_limit` - (Optional) Timeout after first TN. If this field is empty, firewall will not set a limit to other active connections. besides UDP, ICMP and half open TCP connections. Minimum: 1, Maximum: 1000000.
* `tcp_half_open_conn_limit` - (Optional) Active half open TCP connections limit. If this field is empty, firewall will not set a limit to half open TCP connections. Minimum: 1, Maximum: 1000000.
* `udp_active_flow_limit` - (Optional) 	Active UDP connections limit. If this field is empty, firewall will not set a limit to active UDP connections. Minimum: 1, Maximum: 1000000.
* `nat_active_conn_limit` - (Optional) 	Maximum limit of active NAT connections. The maximum limit of active NAT connections. This limit only apply to EDGE components (such as, gateway). If this property is omitted, or set to null, then there is no limit on the specific component. Meanwhile there is an implicit limit which depends on the underlying hardware resource. Minimum: 1, Maximum: 4294967295, Default: 4294967295

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing Gateway Flood Protection Profile can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_gateway_flood_protection_profile.gfpp ID
```

The above command imports the Gateway Flood Protection Profile named `gfpp` with the NSX Policy ID `ID`.

```
terraform import nsxt_policy_gateway_flood_protection_profile.gfpp POLICY_PATH
```
The above command imports the Gateway Flood Protection Profile named `gfpp` with the policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.
