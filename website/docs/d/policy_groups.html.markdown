---
subcategory: "Firewall"
layout: "nsxt"
page_title: "NSXT: policy_groups"
description: A policy groups data source. This data source builds a list representation of the whole groups table, containing `id`, `display_name` and `path` attributes.
---

# nsxt_policy_groups

This data source builds a list of the whole policy Groups table. Such list can be referenced in configuration to obtain object identifier attributes lookup by display name at a cost of single roundtrip to NSX, which improves apply and refresh
time at scale, compared to multiple instances of `nsxt_policy_group` data source.

## Example Usage

```hcl
data "nsxt_policy_groups" "grouplist" {
}

resource "nsxt_policy_predefined_security_policy" "test" {
  path = data.nsxt_policy_security_policy.default_l3.path

  tag {
    scope = "color"
    tag   = "orange"
  }

  rule {
    display_name = "allow_icmp"
    destination_groups = [
      data.nsxt_policy_groups.test.items[index(data.nsxt_policy_groups.test.items.*.display_name, "Cats")].path,
      data.nsxt_policy_groups.test.items[index(data.nsxt_policy_groups.test.items.*.display_name, "Dogs")].path
    ]
    action   = "ALLOW"
    services = [nsxt_policy_service.icmp.path]
    logged   = true
  }

  rule {
    display_name     = "allow_udp"
    source_groups    = [data.nsxt_policy_groups.grouplist.items[index(data.nsxt_policy_groups.grouplist.items.*.display_name, "Fish")].path]
    sources_excluded = true
    scope            = [data.nsxt_policy_groups.grouplist.items[index(data.nsxt_policy_groups.grouplist.items.*.display_name, "Aquarium")].path]
    action           = "ALLOW"
    services         = [nsxt_policy_service.udp.path]
    logged           = true
    disabled         = true
  }

  default_rule {
    action = "DROP"
  }

}
```

## Argument Reference

* `domain` - (Optional) The domain this Group belongs to. For VMware Cloud on AWS use `cgw`. For Global Manager, please use site id for this field. If not specified, this field is default to `default`.
* `context` - (Optional) The context which the object belongs to
  * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `items` - List of policy group entries.
  * `id` - The ID of the group.
  * `display_name` - Display name of the group entry.
  * `path` - The NSX path of the policy resource.
