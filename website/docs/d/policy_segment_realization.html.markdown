---
subcategory: "Realization"
layout: "nsxt"
page_title: "NSXT: policy_segment_realization"
description: State of segment realization on hypervisors.
---

# nsxt_policy_segment_realization

This data source provides information about the realization of a policy segment or policy vlan segment on hypervisor.
This data source will wait until realization is complete with either success, partial success or error. It is recommended
to use this data source in conjunction with vsphere provider, in order to ensure segment is realized on hypervisor before
VM is created on same network.

This data source is applicable to NSX Policy Manager.

## Example Usage

```hcl
resource "nsxt_policy_segment" "s1" {
  display_name        = "segment1"
  transport_zone_path = data.nsxt_policy_transport_zone.tz1.path
}

data "nsxt_policy_segment_realization" "s1" {
  path = nsxt_policy_segment.s1.path
}

# usage in vsphere provider
data "vsphere_network" "net" {
  name          = nsxt_policy_segment_realization.s1.network_name
  datacenter_id = data.vsphere_datacenter.datacenter.id
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

resource "nsxt_policy_segment" "s1" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "segment1"
}

data "nsxt_policy_segment_realization" "s1" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  path = nsxt_policy_segment.s1.path
}

# usage in vsphere provider
data "vsphere_network" "net" {
  name          = nsxt_policy_segment_realization.s1.network_name
  datacenter_id = data.vsphere_datacenter.datacenter.id
}
```

## Argument Reference

* `path` - (Required) The policy path of the segment.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `state` - The realization state of the resource: `success`, `partial_success`, `orphaned`, `failed` or `error`.
* `network_name` - Network name on the hypervisor. This attribute can be used in vsphere provider in order to ensure implicit dependency on segment realization.
