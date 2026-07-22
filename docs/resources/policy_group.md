---
subcategory: "Grouping and Tagging"
page_title: "NSXT: nsxt_policy_group"
description: A resource to configure a Group and its members.
---

# nsxt_policy_group

This resource provides a method for the management of an inventory Group and its members. Groups as often used as sources and destinations, as well as in the Applied To field, in firewall rules.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

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

  conjunction {
    operator = "OR"
  }

  criteria {
    macaddress_expression {
      mac_addresses = ["b2:54:00:98:b0:83"]
    }
  }

  conjunction {
    operator = "OR"
  }

  criteria {
    external_id_expression {
      member_type  = "VirtualMachine"
      external_ids = ["520ba7b0-d9f8-87b1-6f44-15bbeb7935c7", "52748a9e-d61d-e29b-d54b-07f169ff0ee8-4000"]
    }
  }

  extended_criteria {
    identity_group {
      distinguished_name             = "cn=u1,ou=users,dc=example,dc=local"
      domain_base_distinguished_name = "ou=users,dc=example,dc=local"
    }
    identity_group {
      distinguished_name             = "cn=a1,ou=admin,dc=example,dc=local"
      domain_base_distinguished_name = "ou=admin,dc=example,dc=local"
    }
  }
}
```

## Example Usage - Bare Metal Server Group

Note: Bare Metal Server group support requires NSX version 9.0.0 or higher.

```hcl
# Static BMS group — BareMetalServer, BareMetalServerInterface, and Group members
resource "nsxt_policy_group" "bms_static" {
  display_name = "tf-bms-static"
  description  = "Bare Metal Server group with static membership"
  group_type   = "BareMetalServer"

  # Static BareMetalServer members by external ID
  criteria {
    external_id_expression {
      member_type  = "BareMetalServer"
      external_ids = ["bms-external-id-1", "bms-external-id-2"]
    }
  }

  conjunction {
    operator = "OR"
  }

  # Static BareMetalServerInterface members by external ID
  criteria {
    external_id_expression {
      member_type  = "BareMetalServerInterface"
      external_ids = ["bms-intf-external-id-1"]
    }
  }

  conjunction {
    operator = "OR"
  }

  # Include another Group as a static member by policy path
  criteria {
    path_expression {
      member_paths = [nsxt_policy_group.another_bms_group.path]
    }
  }
}

# Dynamic BMS group using conditions  
resource "nsxt_policy_group" "bms_dynamic" {
  display_name = "tf-bms-dynamic"
  description  = "Bare Metal Server group with dynamic membership"
  group_type   = "BareMetalServer"

  criteria {
    condition {
      key         = "Tag"
      member_type = "BareMetalServer"
      operator    = "EQUALS"
      value       = "env|production"
    }
  }

  conjunction {
    operator = "OR"
  }

  criteria {
    condition {
      key         = "OSName"
      member_type = "BareMetalServer"
      operator    = "CONTAINS"
      value       = "Linux"
    }
  }
}

# BMS group selecting management interfaces via ManagementInterface key
resource "nsxt_policy_group" "bms_mgmt_interfaces" {
  display_name = "tf-bms-management-interfaces"
  description  = "All BMS management interfaces"
  group_type   = "BareMetalServer"

  criteria {
    condition {
      key         = "ManagementInterface"
      member_type = "BareMetalServerInterface"
      operator    = "EQUALS"
      value       = "TRUE"
    }
  }
}
```

## Global Manager example

Note: This usage is for Global Manager only using site

```hcl
data "nsxt_policy_site" "paris" {
  display_name = "Paris"
}
resource "nsxt_policy_group" "group1" {
  display_name = "tf-group1"
  description  = "Terraform provisioned Group"
  domain       = data.nsxt_policy_site.paris.id

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

Note: This usage is for Global Manager only using domain

```hcl
resource "nsxt_policy_domain" "france" {
  display_name = "France"
  sites        = ["Paris"]
}

resource "nsxt_policy_group" "group1" {
  display_name = "tf-group1"
  description  = "Terraform provisioned Group"
  domain       = nsxt_policy_domain.france.id

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

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_group" "group1" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
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
* `domain` - (Optional) The domain to use for the Group. This domain must already exist. For VMware Cloud on AWS use `cgw`. For Global Manager, please use site id for this field. If not specified, this field is default to `default`.
* `tag` - (Optional) A list of scope + tag pairs to associate with this Group.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the group resource.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to
* `criteria` - (Optional) A repeatable block to specify criteria for members of this Group. If more than 1 criteria block is specified, it must be separated by a `conjunction`. In a `criteria` block the following membership selection expressions can be used:
    * `ipaddress_expression` - (Optional) An expression block to specify individual IP Addresses, ranges of IP Addresses or subnets for this Group.
        * `ip_addresses` - (Required) This list can consist of a single IP address, IP address range or a subnet. Its type can be of either IPv4 or IPv6. Both IPv4 and IPv6 addresses within one expression is not allowed.
    * `macaddress_expression` - (Optional) An expression block to specify individual MAC Addresses for this Group.
        * `mac_addresses` - (Required) List of MAC addresses.
    * `path_expression` - (Optional) An expression block to specify direct group members by policy path.
        * `member_paths` - (Required) List of policy paths for direct members for this Group (such as Segments, Segment ports, Groups etc).
    * `external_id_expression` - (Optional) An expression block to specify external IDs for the specified member type for this Group.
        * `member_type` - (Optional) External ID member type. Must be one of: `VirtualMachine`, `VirtualNetworkInterface`, `CloudNativeServiceInstance`, `PhysicalServer`, `BareMetalServer`, or `BareMetalServerInterface`. Defaults to `VirtualMachine`.
        * `external_ids` - (Required) List of external IDs for the specified member type.
    * `condition` (Optional) A repeatable condition block to select this Group's members. When multiple `condition` blocks are used in a single `criteria` they form a nested expression that's implicitly ANDed together and each nested condition must used the same `member_type`.
        * `key` (Required) Specifies the attribute to query. Must be one of: `Tag`, `ComputerName`, `OSName`, `Name`, `NodeType`, `GroupType`, `ALL`, `IPAddress`, `PodCidr`, `ManagementInterface`, `OSVersion`. Please note that certain keys are only applicable to certain member types. For BMS-specific supported keys and valid operators per member type, refer to the [Bare Metal Server Management guide](../guides/baremetal_server_management).
        * `member_type` (Required) Specifies the type of resource to query. Must be one of: `IPSet`, `LogicalPort`, `LogicalSwitch`, `Segment`, `SegmentPort`, `VirtualMachine`, `Group`, `DVPG`, `DVPort`, `IPAddress`, `TransportNode`, `Pod`, `Service`, `Namespace`, `KubernetesCluster`, `KubernetesNamespace`, `KubernetesIngress`, `KubernetesService`, `KubernetesNode`, `AntreaEgress`, `AntreaIPPool`, 
          `VpcSubnet`, `VpcSubnetPort`, `BareMetalServer`, `BareMetalServerInterface`. Note that certain member types are only applicable to certain environments.
        * `operator` (Required) Specifies the query operator to use. Must be one of: `CONTAINS`, `ENDSWITH`, `EQUALS`, `NOTEQUALS`, `STARTSWITH`, `IN`, `NOTIN`, `MATCHES`, `TRUE`, `FALSE`. Note that certain operators are only applicable to certain keys/member types.
        * `value` (Required) User specified string value to use in the query. For `Tag` criteria, use 'scope|value' notation if you wish to specify scope in criteria.
* `conjunction` (Required for multiple `criteria`) When specifying multiple `criteria`, a conjunction is used to specify if the criteria should selected using `AND` or `OR`.
    * `operator` (Required) The operator to use. Must be one of `AND` or `OR`. If `AND` is used, then the `criteria` block before/after must be of the same type and if using `condition` then also must use the same `member_type`.
* `extended_criteria` (Optional) A condition block to specify higher level context to include in this Group's members. (e.g. user AD group). This configuration is for Local Manager only. Currently only one block is supported by NSX. Note that `extended_criteria` is implicitly `AND` with `criteria`.
    * `identity_group` (Optional) A repeatable condition block selecting user AD groups to be included in this Group. Note that `identity_groups` are `OR` with each other.
        * `distinguished_name` (Required) LDAP distinguished name (DN). A valid fully qualified distinguished name should be provided here. This value is valid only if it matches to exactly 1 LDAP object on the LDAP server.
        * `domain_base_distinguished_name` (Required) Identity (Directory) domain base distinguished name. This is the base distinguished name for the domain where this identity group resides. (e.g. dc=example,dc=com)
        * `sid` (Optional) Identity (Directory) Group SID (security identifier). A security identifier (SID) is a unique value of variable length used to identify a trustee. This field is only populated for Microsoft Active Directory identity store.
* `group_type` - (Optional) One of `IPAddress`, `ANTREA`, `BareMetalServer`. Empty group type indicates a generic group. Attribute is supported with NSX version 3.2.0 and above. Note that updating this attribute will trigger recreation of the group. `BareMetalServer` group type requires NSX version 9.0.0 or higher.
  * **NOTE**: `group_type = "BareMetalServer"` **must** be set whenever `BareMetalServer` or `BareMetalServerInterface` member types are used. A `BareMetalServer`-typed group can only contain `BareMetalServer` and/or `BareMetalServerInterface` members.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the Group.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing policy Group can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://developer.hashicorp.com/terraform/cli/import

```shell
terraform import nsxt_policy_group.group ID
```

The above command imports the policy Group named `group` with the NSX Policy ID `ID`.

If the Group to import isn't in the `default` domain, the domain name can be added to the `ID` before a slash.

For example to import a Group with `ID` in the `MyDomain` domain:

```shell
terraform import nsxt_policy_group.group MyDomain/ID
```

```shell
terraform import nsxt_policy_group.group POLICY_PATH
```

The above command imports the policy Group named `group` with the policy path `POLICY_PATH`.
