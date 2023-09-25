---
subcategory: "Beta"
layout: "nsxt"
page_title: "NSXT: nsxt_manager_cluster"
description: A resource to configure an NSXT cluster.
---

# nsxt_manager_cluster

This resource provides a method for creating an NSXT cluster with several nodes
This resource is supported with NSX 4.1.0 onwards.
The main node for the cluster is the host in terraform nsxt provider config,
user will need to specify the nodes that will join the cluster in the resource config.
Only one instance of nsxt_manager_cluster resource is supported.

## Example Usage

```hcl
resource "nsxt_manager_cluster" "test" {
  node {
    ip_address = "192.168.240.32"
    username   = "admin"
    password   = "testpassword2"
  }
  node {
    ip_address = "192.168.240.33"
    username   = "admin"
    password   = "testpassword3"
  }
}
```

## Argument Reference

The following arguments are supported:

* `node` - (Required) IP Address of the node that will join the cluster of the host node.
  * `id` - (Computed) Uuid of the cluster node.
  * `ip_address` - (Required) Ip address of the node.
  * `username` - (Required) The username for login to the node.
  * `password` - (Required) The password for login to the node.
  * `fqdn`  - (Computed) Fqdn of the node.
  * `status` - (Computed) Status of the node, value will be one of `JOINING`, `JOINED`, `REMOVING` and `REMOVED`.

## Importing

Importing is not supported for this resource.
