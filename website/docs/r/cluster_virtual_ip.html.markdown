---
subcategory: "Fabric"
layout: "nsxt"
page_title: "NSXT: nsxt_cluster_virtual_ip"
description: A resource to configure virtual IP of NSXT cluster.
---

# nsxt_cluster_virtual_ip

This resource provides a method for configuring the virtual IP of NSXT cluster.
This resource is supported with NSX 4.1.0 onwards.
Only one instance of nsxt_cluster_virtual_ip resource is supported.

## Example Usage

```hcl
resource "nsxt_cluster_virtual_ip" "test" {
  ip_address   = "10.0.0.251"
  ipv6_address = "fd01:1:2:2918:250:56ff:fe8b:7e4d"
  force        = "true"
}
```

## Argument Reference

The following arguments are supported:

* `force` - (Optional) A flag to determine if need to ignore duplicate address detection and DNS lookup validation check. Value can be `true` or `false`. Default value is `false`. This argument is supported for NSX 4.0.0 and above.
* `ip_address` - (Optional) Virtual IP Address of the cluster. Must be in the same subnet as the manager nodes. Default value is `0.0.0.0`.
* `ipv6_address` - (Optional) Virtual IPv6 Address of the cluster. To set ipv6 virtual IP address, IPv6 interface needs to be configured on manager nodes. Default value is `::`. This argument is supported for NSX 4.0.0 and above.

## Importing

Importing is not supported for this resource.
