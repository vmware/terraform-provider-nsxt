---
layout: "nsxt"
page_title: "NSXT: policy_segment_realization"
sidebar_current: "docs-nsxt-datasource-policy-segment-realization"
description: State of segment realization on hypervisors.
---

# nsxt_policy_segment_realization

This data source provides information about the realization of a policy segment or policy vlan segment on hypervisor.
This data source will wait until realization is complete with either success, partial success or error. It is recommended
to use this data source in conjunction with vsphere provider, in order to ensure segment is realized on hypervisor before
VM is created on same network.

## Example Usage

```hcl
resource "nsxt_policy_segment" "s1" {
  display_name        = "segment1"
  transport_zone_path = data.nsxt_transport_zone.tz1.path
}

data "nsxt_policy_segment_realization" "s1" {
  path = data.nsxt_policy_segment.s1.path
}

# usage in vsphere provider
data "vsphere_network" "net" {
  name          = nsxt_policy_segment_realization.s1.network_name
  datacenter_id = data.vsphere_datacenter.datacenter.id
}
```

## Argument Reference

* `path` - (Required) The policy path of the segment.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `state` - The realization state of the resource: `success`, `partial_success`, `orphaned`, `failed` or `error`.
* `network_name` - Network name on the hypervisor. This attribute can be used in vsphere provider in order to ensure implicit dependency on segment realization.
