---
subcategory: "Segments"
page_title: "NSXT: nsxt_policy_segment_ports"
description: A data source to retrieve multiple Segment Ports based on filter criteria.
---

# nsxt_policy_segment_ports

This data source provides a method for retrieving multiple Segment Ports based on filter criteria. Unlike the singular `nsxt_policy_segment_port` data source, this returns a list of segment ports matching the specified filters.

## Example Usage

### Get all segment ports for a specific segment

```hcl
data "nsxt_policy_segment" "segment1" {
  display_name = "segment1"
}

data "nsxt_policy_segment_ports" "all_ports" {
  segment_path = data.nsxt_policy_segment.segment1.path
}

output "port_count" {
  value = length(data.nsxt_policy_segment_ports.all_ports.items)
}

output "first_port_id" {
  value = data.nsxt_policy_segment_ports.all_ports.items[0].id
}
```

### Get segment ports by segment path and display name

```hcl
data "nsxt_policy_segment" "segment1" {
  display_name = "segment1"
}

data "nsxt_policy_segment_ports" "filtered_ports" {
  segment_path = data.nsxt_policy_segment.segment1.path
  display_name = "web-server-port"
}
```

### Get segment ports by VIF ID

```hcl
data "nsxt_policy_segment_ports" "by_vif" {
  vif_id = "50234567-abcd-1234-5678-123456789abc"
}
```

### Get segment ports by VIF ID and display name

```hcl
# This will first filter by VIF ID, then filter the results by display name
data "nsxt_policy_segment_ports" "by_vif_and_name" {
  vif_id       = "50234567-abcd-1234-5678-123456789abc"
  display_name = "web-server-port"
}
```

### Get segment ports by VIF ID and segment path

```hcl
data "nsxt_policy_segment" "segment1" {
  display_name = "segment1"
}

# This will filter by both VIF ID and segment path
data "nsxt_policy_segment_ports" "by_vif_and_segment" {
  vif_id       = "50234567-abcd-1234-5678-123456789abc"
  segment_path = data.nsxt_policy_segment.segment1.path
}
```

### Get segment ports by display name only

```hcl
# When only display_name is provided, behaves like the singular data source
data "nsxt_policy_segment_ports" "by_name" {
  display_name = "web-server-port"
}
```

## Example Usage - Multi-Tenancy

```hcl
data "nsxt_policy_project" "demoproj" {
  display_name = "demoproj"
}

data "nsxt_policy_segment" "segment1" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  display_name = "segment1"
}

data "nsxt_policy_segment_ports" "ports" {
  context {
    project_id = data.nsxt_policy_project.demoproj.id
  }
  segment_path = data.nsxt_policy_segment.segment1.path
}
```

## Argument Reference

The following arguments are supported:

* `context` - (Optional) The context which the object belongs to
    * `project_id` - (Required) The ID of the project which the object belongs to
* `vif_id` - (Optional) Segment Port attachment ID to filter segment ports. This takes priority over other filters.
* `segment_path` - (Optional) Segment path to filter segment ports. Can be combined with `vif_id` or `display_name`.
* `display_name` - (Optional) Display name to filter segment ports. When used with `vif_id`, filters the VIF results by display name. When used with `segment_path`, filters the segment's ports by display name. When used alone, searches all segment ports by display name.

~> **NOTE:** At least one of `vif_id`, `segment_path`, or `display_name` must be specified.

## Filter Priority

The data source applies filters in the following priority order:

1. **VIF ID** - If `vif_id` is set, it's used as the primary filter. Results can be further filtered by `segment_path` and/or `display_name`.
2. **Segment Path** - If `segment_path` is set (without `vif_id`), it filters ports by parent segment. Can be combined with `display_name`.
3. **Display Name** - If only `display_name` is set, it searches all segment ports by display name.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `items` - List of segment ports matching the filter criteria. Each item contains:
    * `id` - Unique identifier of the segment port.
    * `display_name` - Display name of the segment port.
    * `description` - Description of the segment port.
    * `path` - Policy path of the segment port.
    * `vif_id` - VIF attachment ID of the segment port (if attached).
    * `segment_path` - Parent segment path of the segment port.
