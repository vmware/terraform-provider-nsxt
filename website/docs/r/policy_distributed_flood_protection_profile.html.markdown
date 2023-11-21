---
subcategory: "Security"
layout: "nsxt"
page_title: "NSXT: policy_distributed_flood_protection_profile"
description: A resource to configure Policy Distributed Flood Protection Profile on NSX Policy manager.
---

# nsxt_policy_distributed_flood_protection_profile

This resource provides a method for the management of a Distributed Flood Protection Profile.

This resource is applicable to NSX Global Manager and NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_distributed_flood_protection_profile" "test" {
  display_name             = "test"
  description              = "test"
  icmp_active_flow_limit   = 3
  other_active_conn_limit  = 3
  tcp_half_open_conn_limit = 3
  udp_active_flow_limit    = 3
  enable_rst_spoofing      = true
  enable_syncache          = true

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

resource "nsxt_policy_distributed_flood_protection_profile" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name             = "test"
  description              = "test"
  icmp_active_flow_limit   = 3
  other_active_conn_limit  = 3
  tcp_half_open_conn_limit = 3
  udp_active_flow_limit    = 3
  enable_rst_spoofing      = true
  enable_syncache          = true

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
* `enable_rst_spoofing` - (Optional) Flag to indicate rst spoofing is enabled. If set to true, rst spoofing will be enabled. Flag is used only for distributed firewall profiles. Default: false.
* `enable_syncache` - (Optional) Flag to indicate syncache is enabled. If set to true, sync cache will be enabled. Flag is used only for distributed firewall profiles. Default: false.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing Distributed Flood Protection Profile can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_distributed_flood_protection_profile.dfpp ID
```

The above command imports the Distributed Flood Protection Profile named `dfpp` with the NSX Policy ID `ID`.

```
terraform import nsxt_policy_distributed_flood_protection_profile.dfpp POLICY_PATH
```
The above command imports the Distributed Flood Protection Profile named `dfpp` with the policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.
