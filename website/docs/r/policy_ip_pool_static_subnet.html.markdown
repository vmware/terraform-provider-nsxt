---
subcategory: "IPAM"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_ip_pool_static_subnet"
description: A resource to configure IP Pool Static Subnets in NSX Policy.
---

# nsxt_policy_ip_pool_static_subnet

This resource provides a means to configure IP Pool Static Subnets in NSX Policy.

This resource is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_ip_pool_static_subnet" "static_subnet1" {
  display_name = "static-subnet1"
  pool_path    = nsxt_policy_ip_pool.pool1.path
  cidr         = "12.12.12.0/24"
  gateway      = "12.12.12.1"

  allocation_range {
    start = "12.12.12.10"
    end   = "12.12.12.20"
  }
  allocation_range {
    start = "12.12.12.100"
    end   = "12.12.12.120"
  }

  tag {
    scope = "color"
    tag   = "blue"
  }

  tag {
    scope = "env"
    tag   = "test"
  }
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_ip_pool_static_subnet" "static_subnet1" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "static-subnet1"
  pool_path    = nsxt_policy_ip_pool.pool1.path
  cidr         = "12.12.12.0/24"
  gateway      = "12.12.12.1"

  allocation_range {
    start = "12.12.12.10"
    end   = "12.12.12.20"
  }
  allocation_range {
    start = "12.12.12.100"
    end   = "12.12.12.120"
  }

  tag {
    scope = "color"
    tag   = "blue"
  }

  tag {
    scope = "env"
    tag   = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) The display name for the Static Subnet.
* `context` - (Optional) The context which the object belongs to
  * `project_id` - (Required) The ID of the project which the object belongs to
* `pool_path` - (Required) The Policy path to the IP Pool for this Static Subnet.
* `cidr` - (Required) The network CIDR
* `allocation_range` - (Required) One or more IP allocation ranges for the Subnet.
  * `start` - (Required) The start IP address for the allocation range.
  * `end` - (Required) The end IP address for the allocation range.
* `description` - (Optional) Description of the resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this Static Subnet.
* `dns_nameservers` - (Optional) A list of up to 3 DNS nameservers for the Subnet.
* `dns_suffix` - (Optional) The DNS suffix for the Subnet.
* `gateway` - (Optional) The gateway IP for the Subnet.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of this Static Subnet.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the resource.

## Importing

An existing Subnet can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_ip_pool_static_subnet.static_subnet1 pool-id/subnet-id
```

The above would import NSX Static Subnet as a resource named `static_subnet1` with the NSX ID `subnet-id` in the IP Pool `pool-id`, where `subnet-id` is ID of Static Subnet and `pool-id` is the IP Pool ID the Subnet is in.
