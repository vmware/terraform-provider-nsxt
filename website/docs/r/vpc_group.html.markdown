---
subcategory: "VPC"
layout: "nsxt"
page_title: "NSXT: nsxt_vpc_group"
description: A resource to configure a VPC Group and its members.
---

# nsxt_vpc_group

This resource provides a method for the management of an inventory Group and its members in a VPC scope. Groups as often
used as sources and destinations, as well as in the Applied To field, in firewall rules.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_vpc" "demovpc" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "vpc1"
}

resource "nsxt_vpc_group" "group1" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
    vpc_id     = data.nsxt_vpc.demovpc.id
  }

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
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this Group.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the group resource.
* `context` - (Required) The context which the object belongs to
  * `project_id` - (Required) The ID of the project which the object belongs to
  * `vpc_id` - (Required) The ID of the VPC which the object belongs to
* `criteria` - (Optional) A repeatable block to specify criteria for members of this Group. If more than 1 criteria block is specified, it must be separated by a `conjunction`. In a `criteria` block the following membership selection expressions can be used:
  * `ipaddress_expression` - (Optional) An expression block to specify individual IP Addresses, ranges of IP Addresses or subnets for this Group.
    * `ip_addresses` - (Required) This list can consist of a single IP address, IP address range or a subnet. Its type can be of either IPv4 or IPv6. Both IPv4 and IPv6 addresses within one expression is not allowed.
  * `macaddress_expression` - (Optional) An expression block to specify individual MAC Addresses for this Group.
    * `mac_addresses` - (Required) List of MAC addresses.
  * `path_expression` - (Optional) An expression block to specify direct group members by policy path.
    * `member_paths` - (Required) List of policy paths for direct members for this Group (such as Subnets, Subnet ports, Groups etc).
  * `external_id_expression` - (Optional) An expression block to specify external IDs for the specified member type for this Group.
    * `member_type` - (Optional) External ID member type. Must be one of: `VirtualMachine`, `VirtualNetworkInterface`, `CloudNativeServiceInstance`, or `PhysicalServer`. Defaults to `VirtualMachine`.
    * `external_ids` - (Required) List of external IDs for the specified member type.
  * `condition` (Optional) A repeatable condition block to select this Group's members. When multiple `condition` blocks are used in a single `criteria` they form a nested expression that's implicitly ANDed together and each nested condition must used the same `member_type`.
    * `key` (Required) Specifies the attribute to query. Must be one of: `Tag`, `ComputerName`, `OSName`, `Name`. Please note that certain keys are only applicable to certain member types.
    * `member_type` (Required) Specifies the type of resource to query. Must be one of: `Subnet`, `SubnetPort`, `VirtualMachine`
    * `operator` (Required) Specifies the query operator to use. Must be one of: `CONTAINS`, `ENDSWITH`, `EQUALS`, `NOTEQUALS`, `STARTSWITH`, `IN`, `NOTIN`, `MATCHES`. Not that certain operators are only applicable to certain keys/member types.:w
    * `value` (Required) User specified string value to use in the query. For `Tag` criteria, use 'scope|value' notation if you wish to specify scope in criteria.
* `conjunction` (Required for multiple `criteria`) When specifying multiple `criteria`, a conjunction is used to specify if the criteria should selected using `AND` or `OR`.
  * `operator` (Required) The operator to use. Must be one of `AND` or `OR`. If `AND` is used, then the `criteria` block before/after must be of the same type and if using `condition` then also must use the same `member_type`.
* `extended_criteria` (Optional) A condition block to specify higher level context to include in this Group's members. (e.g. user AD group). This configuration is for Local Manager only. Currently only one block is supported by NSX. Note that `extended_criteria` is implicitly `AND` with `criteria`.
  * `identity_group` (Optional) A repeatable condition block selecting user AD groups to be included in this Group. Note that `identity_groups` are `OR` with each other.
    * `distinguished_name` (Required) LDAP distinguished name (DN). A valid fully qualified distinguished name should be provided here. This value is valid only if it matches to exactly 1 LDAP object on the LDAP server.
    * `domain_base_distinguished_name` (Required) Identity (Directory) domain base distinguished name. This is the base distinguished name for the domain where this identity group resides. (e.g. dc=example,dc=com)
    * `sid` (Optional) Identity (Directory) Group SID (security identifier). A security identifier (SID) is a unique value of variable length used to identify a trustee. This field is only populated for Microsoft Active Directory identity store.
* `group_type` - (Optional) One of `IPAddress`, `ANTREA`. Empty group type indicates a generic group. Attribute is supported with NSX version 3.2.0 and above.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the Group.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing VPC Group can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_vpc_group.group1 PATH
```

The above command imports the VPC Group named `group` with the NSX Policy path `PATH`.
