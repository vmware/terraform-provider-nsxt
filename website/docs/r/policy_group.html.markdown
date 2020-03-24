---
layout: "nsxt"
page_title: "NSXT: nsxt_policy_group"
sidebar_current: "docs-nsxt-resource-policy-group"
description: A resource to configure a Group and its members.
---

# nsxt_policy_group

This resource provides a method for the management of an inventory Group and its members. Groups as often used as sources and destinations, as well as in the Applied To field, in firewall rules.

## Example Usage

```hcl
resource "nsxt_policy_group" "group1" {
    display_name = "tf-group1"
    description  = "Terraform provisioned Group"

    criteria {
        condition {
            key         = "Name"
            member_type = "VirtualMachine"
            operator    = "STARTSWITH"
            value       = "public"
        }
        condition {
            key         = "OSName"
            member_type = "VirtualMachine"
            operator    = "CONTAINS"
            value       = "Ubuntu"
        }
    }

    conjunction {
        operator = "OR"
    }

    criteria {
        ipaddress_expression {
            ip_addresses = ["211.1.1.1", "212.1.1.1", "192.168.1.1-192.168.1.100"]
        }
    }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `domain` - (Optional) The domain to use for the Group. This domain must already exist. For VMware Cloud on AWS use `cgw`.
* `tag` - (Optional) A list of scope + tag pairs to associate with this Group.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the group resource.
* `criteria` - (Optional) A repeatable block to specify criteria for members of this Group. If more than 1 criteria block is specified, it must be separated by a `conjunction`. In a `criteria` block the following membership selection expressions can be used:
  * `ipaddress_expression` - (Optional) An expression block to specify individual IP Addresses, ranges of IP Addresses or subnets for this Group.
    * `ip_addresses` - (Required for a `ipaddress_expression`) This list can consist of a single IP address, IP address range or a subnet. Its type can be of either IPv4 or IPv6. Both IPv4 and IPv6 addresses within one expression is not allowed.
  * `path_expression` - (Optional) An expression block to specify direct group members by policy path.
    * `member_paths` - (Required for a `path_expression`) List of policy paths for direct members for this Group (such as Segments, Segment ports, Groups etc).
  * `condition` (Optional) A repeatable condition block to select this Group's members. When multiple `condition` blocks are used in a single `criteria` they form a nested expression that's implicitly ANDed together and each nested condition must used the same `member_type`.
    * `key` (Required for a `condition`) Specifies the attribute to query. Must be one of: `Tag`, `ComputerName`, `OSName` or `Name`. For a `member_type` other than `VirtualMachine`, only the `Tag` key is supported.
    * `member_type` (Required for a `condition`) Specifies the type of resource to query. Must be one of: `IPSet`, `LogicalPort`, `LogicalSwitch`, `Segment`, `SegmentPort` or `VirtualMachine`.
    * `operator` (Required for a `condition`) Specifies the query operator to use. Must be one of: `CONTAINS`, `ENDSWITH`, `EQUALS`, `NOTEQUALS` or `STARTSWITH`.
    * `value` (Required for a `condition`) User specified string value to use in the query. For `Tag` criteria, use 'scope|value' notation if you wish to specify scope in criteria.
* `conjunction` (Required for multiple `criteria`) When specifying multiple `criteria`, a conjunction is used to specify if the criteria should selected using `AND` or `OR`.
  * `operator` (Required for `conjunction`) The operator to use. Must be one of `AND` or `OR`. If `AND` is used, then the `criteria` block before/after must be of the same type and if using `condition` then also must use the same `member_type`.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the Group.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing policy Group can be [imported][docs-import] into this resource, via the following command:

[docs-import]: /docs/import/index.html

```
terraform import nsxt_policy_group.group1 ID
```

The above command imports the policy Group named `group` with the NSX Policy ID `ID`.

If the Group to import isn't in the `default` domain, the domain name can be added to the `ID` before a slash.

For example to import a Group with `ID` in the `MyDomain` domain:

```
terraform import nsxt_policy_group.group1 MyDomain/ID
```
