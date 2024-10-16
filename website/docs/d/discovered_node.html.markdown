---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: discovered_node"
description: An Discovered Node data source.
---

# nsxt_discovered_node

This data source provides information about Discovered Node configured in NSX. A Discovered Node can be used to create a Host Transport Node.

## Example Usage

```hcl
data "nsxt_discovered_node" "test" {
  ip_address = "10.43.251.142"
}
```

## Example Usage when deploying together with Compute Manager

```hcl

resource "nsxt_compute_manager" "vc1" {
  display_name = "test-vcenter"

  server    = "34.34.34.34"
  multi_nsx = false

  credential {
    username_password_login {
      username   = "user1"
      password   = "password1"
      thumbprint = "thumbprint1"
    }
  }
  origin_type = "vCenter"

}

data "nsxt_compute_manager_realization" "vc1" {
  id      = nsxt_compute_manager.vc1.id
  timeout = 1200
}

data "nsxt_discovered_node" "test" {
  compute_manager_state = data.nsxt_compute_manager_realization.vc1.state
  ip_address            = "10.43.251.142"
}
```

## Argument Reference

* `id` - (Optional) External id of the discovered node, ex. a mo-ref from VC.
* `ip_address` - (Optional) IP Address of the discovered node.
* `compute_manager_state` - (Optional) Realized state of compute manager. This argument is only needed to ensure dependency upon `nsxt_compute_manager_realization` data source, so that `nsxt_discovered_node` is not evaluated before compute manager is fully functional.

## Attributes Reference

* `id` - External id of the discovered node, ex. a mo-ref from VC.
* `ip_address` - IP Address of the discovered node.
