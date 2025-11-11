---
subcategory: "Segments"
page_title: "NSXT: nsxt_policy_segment_port"
description: Policy Segment Port data source.
---

# nsxt_policy_segment_port

This resource provides a method for the management of Segments Port.

## Example Usage

```hcl
data "nsxt_policy_segment_port" "segmentport1" {
  display_name = "segport1"
}
```

## Argument Reference

* `id` - (Optional) The ID of Segment Port to retrieve.
* `display_name` - (Optional) The Display Name prefix of the Segment to retrieve.
* `vif_id` - (Optional) Segment Port attachment id.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `description` - The description of the resource.
* `path` - The NSX path of the policy resource.
* `segment_path` - Path of the associated segment.
