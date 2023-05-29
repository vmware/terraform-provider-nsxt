---
subcategory: "Deprecated"
layout: "nsxt"
page_title: "NSXT: transport_zone"
description: A transport zone data source.
---

# nsxt_transport_zone

This data source provides information about Transport Zones (TZ) configured in NSX. A Transport Zone defines the scope to which a network can extend in NSX. For example an overlay based Transport Zone is associated with both hypervisors and logical switches and defines which hypervisors will be able to serve the defined logical switch. Virtual machines on the hypervisor associated with a Transport Zone can be attached to logical switches in that same Transport Zone.

## Example Usage

```hcl
data "nsxt_transport_zone" "overlay_transport_zone" {
  display_name = "1-transportzone-87"
}
```

## Argument Reference

* `id` - (Optional) The ID of Transport Zone to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Transport Zone to retrieve.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the Transport Zone.
* `host_switch_name` - The name of the N-VDS (host switch) on all Transport Nodes in this Transport Zone that will be used to run NSX network traffic.
* `transport_type` - The transport type of this transport zone (OVERLAY or VLAN).
