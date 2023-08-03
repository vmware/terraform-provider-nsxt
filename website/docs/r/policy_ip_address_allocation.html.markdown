---
subcategory: "IPAM"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_ip_address_allocation"
description: A resource to configure a IP Address allocations.
---

# nsxt_policy_ip_address_allocation

This resource provides a method for the management of a IP Address Allocations.
This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_ip_pool" "pool1" {
  context {
    project_id = data.nsxt_policy_project.tenant1.id
  }
  display_name = "ip_pool"
  description  = "Created by Terraform"
}

resource "nsxt_policy_ip_pool_static_subnet" "static_subnet1" {
  display_name = "static-subnet1"
  pool_path    = nsxt_policy_ip_pool.pool1.path
  cidr         = "12.12.12.0/24"
  gateway      = "12.12.12.1"

  allocation_range {
    start = "12.12.12.10"
    end   = "12.12.12.20"
  }
}

resource "nsxt_policy_ip_address_allocation" "test" {
  display_name  = "test"
  description   = "Terraform provisioned IpAddressAllocation"
  pool_path     = nsxt_policy_ip_pool.pool1.path
  allocation_ip = "12.12.12.12"
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_ip_pool" "pool1" {
  context {
    project_id = data.nsxt_policy_project.tenant1.id
  }
  display_name = "ip_pool"
  description  = "Created by Terraform"
}

resource "nsxt_policy_ip_pool_static_subnet" "static_subnet1" {
  context {
    project_id = data.nsxt_policy_project.tenant1.id
  }

  display_name = "static-subnet1"
  pool_path    = nsxt_policy_ip_pool.pool1.path
  cidr         = "12.12.12.0/24"
  gateway      = "12.12.12.1"

  allocation_range {
    start = "12.12.12.10"
    end   = "12.12.12.20"
  }
}
resource "nsxt_policy_ip_address_allocation" "test" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name  = "test"
  description   = "Terraform provisioned IpAddressAllocation"
  pool_path     = nsxt_policy_ip_pool.pool1.path
  allocation_ip = "12.12.12.12"
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to
* `allocation_ip` - (Optional) The IP Address to allocate. If unspecified any free IP in the pool will be allocated.
* `pool_path` - (Required) The policy path to the IP Pool for this Allocation.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the Allocation.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.
* `allocation_ip` - If the `allocation_ip` is not specified in the resource, any free IP is allocated and its value is exported on this attribute.

## Importing

An existing IP Allocation can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_ip_address_allocation.test POOL-ID/ID
```
The above command imports IpAddressAllocation named `test` with the NSX IpAddressAllocation ID `ID` in IP Pool `POOL-ID`.

```
terraform import nsxt_policy_ip_address_allocation.test POLICY_PATH
```
The above command imports IpAddressAllocation named `test` with the NSX policy path `POLICY_PATH`.
Note: for multitenancy projects only the later form is usable.
