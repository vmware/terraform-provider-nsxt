---
subcategory: "Segments"
page_title: "NSXT: policy_segment"
description: Policy Segment data source.
---

# nsxt_policy_segment

This data source provides information about policy Segment configured on NSX.
This data source is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
data "nsxt_policy_segment" "test" {
  display_name = "segment1"
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_segment" "demoseg" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "demoseg"
}
```

## Example Usage - Global infra

```hcl
data "nsxt_policy_segment" "test_global" {
  context {
    from_global = true
  }
  display_name = "test"
}
```

## Argument Reference

* `id` - (Optional) The ID of Segment to retrieve. If ID is specified, no additional argument should be configured.
* `display_name` - (Optional) The Display Name prefix of the Segment to retrieve.
* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Optional) The ID of the project which the object belongs to
    * `from_global` - (Optional) Set to True if the data source will need to search Tier-1 gateway created in a global manager instance (/global-infra)

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
* `transport_zone_path` - The policy path of the transport zone the segment is attached to.
* `connectivity_path` - The policy path of the connected Tier-0 or Tier-1 gateway.
* `vlan_ids` - List of VLAN IDs configured on a VLAN-backed segment.
* `subnet` - List of subnet configurations.
    * `cidr` - Gateway IP address in CIDR format (e.g. `10.0.0.1/24`).
    * `network` - Network CIDR derived from the gateway address and prefix length.
