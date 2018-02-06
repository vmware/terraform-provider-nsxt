---
layout: "nsxt"
page_title: "NSXT: nsxt_alg_type_ns_service"
sidebar_current: "docs-nsxt-resource-alg-type-ns-service"
description: |-
  Provides a resource to configure NS service for algorithm type on NSX-T Manager.
---

# nsxt_alg_type_ns_service

Provides a resource to configure NS service for algorithm type on NSX-T Manager

## Example Usage

```hcl
resource "nsxt_alg_type_ns_service" "S1" {
    description = "S1 provisioned by Terraform"
    display_name = "S1"
    alg = "FTP"
    destination_ports = "21"
    source_ports = [ "9001-9003"]
    tags = [{
        scope = "color"
        tag = "blue"}
    ]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Optional) Display name, defaults to ID if not set.
* `description` - (Optional) Description.
* `destination_ports` - (Required) a single destination port.
* `source_ports` - (Optional) Set of source ports/ranges.
* `alg` - (Required) Algorithm one of "ORACLE_TNS", "FTP", "SUN_RPC_TCP", "SUN_RPC_UDP", "MS_RPC_TCP", "MS_RPC_UDP", "NBNS_BROADCAST", "NBDG_BROADCAST", "TFTP"
* `tags` - (Optional) A list of scope + tag pairs to associate with this ip_set.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the logical switch.
* `default_service` - The default NSServices are created in the system by default. These NSServices can't be modified/deleted.
* `system_owned` - A boolean that indicates whether this resource is system-owned and thus read-only.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
