---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: nsxt_algorithm_type_ns_service"
description: A resource that can be used to configure a networking and security service on NSX.
---

# nsxt_algorithm_type_ns_service

This resource provides a way to configure a networking and security service which can be used with the NSX firewall. A networking and security service is an object that contains the TCP/UDP algorithm, source ports and destination ports in a single entity.

## Example Usage

```hcl
resource "nsxt_algorithm_type_ns_service" "ns_service_alg" {
  description      = "S1 provisioned by Terraform"
  display_name     = "S1"
  algorithm        = "FTP"
  destination_port = "21"
  source_ports     = ["9001-9003"]

  tag {
    scope = "color"
    tag   = "blue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional) Display name, defaults to ID if not set.
* `description` - (Optional) Description.
* `destination_port` - (Required) a single destination port.
* `source_ports` - (Optional) Set of source ports/ranges.
* `algorithm` - (Required) Algorithm one of "ORACLE_TNS", "FTP", "SUN_RPC_TCP", "SUN_RPC_UDP", "MS_RPC_TCP", "MS_RPC_UDP", "NBNS_BROADCAST", "NBDG_BROADCAST", "TFTP"
* `tag` - (Optional) A list of scope + tag pairs to associate with this service.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the NS service.
* `default_service` - The default NSServices are created in the system by default. These NSServices can't be modified/deleted.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.

## Importing

An existing Algorithm type NS service can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_algorithm_type_ns_service.ns_service_alg UUID
```

The above command imports the algorithm based networking and security service named `ns_service_alg` with the NSX id `UUID`.
