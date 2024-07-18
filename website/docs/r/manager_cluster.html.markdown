---
subcategory: "Fabric"
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
If `api_probing` is enabled, this resource will wait for NSX API endpoints to come up
before performing cluster joining.

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

* `node` - (Required) Specification of the node that will join the cluster of the host node.
  * `ip_address` - (Required) Ip address of the node.
  * `username` - (Required) The username for login to the node.
  * `password` - (Required) The password for login to the node.
* `api_probing` - (Optional) Parameters for probing NSX API endpoint connection. Since NSX nodes might have been created during same apply, we might need to wait until the API endpoint becomes available and all required default objects are created.
  * `enabled` - (Optional) Whether API connectivity check is enabled. Default is `true`.
  * `delay` - (Optional) Initial delay before we start probing API endpoint in seconds. Default is 0.
  * `interval` - (Optional) Interval for probing API endpoint in seconds. Default is 10.
  * `timeout` - (Optional) Timeout for probing the API endpoint in seconds. Default is 1800.

## Argument Reference

In addition to arguments listed above, the following attributes are exported:

* `node`
  * `id` - Uuid of the cluster node.
  * `fqdn`  - FQDN of the node.
  * `status` - Status of the node, value will be one of `JOINING`, `JOINED`, `REMOVING` and `REMOVED`.

## Importing

Importing is not supported for this resource.
